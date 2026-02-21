# UI Modification Workflow

Entry Type:
modify_ui_component

Purpose:
Safely modify an existing UI component or page without introducing
state regression, contract breakage, or UX inconsistency.

---

## 1. Scope Identification

- Identify target component/page
- Identify affected routes
- Identify impacted child components
- Identify API contracts consumed

Modification must declare:
- Visual change?
- Behavioral change?
- API contract change?
- State structure change?

If scope unclear → Hard Stop.

---

## 2. Contract Verification

Before modification:

- Validate backend API contract compatibility
- Confirm no DTO shape change without backend alignment
- Confirm response format unchanged

If API contract changes:
  → escalate to create_feature or modify_schema workflow

---

## 3. State Impact Analysis

Evaluate:

- Local state mutation
- Global store mutation
- Derived state recalculation
- Memoization break risk

Must confirm:
- No orphan state
- No stale closure risk
- No unintended re-render cascade

---

## 4. UX Regression Guard

Modified UI must preserve:

- Loading state
- Error state
- Empty state
- Disabled state (if form)

If any state handling removed → Regenerate.

---

## 5. Access Control Validation

- Route guard still applied
- Role-based rendering still enforced
- Sensitive data still masked if required

---

## 6. Refactor Safety

If refactoring component structure:

- Preserve prop contract
- Preserve public interface
- Update tests accordingly

---

## 7. Testing Requirements

Minimum:

- Existing test must pass
- Update snapshot if structure changed
- Add regression test if behavior changed

---

## 8. Post-Modification Validation

Ensure:

- No business logic introduced into component
- API calls still abstracted
- No direct side-effects in render body
- No state mutation outside setter

Failure → Regenerate
Repeated failure → Hard Stop