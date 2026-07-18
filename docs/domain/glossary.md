# Domain Glossary

## Personnel

A person who can receive assets under custody.

In the user interface, this may be presented as "militar", but the internal domain term is `Personnel`.

## Asset

A material or item that can be checked out to personnel.

## Asset Name

Asset name is the current primary user-facing identifier for an asset record.

Because Cordell does not yet model serial numbers, patrimony numbers, stock, or inventory identity, asset names must be unique.

Asset names are normalized by trimming surrounding whitespace and collapsing repeated internal whitespace.

## Personnel Lifecycle

Personnel lifecycle describes whether a personnel record is available for normal operational workflows.

Active personnel can receive new checkout.

Inactive personnel should not receive new checkout but remain available for historical custody records, receipts, and auditability.

## Asset Lifecycle

Asset lifecycle describes whether an asset record is available for normal operational workflows.

Active assets can be used in new checkout.

Inactive assets should not be used in new checkout but remain available for historical custody records, receipts, and auditability.

## Deactivation

Deactivation removes a record from normal operational workflows without deleting it.

Deactivation is different from deletion.

A deactivated record remains part of historical custody records.

## Reactivation

Reactivation makes a previously inactive record available for normal operational workflows again.

## Confirmation Modal

A confirmation modal is a contextual confirmation prompt shown before a sensitive state-changing action.

It reduces accidental clicks while keeping the user on the current page.

A confirmation modal does not replace backend validation, authentication, authorization, CSRF protection, application services, or repository/database rules.

## Search Status Filter

Search status filter controls whether global search returns active records, inactive records, or all records.

It prevents inactive records from polluting daily operational search while preserving explicit access to archived records.

## Custody

The state of responsibility over one or more assets.

## Checkout

A business event where assets are assigned to personnel.

## Return

A business event where assets previously assigned to personnel are returned.

## Custody Transaction

An immutable business event representing either a checkout, a return, or a future correction/adjustment.

## Custody Line

A line inside a custody transaction representing an asset and a quantity.

## Custody Balance

The current quantity of a given asset under the custody of a given person.

## Atomic Custody Balance Update

An atomic custody balance update validates and changes a custody balance in one database operation.

For returns, Cordell decreases balance only when the current quantity is still sufficient.

This prevents concurrent returns from consuming the same balance twice.

## Custody Receipt

A custody receipt is a read-only view of a checkout or return transaction.

It shows who received or returned assets, which authenticated operator registered the event, which assets were involved, quantities, notes, and timestamp.

## Custody Correction

A custody correction is a future explicit event used to correct a previously registered custody transaction.

Corrections must not edit or delete the original transaction.

A correction should explain what was wrong, who registered the correction, when it happened, and how it affects custody balances.

## Custody Ledger

The custody ledger is the append-only history of checkout, return, and future correction events.

Current custody state should be derived from the ledger rather than manually overwritten.

## Current Custody Balance

Current custody balance is the current quantity of a given asset under the responsibility of a given personnel record.

Checkout transactions increase current custody balance.

Return transactions decrease current custody balance.

Current custody views hide zero-balance rows.

## Custody Line Consolidation

Custody line consolidation is the application-level normalization that combines duplicate asset lines in the same checkout or return command before the custody transaction is created.

For example, two lines for the same asset with quantities `1` and `2` become one line with quantity `3`.

## Return UI Guardrail

A return UI guardrail is a user-interface restriction that helps prevent invalid return attempts before submission.

For example, Cordell only shows assets currently under a selected personnel record's custody in the return form.

Guardrails improve usability but do not replace backend validation.

## Active Record Guardrail

An active record guardrail prevents new operational responsibility from being created with inactive records.

In Cordell, checkout requires active personnel and active assets.

Return may still be allowed for inactive personnel or inactive assets when current custody exists, because return reduces or clears an existing responsibility.

## Operator

A system user who performs actions inside Cordell.

Operators authenticate with a unique operator registration ID and password.

Operator display names use rank and alias, such as "sergeant silva".

## Operator Attribution

Operator attribution is the link between a custody transaction and the authenticated operator who registered it in Cordell.

It is different from the personnel receiving or returning assets.

## Audit Log

A technical append-only record of relevant actions performed in the system.

## Audit Event

An audit event is an append-only application record describing a high-value action performed in Cordell.

It records the actor, event type, entity, outcome, timestamp, and minimal metadata.

Audit events are not debug logs and must not contain secrets.

Audit events are currently recorded on a best-effort basis after the main application action succeeds.

This means audit events are useful for traceability, but they are not yet guaranteed to be committed atomically with the action that caused them.

The transaction boundary decision is documented in ADR 0005.

## Status Filter

A status filter controls whether list pages show active records, inactive records, or all records.

Cordell uses status filters to avoid polluting daily operational screens while preserving access to historical and inactive records.

## Inactive Pending Custody

Inactive pending custody means an inactive personnel or inactive asset is still involved in a current custody balance.

This is allowed because deactivation does not settle custody.

Inactive pending custody must remain visible until it is cleared through return or future correction workflows.

## Asset Deactivation and Reactivation

Asset records can be deactivated and reactivated by authenticated operators.

Deactivation:

- does not delete the asset record
- does not delete custody history
- does not settle current custody
- does not create a return transaction
- prevents new checkout for that asset
- keeps receipts and history readable
- records an audit event

Reactivation makes the asset available for normal checkout workflows again.

Asset deactivation does not require a reason in the current implementation.

If a reason is added later, it should preferably be selected from a controlled list instead of arbitrary free text.

Inactive assets with current custody should remain visible in current custody and detail views with clear warnings.

Inactive asset names remain reserved. Deactivation does not allow another asset with the same name to be created.
