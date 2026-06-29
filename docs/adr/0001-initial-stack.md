# ADR 0001 - Initial Stack

## Context

Cordell must be simple enough to be developed and maintained by a small team, but serious enough to support auditability, traceability, and secure internal usage.

The developer also wants the project to be useful as a portfolio project and as a learning opportunity for Go, web development, architecture, and backend fundamentals.

## Decision

The initial technical direction is:

- Go for the backend
- Server-rendered HTML for the frontend
- Tailwind CSS for styling
- HTMX for progressive enhancement
- PostgreSQL for production persistence
- SQL-first persistence, likely using sqlc later
- Modular monolith architecture

## Consequences

### Positive

- Keeps the system simple.
- Avoids unnecessary frontend complexity.
- Makes HTTP, SQL, and backend architecture explicit.
- Supports a clean portfolio narrative.
- Allows gradual evolution.

### Negative

- Requires more explicit decisions than a full-stack framework.
- Requires careful implementation of authentication, sessions, CSRF, and authorization.
- Requires discipline to keep handlers, domain, application, and infrastructure separated.

## Alternatives Considered

### Laravel

Laravel would provide more built-in infrastructure, but it would shift the learning focus toward the framework and PHP ecosystem.

### SPA Frontend

A SPA frontend would add unnecessary complexity for the initial version of this internal system.

### Microservices

Microservices are unnecessary for the current size and deployment context of the project.
