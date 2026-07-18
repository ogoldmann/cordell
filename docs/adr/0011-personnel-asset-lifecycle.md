# ADR 0011: Personnel and Asset Lifecycle

## Status

Accepted for future implementation

## Context

Cordell manages custody and return of assets linked to personnel records.

Personnel and asset records are operational records, but they are also part of historical custody evidence.

A military unit may have significant yearly turnover. Temporary military personnel may leave the unit after a short period. Assets may also leave regular operational use, become unavailable, be replaced, or stop being used for new checkout workflows.

The system needs a way to prevent old records from polluting daily operational screens while preserving custody history, receipts, and auditability.

The current model already includes an `active` flag for personnel and assets.

Cordell currently follows these principles:

* custody transactions are historical records
* existing custody transactions should not be edited or deleted
* custody receipts should remain readable
* operators are deactivated instead of deleted
* checkout creates new custody responsibility
* return reduces or clears existing custody responsibility
* auditability is more important than hiding history

The project favors:

* explicit domain rules
* low coupling
* modular monolith
* SQL clarity
* server-rendered simplicity
* security and auditability
* maintainability
* correctness
* incremental verified milestones
* no premature complexity

## Decision

Cordell will not hard-delete personnel or asset records through normal application workflows.

Personnel and assets will use lifecycle deactivation/reactivation.

The existing `active` flag represents whether a record is currently available for normal operational workflows.

A deactivated record is not deleted. It is archived from daily operational use.

Default operational screens should avoid pollution by showing active records by default.

Historical screens, receipts, audit events, and custody history must remain able to display deactivated records when they are part of historical data.

## Lifecycle Language

The application may use user-facing language such as:

* active
* inactive
* deactivated
* archived

The database currently uses:

```txt
active BOOLEAN
```

This is sufficient for the current stage.

A future migration may replace or complement `active` with a richer lifecycle status if needed, such as:

```txt
active
inactive
transferred
discharged
retired
lost
disposed
```

This is deferred until there is a proven operational need.

## Personnel Lifecycle Rules

Personnel records should not be hard-deleted.

A personnel record may be deactivated when the person should no longer appear in daily operational flows.

Examples:

* temporary military personnel left the unit
* personnel transferred
* personnel should no longer receive new custody
* duplicate or obsolete record was resolved

Deactivated personnel:

* must not appear in normal checkout selectors
* must not receive new checkout
* may still appear in historical custody records
* may still appear in custody receipts
* may still appear in audit-related views
* may be searchable through explicit inactive/archive filters
* may be reactivated if needed

Initial implementation should prevent deactivation of personnel with current custody balance greater than zero.

Reason:

* hidden pending custody is operationally dangerous
* deactivation should not hide unresolved responsibility
* outstanding material should be returned or corrected before deactivation

Return behavior remains different from checkout.

If an inactive personnel record has current custody due to legacy data or exceptional future workflow, return should still be allowed to clear the responsibility.

## Asset Lifecycle Rules

Asset records should not be hard-deleted.

An asset may be deactivated when it should no longer be issued in new checkout workflows.

Examples:

* asset removed from regular use
* asset replaced
* asset no longer issued
* asset record is obsolete
* asset should be preserved only for history

Deactivated assets:

* must not appear in normal checkout selectors
* must not be used for new checkout
* may still appear in historical custody records
* may still appear in custody receipts
* may still appear in current custody if it was already checked out before deactivation
* may be searchable through explicit inactive/archive filters
* may be reactivated if needed

Initial implementation should prevent deactivation of assets with current custody balance greater than zero.

Reason:

* a currently custodied asset should remain operationally visible until returned or corrected
* deactivation should not hide unresolved custody
* correction workflows are not implemented yet

Return behavior remains different from checkout.

If an inactive asset has current custody due to legacy data or exceptional future workflow, return should still be allowed to clear the responsibility.

## Default Visibility Rules

To avoid operational pollution, normal screens should prefer active records by default.

Default behavior:

```txt
personnel list        -> active records by default
asset list            -> active records by default
checkout form         -> active personnel and active assets only
return form           -> current custody items only
custody history       -> includes historical records regardless of active state
custody receipts      -> include historical records regardless of active state
admin/lifecycle views -> can show active, inactive, or all records
```

