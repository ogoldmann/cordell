# ADR 0005 - Application Composition

## Status

Accepted.

## Context

Cordell has separate packages for domain models, application use cases, ports, PostgreSQL repositories, and HTTP delivery.

The system needs a clear place where concrete dependencies are created and connected.

## Decision

The `cmd/cordell/main.go` file will act as the application composition root.

It is responsible for:

- loading configuration
- creating the logger
- opening the PostgreSQL connection pool
- creating sqlc queries
- creating PostgreSQL repositories
- creating application services
- creating the HTTP server
- starting the HTTP listener

## Consequences

### Positive

- Dependency wiring is explicit.
- Domain and application layers remain independent from PostgreSQL and HTTP details.
- The executable entrypoint clearly shows how the application is assembled.
- Future dependencies can be added in one visible place.

### Negative

- `main.go` can grow if composition is not kept organized.
- Care must be taken to keep business logic out of the composition root.

## Alternatives Considered

### Global dependencies

Global dependencies were avoided because they make tests and reasoning harder.

### Dependency injection framework

A dependency injection framework was avoided because the current project size does not justify the additional abstraction.
