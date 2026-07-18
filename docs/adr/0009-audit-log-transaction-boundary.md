# ADR 0009: Audit Log Transaction Boundary

## Status

Accepted for current development stage

## Context

Cordell stores application-level audit events for high-value actions such as operator administration and custody transaction creation.

The audit log currently records events after the main application action succeeds.

Examples:

- an operator is created, then `operator.created` is recorded
- a checkout is registered, then `custody.checkout_created` is recorded
- a return is registered, then `custody.return_created` is recorded

This means the main action and the audit event are not currently part of the same database transaction boundary.

If the main action succeeds but audit recording fails, Cordell logs the audit failure as a server-side error and still returns the successful result to the user.

This avoids a misleading user-facing failure after the action has already been persisted.

The project currently favors:

- explicit domain rules
- low coupling
- modular monolith
- SQL clarity
- server-rendered simplicity
- security and auditability
- maintainability
- correctness
- incremental verified milestones
- no premature complexity

## Decision

Cordell will keep audit event recording as best-effort for now.

The current behavior is:

1. Execute the main application action.
2. If the action fails, do not record a success audit event.
3. If the action succeeds, attempt to record the audit event.
4. If the audit event fails after the action already succeeded, log the failure on the server.
5. Do not return a misleading failure to the user after the main action already completed.

The audit log table remains append-only at the database level.

The system rejects:

- `UPDATE` on `audit_events`
- `DELETE` on `audit_events`
- `TRUNCATE` on `audit_events`

Audit event metadata must remain minimal and must not include:

- plaintext passwords
- password hashes
- session tokens
- CSRF tokens
- cookies
- raw form payloads
- secrets

## Alternatives Considered

### Option A: Keep best-effort audit logging

The application action and audit event are written separately.

Benefits:

- simplest option
- keeps current architecture intact
- avoids coupling application services to PostgreSQL transaction mechanics
- avoids premature transaction manager abstraction
- easy to understand and maintain
- sufficient for current local/pre-release development stage

Drawbacks:

- action and audit event are not perfectly atomic
- audit write failure can leave a successful action without a corresponding audit event
- requires server logs to detect audit write failures

Decision:

Accepted for now.

### Option B: Application transaction manager

Introduce an application-level transaction boundary so the main action and audit event are committed together.

Possible shape:

```txt
WithTransaction(ctx, func(ctx context.Context, tx Repositories) error {
    perform main action
    record audit event
})
```

Benefits:

- main action and audit event can commit atomically
- improves correctness for audit-critical operations
- keeps domain events in application code instead of database triggers

Drawbacks:

- introduces a new abstraction across app and infra
- requires transactional repository variants or context-bound repositories
- touches many services and repositories
- increases architectural complexity
- can reduce clarity if introduced too early

Decision:

Deferred.

This is the most likely future direction when Cordell needs stronger production-grade audit guarantees.

### Option C: Outbox pattern

Write an outbox record in the same transaction as the main action, then process it separately into audit events.

Benefits:

- strong reliability pattern
- useful for asynchronous processing
- useful if audit events later need to be published externally

Drawbacks:

- overkill for the current monolithic, server-rendered architecture
- requires worker/process management
- introduces retry semantics
- introduces extra operational complexity
- unnecessary while everything runs inside one application and one database

Decision:

Rejected for now.

This can be reconsidered if Cordell later integrates external systems or background workers.

### Option D: Database triggers for audit events

Use database triggers to automatically write audit events when tables change.

Benefits:

- very strong database-level consistency
- audit event can be written in the same database transaction
- difficult for application code to forget audit writes

Drawbacks:

- pushes domain semantics into SQL triggers
- requires passing actor context into the database session
- can make audit behavior less visible from application code
- harder to test from the app layer
- can become brittle as domain language evolves
- not all audit events are simple row-level changes

Decision:

Rejected for now.

Cordell should keep audit semantics explicit in application code.

### Consequences

The current audit log is useful for traceability but is not yet perfectly atomic with the main action.

This is acceptable for the current pre-release stage.

If an audit write fails after a successful action, the application logs the failure. The user is not shown a misleading failure response for an action that already succeeded.

The system must continue to treat audit write failures as important server errors.

Before production or real operational deployment, Cordell should revisit this ADR and likely implement stronger transaction boundaries.

### Future Work

Future audit hardening may include:

- application transaction manager
- transactional repository boundaries
- domain event recording inside the same transaction
- canonical audit event payloads
- tamper-evident hash chain
- audit event verification command
- stronger database permissions around audit tables
- operational alerting for failed audit writes