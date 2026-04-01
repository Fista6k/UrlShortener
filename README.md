# URL Shortener

URL shortener service built with Go and Gin.  
The application stores links in Postgres, uses Redis as cache, and applies DB migrations automatically on startup.

## Features

- Shortens links via `POST /shorten`
- Redirects by short code via `GET /:shortUrl`
- Main page at `GET /`
- Redis cache for fast short-code lookup
- Postgres as source of truth
- Built-in rate limiting middleware

## Tech Stack

- Go
- Gin
- Postgres
- Redis
- Goose migrations (`go:embed`)
- Docker Compose

## Project Structure

- `cmd/` - application entrypoint
- `internal/service/` - business logic and handlers
- `internal/adapter/` - storage and migrations
- `internal/controller/http/` - router and middleware
- `static/` - HTML templates

## Prerequisites

- Go `1.25+`
- Docker + Docker Compose (recommended for full local run)

## Environment Variables

Create `.env` from `.env.example`:

```bash
cp .env.example .env
```

Minimal example:

```env
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_PORT=5432
DB_HOST=postgres

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=your_database_connstring
GOOSE_MIGRATION_DIR=./migrations

DB_NAME_REDIS=0
REDIS_HOST=redis
REDIS_PORT=6379

APP_PORT=8080
```

## Run with Docker (recommended)

```bash
docker compose up -d --build
```

Open: [http://localhost:8080](http://localhost:8080)

Stop:

```bash
docker compose down
```

## Run locally (without Docker app container)

1. Ensure Postgres and Redis are running and reachable.
2. Set `.env` values for your local environment.
3. Start app:

```bash
go run ./cmd
```

## Makefile Commands

This project includes `makefile` with useful shortcuts:

- `make help` - list all targets
- `make run` - run app locally
- `make test` - run all tests
- `make test-v` - run tests in verbose mode
- `make fmt` - format Go code
- `make vet` - run `go vet`
- `make tidy` - cleanup dependencies
- `make up` - start docker services
- `make down` - stop docker services
- `make logs` - tail docker logs
- `make rebuild` - rebuild and restart containers

## Testing

Run all tests:

```bash
go test ./...
```

Verbose mode:

```bash
go test ./... -v
```
