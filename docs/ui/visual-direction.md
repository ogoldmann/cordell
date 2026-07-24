# Cordell UI/UX Direction Brief

## Core premise

This redesign is a user interface reformulation.

All existing behavior, business rules, validations, security rules, custody flows, audit flows, search behavior, correction behavior, and data persistence rules must remain intact.

The redesign changes presentation, hierarchy, responsiveness, interaction polish, and visual consistency.

## Product language

Cordell is not presented as a "system".

Cordell is presented as a platform.

The word "system" should be avoided in user-facing copy because it carries a negative association with slow, ugly, corporate software.

Preferred language:

- platform
- Cordell
- ambiente
- operação
- painel
- fluxo
- registro
- controle

Avoid:

- sistema

## Brand principles

Cordell must feel:

- clean
- organized
- practical
- operational
- premium
- serious, but not lifeless
- stimulating, but not playful
- modern, but not generic

The interface should make the user want to navigate, search, register, and operate.

The desired feeling is controlled dopamine: enough energy to be pleasant, but not enough to compromise seriousness.

## UX principles

### Extreme responsiveness

Cordell must support different device sizes extremely well.

The layout must work from small screens to large desktop screens.

Responsive behavior must be intentional, not accidental.

Important responsive rules:

- no horizontal overflow
- no cramped forms
- tables must degrade gracefully on small screens
- navigation must remain usable on narrow screens
- important actions must remain reachable
- search must stay useful on mobile
- fixed/floating components must not cover critical content

### Pattern consistency

The same type of information must look the same across different pages.

Examples:

- personnel identity display
- asset display
- transaction cards
- search bars
- filters
- status badges
- action buttons
- empty states
- timeline entries

Consistency is a UX requirement.

Users should learn a pattern once and recognize it everywhere.

## Visual style

The visual style should use:

- clean surfaces
- comfortable spacing
- rounded cards
- soft shadows
- strong but controlled visual hierarchy
- clear hover states
- recognizable color semantics
- subtle premium animation
- low visual noise

Cards should have sufficiently rounded borders.

The UI should avoid the feeling of old corporate systems.

## Themes

Cordell supports three themes:

- light
- dark
- sepia

### Brand colors

Primary color:

- comfortable orange

Secondary color:

- brown / warm neutral

The dark and light themes should include well-placed brown tones when appropriate.

The dark theme especially should include warm brown tones so it does not become a generic gray/black interface.

### Semantic colors

Positive / input / checkout:

- lime green

Used for:

- success feedback
- checkout/cautela
- positive movement
- entry-like operation

Negative / output / return:

- muted red

Used for:

- destructive or negative feedback
- return/descautela
- output-like operation

The red should not be overly saturated or aggressive.

## Icons

Cordell will use Lucide icons.

Lucide icons should be used consistently, with a clear icon size scale.

Preferred icon style:

- stroke icons
- simple
- consistent
- not decorative-only
- paired with labels when the action is important

Initial icon candidates:

- search
- log-out
- external-link
- chevron-down
- user
- package
- clipboard-list
- home
- shield
- plus
- receipt
- filter
- calendar
- moon
- sun
- monitor

## Icon library decision

Cordell uses Lucide icons.

Reasoning:

- open-source
- permissive license
- consistent stroke style
- broad icon set
- works well for clean operational interfaces

Implementation rule:

- icons should support meaning, not replace important labels
- critical actions must keep text labels
- icon usage must be consistent across the platform

## Design tokens

Cordell's visual design is driven by semantic tokens.

Tokens should describe meaning, not raw color names.

Preferred examples:

- background
- surface
- surface-raised
- border
- text
- muted
- primary
- secondary
- positive
- negative
- warning
- focus

Avoid using raw color concepts in component names when the color has semantic meaning.

Good:

- positive
- negative
- primary

Avoid:

- green
- red
- orange

The checkout/cautela color is semantically positive.

The return/descautela color is semantically negative.

## Theme strategy

Cordell has three product themes:

- light
- dark
- sepia

The dark theme should use warm brown undertones.

The sepia theme should feel premium and calm, not old or visually dirty.

The primary brand color is comfortable orange.

The secondary brand tone is warm brown.
