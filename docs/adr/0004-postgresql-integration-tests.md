# ADR 0004 - PostgreSQL Integration Tests

## Status

Accepted.

## Context

Cordell uses PostgreSQL repositories implemented with pgx and sqlc. Unit tests with fake repositories validate application use cases, but they do not prove that SQL queries, migrations, repositories, and database constraints work together.

## Decision

Cordell will include PostgreSQL integration tests for persistence behavior.

Integration tests will be opt-in through the `CORDELL_INTEGRATION_TESTS=1` environment variable and will use the local PostgreSQL database configured by `CORDELL_DATABASE_URL`.

## Consequences

### Positive

- Repository behavior is tested against a real PostgreSQL database.
- SQL queries and migrations are validated together.
- Custody balance updates are tested through real database state.
- The normal unit test workflow remains fast and independent from PostgreSQL.

### Negative

- Integration tests require a running PostgreSQL instance.
- Test cleanup must be handled carefully.
- The local test database must not contain real or important data.

## Alternatives Considered

### Only unit tests with fake repositories

This would be faster, but would not validate SQL, constraints, transactions, or generated query code.

### Testcontainers

Testcontainers would provide stronger isolation, but it adds complexity and an extra dependency. It may be reconsidered later.

### Running integration tests by default

This was avoided because regular development tests should not require PostgreSQL to be running.
    