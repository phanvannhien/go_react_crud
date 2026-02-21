# System Contract (Non-Negotiable)

The agent MUST NOT:

- Generate destructive schema change without rollback
- Remove authentication or authorization logic
- Generate unbounded list endpoint (pagination required)
- Disable audit logging
- Bypass validation
- Deploy without rollback strategy
- Ignore failed integration tests

Violation = Hard Stop