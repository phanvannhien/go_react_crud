# Prompt Guide: Generate User, Auth & Product Modules

Purpose:
Standardized prompt examples for generating production-grade modules
using CRUD Generation Skill (Enterprise Edition).

---

# 1️⃣ CREATE USER MODULE (Auth Enabled)

Use case:
- User management
- Authentication required
- Email login
- Searchable email
- Role-based extension possible

Command:

/generate-crud user \
  email:string:required:unique:index:searchable:sortable \
  password_hash:string:required \
  role:string:required:filterable:sortable \
  is_active:bool:default=true:filterable \
  created_at:timestamp \
  updated_at:timestamp \
  --auth_required=true \
  --pagination_mode=offset

Expected generation:

Backend:
- Migration with UNIQUE email
- Index on email
- Filter support on role, is_active
- Search support on email
- Offset pagination
- Validation for:
  - email format
  - required fields
- OpenAPI spec
- Integration test

Frontend:
- Zod validation schema
- Search input (email)
- Filter dropdown (role, is_active)
- Pagination UI
- Disable submit if invalid

Security:
- Do not return password_hash in API response
- Whitelist sort fields
- Enforce auth middleware

---

# 2️⃣ CREATE AUTH MODULE (Login + Register)

Use case:
- Authentication system
- JWT-based
- No generic CRUD listing

Command:

/generate-auth user \
  email:string:required:unique:index \
  password:string:required \
  --jwt=true

Expected generation:

Backend:
- Register endpoint
- Login endpoint
- Password hashing (bcrypt)
- JWT generation
- Token expiration
- Middleware for protected routes
- No SELECT *

Security enforcement:
- Hash password before storing
- Never expose password
- Validate email format
- Validate password length

Frontend:
- Login page
- Register page
- Token storage strategy
- Global auth context
- 401 handler redirect

---

# 3️⃣ CREATE PRODUCT MODULE (With Filters + Cursor Mode)

Use case:
- Large dataset
- Searchable product name
- Filter by category
- Price range filter
- High-scale ready

Command:

/generate-crud product \
  name:string:required:searchable:sortable:index \
  description:text \
  price:decimal:required:filterable:sortable:index \
  category_id:uuid:required:fk=categories.id:index:filterable \
  stock:int:required:filterable \
  is_active:bool:default=true:filterable \
  created_at:timestamp:cursor \
  updated_at:timestamp \
  --auth_required=true \
  --pagination_mode=cursor \
  --high_scale=true

Expected generation:

Backend:

- Migration:
  - FK constraint
  - Index on category_id
  - Composite index suggestion (category_id, is_active)
- Cursor pagination (created_at or id)
- Filter builder using typed sqlc params
- Search on name (ILIKE)
- Price filter
- Stock filter
- Integration test
- OpenAPI spec

Frontend:

- Infinite scroll OR Load More
- Filter panel:
  - Category select
  - Price range
  - Stock filter
  - Active toggle
- Debounced search input
- Whitelisted sorting
- Cursor-based loading

Performance:

- No OFFSET
- Indexed cursor field
- Composite index recommendation
- No unbounded queries

---

# 4️⃣ VALIDATION EXPECTATIONS

Every generation must include:

Backend:
- DTO validation
- UUID validation
- Email validation
- Numeric range validation
- Boolean strict parsing
- RFC3339 timestamps

Frontend:
- Zod schema
- Inline errors
- Disabled submit if invalid
- Server error mapping

---

# 5️⃣ PAGINATION MODES

Offset Mode:

/generate-crud resource ... --pagination_mode=offset

Cursor Mode (recommended for large tables):

/generate-crud resource ... --pagination_mode=cursor --high_scale=true

---

# 6️⃣ FAILURE CONDITIONS

Generation must be rejected if:

- SELECT *
- No LIMIT
- No ORDER BY whitelist
- Handler calls DB
- Missing validation
- Missing index for sortable field
- Missing pagination
- Missing OpenAPI
- Missing integration test

---

# 7️⃣ DESIGN PRINCIPLES

Generated module must be:

- Clean architecture compliant
- Security hardened
- Performance aware
- Fully validated
- Index optimized
- Test covered
- OpenAPI documented
- SPA safe

---

# FINAL NOTE

Use:

- Offset mode for small tables
- Cursor mode for high-scale tables
- Always index searchable and filterable fields
- Never expose sensitive fields
- Always enforce structured response contract