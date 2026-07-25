# Cordell UI Redesign Checklist

## Foundation

- [ ] Existing behavior remains unchanged.
- [ ] Existing routes remain unchanged unless explicitly planned.
- [ ] Existing validations remain unchanged.
- [ ] Existing custody rules remain unchanged.
- [ ] Existing audit rules remain unchanged.
- [ ] Existing search backend remains authoritative.
- [ ] UI copy avoids the word "sistema".
- [ ] UI presents Cordell as a platform.

## Design Tokens

- [ ] Semantic CSS variables are defined.
- [ ] Light theme tokens are defined.
- [ ] Dark theme tokens are defined.
- [ ] Sepia theme tokens are defined.
- [ ] Primary orange token is defined.
- [ ] Secondary brown token is defined.
- [ ] Positive lime token is defined.
- [ ] Negative muted red token is defined.
- [ ] Focus token is defined.
- [ ] Radius tokens are defined.
- [ ] Shadow tokens are defined.
- [ ] Existing legacy Cordell color utilities remain temporarily compatible.

## Responsiveness

- [ ] Layout works on small screens.
- [ ] Layout works on tablets.
- [ ] Layout works on desktop.
- [ ] No horizontal overflow.
- [ ] Tables have mobile strategy.
- [ ] Fixed components do not hide important content.
- [ ] Navbar remains usable on narrow screens.

## Themes

- [ ] Light theme exists.
- [ ] Dark theme exists.
- [ ] Sepia theme exists.
- [ ] Primary orange is consistently applied.
- [ ] Secondary warm brown is consistently applied.
- [ ] Checkout/cautela uses lime positive tone.
- [ ] Return/descautela uses muted red tone.
- [ ] Theme selection is available from operator dropdown.
- [ ] Theme selection is available on login page.

## Theme System

- [ ] Theme preference is stored locally.
- [ ] Theme is applied through `data-theme`.
- [ ] Theme is applied before full page render to avoid flash.
- [ ] Light theme can be selected.
- [ ] Dark theme can be selected.
- [ ] Sepia theme can be selected.
- [ ] Theme selector component exists.
- [ ] Theme selector works before login.
- [ ] Theme selector is ready for operator dropdown.

## Components

- [ ] Search bar component redesigned.
- [ ] Navbar redesigned.
- [ ] Operator dropdown redesigned.
- [ ] Transaction card redesigned.
- [ ] Data table component created.
- [ ] Status badges standardized.
- [ ] Filter boxes standardized.
- [ ] Pagination standardized.
- [ ] Operational dock created for dashboard.
- [ ] Lucide icons integrated.

## Icon System

- [ ] Lucide is installed locally.
- [ ] Initial icon set is copied into the repository.
- [ ] Icon sync script exists.
- [ ] Icons render through Go templates.
- [ ] Icons use currentColor.
- [ ] Icons support decorative and labeled modes.
- [ ] Critical actions keep visible text labels.
- [ ] No CDN is used for icons.

## Search Bar Component

- [ ] Shared search bar view model exists.
- [ ] Shared search bar template exists.
- [ ] Search button uses Lucide search icon.
- [ ] Right-side button is square and rounded.
- [ ] Default variant exists.
- [ ] Hero variant exists.
- [ ] Compact variant exists.
- [ ] Dashboard uses hero search bar.
- [ ] `/search` uses shared search bar.
- [ ] Personnel list uses shared search bar.
- [ ] Asset list uses shared search bar.
- [ ] Custody ledger uses shared search bar.
- [ ] Admin operators uses shared search bar.
- [ ] Search behavior remains outside the component.
- [ ] Live search still works.

## Navbar Redesign

- [ ] Brand block shows cordell.
- [ ] Brand block shows PELOTÃO DE SEGURANÇA.
- [ ] Developer attribution appears subtly at the top.
- [ ] GitHub attribution uses external-link icon.
- [ ] Navbar search is hidden on dashboard.
- [ ] Navbar search appears on internal pages.
- [ ] Navbar search uses compact shared search bar.
- [ ] Home link is hidden on dashboard.
- [ ] Navigation order is Home, Militares, Materiais, Transações, Admin.
- [ ] Active navigation link is visually clear.
- [ ] Admin link is hidden for non-admin operators.
- [ ] Operator card shows display name.
- [ ] Operator card shows role.
- [ ] Operator dropdown opens and closes.
- [ ] Operator dropdown contains theme selector.
- [ ] Operator dropdown contains logout.
- [ ] Logout uses log-out icon.
- [ ] Navbar remains usable on tablet.
- [ ] Navbar remains usable on mobile.

## Pages

- [ ] Login page redesigned.
- [ ] Dashboard redesigned.
- [ ] Personnel list redesigned.
- [ ] Asset list redesigned.
- [ ] Custody transactions redesigned.
- [ ] Personnel detail redesigned.
- [ ] Asset detail redesigned.
- [ ] Checkout form redesigned.
- [ ] Return form redesigned.
- [ ] Correction form redesigned.
- [ ] Admin pages redesigned.
- [ ] 404 page aligned with new visual style.

## Login Page

- [ ] Login page uses redesigned visual shell.
- [ ] Cordell brand appears at top.
- [ ] Theme selector appears before login.
- [ ] Login card is centered.
- [ ] Login card uses subtle blur.
- [ ] Background has subtle particles.
- [ ] Identidade field remains functional.
- [ ] Password field remains functional.
- [ ] CSRF remains functional.
- [ ] Error message still appears.
- [ ] return_to is preserved if present.
- [ ] Theme persists after reload.
- [ ] Design remains serious and minimal.
