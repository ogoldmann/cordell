# Cordell Page Redesign Spec

## Login page

Status: redesigned.

The login page uses a minimalist operational visual style.

Design:

- small centered Cordell brand at top
- section label below brand
- theme selector on the top right
- centered vertical login card
- moderate rounded corners
- subtle blur surface
- title: login
- subtitle: faça login na sua conta para continuar
- fields:
  - Identidade
  - Senha
- solid theme-aware background
- subtle rising particles
- subtle primary glow

Rules:

- authentication behavior must remain unchanged
- CSRF must remain enabled
- return_to must remain preserved if present
- theme can be selected before login
- no oversized card
- no playful particles
- no excessive rounding

## Home / Dashboard

Status: redesigned.

The dashboard is centered around global search.

Navbar search does not appear on the dashboard.

### Main search

The first main element below the navbar is the global search experience.

It includes:

- subtle welcome phrase
- hero search bar
- live search results below

The search bar is not wrapped in a large card.

It sits directly over the page background.

Behavior:

- input is focused automatically on dashboard load
- live search results render inline
- without JavaScript, form submits to `/search`
- dashboard URL does not change during live search

### Welcome phrase

The welcome phrase is selected randomly from the approved phrase list.

It uses subtle typewriter animation.

The animation should not loop aggressively.

It must respect reduced motion preferences.

### Welcome phrase behavior

The welcome phrase appears only once after login.

After the user leaves the dashboard and returns during the same authenticated browser session, the phrase should not appear again.

The phrase is reset when the user returns to the login page.

### Scroll behavior

The dashboard should not show a vertical scrollbar when there are no search results and no content overflow.

Scrolling should appear only when live search results or viewport constraints require it.

### Operational action dock

The operational dock is fixed near the bottom of the viewport.

It contains:

- CADASTRAR
  - Militar
  - Material

- REGISTRAR
  - Cautela
  - Descautela

The dock is visually quiet when the search is inactive.

When search is active, it gains subtle elevation so it separates from the results area.

### Operational dock personality

The operational dock may use restrained icons for its four actions.

Icons must be subtle and functional.

They should add recognition, not visual noise.

Future possibility:

- use this dock globally as an operational shortcut.

## Personnel list

Status: polished.

The personnel list uses the shared data table component.

Structure:

- breadcrumb
- title and primary action
- short description
- filters on the left
- search on the right
- data table
- pagination

Columns:

- Nome de Guerra + Nome
- Seção
- Materiais sob custódia
- Status

Rules:

- navigation links do not use icons
- rows are clickable
- the entity name remains a real link
- filters and search preserve each other
- live search updates only the table region
- table scrolls horizontally on small screens

### Current custody summaries

Personnel list shows the current total quantity of materials under each person's custody.

Asset list shows the current number of personnel with positive custody balance for each asset.

These summaries must use effective custody state.

They must account for transaction corrections.

The list must not compute custody state from original transactions only.

## Assets list page

Route:

- `/assets`

Same structure as personnel list, adapted to materials.

Columns:

1. Nome do material
2. Quantidade total de militares que custodiam esse material atualmente
3. Status badge:
   - ATIVO
   - INATIVO

Below:

- pagination

## Custody transactions page

Route:

- `/custody/transactions`

Layout:

- breadcrumb
- heading: Transações de Custódia
- right of heading:
  - Registrar Cautela
  - Registrar Descautela

Filters:

- external box with heading: FILTROS
- inside, two stacked sub-boxes

Sub-box 1:

- heading: Período
- year/month controls
- no separate "Aplicar período" button

Sub-box 2:

- transaction type filter
- edit status filter
- no search field here
- no internal apply/clear buttons
- no heading required

Below filters:

- left:
  - Aplicar Filtros
  - Limpar
- right:
  - search bar

Important behavior:

- there must be only one "Aplicar Filtros" button
- this button applies both period and other filters
- avoid duplicate apply buttons

Results:

- transaction cards
- not table view
- transaction cards follow the shared transaction card spec

Below:

- pagination
