# UI Component Guidelines

Cordell uses small server-rendered UI components.

Components should be created when they represent a repeated visual pattern with stable behavior.

Good component candidates:

- page headers
- breadcrumbs
- empty states
- feedback messages
- form action bars
- detail fields
- simple cards
- reusable custody line items

Avoid componentizing too early:

- complex forms with unique validation rules
- search and filter panels that are still evolving
- repository or query behavior inside transaction timeline cards
- highly page-specific layouts

## Custody Timeline

The custody timeline component is a shared visual component.

It should receive an already-prepared view model and should not perform business decisions.

Good uses:

- transaction ledger
- personnel custody history
- future asset custody history

The custody timeline is currently used by:

- global custody transaction ledger
- personnel custody history
- asset custody history

Avoid placing query/filter/pagination logic inside the component.

## Live Search

Live search should be progressive.

Rules:

- forms must work without JavaScript
- backend search remains authoritative
- JavaScript may only debounce, request partial HTML, update targets, and update the URL
- no business filtering rules should be implemented only in JavaScript
- partial templates should render the same cards/components as the full page
- live search can apply to complex pages when the result area is isolated into a partial
- good candidates include entity lists, admin lists, and ledgers
- pagination and filters must remain server-rendered and work without JavaScript

Internal code remains in English.

User-facing labels are presented in Portuguese through the web presentation layer.
