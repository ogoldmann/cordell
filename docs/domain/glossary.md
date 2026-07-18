# Domain Glossary

## Personnel

A person who can receive assets under custody.

In the user interface, this may be presented as "militar", but the internal domain term is `Personnel`.

## Asset

A material or item that can be checked out to personnel.

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

## Custody Receipt

A custody receipt is a read-only view of a checkout or return transaction.

It shows who received or returned assets, which authenticated operator registered the event, which assets were involved, quantities, notes, and timestamp.

## Operator

A system user who performs actions inside Cordell.

Operators authenticate with a unique operator registration ID and password.

Operator display names use rank and alias, such as "sergeant silva".

## Operator Attribution

Operator attribution is the link between a custody transaction and the authenticated operator who registered it in Cordell.

It is different from the personnel receiving or returning assets.

## Audit Log

A technical append-only record of relevant actions performed in the system.
