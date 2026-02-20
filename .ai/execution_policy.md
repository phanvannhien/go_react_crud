# Execution Policy

Purpose:
Prevent unsafe autonomous behavior.

---

## Hard Stops

Agent MUST stop if:

- Data deletion without backup
- Migration without rollback
- Removing auth guard
- Disabling audit log
- Bypassing feature flag

---

## Soft Warnings

Agent should warn if:

- No pagination
- No filter limit
- No integration test
- No error handling

---

## Deployment Gate

Block deployment if:

- SLO violated
- Migration unsafe
- Rollback missing
- Feature flag missing