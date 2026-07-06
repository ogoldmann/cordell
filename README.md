# Cordell

Cordell is a server-rendered Go web application for managing custody, checkout, and return of assets assigned to personnel designed for extreme clarity and practicality.

## Tech Direction

- Go
- PostgreSQL
- SQL-first persistence
- Server-rendered HTML
- Tailwind CSS
- HTMX
- Modular monolith architecture
- ULID-based identifiers

## Current Endpoint

```bash
GET /health
```

Expected response: 

```json
{"status":"ok","service":"cordell"}
```

## Running Locally

```bash
go run ./cmd/cordell
```

Then test:

```bash
curl http://localhost:8080/health
```

## Architecture Direction

The project is organized around explicit responsibilities:
- `cmd/cordell`: application entrypoint
- `internal/web`: HTTP server, routes, handlers and middleware
- `internal/domain`: core business concepts and rules
- `internal/app`: use cases and application services
- `internal/ports`: interfaces/contracts
- `internal/infra`: infrastructure implementations
- `internal/security`: security-related services
- `migrations`: database migrations
- `docs`: project documentation and architecture decisions

## Local Database

Start PostgreSQL:

```bash
make db-up
```

Run migrations:

```bash
make migrate-up
```

Check migration status:

```bash
make migrate-status
```

Stop PostgreSQL:

```bash
make db-down
```

## SQLC

Generate PostgreSQL query code:

```bash
make sqlc-generate
```

SQL queries live in:

```bash
internal/infra/postgres/queries
```

Generated code lives in:

```bash
internal/infra/postgres/db
```

## Running the Application

Start PostgreSQL:

```bash
make db-up
```

Run migrations:

```bash
make migrate-up
```

Start the application:

```bash
make run
```

Health check:

```bash
curl http://localhost:8080/health
```

Expected response:

```bash
{"status":"ok","service":"cordell"}
```

## Identifiers

Cordell uses ULIDs for application-generated identifiers.

ULIDs are stored as text and are used for records such as personnel, assets, and custody transactions.

## First Web Flow

Start the application:

```bash
make run
```

Open:

```bash
http://localhost:8080/personnel/new
```

This page creates a real personnel record in PostgreSQL and redirects to the created personnel detail page.
