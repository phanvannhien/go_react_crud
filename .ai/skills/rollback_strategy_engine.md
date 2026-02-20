# Rollback Strategy Engine

Purpose:
Ensure every change can be safely reverted.

---

## Migration Rollback Rules

- Down migration must restore previous schema
- Never drop column without shadow backup
- For destructive changes:
  - Create new column
  - Backfill data
  - Switch application
  - Remove old column in later migration

---

## Service Rollback Rules

If logic changes:
- Maintain old API behavior behind flag
- Avoid breaking response structure

---

## Deployment Rollback

If deployment fails:
- Revert container image
- Rollback migration (if safe)
- Restore previous config

---

## Data Safety Rule

Before:
- Changing column type
- Removing constraint

Require:
- Data backup plan
- Migration dry-run

---

## Zero-Downtime Migration Pattern

1. Add new column nullable
2. Backfill
3. Deploy code writing both
4. Switch read path
5. Remove old column later