# ADR 0004: Authentication and RBAC

## Status

Accepted

## Context

Cordell is a server-rendered custody management system with administrative and operational workflows.

The system needs authentication before it can be safely used beyond local development. It also needs basic authorization so administrative functionality is not available to regular operators.

The project favors explicit domain rules, low coupling, SQL clarity, server-rendered simplicity, security, auditability, and incremental implementation.

## Decision

Cordell uses application-managed operator accounts.

Operators have:

- an ID
- a username
- a role
- a password hash
- an active flag

Current roles:

- `admin`
- `operator`

Authentication uses:

- Argon2id password hashing
- server-side sessions
- opaque random session tokens
- hashed session tokens stored in PostgreSQL
- HTTP-only session cookies
- CSRF tokens for authenticated unsafe methods

Authorization uses:

- middleware that reads the authenticated operator from request context
- admin-only middleware for administrative routes
- explicit domain role checks

Operator administration supports:

- listing operators
- viewing operator details
- creating operators
- deactivating operators
- reactivating operators
- changing roles
- resetting passwords

Operator deletion is intentionally not implemented. Operators are identities and should remain available for future audit history.

## Consequences

Benefits:

- password hashes are never exposed in admin read models
- session tokens are not stored in plaintext
- administrative routes are separated from regular workflows
- admin actions are protected by CSRF
- role changes and password resets invalidate existing sessions
- inactive operators cannot authenticate
- the system prevents deactivating or demoting the last active admin

Tradeoffs:

- login rate limiting is currently in-memory and suitable only for the current single-node architecture
- roles are fixed in code and constrained in SQL instead of being stored in a separate roles table
- operator management UI is intentionally simple and server-rendered
- no full audit log exists yet

## Follow-ups

Future milestones should add:

- append-only audit log
- tamper-evident audit chain
- audit events for operator administration
- finer-grained authorization if needed
- stronger production deployment guidance