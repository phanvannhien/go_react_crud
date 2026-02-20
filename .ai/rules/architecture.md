# Architecture Rules

Backend structure:

app/
  cmd/
  internal/
    handler/
    service/
    repository/
    middleware/
  db/
    migrations/
    query/

Rules:

- handler → service → repository (sqlc)
- handler must not access db directly
- service must not know HTTP layer
- repository wraps sqlc generated code
- transaction logic stays in service
- dependency direction is inward only