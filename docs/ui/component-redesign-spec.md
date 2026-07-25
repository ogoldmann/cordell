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
- navbar search, after navbar redesign

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
