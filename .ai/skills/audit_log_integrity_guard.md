# Audit Log Integrity Guard

Purpose:
Ensure traceability of critical operations.

---

## When Audit Log Required

- User creation
- Role change
- Product price change
- Stock update
- Permission change
- Deletion operation

---

## Audit Structure

Fields:
- id
- entity_type
- entity_id
- action
- actor_id
- timestamp
- metadata (jsonb)

---

## Integrity Rules

- Never allow delete of audit log
- Append-only table
- No update to audit records
- Use transaction when writing audit log

---

## Sensitive Operation Rule

If service:
- Changes role
- Changes permission
- Changes price
- Deletes entity

Then:
- Must write audit log
- Must be in same transaction

---

## Forbidden

- Silent destructive change
- Logging partial data
- Logging outside transaction