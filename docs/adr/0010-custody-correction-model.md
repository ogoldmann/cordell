# ADR 0010: Custody Correction Model

## Status

Accepted for future implementation

## Context

Cordell records custody transactions for checkout and return workflows.

A custody transaction currently represents an operational event:

* a checkout registers that assets were issued to personnel
* a return registers that assets were returned by personnel
* each transaction records the authenticated operator who registered it
* each transaction has immutable lines with assets and quantities
* each transaction has a timestamp
* each transaction can be viewed through a custody receipt
* high-value actions are recorded in the application audit log

This creates a ledger-like custody history.

However, real users can make mistakes:

* wrong asset selected
* wrong quantity entered
* wrong personnel selected
* checkout registered when it should not have been
* return registered with incorrect quantity
* notes missing or incorrect
* duplicated transaction
* operational correction requested after review

Cordell needs a correction model that preserves accountability and avoids silent history mutation.

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

Cordell will not edit or delete existing custody transactions to correct mistakes.

Custody corrections will be represented by new explicit correction events.

The custody ledger must remain append-only from the application perspective.

Existing checkout and return transactions remain part of history even if they were wrong. A later correction event explains and compensates for the mistake.

The future correction model should support:

* linking a correction to the original transaction
* recording the operator who registered the correction
* recording a required correction reason
* recording correction lines
* preserving the original transaction receipt
* creating a new correction receipt
* showing correction status in history
* recording audit events for corrections

The current likely transaction types are:

* `checkout`
* `return`
* `correction`

A correction transaction should not silently modify the original transaction. It should be a new transaction that affects custody balances according to explicit correction lines.

## Alternatives Considered

### Option A: Edit the original transaction

Allow operators or admins to edit a previously saved checkout or return.

Benefits:

* simple user experience
* the visible transaction can be made to look correct
* fewer records in history

Drawbacks:

* destroys historical truth
* makes receipts unreliable
* makes auditability weaker
* makes it harder to know what actually happened
* risks hiding mistakes or abuse
* conflicts with ledger-like custody history

Decision:

Rejected.

Custody transactions must not be edited after creation.

### Option B: Delete the original transaction

Allow admins to delete an incorrect custody transaction.

Benefits:

* simple to implement
* removes wrong records from normal views
* keeps current balances simple

Drawbacks:

* destroys evidence
* breaks transaction receipts
* makes audit trails harder to trust
* can hide operational mistakes
* creates ambiguity around historical balances
* conflicts with accountability requirements

Decision:

Rejected.

Custody transactions must not be deleted as part of normal application workflows.

### Option C: Cancel the original transaction

Mark a transaction as canceled and remove its effect from current custody balances.

Benefits:

* keeps original transaction visible
* makes intent clearer than deletion
* useful for fully invalid transactions
* simpler than line-level correction in some cases

Drawbacks:

* cancellation alone is not enough for partial corrections
* requires careful balance recalculation
* can become ambiguous if cancellation and returns interact
* still needs reason, actor, timestamp, and receipt
* may not handle wrong quantity or wrong asset elegantly

Decision:

Deferred as a possible future specialized correction type.

Cancellation may be useful later, but it should still be implemented as an explicit event, not as mutation of the original transaction.

### Option D: Manual adjustment

Allow admins to directly adjust current custody balances.

Benefits:

* flexible
* can fix almost any current-balance problem
* easy to build as an administrative tool

Drawbacks:

* weak domain language
* high risk of abuse or confusion
* separates current balance from transaction history
* makes reconciliation harder
* creates hidden state changes if not modeled carefully

Decision:

Rejected for now.

Cordell should avoid generic balance adjustment until there is a proven operational need.

### Option E: Correction event

Create a new correction transaction linked to the original transaction.

Benefits:

* preserves original history
* makes correction explicit
* supports auditability
* supports correction receipts
* supports operator attribution
* keeps ledger append-only
* can support partial corrections
* can be shown clearly in personnel history

Drawbacks:

* more complex than editing or deleting
* requires careful balance rules
* requires new UI and read models
* requires clear correction reason
* requires careful user education

Decision:

Accepted for future implementation.

This is the preferred model.

## Correction Principles

Custody correction implementation must follow these principles:

1. Never edit original custody transactions.
2. Never delete original custody transactions.
3. Every correction must be a new event.
4. Every correction must have an authenticated operator.
5. Every correction must have a required reason.
6. Every correction should link to the original transaction when applicable.
7. Every correction should have a receipt.
8. Every correction should produce an audit event.
9. Current custody balances should be derived from the ledger, not manually overwritten.
10. The UI must clearly show corrected transactions and correction events.

## Future Domain Shape

A future implementation may add:

```txt
custody_transactions
  id
  transaction_type
  personnel_id
  operator_id
  corrected_transaction_id NULL
  correction_reason
  notes
  created_at
```

Possible transaction types:

```txt
checkout
return
correction
```

A correction transaction may contain positive or negative line effects depending on the final domain design.

Alternatively, the system may introduce a separate table:

```txt
custody_corrections
```

This ADR does not yet decide the final database shape. It decides the domain rule:

```txt
corrections are explicit append-only events
```

## Balance Rules To Decide Later

Before implementation, Cordell must define exactly how correction lines affect current custody balances.

Open questions:

* Should correction lines use signed quantities?
* Should correction lines use a direction field?
* Should correction transactions be restricted to the same personnel as the original transaction?
* Can a correction change personnel?
* Can one correction correct multiple original transactions?
* Should correction be allowed only for admins?
* Should regular operators request correction while admins approve it?
* Should a corrected transaction be visually marked as corrected?
* Should the original receipt link to the correction receipt?
* Should the correction receipt link back to the original receipt?

These questions must be resolved in a future implementation milestone.

## Authorization Direction

The likely initial policy is:

* admins can register corrections
* regular operators cannot register corrections
* all authenticated operators can view correction history if they can view custody history

This policy may be adjusted later based on operational needs.

## Audit Direction

Future correction implementation should record audit events such as:

```txt
custody.correction_created
custody.transaction_corrected
```

Audit metadata may include:

* original transaction ID
* correction transaction ID
* personnel ID

Audit metadata must not include secrets, session tokens, CSRF tokens, password hashes, or raw form payloads.

## Consequences

Cordell will preserve historical custody records even when mistakes occur.

Corrections will increase the number of records in custody history, but the history will be more trustworthy.

Users may need clearer UI language to understand original transactions, correction events, and corrected status.

Implementation is deferred until the exact correction workflow is designed.

## Future Work

Future milestones should include:

* correction balance model ADR
* database schema for correction events
* correction use cases
* correction receipts
* correction audit events
* admin-only correction routes
* UI indicators for corrected transactions
* tests for current balance after corrections
* tests for correction history and receipts
