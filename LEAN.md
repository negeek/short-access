# How it stays lean

Short-Access treats every dependency and every moving part as a cost — in memory,
image size, attack surface, and the time it takes someone to understand the code.
The whole design leans on that idea: do the job well with as little as possible,
so running a URL shortener never turns into running a platform. Here is what that
looks like in practice.

## One static binary, no runtime

The app compiles to a single, statically linked Go binary (`CGO_ENABLED=0`). There
is no interpreter, no VM, and no language runtime to install — the binary *is* the
program. The Docker image is a [multi-stage build](Dockerfile): a build stage
compiles it, and the final image is just a minimal Alpine base plus that one
binary, running as a non-root user. The result is a small image that starts
almost instantly, because starting it is just executing a file.

## Almost no dependencies

There is no web framework and no ORM. The entire direct dependency list is:

- **Go's standard library** — HTTP server, JSON, crypto, `log/slog` logging.
- **`gorilla/mux`** — routing, and nothing else.
- **`jackc/pgx`** — the Postgres driver and connection pool.
- **`golang-jwt`**, **`google/uuid`**, **`golang.org/x/crypto/bcrypt`** — tokens,
  ids, password hashing.

That's the whole tree. A small dependency set means a small binary, fast builds,
a smaller attack surface, and a codebase you can actually read end to end in an
afternoon.

## No framework tax on each request

Requests flow through plain `net/http` handlers over a thin layering of
handler → service → repository. There is no framework wrapping every request in
layers of reflection-driven middleware. Routing is a lightweight mux match; the
handler decodes the body, the service applies the rules, the repository runs a
query. Little happens on the way in and out that you didn't ask for.

## An efficient hot path

Redirects are the busy path in any shortener, so they are kept cheap. Following a
short link is a single indexed lookup by the unique `short_url`, then a redirect.
Counting the visit is one atomic `UPDATE ... access_count + 1` (and it's
best-effort — a counting error is logged, never blocks the redirect).

Creating short links is careful with the database too. Slugs come from a
sequential counter that reserves a *block* of numbers at a time and hands them out
from memory, so most shorten requests don't touch the counter table at all — see
the counter in [`service/v1/url/service.go`](service/v1/url/service.go). Database
connections are pooled by `pgx`, so there's no per-request connect cost.

## Migrations embedded in the binary

The schema migrations are baked into the binary with `embed.FS` (see
[`db/migrate.go`](db/migrate.go)). The same image both applies migrations
(`short-access migrate up`) and serves traffic — there is no separate migration
tool, image, or mounted file to manage. One artifact does everything.

## Stateless, so it scales sideways

The server keeps no session state. Auth is carried by JWTs and API keys, and all
real state lives in Postgres. That means you scale by running more identical
containers behind a load balancer, not by growing one big process. A single
container is plenty for most self-hosters; if you ever need more, you add
replicas rather than resources.

## What that adds up to

Small image, low memory, near-instant cold start, and few surprises. It sits
comfortably next to your existing app on a small VM or a free tier, rather than
demanding a box of its own. Leanness here isn't austerity for its own sake — it's
less to run, less to pay for, and less that can break.
