# Feature Flag Enforcement

Purpose:
Control rollout of risky features.

---

## When Required

Use feature flag if:

- New business logic
- New auth rule
- Pricing logic change
- Schema-dependent feature
- Critical performance change

---

## Flag Pattern

- Define flag in config
- Inject flag into service
- Branch logic inside service layer only

---

## Forbidden

- Feature flag in handler
- Feature flag in repository
- Flag without default value
- Flag without removal plan

---

## Removal Policy

After stable release:
- Remove dead flag
- Remove dual logic
- Clean code path

---

## Testing Rule

- Test both flag ON and OFF