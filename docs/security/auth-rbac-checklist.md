# Auth and RBAC Manual Checklist

Use this checklist after changes to authentication, sessions, CSRF, operators, or admin routes.

## Bootstrap

- [ ] Start database.
- [ ] Run migrations.
- [ ] Create an admin operator through `cordell-admin`.
- [ ] Confirm admin can log in.
- [ ] Confirm regular operator can be created.

## Authentication

- [ ] Valid registration ID and password logs in successfully.
- [ ] Invalid registration ID returns a generic invalid credentials message.
- [ ] Invalid password returns a generic invalid credentials message.
- [ ] Inactive operator cannot log in.
- [ ] Already authenticated operator visiting `/login` is redirected to `/`.
- [ ] Requested private URL is preserved through `return_to`.
- [ ] External `return_to` URLs are rejected.

## Sessions

- [ ] Login creates a `cordell_session` cookie.
- [ ] Cookie is `HttpOnly`.
- [ ] Cookie uses `SameSite=Lax`.
- [ ] Cookie `Secure` flag follows environment config.
- [ ] Logout deletes the session.
- [ ] Expired sessions are removed.
- [ ] Deactivated operators lose access.
- [ ] Password reset invalidates target operator sessions.
- [ ] Role change invalidates target operator sessions.
- [ ] Reactivation removes stale target operator sessions.

## CSRF

- [ ] Authenticated POST without CSRF token returns `403 Forbidden`.
- [ ] Authenticated POST with invalid CSRF token returns `403 Forbidden`.
- [ ] Authenticated POST with valid CSRF token works.
- [ ] GET routes do not require CSRF token.
- [ ] Login route remains public.

## Authorization

- [ ] Admin can access `/admin`.
- [ ] Regular operator receives `403 Forbidden` on `/admin`.
- [ ] Anonymous user is redirected to `/login` for private GET routes.
- [ ] Admin link is visible only to admins.
- [ ] Regular private routes remain available to regular operators.

## Operator Administration

- [ ] Admin can list operators.
- [ ] Admin can view operator detail.
- [ ] Admin can create operator.
- [ ] Admin can deactivate another active operator.
- [ ] Admin cannot deactivate self.
- [ ] Admin cannot deactivate the last active admin.
- [ ] Admin can reactivate inactive operator.
- [ ] Admin can change another operator role.
- [ ] Admin cannot change own role.
- [ ] Admin cannot demote the last active admin.
- [ ] Admin can reset another active operator password.
- [ ] Admin cannot reset own password through admin action.
- [ ] Password hashes never appear in UI.

## Headers and Cache

- [ ] Private authenticated routes return `Cache-Control: no-store`.
- [ ] Security headers are present.
- [ ] HSTS is disabled in local HTTP development.
- [ ] HSTS is enabled only behind HTTPS.

## Data Exposure

- [ ] Operator list does not expose `password_hash`.
- [ ] Operator detail does not expose `password_hash`.
- [ ] Logs do not print plaintext passwords.
- [ ] Logs do not print session tokens.
- [ ] Logs do not print CSRF tokens.

## Custody Attribution

- [ ] Checkout history shows the operator that registered the transaction.
- [ ] Return history shows the operator that registered the transaction.
- [ ] Checkout cannot be registered without an authenticated operator.
- [ ] Return cannot be registered without an authenticated operator.
- [ ] Custody transaction rows include `operator_id`.

## Custody Receipts

- [ ] Checkout history links to a receipt page.
- [ ] Return history links to a receipt page.
- [ ] Receipt shows personnel.
- [ ] Receipt shows registering operator.
- [ ] Receipt shows assets and quantities.
- [ ] Receipt does not expose password hashes, session tokens, or CSRF tokens.

## Audit Events

- [ ] Creating an operator records `operator.created`.
- [ ] Deactivating an operator records `operator.deactivated`.
- [ ] Reactivating an operator records `operator.reactivated`.
- [ ] Changing an operator role records `operator.role_changed`.
- [ ] Resetting an operator password records `operator.password_reset`.
- [ ] Creating a checkout records `custody.checkout_created`.
- [ ] Creating a return records `custody.return_created`.
- [ ] Audit events do not include passwords.
- [ ] Audit events do not include password hashes.
- [ ] Audit events do not include session tokens.
- [ ] Audit events do not include CSRF tokens.

## Audit Hardening

- [ ] `audit_events` rejects `UPDATE`.
- [ ] `audit_events` rejects `DELETE`.
- [ ] `audit_events` rejects `TRUNCATE`.
- [ ] `audit_events.metadata` only accepts JSON objects.
- [ ] Audit metadata does not contain passwords.
- [ ] Audit metadata does not contain password hashes.
- [ ] Audit metadata does not contain session tokens.
- [ ] Audit metadata does not contain CSRF tokens.
- [ ] A failed audit write after a successful action is logged as a server error.