A future implementation may add filters:

```txt
status=active
status=inactive
status=all
```

or UI controls:

```txt
Active
Inactive
All
```

## Search Direction

Search should not permanently hide inactive records.

However, to avoid pollution, broad operational browsing should favor active records.

Possible future behavior:

* default list pages show active only
* search results clearly mark inactive records
* admin views can include inactive records
* exact registration ID search may return inactive records with a visible inactive badge

This ADR does not require immediate search changes. It defines the direction.

## Deactivation Authorization

Initial lifecycle actions should likely be admin-only.

Regular operators may view records and register normal custody operations, but deactivation/reactivation changes the operational availability of records and should be restricted.

Future policy may allow more granular permissions if needed.

## Audit Direction

Future lifecycle implementation should record audit events such as:

```txt
personnel.deactivated
personnel.reactivated
asset.deactivated
asset.reactivated
```

Audit metadata may include:

* target record ID
* reason
* whether current custody existed
* previous active state
* new active state

Audit metadata must not include secrets, session tokens, CSRF tokens, password hashes, or raw form payloads.

## Reason Requirement

Deactivation should require a reason.

Reason examples:

* transferred
* temporary service ended
* duplicate record
* asset removed from use
* asset replaced
* administrative correction

The reason should be stored in an audit event at minimum.

A future schema may add explicit lifecycle event tables if richer reporting is needed.

## Alternatives Considered

### Option A: Hard delete

Delete personnel or asset rows from the database.

Benefits:

* removes clutter
* simple mental model for users
* reduces visible records

Drawbacks:

* breaks historical traceability
* risks breaking receipts
* conflicts with custody history
* conflicts with auditability
* may violate operational evidence requirements
* makes past transactions harder to understand

Decision:

Rejected.

Hard delete should not exist as a normal application workflow.

### Option B: Soft delete with `deleted_at`

Add a `deleted_at` timestamp and treat records as deleted but recoverable.

Benefits:

* common application pattern
* preserves rows
* easy to filter out from normal views

Drawbacks:

* "deleted" is the wrong domain language
* personnel who left the unit are not deleted people
* assets removed from issue are not necessarily deleted assets
* can confuse lifecycle with deletion
* does not express operational availability clearly

Decision:

Rejected for now.

The existing `active` flag better expresses operational availability.

A future `deactivated_at` or lifecycle event table may be considered if needed.

### Option C: Active/inactive lifecycle

Use `active` to decide whether records participate in normal operational workflows.

Benefits:

* already exists in the schema
* simple
* clear enough for current needs
* avoids clutter in daily screens
* preserves history
* supports reactivation
* matches current checkout guardrails

Drawbacks:

* does not capture detailed reasons by itself
* does not distinguish transfer, discharge, disposal, duplicate, or retirement
* may need richer lifecycle reporting later

Decision:

Accepted for current implementation.

### Option D: Rich status enum immediately

Replace `active` with a status enum.

Possible statuses:

```txt
active
inactive
transferred
discharged
disposed
lost
retired
duplicate
```

Benefits:

* more expressive
* better reporting
* richer operational semantics

Drawbacks:

* premature for current stage
* requires more UI and policy decisions
* can overfit to one unit's process too early
* increases migration and validation complexity

Decision:

Deferred.

Cordell should keep `active` for now and evolve only when real operational need appears.

## Consequences

Cordell will avoid polluting daily operational screens without deleting historical records.

The system will preserve custody history, receipts, and auditability.

Inactive personnel and assets will be hidden from normal checkout workflows but remain available where history requires them.

Before implementing deactivation UI, Cordell must ensure current custody is not accidentally hidden.

The initial policy should block deactivation when current custody exists.

## Future Work

Future milestones should include:

* personnel deactivation use case
* personnel reactivation use case
* asset deactivation use case
* asset reactivation use case
* admin-only lifecycle routes
* lifecycle reason field
* audit events for lifecycle actions
* list filters for active/inactive/all
* inactive badges in list/search/detail views
* tests preventing deactivation with current custody
* tests allowing return for inactive records when current custody exists
