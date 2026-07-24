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
