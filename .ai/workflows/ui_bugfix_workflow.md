# UI Bugfix Workflow

Entry Type:
fix_ui_bug

Purpose:
Fix UI defect without introducing new regressions,
side effects, or contract violations.

---

## 1. Bug Classification

Identify bug type:

- Rendering issue
- State inconsistency
- Async handling bug
- Form validation issue
- Access control leak
- Performance regression

Bug must be reproducible.

If reproduction unclear → Hard Stop.

---

## 2. Root Cause Analysis

Determine:

- Is it UI-only?
- Is it API response mismatch?
- Is it race condition?
- Is it state mutation bug?

If backend-related → escalate to appropriate backend workflow.

---

## 3. Fix Isolation

Fix must:

- Modify minimal surface area
- Avoid large refactor unless required
- Preserve component contract

If fix requires structural rewrite:
  → escalate to modify_ui_component

---

## 4. State Integrity Check

After fix:

- Ensure loading state intact
- Ensure error state intact
- Ensure no memory leak
- Ensure no duplicated event handler

---

## 5. Async Safety Check

If async involved:

- Confirm no unmounted state update
- Confirm no double-submit
- Confirm proper error propagation

---

## 6. Visual Regression Guard

If UI layout affected:

- Confirm responsive layout preserved
- Confirm no layout shift introduced
- Confirm accessibility not degraded

---

## 7. Testing Requirements

Must include:

- Reproduction test
- Fix validation test
- No regression in related components

---

## 8. Post-Fix Validation

Ensure:

- No new console warnings
- No unhandled promise rejection
- No bypassed access control
- No direct API call inside render

If validation fails → Regenerate
Repeated failure → Hard Stop