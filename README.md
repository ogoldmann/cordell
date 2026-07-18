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
{ "status": "ok", "service": "cordell" }
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

## Custody Operator Attribution

Every custody transaction records the operator that registered it.

Custody transactions store:

- transaction type
- personnel
- operator
- notes
- timestamp
- custody lines

This means checkout and return history can answer:

- who received or returned the asset
- which operator registered the event
- when the event happened
- which assets and quantities were involved

Operator records are not deleted, so historical custody records keep their administrative attribution.

## Custody Receipts

Each custody transaction has a receipt page:

```txt
/custody/transactions/{id}
```

A custody receipt displays:

- transaction ID
- transaction type
- timestamp
- personnel
- operator that registered the transaction
- asset lines
- quantities
- notes

Receipts are read-only and generated from persisted custody transaction data.

## Custody Receipt

A custody receipt is a read-only view of a checkout or return transaction.

It shows who received or returned assets, which authenticated operator registered the event, which assets were involved, quantities, notes, and timestamp.

## Custody Corrections

Cordell treats custody history as a ledger-like record.

Existing custody transactions should not be edited or deleted to correct mistakes.

The correction model is documented in:

```txt
docs/adr/0006-custody-correction-model.md
```

The accepted direction is:

- original custody transactions remain immutable
- corrections are explicit new events
- corrections must preserve operator attribution
- corrections must require a reason
- corrections should generate receipts and audit events

Correction implementation is deferred until the balance rules and workflow are designed.

## Custody Current State

Cordell maintains current custody balances from custody transactions.

Current balance behavior:

- checkout increases the personnel custody balance for each asset
- return decreases the personnel custody balance for each asset
- return cannot reduce a balance below zero
- a full return leaves the balance at zero
- zero-balance items are omitted from current custody views
- duplicate asset lines in a checkout or return command are consolidated before the transaction is created

Custody transaction persistence is atomic inside the PostgreSQL repository. If balance update fails, transaction rows and line rows are rolled back.

## Return UI Guardrails

The return form is guided by current custody state.

Return workflow:

1. Select personnel.
2. Cordell lists only assets currently under that personnel's custody.
3. Each return card shows the available quantity.
4. The quantity input uses the available quantity as the HTML maximum.
5. The backend still validates the real current balance before saving.

The UI reduces operational mistakes, but the application service and repository remain the source of truth.

## Checkout Active Record Guardrails

Checkout creates a new custody responsibility.

Because of that, Cordell only allows checkout when:

- the personnel record is active
- the asset record is active
- the operator is authenticated and active

The checkout form lists only active personnel and active assets.

Return behaves differently. Return can still be registered for inactive personnel or inactive assets when there is current custody to clear.

This allows the system to close existing custody responsibilities even when a personnel or asset record becomes inactive later.

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

## Personnel and Asset Lifecycle

Cordell does not hard-delete personnel or asset records through normal application workflows.

Personnel and assets use an active/inactive lifecycle.

The accepted direction is documented in:

```txt
docs/adr/0007-personnel-asset-lifecycle.md
```

Lifecycle direction:

- active records appear in normal operational workflows
- inactive records are hidden from daily operational lists by default
- inactive records remain available for history, receipts, and auditability
- checkout requires active personnel and active assets
- return may still clear existing custody if inactive records are involved
- future deactivation/reactivation actions should be admin-only
- future deactivation should require a reason
- future deactivation should produce audit events
- initial lifecycle implementation should prevent deactivation while current custody exists

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

Personnel and assets support basic server-rendered tokenized search.

Personnel search:

```bash
http://localhost:8080/personnel?q=sergeant%20doe
```

Personnel search matches:

- full name
- alias
- registration ID
- rank
- section
- organization unit

Search queries are split into tokens. All tokens must match at least one searchable field in the same record.

Examples:

sergeant doe
operations doe
529.982 doe

Asset search:

```bash
http://localhost:8080/assets?q=radio%20battery
```

Asset search matches:

- name

Asset search also uses tokenized matching, so all query tokens must match the asset name.

## Global Search

Cordell supports global server-rendered search through:

```bash
/search?q=doe
```

Global search currently includes:

- personnel
- assets

Results are grouped by record type.

Global search uses the same tokenized search behavior as personnel and asset listing pages.

## CSS

Cordell uses Tailwind CSS through the Tailwind CLI.

The source CSS file is:

```bash
static/css/input.css
```

The compiled CSS file served by the application is:

