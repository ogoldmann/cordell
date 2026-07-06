# ADR 0007 - Server-Rendered HTML

## Status

Accepted.

## Context

Cordell is an internal web application with practical workflows, forms, dashboards, and simple interactions.

The project should avoid unnecessary frontend complexity while keeping the UI maintainable and progressively enhanceable.

## Decision

Cordell will initially use server-rendered HTML with Go's `html/template`.

Templates are embedded into the Go binary using `embed`.

Static assets are served from the `static` directory during development.

## Consequences

### Positive

- Keeps the frontend simple.
- Avoids a separate SPA build pipeline at this stage.
- Keeps form-based workflows straightforward.
- Fits well with future HTMX progressive enhancement.
- Allows the application to work with minimal JavaScript.

### Negative

- Rich client-side interactions require either HTMX or custom JavaScript.
- Template organization must be kept disciplined as pages grow.
- Static asset handling may need to be revisited for production packaging.

## Alternatives Considered

### SPA frontend

A SPA frontend was avoided because it would add significant complexity before the core workflows are stable.

### Laravel Blade-style full-stack framework

Cordell is a Go project, so Go's `html/template` is the initial server-side rendering tool.
