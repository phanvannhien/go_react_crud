# Agent Execution Flow Contract

This document defines the mandatory execution lifecycle for the AI agent.

The agent MUST follow this flow.
Skipping any phase is considered a violation.

---

# 1. Boot Phase (Mandatory)

Before any action:

1. Load `.ai/master_policy.md`
2. Load `.ai/project_structure.md`
3. Load all relevant rules:
   - rules/backend.md
   - rules/frontend.md
   - rules/architecture.md
   - rules/security.md

If any file is missing → STOP execution.

---

# 2. Skill Resolution Phase

Agent must determine the correct skill based on user intent.

| Intent Type | Skill |
|------------|-------|
| Generate feature | CRUD Generation Skill |
| Improve structure | Refactor Skill |
| Validate structure | Architecture Constraint Skill |
| Optimize performance | Performance Tuning Skill |
| Audit vulnerabilities | Security Audit Skill |

If skill is ambiguous → request clarification.
Agent must NOT guess.

---

# 3. Pre-Execution Validation

Before writing or modifying code:

Agent MUST validate:

- Directory structure compliance
- Layer dependency direction
- No raw SQL outside /db/query
- No DB access inside handler
- No business logic in handler
- Pagination required for list endpoints
- No SELECT *

If violation detected → abort and explain.

---

# 4. Execution Phase

Agent executes selected skill.

Possible artifact generation:

Backend:
- dbmate migration
- sqlc query file
- service layer
- handler layer
- route registration

Frontend:
- feature folder
- api.ts
- types.ts
- pages
- components

Agent must respect project_structure.md strictly.

---

# 5. Post-Execution Enforcement Gates

After generation, agent MUST run:

1. Architecture Constraint Skill
2. Security Audit Skill
3. Performance Tuning Skill (basic level)

If any gate fails:
- Rollback generated output
- Provide violation report
- Do not finalize response

---

# 6. Dependency Direction Enforcement

Allowed dependency direction:

handler → service → sqlc
service → domain
domain → stdlib only

Forbidden:
- handler → sqlc
- handler → pgx
- service → echo
- domain → sqlc
- circular imports

Violation = reject.

---

# 7. Security Enforcement

Agent must ensure:

- JWT validation if auth_required = true
- No hardcoded secrets
- No localStorage token storage (frontend)
- No raw SQL
- Input validation exists
- Rate limiting middleware considered
- Proper CORS handling

---

# 8. Performance Baseline Enforcement

Agent must check:

Backend:
- No SELECT *
- LIMIT + OFFSET required
- ORDER BY required
- No N+1 pattern
- Index recommendation if needed

Frontend:
- No unnecessary re-renders
- No API calls in components
- Lazy load pages if applicable

---

# 9. Failure Conditions

Agent MUST reject output if:

- Clean architecture violated
- Security violation found
- SQL outside sqlc
- Handler accessing DB directly
- Pagination missing
- Project structure broken

Agent must explain why.

---

# 10. Execution Modes

Controlled Mode:
User explicitly selects skill.

Autonomous Mode:
Agent selects skill automatically but must still run all gates.

Default mode: Controlled Mode unless specified.

---

# 11. Final Output Contract

Final output must:

- Respect project directory structure
- Follow clean architecture
- Pass security gate
- Pass performance gate
- Be consistent with stack:
  - Go 1.24
  - Echo
  - PostgreSQL (pgx)
  - sqlc
  - dbmate
  - React 19
  - Vite
  - Tailwind
  - SPA

If not compliant → do not output code.