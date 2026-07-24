# Cordell Component Redesign Spec

## Search bars

Search bars are a core interaction pattern in Cordell.

Design:

- clean horizontal input
- rounded borders
- comfortable but not oversized height
- search icon inside a square rounded box on the right side
- strong focus state
- responsive behavior

Responsive behavior:

- on large screens, search bars are inline
- on smaller screens, search interaction may become a side panel or stacked block
- no search bar should break layout width

Search bars must remain progressive:

- without JavaScript, normal GET forms work
- with JavaScript, live search can request partial HTML
- backend remains the authority for search rules

## Navbar

The navbar is global.

Layout order:

1. Brand block
2. Global search, only outside the dashboard
3. Navigation links
4. Connected operator card

### Brand block

Shows:

- cordell
- PELOTÃO DE SEGURANÇA

The section name appears below "cordell", smaller, bold, and visually muted.

### Developer attribution

At the absolute top center of the navbar area:

- desenvolvido por sd galliac • github [external-link icon]

This must be very small and visually subtle.

### Navbar search

The navbar search uses the shared search bar pattern.

Behavior:

- hidden on dashboard/home
- visible on internal pages
- lightly animated when appearing after leaving home
- Enter opens full search page
- suggestions appear while typing
- no business search logic in JavaScript

### Navigation links

Order:

- Home
- Militares
- Materiais
- Transações
- Admin

Home link behavior:

- visible only outside the dashboard
- hidden on dashboard
- hiding/showing should not visually disturb the navbar layout

### Connected operator card

Shows:

- operator rank abbreviation and alias
- role below
- chevron/down button area

The whole card is clickable.

Dropdown options:

- Theme
- Sair

Logout should use a door/log-out icon.

Theme selection should support:

- light
- dark
- sepia

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
