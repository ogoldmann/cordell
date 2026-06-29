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

## Operator

A system user who performs actions inside Cordell.

## Audit Log

A technical append-only record of relevant actions performed in the system.
