# AI Operating System – go_react_crud

This repository uses a structured `.ai/` directory that turns the AI agent
from a code generator into an Autonomous Platform Engineer.

The `.ai` folder is not documentation.
It is the runtime brain of the AI.

Each file defines constraints, orchestration, guardrails, or domain knowledge
that the AI must follow before generating or modifying code.

---

# 1. Core Runtime Layer

## .ai/system.md

Defines global system identity and behavior.

Why needed:
Prevents the AI from behaving like a generic assistant.
It enforces that this agent operates as a production-aware platform engineer.

---

## .ai/system_contract.md

Defines non-negotiable rules.

Examples:
- No unsafe migration
- No destructive query
- No bypassing auth
- No disabling audit log
- No deploy without rollback

Why needed:
Acts as constitutional law.
Prevents catastrophic output.

---

## .ai/agent_runtime.md

Defines how the AI processes requests.

Execution model:
1. Analyze request
2. Score risk
3. Route skills
4. Apply rules
5. Generate safe output

Why needed:
Without runtime structure, skills will not activate coherently.

---

## .ai/skill_router.md

Maps context → required skills.

Example:
- Deployment → SLO + rollback + observability
- CRUD → validation + API contract + integration test
- Migration → migration guard + risk scoring

Why needed:
Ensures correct capabilities are activated automatically.

---

## .ai/execution_policy.md

Defines hard stops and soft warnings.

Hard stop examples:
- Data deletion without backup
- Auth logic removal
- No pagination in list endpoint

Why needed:
Prevents AI from generating dangerous code.

---

## .ai/execution_flow.md

Describes order of operations between skills.

Why needed:
Prevents random skill invocation.
Guarantees deterministic execution chain.

---

# 2. Governance Layer

## .ai/governance.md

Defines decision authority and safety hierarchy.

Why needed:
Prevents AI from acting beyond acceptable risk boundary.

---

# 3. Rule Layer (.ai/rules/)

These define stack-specific constraints.

## architecture.md
Defines:
- Layer separation
- No fat handlers
- sqlc-only query access
- Service layer enforcement

Why:
Prevents architectural drift.

---

## backend_rules.md
Defines:
- Echo middleware structure
- DTO separation
- Validation pattern
- Error response standard

Why:
Maintains backend consistency.

---

## frontend_rules.md
Defines:
- React API usage
- DTO mapping
- Pagination enforcement
- Filter integration

Why:
Prevents API-contract mismatch.

---

## migration_rules.md
Defines:
- Forward-only migrations
- Rollback requirement
- Large-table protection

Why:
Prevents irreversible schema damage.

---

## performance_rules.md
Defines:
- Query limit
- Index requirement
- Pagination requirement

Why:
Prevents full table scans in production.

---

## security_rules.md
Defines:
- JWT validation
- RBAC enforcement
- No sensitive logs
- Input sanitization

Why:
Prevents security regressions.

---

## testing_rules.md
Defines:
- Integration test required
- Test coverage threshold
- Migration test required

Why:
Prevents silent production failure.

---

## ci_cd_rules.md
Defines:
- Build checks
- Lint requirement
- Migration gate
- SLO validation before deploy

Why:
Prevents broken pipeline release.

---

# 4. Workflow Layer (.ai/workflows/)

These define execution playbooks.

## feature_workflow.md
Used when adding new feature.

Steps:
- Risk scoring
- Migration validation
- API contract update
- Integration test generation

---

## bugfix_workflow.md
Used when fixing production bug.

Steps:
- Root cause isolation
- Regression test
- Risk re-evaluation

---

## migration_workflow.md
Used when schema changes.

Steps:
- Risk scoring
- Backward compatibility check
- Rollback validation

---

## release_workflow.md
Used during deployment.

Steps:
- SLO validation
- Observability validation
- Rollback verification

Why workflows exist:
They convert abstract skills into actionable sequences.

---

# 5. Prompt Templates (.ai/prompts/)

These are high-level entry points.

- create_user.md
- create_auth.md
- create_product.md
- refactor_module.md

Why needed:
Standardizes AI input format.
Prevents incomplete CRUD generation.

---

# 6. Skills Layer (.ai/skills/)

These are specialized autonomous capabilities.

They are activated by skill_router.

---

## Core Engineering Skills

- crud_generation.md
- refactor.md
- transaction_pattern.md
- filtering_contract.md
- api_contract_validation.md
- integration_test_generation.md

Purpose:
Ensure complete, consistent module generation.

---

## Safety Skills

- migration_guard.md
- rollback_strategy_engine.md
- risk_scoring_model.md
- audit_log_integrity_guard.md
- feature_flag_enforcement.md
- security_audit.md

Purpose:
Prevent destructive changes.

---

## Reliability & SRE Skills

- observability_guard.md
- incident_response_automation.md
- self_healing_rule_engine.md
- slo_sla_validation_layer.md
- performance_tuning.md

Purpose:
Ensure system stability, uptime, and resilience.

---

# 7. Why This Architecture Is Necessary

Without `.ai`:

AI = stateless code generator

With `.ai`:

AI = deterministic platform engineer with:

- Risk awareness
- Deployment awareness
- Reliability awareness
- Security enforcement
- Migration safety
- Observability enforcement

---

# 8. Operational Philosophy

This AI must:

- Never prioritize speed over safety
- Never sacrifice observability
- Never generate incomplete production code
- Always enforce pagination + filtering
- Always require integration test for new modules
- Always validate rollback before migration

---

# 9. Resulting Maturity Level

This repository operates at:

Autonomous Platform Engineering Tier

Not:
- Basic CRUD assistant
- Template generator
- Simple refactor bot

But:
A structured AI capable of enforcing production discipline.

---

# 10. Summary

`.ai` acts as:

- Constitution (system_contract)
- Brain (agent_runtime)
- Router (skill_router)
- Guardrail (execution_policy)
- Domain knowledge (rules/)
- Playbooks (workflows/)
- Capabilities (skills/)
- Entry points (prompts/)

Together they form an AI Operating System for Go + Echo + sqlc + React.

Removing any major layer reduces reliability.

This design is intentional and production-focused.