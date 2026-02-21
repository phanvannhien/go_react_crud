# Skill Dependency Matrix

rollback_strategy_engine
  requires:
    - migration_safety_gate

integration_test_generation
  requires:
    - api_contract_validation

performance_tuning
  requires:
    - observability_guard

slo_sla_validation_layer
  requires:
    - observability_guard