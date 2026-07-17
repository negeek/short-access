# Short-Access

A self-hosted URL shortener with a small HTTP API. You run it yourself with
Docker: point your app at it, create an API key, and start shortening links.

It is deliberately lean — Go standard library, `gorilla/mux` for routing, and
Postgres. No web framework.

## How auth works

There are two credentials:

- **JWT** — sent as `Authorization: Bearer <token>`. You get one when you sign up
  or log in. Required for managing your account and your API keys.
- **API keys** — sent as `X-API-Key: <key>`. You create a key once (with your
  JWT) for your application to use.

The **URL endpoints accept either** — an API key (typical for an app) or a JWT
(handy for a signed-in user). Account and key-management endpoints require the
JWT.

The usual flow is: **sign up → get a token → create an API key → use that key
from your app.**

## Interactive docs

Once it's running, browse the API at **`http://localhost:8080/docs`** (Swagger
UI). The raw OpenAPI spec is at `/openapi.yaml`.

## Run it (self-host)

You only need Docker. You do not build anything — you run the published image.

1. Grab the template compose file [`docker-compose.sample.yml`](docker-compose.sample.yml)
   and save it in your project as `docker-compose.yml`.
2. Create a `.env` next to it (see [`.env.example`](.env.example) for the full
   list):

   ```env
   POSTGRES_USER=sauser
   POSTGRES_PASSWORD=sapass
   POSTGRES_DB=sadb
   BASE_URL=http://localhost:8080
   AUTH_KEY=change-me-to-a-long-random-secret
   ```

3. Start it:

   ```bash
   docker compose up -d
   ```

   Postgres comes up, `short-access-migrate` runs the migrations and exits, then
   `short-access-engine` starts serving on `http://localhost:8080`.

### Image tags

Images are published to Docker Hub as `negeek/short-access`. `:latest` tracks
the main branch; released versions are tagged with semver (`:1.0.0`, `:1.0`,
`:1`). Pin a version in your compose file for anything you care about.

### Run without Docker

If you have Go, you can install the binary straight from the repo:

```bash
go install github.com/negeek/short-access/cmd/short-access@latest
```

You still need a Postgres. Set the same environment variables (see
[`.env.example`](.env.example)), apply migrations with `short-access migrate up`,
then run `short-access` to serve.

## API

Base URL below is `http://localhost:8080`.

### Sign up

```bash
curl -X POST 'http://localhost:8080/api/v1/user_mgt/join/' \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"a-good-password"}'
```

Returns an `access_token` (JWT).

### Get a new token

```bash
curl -X POST 'http://localhost:8080/api/v1/user_mgt/new_token/' \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"a-good-password"}'
```

### Create an API key (JWT required)

```bash
curl -X POST 'http://localhost:8080/api/v1/user_mgt/api_keys/' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"my app"}'
```

The response includes `api_key` — the raw key. **Copy it now; it is never shown
again.** You can also pass an optional `expire_at` (RFC 3339 timestamp); leave it
out for a key that never expires.

Manage keys:

- `GET /api/v1/user_mgt/api_keys/` — list your keys
- `POST /api/v1/user_mgt/api_keys/{id}/revoke` — revoke a key
- `DELETE /api/v1/user_mgt/api_keys/{id}` — delete a key

All URL endpoints below accept **either** `X-API-Key: <key>` **or**
`Authorization: Bearer <token>`.

### Shorten a URL

```bash
curl -X POST 'http://localhost:8080/api/v1/url_mgt/shorten/' \
  -H 'X-API-Key: <key>' \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://pkg.go.dev/net/http"}'
```

The response includes `short_url` (the slug) and `short_access` (the full link,
built from `BASE_URL`).

You can optionally set an expiry at creation with `time_unit` + `time_value`
(pass both or neither — see the units under [Set expiry](#set-expiry)):

```bash
curl -X POST 'http://localhost:8080/api/v1/url_mgt/shorten/' \
  -H 'X-API-Key: <key>' \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://pkg.go.dev/net/http","time_unit":"d","time_value":7}'
```

### Custom slug

```bash
curl -X POST 'http://localhost:8080/api/v1/url_mgt/custom/' \
  -H 'X-API-Key: <key>' \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp"}'
```

Custom URLs accept the same optional `time_unit` + `time_value`.

### Set expiry

`time_unit` is one of `y`, `mo`, `d`, `h`, `m`, `s` (year, month, day, hour,
minute, second). This one expires 40 seconds from now:

```bash
curl -X POST 'http://localhost:8080/api/v1/url_mgt/url_expiry/' \
  -H 'X-API-Key: <key>' \
  -H 'Content-Type: application/json' \
  -d '{"time_unit":"s","time_value":40,"url_id":1}'
```

### List / filter your URLs

Paginated with `limit` (default 20, max 100) and `offset`. You can also filter
by any column, e.g. `id` or `short_url`.

```bash
curl 'http://localhost:8080/api/v1/url_mgt/?limit=20&offset=0' -H 'X-API-Key: <key>'
curl 'http://localhost:8080/api/v1/url_mgt/?id=1&short_url=nethttp' -H 'X-API-Key: <key>'
```

The `data` is a page object: `items` (the URLs), plus `limit`, `offset`, `count`
and `has_more`.

### Update / delete a URL

```bash
curl -X PATCH  'http://localhost:8080/api/v1/url_mgt/1' -H 'X-API-Key: <key>' -H 'Content-Type: application/json' -d '{"original_url":"https://go.dev"}'
curl -X PUT    'http://localhost:8080/api/v1/url_mgt/1' -H 'X-API-Key: <key>' -H 'Content-Type: application/json' -d '{"original_url":"https://go.dev"}'
curl -X DELETE 'http://localhost:8080/api/v1/url_mgt/1' -H 'X-API-Key: <key>'
```

### Follow a short link

```bash
curl -i 'http://localhost:8080/<slug>'
```

Redirects to the original URL and counts the visit.

### Health check

```bash
curl -i 'http://localhost:8080/healthz'
```

Returns `200` when the service can reach its database, `503` otherwise. The
container healthcheck uses this.

## Local development

You need Go and a Postgres you can point at. Common tasks are in the
[`Makefile`](Makefile):

```bash
make run          # run the server
make format       # format the code
make test         # run tests
make migrate-up   # apply migrations
make migrate-down # roll back the last migration
make docker-up    # build and run the whole stack locally
```

Run `make help` to see everything.

### Tests

Tests run against a real database. The compose file ships a throwaway one behind
a `test` profile, so the easiest path is:

```bash
make test-db-up        # start the throwaway Postgres (host port 5444)
make test-integration  # run the suite against it
make test-db-down      # tear it down when you're done
```

`make test` on its own runs everything too, but the database-backed tests skip
themselves unless `TEST_DATABASE_URL` is set. The harness drops and recreates
the schema, runs migrations, and clears the tables between tests.
