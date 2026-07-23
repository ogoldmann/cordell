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
- transaction timeline cards before the final read model is stable
- highly page-specific layouts

Internal code remains in English.

User-facing labels are presented in Portuguese through the web presentation layer.
