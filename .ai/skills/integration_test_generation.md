# Integration Test Generation Skill

Purpose:
Ensure every feature is production-safe.

---

## Integration Test Requirements

For every endpoint:

- Setup test DB
- Apply migration
- Seed necessary data
- Call HTTP handler
- Assert response code
- Assert JSON response
- Assert DB state

---

## Required Cases

1. Success case
2. Validation failure
3. Unauthorized access (if protected)
4. Pagination boundary case
5. Not found case

---

## Setup Pattern

- Use test container or local test DB
- Use transaction rollback between tests
- Avoid mocking repository

---

## Table Driven Structure

tests := []struct {
    name string
    payload interface{}
    expectedStatus int
    expectedResponse interface{}
}{}

---

## Performance Check

If list endpoint:
- Insert multiple rows
- Ensure pagination works
- Ensure no N+1 query

---

## Coverage Target

- Service layer: >= 80%
- Handler layer: >= 70%

---

## Forbidden

- No integration test
- Mocking DB for integration test
- Skipping validation test