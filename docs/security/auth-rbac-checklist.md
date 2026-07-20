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

## Custody Receipt Integrity

- [ ] Receipt shows transaction ID.
- [ ] Receipt shows transaction type.
- [ ] Receipt shows created timestamp.
- [ ] Receipt shows personnel identity.
- [ ] Receipt shows operator identity.
- [ ] Receipt shows asset lines and quantities.
- [ ] Receipt shows notes or empty notes state.
- [ ] Receipt remains readable after personnel deactivation.
- [ ] Receipt remains readable after asset deactivation.
- [ ] Receipt marks current Active/Inactive status clearly.
- [ ] Receipt explains that current status badges are not historical snapshots.

## Custody Correction Model

- [ ] Existing custody transactions are not edited for correction.
- [ ] Existing custody transactions are not deleted for correction.
- [ ] Correction implementation follows ADR 0006.
- [ ] Future correction events preserve operator attribution.
- [ ] Future correction events require a reason.
- [ ] Future correction events produce audit events.

## Custody Current State

- [ ] Checkout increases current custody balance.
- [ ] Return decreases current custody balance.
- [ ] Return cannot make current custody balance negative.
- [ ] Full return leaves current custody balance at zero.
- [ ] Zero-balance items are not shown in current custody views.
- [ ] Duplicate asset lines are consolidated before transaction creation.
- [ ] Invalid return rolls back transaction rows and line rows.

## Custody Concurrency

- [ ] Return balance decrease uses atomic conditional update.
- [ ] Return does not rely only on read-before-update validation.
- [ ] Concurrent returns cannot overdraw custody balance.
- [ ] Failed concurrent return does not save transaction rows.
- [ ] Failed concurrent return does not save line rows.
- [ ] Concurrent checkouts accumulate balance correctly.
- [ ] Repository/database remains the source of truth for custody balance.

## Custody Form Idempotency

- [ ] Checkout form includes generated transaction ID.
- [ ] Return form includes generated transaction ID.
- [ ] Register checkout service uses command transaction ID.
- [ ] Register return service uses command transaction ID.
- [ ] Duplicate custody transaction ID is treated as duplicate submit.
- [ ] Duplicate submit does not update custody balance twice.
- [ ] Duplicate submit does not create duplicate custody lines.
- [ ] Duplicate submit does not create duplicate audit events.
- [ ] Backend/database uniqueness is the source of truth.

## Return UI Guardrails

- [ ] Return form requires personnel selection first.
- [ ] Return form lists only assets currently under selected personnel custody.
- [ ] Return form shows available quantity.
- [ ] Return quantity input has `max` equal to available quantity.
- [ ] Personnel with no current custody shows an empty state.
- [ ] Backend still rejects return quantity greater than current balance.

## Checkout Active Record Guardrails

- [ ] Checkout form lists only active personnel.
- [ ] Checkout form lists only active assets.
- [ ] Backend rejects checkout for inactive personnel.
- [ ] Backend rejects checkout for inactive assets.
- [ ] Return still works for inactive personnel when current custody exists.
- [ ] Return still works for inactive assets when current custody exists.

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

## Audit Transaction Boundary

- [ ] Audit event recording behavior matches ADR 0005.w
- [ ] A successful main action does not return a misleading failure if audit recording fails afterward.
- [ ] Audit write failures are logged as server-side errors.
- [ ] Audit events are still treated as best-effort until a transaction boundary is implemented.
- [ ] Production deployment review includes revisiting audit transaction boundaries.

## Personnel and Asset Lifecycle

- [ ] Personnel and assets are not hard-deleted through normal workflows.
- [ ] Active records appear in normal operational workflows.
- [ ] Inactive records are hidden from checkout workflows.
- [ ] Inactive records remain visible in historical custody records.
- [ ] Inactive records remain visible in custody receipts.
- [ ] Checkout rejects inactive personnel.
- [ ] Checkout rejects inactive assets.
- [ ] Return can clear current custody involving inactive records.
- [ ] Future deactivation actions require admin authorization.
- [ ] Future deactivation actions require a reason.
- [ ] Future deactivation actions produce audit events.
- [ ] Future deactivation is blocked when current custody exists.

## Active/Inactive Filtering

- [ ] Personnel list defaults to active records.
- [ ] Asset list defaults to active records.
- [ ] Personnel list supports inactive filter.
- [ ] Asset list supports inactive filter.
- [ ] Personnel list supports all filter.
- [ ] Asset list supports all filter.
- [ ] Lists show Active/Inactive badges.
- [ ] Direct detail access still works for inactive records.
- [ ] Checkout selectors use active records only.

## Inactive Visibility

