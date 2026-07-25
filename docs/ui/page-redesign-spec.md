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

The home page is centered around global search.

Navbar search must not appear on the dashboard.

### Main search

The first main element below the navbar is the global search bar.

Design:

- horizontally centered
- comfortable distance from navbar
- not inside a card
- directly over the solid background
- long but not vertically heavy
- clean and powerful
- result area appears below it

Behavior:

- search input is focused automatically when entering home
- dashboard live search remains progressive
- without JavaScript, form submits to `/search`

### Welcome phrase

When the user enters the dashboard, a phrase appears above the search bar, aligned with the search bar's left edge.

The phrase should feel premium, with subtle typing/deleting animation.

Possible phrases:

- Segura, aqui é Pelotão de Segurança.
- Excelente, vamos ao trabalho.
- Excelente.
- Mais um dia. Tudo sob controle.
- Vamos colocar tudo em ordem.
- Hoje também, no controle.

### Operational action dock

A fixed operational action component appears near the bottom of the viewport.

It contains two groups:

- CADASTRAR
  - Militar
  - Material

- REGISTRAR
  - Cautela
  - Descautela

Design:

- horizontally centered
- slightly narrower than the home search bar
- not too tall
- external card/box
- two internal sub-boxes
- each sub-box has a small uppercase heading
- each action is a rectangular list-like button

Behavior:

- when no search is active and content does not overflow, the dock blends with the background
- when search results appear or page scroll becomes relevant, the dock gains visual elevation/shadow
- the dock is fixed near the bottom with a small gap from the viewport edge
- future possibility: use this dock globally as an operational shortcut

## Personnel list page

Route:

- `/personnel`

Layout:

- breadcrumb
- heading: Militares
- right of heading: Cadastrar Militar button
- below: description

Description:

- Militares que podem cautelar e descautelar materiais

Filters/search row:

- left: status filters
- right: search bar

List presentation:

- data table, not cards
- each row is fully clickable
- hover state must clearly indicate clickability

Columns:

1. Nome de Guerra + Nome
   - display name on top
   - full name below
2. Seção
3. Quantidade total de materiais sob custódia
4. Status badge:
   - ATIVO
   - INATIVO

Below:

- pagination

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
