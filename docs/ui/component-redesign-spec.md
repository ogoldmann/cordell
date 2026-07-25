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

### Recalibration notes

The search bar should be minimalist and precise.

The search icon box should read as a rounded square, not a pill.

Rounding should be moderate.

The primary color should be used subtly.

The component should not feel oversized.

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

### Navbar search placement

The navbar search is not absolutely centered.

It occupies the available space between the Cordell brand block and the navigation/operator area.

Desktop structure:

```txt
[brand] [search] [navigation + operator]
```

The search column remains reserved even when search is hidden on the dashboard.

This prevents layout movement when entering or leaving the dashboard.

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

### Recalibration notes

Navbar navigation links do not use icons.

The navbar must not shift when the Home link is hidden on the dashboard.

The navbar must not shift when the navbar search is hidden on the dashboard.

The navbar search should stay centered.

The operator profile card should be compact.

The operator role appears as a small rectangular badge.

The operator dropdown should match the profile card width.

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

## Custody transaction cards

Custody transaction cards are used in ledger and custody history timelines.

### Structure

Each card contains:

- sequence badge
- transaction type badge
- optional edited badge
- receipt link
- registration metadata
- personnel block
- material lines block

### Transaction type

Checkout and return are differentiated by:

- timeline dot color
- card left border color
- transaction type badge

Checkout uses the positive token.

Return uses the negative token.

Color must remain subtle.

### Timeline

The timeline uses:

- small dot
- subtle vertical line
- compact spacing between cards

### Rules

- no oversized cards
- no unnecessary icons
- receipt action remains text-based
- metadata and operational data are visually separated
- material lines use compact tabular layout

## Data tables

Cordell uses shared data table styling for operational lists.

Tables are used when the user needs to scan many records quickly.

Initial use cases:

- personnel list
- asset list

### Visual rules

Data tables should be:

- compact
- readable
- serious
- minimally rounded
- horizontally scrollable on small screens
- consistent across pages

Avoid:

- oversized rows
- excessive shadows
- navigation icons inside every row
- colorful decoration
- card-per-record layouts for dense lists

### Row behavior

Rows may be clickable for convenience.

The primary entity name remains a real link for accessibility.

### Personnel columns

- Nome de Guerra + Nome
- Seção
- Materiais sob custódia
- Status

### Asset columns

- Material
- Militares com custódia atual
- Status

### Responsiveness

Tables may scroll horizontally on small screens.

The page itself must not break or overflow horizontally.
