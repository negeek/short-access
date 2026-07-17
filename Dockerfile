# syntax=docker/dockerfile:1

# --- build stage ---
FROM golang:1.24-alpine AS build

WORKDIR /app

# Download modules first so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static binary. -s -w strip the symbol and debug tables to keep it small.
# Migrations are embedded, so the image needs nothing else.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /short-access ./cmd/short-access

# --- run stage ---
FROM alpine:3.20

# Run as a non-root user.
RUN adduser -D -u 10001 app
USER app

COPY --from=build /short-access /usr/local/bin/short-access

EXPOSE 8080

# Default is to serve. Override the command with "migrate up" to run migrations.
ENTRYPOINT ["short-access"]
