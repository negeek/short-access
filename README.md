# Short-Access

Short-Access is a self-hostable URL shortener with an easy-to-use API.
You run it on your own machine or server, point your application at it, and it
hands back short links — with expiry, custom slugs, visit counting, and API keys.

It is built to be lean, and that is the whole point: wanting to shorten a few
URLs shouldn't mean standing up a heavy service that eats into your
infrastructure. Staying lean keeps the compiled binary small, startup close to
instant, and memory use low, so a single modest container — the whole image comes
in under 20 MB — comfortably serves a lot of redirects without asking much of the
box it runs on. Curious how it pulls that off? See [how it stays lean](LEAN.md).

## How auth works

There are two kinds of credential:

- **JWT** — sent as `Authorization: Bearer <token>`. You receive one when you
  sign up or log in. It is what you use to manage your account and your API keys.
- **API keys** — sent as `X-API-Key: <key>`. You create a key once (using your
  JWT) and hand it to your application.

The **action endpoints** (everything under `/url_mgt`) accept **either**
credential, so you can call them with a JWT or an API key. In practice, **prefer
an API key**: unlike the JWT, you can give a key its own expiry, and you can
revoke it at any time without touching your password — which is exactly what you
want for something living inside an application.

So the usual journey is: **sign up → get a token → create an API key → use that
key from your app.**

## Self-host guide

You only need Docker. You are not building anything from source here — you run
the image that's already published.

**1. Get the compose file.** Grab
[`docker-compose.sample.yml`](docker-compose.sample.yml) and save it in your
project as `docker-compose.yml`. It pulls the published image and wires up
Postgres, a one-shot migration step, and the server.

**2. Create a `.env` next to it.** See [`.env.example`](.env.example) for the
full list. The important ones:

```env
POSTGRES_USER=sauser
POSTGRES_PASSWORD=sapass
POSTGRES_DB=sadb
BASE_URL=http://localhost:8080
AUTH_KEY=change-me-to-a-long-random-secret
```

- **`BASE_URL`** is the public address people will reach your instance at. It is
  used to build the *full* short link that gets returned to you — the slug is
  tacked onto the end of it (so `BASE_URL` of `https://sho.rt` turns slug `abc123`
  into `https://sho.rt/abc123`). Set it to wherever your instance actually lives.
- **`AUTH_KEY`** is the secret used to sign your JWT tokens. Make it a long,
  random string and keep it private — anyone who has it can mint valid tokens.
  If you ever change it, existing tokens stop working (everyone re-logs in).

**3. Start it.**

```bash
docker compose up -d
```

Postgres comes up, the migration container runs and exits, and the server starts
serving on `http://localhost:8080`. Open **`http://localhost:8080/docs`** to
browse and try the API right away (interactive Swagger docs).

### Image tags

Images are published to Docker Hub as `negeek/short-access`. `:latest` follows
the main branch; releases are tagged with semver (`:1.0.0`, `:1.0`, `:1`). Pin a
version in your compose file for anything you care about keeping stable.

### Installing with Go (no Docker)

If you already have Go, you can skip Docker entirely and install the binary
straight from the source:

```bash
go install github.com/negeek/short-access/cmd/short-access@latest
```

Here's what that actually does, step by step. `go install` reads that import
path, downloads this project's code from GitHub, compiles the `short-access`
command, and drops the finished binary into your Go bin directory
(`$(go env GOPATH)/bin`, usually `~/go/bin`). If that directory is on your
`PATH`, you can then just type `short-access` to run it. The `@latest` on the end
tells Go to use the newest released version — that is, the highest `vX.Y.Z` git
tag on the repo; you could pin a specific one instead, like `@v1.0.0`.

You still need a Postgres for it to talk to. Set the same environment variables
as above, apply the schema once with `short-access migrate up`, then run
`short-access` to serve. Once it's up, the interactive API docs are at
**`http://localhost:8080/docs`**.

## What the API can do

Base URL below is `http://localhost:8080`. Here is the whole onboarding flow you
need to get started — sign up, create a key, and check the service is healthy:

```bash
# 1. Sign up -> returns a JWT in "access_token"
curl -X POST 'http://localhost:8080/api/v1/user_mgt/join/' \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"a-good-password"}'

# 2. Create an API key with that JWT -> returns the raw key ONCE, save it
curl -X POST 'http://localhost:8080/api/v1/user_mgt/api_keys/' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"my app"}'

# 3. Health check -> 200 when the service and its database are up
curl -i 'http://localhost:8080/healthz'
```

From there, using your API key, the URL API lets you:

- **Shorten a URL** and get an auto-generated slug back.
- **Choose a custom slug** instead of the generated one.
- **Set an expiry** so a link stops working after a chosen amount of time
  (seconds through years) — either at creation or later.
- **List and filter** your links, paginated, filtering by fields like id or slug.
- **Update or delete** a link.
- **Follow a link** — visiting the short URL redirects to the original and counts
  the visit.

And with your JWT you can **manage API keys** — list them, revoke one, or delete
one.

You don't have to memorize request shapes: the running service ships **interactive
API docs**. Open **`http://localhost:8080/docs`** in a browser to try every
endpoint (Swagger UI), or read the raw spec at `/openapi.yaml`. That is the
source of truth for exact fields, parameters and responses — this README just
gets you moving.

## Local development

Most people should just run the Docker image; this section is only if you want to
hack on the code.

```bash
make run     # run the server locally
make test    # run tests (see below)
make help    # list every available command
```

Tests run against a real database. The compose file includes a throwaway one
behind a `test` profile:

```bash
make test-db-up        # start the throwaway Postgres
make test-integration  # run the suite against it
make test-db-down      # tear it down
```
