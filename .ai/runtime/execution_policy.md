# Execution Escalation Policy

If risk_level == High:
    - Require explicit confirmation
    - Require rollback plan

If risk_level == Critical:
    - Deny execution
    - Require override justification

Override must:
    - Log reason
    - Attach risk impact note