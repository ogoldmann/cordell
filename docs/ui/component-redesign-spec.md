# Cordell Component Redesign Spec

## Search bars

Search bars are a shared Cordell component.

They are used in:

- dashboard global search
- full global search page
- personnel list
- asset list
- custody transaction ledger
- admin operators list
- navbar search

### Visual rules

The search bar has:

- rounded outer frame
- clean input area
- square rounded search button on the right
- Lucide search icon
- strong focus state
- token-based theme colors
- no raw hardcoded colors

### Variants

The component supports three visual variants:

- default
- hero
- compact

#### Default

Used in normal pages and filters.

#### Hero

Used in the dashboard/home global search.

Characteristics:

- wider visual presence
- slightly taller
- more elevated
- autofocus allowed

#### Compact

Used in dense spaces such as navbar or future compact filters.

### Behavioral rules

The component does not own search behavior.

The form around the component owns:

- action
- method
- live search attributes
- hidden fields
- target containers

JavaScript remains generic.

The backend remains the authority for search behavior.

### Responsiveness

The search bar must be full-width inside its container.

On small screens, parent layouts may stack it or move it into a future mobile search panel.

## Navbar

The global authenticated navbar contains:

1. Cordell brand block
2. Navbar global search, except on dashboard
3. Navigation links
4. Connected operator card

### Brand block

The brand block shows:

- cordell
- PELOTÃO DE SEGURANÇA

The section name is small, bold, and muted.

### Developer attribution

The navbar shows a very small top-centered attribution:

- desenvolvido por sd galliac • github

The GitHub link uses the Lucide external-link icon.

### Search behavior

Navbar search is visible on internal pages and hidden on the dashboard.

It uses the shared compact search bar component.

It shows suggestions while typing.

Pressing Enter opens `/search?q=...`.

On very small screens, navbar search may be hidden to preserve navigation usability.

### Navigation links

Order:

- Home
- Militares
- Materiais
- Transações
- Admin

Home is hidden on the dashboard.

Admin is shown only to admin operators.

### Operator card

The connected operator appears as a rectangular card.

It shows:

- operator rank and alias
- role
- dropdown chevron

The whole card is clickable.

The dropdown contains:

- theme selector
- logout action

Logout uses the Lucide log-out icon.

## Transaction card

Transaction cards represent custody events.

They are used by:

- custody ledger
- personnel history
- asset history

Design:

- rectangular card
- external border color indicates transaction type
- checkout/cautela uses positive lime tone
- return/descautela uses muted red tone
- internal visual separation between metadata and transaction data

Metadata area:

- left: sequence badge
- below/near it: strong CAUTELA/DESCAUTELA badge
- right: Abrir Recibo button
- below: "registrado por [operator] às [date/time]"

Transaction data area:

- personnel block:
  - display name on top
  - full name below
- materials table:
  - asset name
  - quantity

Timeline detail:

- left side has a colored dot
- dot connects vertically to adjacent cards through a subtle line
- dot color follows transaction type

The transaction card must not look like isolated boxes stacked randomly.

It should feel like a timeline/history.
