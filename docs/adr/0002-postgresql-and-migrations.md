# ADR 0002 - PostgreSQL and Database Migrations

## Status

Accepted.

## Context

Cordell needs reliable persistence for personnel, assets, custody transactions, custody lines, and current custody balances.

The system must prioritize correctness, traceability, and future auditability.

## Decision

Cordell will use PostgreSQL as its database and goose for database migrations.

The initial schema uses string-based identifiers stored as `TEXT`, matching the current domain identifier types. This may be revisited later if the project adopts UUID-specific database types.

## Consequences

### Positive

- PostgreSQL supports relational integrity, constraints, indexes, and transactions.
- Database changes are versioned through migrations.
- The schema is explicit and reviewable.
- The project avoids ORM-generated hidden schema changes.

### Negative

- The project requires explicit SQL design.
- Schema evolution must be carefully managed.
- Application repositories must keep domain models and database rows consistent.

## Alternatives Considered

### SQLite

SQLite would be simpler for local development, but PostgreSQL better matches the intended multi-user intranet environment.

### ORM-managed migrations

ORM-managed migrations were avoided because Cordell favors explicit SQL and clear database control.