# Guard Enforcement Policy

Hard Stops:

- destructive_migration_without_backup
- drop_table_without_rollback
- auth_logic_removal
- sensitive_data_logging
- unpaginated_list_endpoint
- missing_rollback_plan_on_deploy

Soft Warnings:

- missing_index
- missing_observability
- missing_test_coverage
- no_feature_flag_on_high_risk_change

# Output-Type Guard Mapping

Post-Generation Validation MUST depend on output type.

If output_type == "migration_script":
  - Must include rollback
  - Must not contain irreversible DROP without backup
  - Must not modify auth-related tables without migration_safety_gate

If output_type == "api_endpoint":
  - Must enforce pagination for list endpoints
  - Must include input validation
  - Must return standardized error response

If output_type == "deployment_plan":
  - Must include rollback strategy
  - Must include SLO validation step
  - Must include observability confirmation

If output_type == "security_patch":
  - Must not expose sensitive logs
  - Must include regression test

Failure in output-type validation:
  → Regenerate
  → If still invalid → Hard Stop