- [ ] Inactive personnel with current custody remains visible in personnel detail.
- [ ] Inactive personnel with current custody shows a warning.
- [ ] Asset detail marks inactive personnel holders.
- [ ] Current custody marks inactive assets.
- [ ] Global search marks inactive personnel.
- [ ] Global search marks inactive assets.
- [ ] Checkout backend rejects inactive personnel.
- [ ] Checkout backend rejects inactive assets.
- [ ] Return backend can clear custody involving inactive records.
- [ ] UI visibility rules do not replace backend validation.

## Asset Name Uniqueness

- [ ] Asset name is required.
- [ ] Asset name is normalized before saving.
- [ ] Duplicate asset names are rejected.
- [ ] Duplicate asset names are rejected case-insensitively.
- [ ] Duplicate normalized asset names are rejected.
- [ ] Duplicate asset names are rejected across active and inactive assets.
- [ ] Asset creation shows a human-readable duplicate-name error.

## Asset Lifecycle Actions

- [ ] Authenticated operator can deactivate active asset.
- [ ] Authenticated operator can reactivate inactive asset.
- [ ] Asset deactivation does not require admin role.
- [ ] Asset deactivation does not require reason.
- [ ] Asset deactivation does not delete custody history.
- [ ] Asset deactivation does not settle current custody.
- [ ] Inactive asset cannot be used in new checkout.
- [ ] Inactive asset can still be returned when current custody exists.
- [ ] Deactivation records `asset.deactivated`.
- [ ] Reactivation records `asset.reactivated`.
- [ ] Inactive asset appears under inactive filter.
- [ ] Reactivated asset appears under active filter.
- [ ] Deactivated asset name remains reserved.

## Dangerous Action Confirmation Modals

- [ ] Personnel deactivation uses a confirmation modal.
- [ ] Personnel reactivation uses a confirmation modal.
- [ ] Asset deactivation uses a confirmation modal.
- [ ] Asset reactivation uses a confirmation modal.
- [ ] Confirmation modals submit through POST.
- [ ] Confirmation POST requests include CSRF token.
- [ ] Direct POST without CSRF is rejected.
- [ ] Confirmation modals do not replace backend validation.
- [ ] No lifecycle `/confirm` routes remain.

## Global Search Inactive Awareness

- [ ] Global search defaults to active records.
- [ ] Global search supports inactive filter.
- [ ] Global search supports all filter.
- [ ] Search status filter is applied in backend queries.
- [ ] Search form preserves selected status filter.
- [ ] Personnel search results show Active/Inactive badge.
- [ ] Asset search results show Active/Inactive badge.
- [ ] Inactive records are not hidden when explicitly requested.

## Custody Transaction Edit Form

- [ ] Receipt exposes Edit checkout/Edit return action.
- [ ] Edit form loads original transaction.
- [ ] Edit form uses latest correction as current effective base.
- [ ] Edit form includes generated correction ID.
- [ ] Edit form allows corrected personnel.
- [ ] Edit form allows corrected asset lines.
- [ ] Edit form allows corrected quantities.
- [ ] Edit form allows corrected notes.
- [ ] POST creates append-only correction.
- [ ] POST does not mutate original custody transaction.
- [ ] POST applies balance deltas through correction service.
- [ ] POST records audit only when correction is newly created.
- [ ] POST redirects to original receipt.

## Custody Correction UX Hardening

- [ ] Edit form uses submit confirmation modal.
- [ ] Submit confirmation modal submits the original form.
- [ ] Submit confirmation modal preserves corrected personnel.
- [ ] Submit confirmation modal preserves corrected lines.
- [ ] Submit confirmation modal preserves corrected notes.
- [ ] Submit confirmation modal preserves correction ID.
- [ ] Receipt distinguishes original transaction from latest edit.
- [ ] Edit form explains whether base is original transaction or latest edit.
- [ ] Insufficient balance error explains possible later custody activity.

## Custody Receipt Effective State

- [ ] Receipt main cards show latest effective state.
- [ ] Original transaction remains visible in edit history.
- [ ] Receipt shows edit count.
- [ ] Receipt links to edit history.
- [ ] Edit history lists original transaction.
- [ ] Edit history lists all corrections chronologically.
- [ ] Edit history shows only changed fields.
- [ ] Edit metadata is visually separated from changed content.
- [ ] Latest correction is not rendered as a full duplicate receipt snapshot.

## Custody Correction Line UX

- [ ] Edit form does not render arbitrary blank line rows.
- [ ] Add line button appends one empty line.
- [ ] Remove button removes a line.
- [ ] At least one line remains required.
- [ ] Each visible line requires asset and quantity.
- [ ] Backend validates all submitted lines.
