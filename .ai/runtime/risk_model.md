# Quantitative Risk Model

Base Scores:

- code_generation → 20
- migration_change → 60
- deployment_request → 50
- production_incident → 70
- performance_regression → 40
- security_finding → 80

Modifiers:

+30 → drop_table
+20 → no_rollback
+15 → production_environment
+25 → auth_modification
+20 → data_loss_risk
+15 → no_integration_test
+10 → no_index_on_large_table

Environment Multiplier:

production → x1.5
staging → x1.2
development → x1.0

Critical Rule:

Final Score ≥ 86 → Hard Stop Required

# Compound Risk Escalation Rules

Certain risk combinations automatically escalate to Critical
regardless of numeric total.

The following combinations MUST override calculated score:

1. migration_change + drop_table + production_environment
   → Force Risk Score = 95 (Critical)

2. migration_change + no_rollback
   → Minimum Risk Score = 85 (High boundary)

3. migration_change + auth_modification
   → Force Risk Score ≥ 90

4. deployment_request + no_rollback_plan + production_environment
   → Force Risk Score ≥ 90

5. security_finding + production_environment
   → Force Risk Score ≥ 92

Escalation Rules:

- Escalation overrides multipliers.
- Escalation applies before final threshold evaluation.
- Escalation cannot be bypassed without explicit override justification.