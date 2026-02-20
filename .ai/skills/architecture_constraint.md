# Skill: Architecture Constraint Enforcement

Purpose:
Ensure generated or modified code follows defined architecture rules.

When to Trigger:
- New feature generation
- Pull request review
- Refactor operation

Actions:
1. Detect layer violations:
   - Handler calling DB directly
   - Business logic in handler
   - Direct SQL usage outside sqlc

2. Detect circular dependencies.

3. Verify:
   - Dependency injection usage
   - No global state

4. Reject code if:
   - Raw SQL outside sqlc
   - Inline DB calls in HTTP handler
   - UI calling backend outside API layer

Output:
- List violations
- Suggested refactor steps
- Code patch if possible

Additional Checks:

- List endpoint must implement pagination
- Filterable fields must not be hardcoded
- No SELECT *
- Validation must exist for Create and Update DTO
- No direct DB call inside handler
- Sort field must be validated against whitelist