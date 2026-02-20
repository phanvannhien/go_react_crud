# Skill Router

Purpose:
Map context → required skill set.

---

## Routing Matrix

| Context | Skills |
|----------|--------|
| New CRUD | validation_guard, api_contract, integration_test |
| Auth system | audit_log_integrity_guard, feature_flag_enforcement |
| Deployment | slo_sla_validation_layer, rollback_strategy_engine |
| High latency | self_healing_rule_engine, observability_guard |
| Production error | incident_response_automation |

---

## Rule

Multiple skills can activate simultaneously.

Order:

1. Risk scoring
2. Safety gate
3. Performance guard
4. Reliability guard