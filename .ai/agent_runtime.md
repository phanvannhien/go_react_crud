# Agent Runtime Model

Purpose:
Define when and how skills are activated.

---

## Event Sources

1. Code generation request
2. Migration change
3. Deployment request
4. Production failure report
5. Performance regression
6. Security finding

---

## Execution Phases

Phase 1: Analyze Request
Phase 2: Risk Scoring
Phase 3: Skill Selection
Phase 4: Execution Plan
Phase 5: Guard Validation
Phase 6: Output

---

## Skill Activation Rules

- If deployment mentioned → activate:
  - slo_sla_validation_layer
  - observability_guard
  - rollback_strategy_engine

- If performance mentioned → activate:
  - self_healing_rule_engine
  - observability_guard

- If error spike / outage mentioned → activate:
  - incident_response_automation
  - rollback_strategy_engine

- If database schema changed → activate:
  - migration_safety_gate
  - dynamic_filter_builder validation