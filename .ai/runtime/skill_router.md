# Deterministic Skill Routing

Event: migration_change
Activate:
- migration_safety_gate
- rollback_strategy_engine
- integration_test_generation

Event: deployment_request
Activate:
- slo_sla_validation_layer
- observability_guard
- rollback_strategy_engine

Event: performance_regression
Activate:
- performance_tuning
- observability_guard

Event: security_finding
Activate:
- security_audit
- audit_log_integrity_guard

Event: code_generation
Activate:
- crud_generation
- api_contract_validation
- integration_test_generation