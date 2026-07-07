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

## Personnel Flow

List personnel:

```bash
http://localhost:8080/personnel
```

Create personnel:

```bash
http://localhost:8080/personnel/new
```

After creation, Cordell redirects to the created personnel detail page.

## Asset Flow

List assets:

```bash
http://localhost:8080/assets
```

Create asset:

```bash
http://localhost:8080/assets/new
```

After creation, Cordell redirects to the created asset detail page.

## Checkout Flow

Register checkout:

```bash
http://localhost:8080/custody/checkouts/new
```

The checkout flow assigns an asset quantity to personnel and redirects to the personnel detail page, where current custody balances are displayed.

## Return Flow

Register return:

```bash
http://localhost:8080/custody/returns/new
```

The return flow subtracts an asset quantity from a personnel custody balance.

Cordell prevents returns greater than the current custody quantity.

## Custody History

Personnel detail pages display:

- current custody balances
- custody transaction history
- checkout and return events
- asset quantities
- transaction notes

## Asset Current Holders

Asset detail pages display current custody holders.

This provides the inverse view of personnel current custody:

- Personnel page: assets currently assigned to that personnel
- Asset page: personnel currently holding that asset

## Dashboard

The dashboard is available at:

```bash
http://localhost:8080/
```

It provides quick access to core workflows:

- personnel
- assets
- checkout
- return

It also displays recent personnel and asset records.

## Personnel Profile

Personnel records include:

- full name
- alias
- rank
- registration ID
- section
- organization unit
- active status

`registration_id` is the conceptual personnel registration identifier.

The current validation mechanism accepts CPF-like identifiers, but the domain model does not expose a CPF-specific field name.

## Development Database Reset

During pre-release development, early migrations may still be edited to keep the baseline schema clean.

When an already-applied migration is edited, reset the local development database:

```bash
docker compose down -v
make db-up
make migrate-up
```

Do not use this process after real data exists.

## Search

Personnel and assets support basic server-rendered search.

Personnel search:

```bash
http://localhost:8080/personnel?q=doe
```

Personnel search matches:

- full name
- alias
- registration ID
- rank
- section
- organization unit

Asset search:

```bash
http://localhost:8080/assets?q=radio
```

Asset search matches:

- name