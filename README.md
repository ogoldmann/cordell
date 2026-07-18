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