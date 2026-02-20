# Risk Scoring Model

Purpose:
Evaluate operational risk before executing changes.

---

## Risk Levels

LOW:
- New endpoint
- Add non-breaking field
- Add index

MEDIUM:
- Modify service logic
- Add migration (non-destructive)
- Add filter behavior

HIGH:
- Auth logic change
- Transaction logic change
- Foreign key change
- Index removal
- Column type change

CRITICAL:
- Drop table
- Drop column
- Change primary key
- Modify permission logic
- Remove audit logging

---

## Risk Evaluation Matrix

Score factors:

1. Data Impact (0-3)
2. Auth Impact (0-3)
3. Backward Compatibility (0-3)
4. Migration Complexity (0-3)

Total:
0-3 → Low
4-6 → Medium
7-9 → High
10-12 → Critical

---

## Required Actions

Low:
- Proceed with test

Medium:
- Add integration test
- Add performance check

High:
- Add rollback strategy
- Add migration safety explanation
- Add compatibility test

Critical:
- Require explicit confirmation
- Require data backup strategy
- Require staged deployment