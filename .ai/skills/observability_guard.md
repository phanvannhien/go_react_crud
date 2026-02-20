# Observability Guard

Purpose:
Ensure full visibility of system behavior.

---

## Required Metrics

For every endpoint:

- Request count
- Error rate
- P95 latency
- P99 latency

For DB:

- Query duration
- Slow query count
- Connection pool usage

---

## Required Logging

- Structured logging (JSON)
- Correlation ID per request
- Log level separation
- No sensitive data in logs

---

## Tracing

- Each request must propagate trace ID
- Service layer spans
- DB query spans

---

## Health Endpoints

- /health (liveness)
- /ready (readiness)

---

## Forbidden

- Println debugging
- Logging secrets
- Missing trace ID
- No latency measurement