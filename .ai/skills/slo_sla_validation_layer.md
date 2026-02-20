# SLO / SLA Validation Layer

Purpose:
Enforce service reliability targets.

---

## Definitions

SLA:
Contract with users (e.g., 99.9% uptime)

SLO:
Internal reliability target (e.g., P95 < 200ms)

---

## Default SLO

API:
- Error rate < 1%
- P95 latency < 200ms
- P99 latency < 500ms

DB:
- Slow query < 2%
- Pool utilization < 80%

---

## Validation Rules

If SLO violated:

- Trigger alert
- Enable performance mode
- Activate circuit breaker
- Reduce traffic (if possible)

---

## Deployment Gate

Do not deploy if:

- Test latency > threshold
- Error rate increased
- Migration time too long

---

## Error Budget Policy

If error budget exhausted:

- Freeze new feature release
- Focus on stabilization