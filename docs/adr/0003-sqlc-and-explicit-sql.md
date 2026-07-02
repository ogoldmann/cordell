# ADR 0003 - SQLC and Explicit SQL

## Status

Accepted.

## Context

Cordell needs reliable PostgreSQL persistence while keeping database behavior explicit and reviewable.

The project intentionally avoids hiding important persistence behavior behind an ORM because custody transactions, balances, constraints, and auditability require clear SQL.

## Decision

Cordell will use sqlc to generate type-safe Go code from handwritten SQL queries.

PostgreSQL access will use pgx.

## Consequences

### Positive

- SQL remains explicit and reviewable.
- Generated Go code is type-safe.
- The project avoids ORM magic.
- Database behavior is easier to reason about.
- Query code remains close to the schema.

### Negative

- Developers must understand SQL.
- Schema changes may require query and generated code updates.
- The generated code should not be edited manually.

## Alternatives Considered

### GORM

GORM would provide ORM-style productivity, but it would hide some persistence details that Cordell needs to keep explicit.

### database/sql manually

Writing all scanning and query boilerplate manually would be explicit, but more repetitive and error-prone.

### Ent

Ent is powerful and type-safe, but it adds a larger framework-like layer than Cordell needs at this stage.