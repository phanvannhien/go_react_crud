# Self-Healing Rule Engine

Purpose:
Automatically mitigate recurring operational failures.

---

## Trigger Conditions

- Repeated 5xx from endpoint
- Transaction deadlock
- Slow query > threshold
- DB pool exhaustion
- Memory spike

---

## Healing Actions

If DB pool exhaustion:
- Reduce max open connections
- Reject non-critical requests

If slow query:
- Log EXPLAIN plan
- Suggest index

If deadlock:
- Retry transaction (limited attempts)

If panic:
- Recover middleware
- Log stack trace
- Continue serving

---

## Retry Policy

- Max retry attempts: 3
- Backoff strategy: exponential
- Retry only idempotent operations

---

## Circuit Breaker Rule

If endpoint failure rate > threshold:
- Open circuit
- Return fallback response
- Retry after cooldown

---

## Forbidden

- Infinite retry loop
- Retrying non-idempotent write
- Swallowing critical error