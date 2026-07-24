# Cordell Icon System

Cordell uses Lucide icons.

Icons are treated as part of the design system, not decorative fragments.

## Principles

- Icons must support meaning.
- Critical actions must keep text labels.
- Icons should use `currentColor`.
- Icons should follow the same size scale.
- Icons should not introduce new colors by themselves.
- Icons should not be loaded from a CDN.
- Icons should work without JavaScript.

## Source

Lucide icons are installed through NPM and copied into the project.

Source directory:

```txt
node_modules/lucide-static/icons
```

Cordell icon directory:

```txt
internal/web/views/icons
```

Sync script:

```txt
scripts/sync-icons.sh
```

## Initial Icon Set

- search
- log-out
- external-link
- chevron-down
- home
- users
- package
- clipboard-list
- shield
- plus
- receipt-text
- filter
- calendar
- sun
- moon
- monitor
- user
- key-round

## Size Scale

Recommended Tailwind sizes:

- size-3: tiny metadata icon
- size-4: inline text icon
- size-5: normal action icon
- size-6: prominent action icon
- size-8: empty state or hero icon

## Accessibility

Decorative icons should be hidden from screen readers.

Icons that carry meaning without text must provide accessible labels.

Preferred pattern:

- icon + visible text for important actions
- icon-only only for secondary UI controls with aria-label
