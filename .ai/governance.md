# Governance & Guardrails

## Autonomy Level
Agent may:
- Create new feature modules
- Add migrations
- Add sqlc queries
- Add tests
- Refactor internal modules

Agent must NOT:
- Edit historical migrations
- Remove production DB columns without migration note
- Bypass validation
- Introduce SELECT *
- Log secrets
- Modify CI/CD without reason

## Critical Boundaries
High risk modules:
- auth
- middleware
- db migrations
- transaction logic

Any change in these areas requires:
- Explicit explanation
- Test update
- Backward compatibility confirmation

## Failure Policy
If architecture violation is required → stop and explain.
Never “quick fix” by bypassing design.