```bash
static/css/app.css
```

Build CSS:

```bash
make css-build
```

Watch CSS during UI development:

```bash
make css-watch
```

The compiled CSS is committed so the Go application can run without requiring a CSS build step at runtime.

## Tailwind Migration

Cordell has migrated the current server-rendered UI to Tailwind utilities.

Migrated areas:

- base layout
- dashboard
- personnel listing
- asset listing
- search forms
- personnel creation form
- asset creation form
- checkout form
- return form
- personnel detail page
- asset detail page
- current custody sections
- custody history timeline

The Tailwind source file should remain small and focused on:

- Tailwind import
- source registration
- dark mode variant
- design tokens

## Template Components

Cordell uses small server-rendered template components for repeated UI fragments.

Current shared components:

- `personnel_card`
- `asset_card`

Shared components live in:

```bash
internal/web/views/components.html
```

Template helper functions are registered by the web renderer before templates are parsed.

Current template helpers:

- queryEscape: escapes values for query string links.

## Theme

Cordell supports light and dark themes through a `data-theme` attribute on the root HTML element.

The current implementation stores the user's visual preference in the browser through `localStorage`.

No server-side user preference is persisted yet.

Theme behavior:

- `data-theme="light"` enables the light theme.
- `data-theme="dark"` enables the dark theme.
- When no preference is stored, the browser/system preference is used.
- The theme toggle is implemented in `static/js/theme.js`.

Server-side persistence may be added later when authentication and user sessions exist.

## Operators

Operators represent authenticated users of the system.

Each operator has:

- an internal ID
- a unique registration ID used for login
- an alias
- a rank
- a role
- a password hash
- an active flag

Current roles:

- `admin`: administrative operator
- `operator`: regular custody workflow operator

The current RBAC foundation includes:

- operator role value object
- role validation in the domain
- role persistence in PostgreSQL
- role support in the admin CLI
- layout identity display with `Signed in as` followed by rank and alias

Route-level authorization is implemented in a later milestone.

## Sessions

Cordell uses server-side operator sessions.

The browser receives an opaque session token in a cookie.

The database stores only a hash of the session token.

Session cookie attributes:

- `HttpOnly`
- `SameSite=Lax`
- `Secure` configurable through `CORDELL_SESSION_COOKIE_SECURE`

In local development, `CORDELL_SESSION_COOKIE_SECURE=false`.

In HTTPS production environments, set:

```env
CORDELL_SESSION_COOKIE_SECURE=true
```

Login, logout, and session creation are implemented.
Route protection and CSRF protection are implemented.

Expired sessions are deleted during login and when expired sessions are encountered while loading the current operator.

## Admin CLI

Cordell includes a local administrative CLI for bootstrap tasks.

The admin CLI reads the same environment configuration as the main server.

When running it manually, export `.env` variables first:

```bash
set -a
source .env
set +a
```

Create an admin operator:

```bash
go run ./cmd/cordell-admin create-operator -registration-id 52998224725 -alias silva -rank sergeant -role admin
```

Or through Make:

```bash
make admin-create-operator REGISTRATION_ID=52998224725 ALIAS=silva RANK=sergeant ROLE=admin
```

The Make target automatically loads .env when the file exists.

Create a regular operator:

```bash
make admin-create-operator REGISTRATION_ID=93541134780 ALIAS=costa RANK=corporal ROLE=operator
```

The command prompts for the password interactively and does not echo it in the terminal.

Operator creation is intentionally not exposed as a public registration flow.

## CSRF Protection

Cordell protects authenticated state-changing requests with CSRF tokens.

The current implementation uses a synchronizer token pattern:

- a CSRF token is generated when an operator session is created
- the token is stored server-side in the operator session
- private HTML forms include the token as a hidden field
- unsafe authenticated methods require a valid token
- CSRF tokens are not placed in URLs
- CSRF tokens are not stored in cookies

Safe methods such as `GET`, `HEAD`, and `OPTIONS` do not require CSRF tokens.

Public login CSRF protection may be revisited later with a pre-session flow if needed.

## Security Headers

Cordell sends baseline HTTP security headers.

Current headers:

