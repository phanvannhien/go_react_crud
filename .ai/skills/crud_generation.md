# Skill: CRUD Generation (Enterprise Edition)

Purpose:
Generate a full production-grade CRUD module aligned with project stack,
architecture constraints, validation policy, performance rules, and security standards.

---

# STACK

Backend:
- Go 1.24
- Echo
- PostgreSQL (pgx)
- sqlc
- dbmate

Frontend:
- React 19
- Vite
- TailwindCSS
- SPA architecture
- Zod (validation)

Architecture:
- Clean architecture
- Handler → Service → sqlc Repository
- No DB call in handler
- No business logic in handler
- SQL only inside /db/query
- Explicit column selection
- Context-aware operations

---

# INPUT CONTRACT

Required:
- resource_name (singular, lowercase)
- fields (name + type + modifiers)
- auth_required (true/false)

Optional:
- pagination_mode: offset | cursor
- high_scale: true/false

---

# FIELD DSL FORMAT

field:type[:modifier[:modifier...]]

Supported types:

| DSL       | PostgreSQL    | Go        | React |
|-----------|---------------|----------|-------|
| string    | TEXT          | string   | string |
| text      | TEXT          | string   | string |
| int       | INTEGER       | int32    | number |
| bigint    | BIGINT        | int64    | number |
| decimal   | NUMERIC       | float64  | number |
| bool      | BOOLEAN       | bool     | boolean |
| uuid      | UUID          | uuid.UUID| string |
| timestamp | TIMESTAMPTZ   | time.Time| string |

Modifiers:

- required
- unique
- index
- nullable
- default=value
- searchable
- sortable
- filterable
- fk=table.column
- cursor

---

# OUTPUT STRUCTURE

Backend:

/db/migrations/
/db/query/
/internal/service/
/internal/handler/
/internal/routes/
/internal/test/
/docs/openapi/

Frontend:

/ui/src/features/<resource>/

---

# STEP 1: Generate Migration (dbmate)

Rules:

- UUID primary key
- TIMESTAMPTZ
- created_at DEFAULT now()
- updated_at DEFAULT now()
- Add FK constraints
- Add UNIQUE constraints
- Add index for:
  - unique fields
  - searchable fields
  - filterable fields
  - sortable fields
  - foreign keys

Auto index suggestion block must be included.

Example:

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

---

# STEP 2: Generate sqlc Query File

Location:
- /db/query/<resource>.sql

Required queries:

- Create<Resource>
- Get<Resource>ByID
- List<Resources>
- Update<Resource>
- Delete<Resource>
- Count<Resources>

Rules:

- Never SELECT *
- Always explicit columns
- Always LIMIT
- Always ORDER BY whitelist
- No raw string concat
- Parameterized SQL only

---

# FILTER BUILDER (MANDATORY)

For filterable fields:

Use sqlc-compatible dynamic pattern:

WHERE
  ($1::text IS NULL OR email ILIKE '%' || $1 || '%')
  AND ($2::uuid IS NULL OR role_id = $2)

All filters must be typed.
No dynamic SQL building in Go.

---

# PAGINATION

## Offset Pagination (default)

Query params:
- page
- limit
- sort
- order

SQL:
LIMIT $limit OFFSET $offset

Response:

{
  "data": [],
  "meta": {
    "page": number,
    "limit": number,
    "total": number
  }
}

---

## Cursor Pagination (if modifier cursor OR high_scale=true)

Query params:
- cursor
- limit

SQL pattern:

WHERE id > $cursor
ORDER BY id ASC
LIMIT $limit

Response:

{
  "data": [],
  "next_cursor": "uuid"
}

Cursor must be stable and indexed.

---

# STEP 3: Run sqlc generate

Assume sqlc.yaml exists.

---

# STEP 4: Generate Service Layer

Location:
/internal/service/<resource>_service.go

Rules:

- Accept context.Context
- Use injected sqlc Queries
- Wrap errors using %w
- No HTTP logic
- No Echo dependency
- Validate defensive checks if critical

---

# STEP 5: Generate Handler Layer

Location:
/internal/handler/<resource>_handler.go

Responsibilities:

- Parse input
- Bind DTO
- Validate DTO
- Call service
- Return structured JSON

Response format:

{
  "data": ...,
  "error": null
}

Validation error:

{
  "error": "validation_error",
  "details": {...}
}

Prohibited:

- DB access
- Business logic
- Raw SQL

---

# BACKEND VALIDATION (MANDATORY)

For each field:

- required → validate:"required"
- string → trim
- email → format validation
- uuid → uuid validation
- int → range if provided
- bool → strict parse
- timestamp → RFC3339 only

Service must not assume handler validated correctly.

---

# FRONTEND FEATURE GENERATION

Location:
/ui/src/features/<resource>/

Structure:

- api.ts
- types.ts
- schema.ts
- components/
- pages/

Rules:

- No fetch inside component
- All fetch inside api.ts
- Tailwind only
- No inline styles
- Use Zod schema
- Reuse backend validation rules
- Display server validation errors

---

# AUTO GENERATE OPENAPI SPEC

Location:
/docs/openapi/<resource>.yaml

Must include:

- Paths
- Query parameters
- Request schema
- Response schema
- Pagination schema
- Error responses (400,401,404,500)

---

# AUTO GENERATE INTEGRATION TEST

Location:
/internal/test/<resource>_integration_test.go

Test must:

- Spin test DB
- Insert fixture
- Call handler
- Validate response JSON
- Validate pagination
- Validate filter behavior
- Validate validation error

---

# PERFORMANCE ENFORCEMENT

Agent must:

1. Detect high cardinality fields
2. Suggest composite index if multiple filterable fields
3. Ensure FK indexed
4. Ensure ORDER BY field indexed
5. Warn if missing index on searchable field

---

# SECURITY ENFORCEMENT

Mandatory:

- No SELECT *
- No unbounded query
- Whitelisted ORDER BY
- Parameterized SQL
- No raw string building
- Validate all query params
- Protect against mass assignment
- Enforce auth middleware if auth_required=true

---

# VALIDATION CHECKLIST

Before finishing generation:

- No SQL outside /db/query
- No DB call in handler
- No business logic in UI
- UUID used
- Timestamps RFC3339
- Pagination implemented
- Filter implemented
- Validation implemented
- Index suggested
- OpenAPI generated
- Integration test generated
- Error wrapping used
- No SELECT *
- No missing LIMIT

---

# FAILURE CONDITIONS

Reject generation if:

- SELECT *
- Missing LIMIT
- Missing pagination
- Missing validation
- Handler accessing DB
- React component calling fetch
- Unindexed sortable field
- Missing context.Context
- Raw SQL string concat

---

# HIGH SCALE MODE (Optional)

If high_scale=true:

- Force cursor pagination
- Force created_at index
- Force composite filter index
- Enforce strict sorting whitelist
- Recommend read replica usage
- Recommend soft delete support

---

# FINAL OBJECTIVE

Generated CRUD must be:

- Production-safe
- Performance-aware
- Security-hardened
- Test-covered
- OpenAPI-documented
- Clean-architecture compliant