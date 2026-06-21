# syntax=docker/dockerfile:1

# ---- Build stage ----------------------------------------------------------
# CGO must be enabled because the SQLite driver (mattn/go-sqlite3) compiles
# C code, so we build on a Debian-based image that has a C toolchain.
FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /todo-api .

# ---- Runtime stage --------------------------------------------------------
# debian-slim has glibc (needed by the CGO binary) and we add CA certs for TLS.
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /todo-api /app/todo-api
# config.yaml is gitignored (holds local DB creds) and unused in prod —
# the container loads config.production.yaml via CONFIG_FILE (see fly.toml).
COPY config.production.yaml /app/config.production.yaml

# The SQLite file lives on a mounted volume; see fly.toml + SQLITE_PATH.
EXPOSE 8080

CMD ["/app/todo-api"]