- `X-Content-Type-Options: nosniff`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy`
- `Permissions-Policy`

HSTS is configurable through:

```env
CORDELL_ENABLE_HSTS=true
```

HSTS should only be enabled in HTTPS environments.

The current CSP allows inline scripts and inline styles because the UI still uses an inline theme bootstrap script and server-rendered Tailwind output. This may be tightened later.

## Authentication Hardening

Cordell includes basic authentication hardening:

- failed login rate limiting
- expired session cleanup during login
- redirect back to the originally requested private URL after login
- local-path validation for post-login redirects to prevent open redirects

The login rate limiter is currently in-memory and intended for the current single-node application architecture.

If Cordell is deployed across multiple instances later, login throttling should move to shared storage such as PostgreSQL or Redis.

## Authorization

Cordell includes a role-based authorization foundation.

Current roles:

- `admin`
- `operator`

Current authorization behavior:

- authenticated operators can access regular private routes
- only admins can access `/admin`
- regular operators receive `403 Forbidden` when accessing admin-only routes
- unauthenticated users are redirected to login for private `GET` routes

Authorization is implemented through web middleware using the authenticated operator stored in the request context.

Operator management pages and finer-grained permissions are implemented in later milestones.

## Operator Administration

Admins can access a read-only operator management page at:

```txt
/admin/operators
```

The operator list displays:

- registration ID
- rank
- alias
- role
- active status
- creation timestamp
- operator ID

Password hashes are intentionally excluded from this read model and are never displayed in the admin UI.

Operator creation, deactivation, and role changes are implemented in later milestones.

Admins can create operators from the web admin area: 

```
txt /admin/operators/new
```

Operator creation through the web admin area:

- requires an authenticated admin
- is protected by CSRF
- supports the admin and operator roles
- hashes passwords through the application service
- never displays password hashes
- defaults new web-created accounts to the operator role in the form

The local admin CLI remains available for bootstrap tasks.

Admins can deactivate operators from the web admin area.

Operator deactivation:

- marks the operator as inactive
- does not delete the operator record
- invalidates the operator's active sessions
- prevents deactivating the currently authenticated operator
- prevents deactivating the last active admin
- is protected by admin-only authorization
- is protected by CSRF

Inactive operators cannot authenticate.

Admins can change operator roles from the web admin area.

Operator role changes:

- require an authenticated admin
- are protected by CSRF
- cannot target the currently authenticated operator
- cannot demote the last active admin
- invalidate active sessions for the changed operator
- support the `admin` and `operator` roles

This keeps role changes effective immediately and avoids stale privileged sessions.

Admins can reset operator passwords from the web admin area.

Password reset:

- requires an authenticated admin
- is protected by CSRF
- cannot target the currently authenticated operator
- requires password confirmation
- uses the application password hashing service
- invalidates active sessions for the changed operator
- never displays password hashes
- only applies to active operators

This is an administrative reset flow, not a self-service password change flow.

The operator administration UI separates listing from sensitive actions.

Routes:

```txt
/admin/operators
/admin/operators/{id}
```

The operators index is a read-oriented table.

Sensitive operator actions are handled on the operator detail page:

- role change
- password reset
- deactivation

This keeps the admin list readable and gives each sensitive action more context.

Admins can reactivate inactive operators.

Operator reactivation:

- marks the operator as active
- does not change the operator password
- does not create a session automatically
- removes any stale sessions for the operator
- requires an authenticated admin
- is protected by CSRF

Reactivated operators can authenticate again using their current password.

## Audit Events

Cordell stores application-level audit events for high-value actions.

Current audited events include:

- operator creation
- operator deactivation
- operator reactivation
- operator role changes
- operator password resets
- checkout creation
- return creation

Audit events are append-only application records and are visible to admins at:

```txt
/admin/audit-events
```

Audit events must not include passwords, password hashes, session tokens, CSRF tokens, or raw form payloads.

Audit events are append-only at the database level.

The `audit_events` table rejects:

- updates
- deletes
- truncation

This protects against accidental mutation and reinforces the audit model. It is not yet a cryptographic tamper-evident chain.

Audit recording currently happens after the main application action succeeds. If audit recording fails after the action has already completed, Cordell logs the audit failure as a server error instead of returning a misleading failure to the user.

Future audit hardening should add stronger transaction boundaries and tamper-evident hashing.

### Audit transaction boundary

Audit event recording is currently best-effort.

Cordell first executes the main application action. If the action succeeds, Cordell attempts to record an audit event.

If audit recording fails after the main action has already succeeded, Cordell logs the audit failure as a server-side error instead of returning a misleading failure to the user.

This is documented in:

```txt
docs/adr/0005-audit-log-transaction-boundary.md
```

Future production hardening should revisit this decision and likely introduce an application transaction boundary for audit-critical operations.