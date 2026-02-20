# Security Rules (Enterprise Edition)

Purpose:
Define mandatory security standards for backend and frontend.
All generated modules MUST comply.

Security is not optional.

---

# 1️⃣ CORE PRINCIPLES

- Zero trust input
- Explicit validation
- Least privilege
- Whitelisted operations
- No implicit behavior
- No hidden magic

---

# 2️⃣ INPUT VALIDATION (MANDATORY)

All external input must be validated:

Sources:
- JSON body
- Query params
- Path params
- Headers

Rules:

- UUID must validate format
- Email must validate format
- Numeric fields must enforce range
- Boolean must be strict parse
- Timestamps must be RFC3339
- Required fields must not be empty
- Strings must be trimmed
- Enforce max length for strings

Reject invalid input with:

{
  "error": "validation_error",
  "details": {...}
}

Never rely on frontend validation alone.

---

# 3️⃣ SQL INJECTION PROTECTION

Mandatory:

- No dynamic SQL building
- No string concatenation for SQL
- Use sqlc parameter binding only
- No dynamic table names
- No raw WHERE clause building

Prohibited:

fmt.Sprintf("SELECT ... %s", userInput)

All queries must live in /db/query.

---

# 4️⃣ ORDER BY & FILTER PROTECTION

- Only allow whitelisted sortable fields
- Only allow whitelisted filterable fields
- Never inject raw query param into SQL
- Validate sort direction:
  - asc
  - desc only

Reject unknown sort field.

---

# 5️⃣ AUTHENTICATION

If auth_required=true:

- All protected routes must attach middleware
- JWT must:
  - Have expiration
  - Be signed with strong secret
  - Be validated on each request
- Reject expired token
- Reject malformed token

Never trust user ID from client.
Always extract from validated token.

---

# 6️⃣ AUTHORIZATION (RBAC READY)

If role field exists:

- Validate role in middleware
- Do not trust role from request body
- Enforce permission at service layer
- Deny by default

Never allow privilege escalation via payload.

---

# 7️⃣ PASSWORD SECURITY

Mandatory:

- Use bcrypt for hashing
- Never store plain password
- Never log password
- Never return password_hash
- Enforce minimum password length
- Optional: rate limit login

---

# 8️⃣ SENSITIVE DATA EXPOSURE

Never expose:

- password_hash
- internal IDs if not required
- internal error stack
- database errors
- private tokens

Response must sanitize sensitive fields.

---

# 9️⃣ MASS ASSIGNMENT PROTECTION

Never bind request body directly to DB model.

Use separate:

- Request DTO
- DB model

Explicit field mapping required.

---

# 🔟 RATE LIMITING

Sensitive endpoints:

- Login
- Register
- Password reset

Must support:

- Rate limiting middleware
- IP throttling

---

# 11️⃣ CORS POLICY

- Do not allow wildcard origin in production
- Explicitly define allowed origins
- Allow credentials only if necessary

---

# 12️⃣ CSRF

If using cookie auth:

- Enable CSRF protection
- Use SameSite=strict

If using JWT header auth:
- CSRF less relevant but still enforce origin checks if needed

---

# 13️⃣ XSS PROTECTION (Frontend)

- Escape unsafe HTML
- Never use dangerouslySetInnerHTML without sanitization
- Sanitize rich text input
- Encode output properly

---

# 14️⃣ LOGGING RULES

Never log:

- Password
- JWT
- Authorization header
- Sensitive personal data

Log only:

- Request ID
- Safe metadata
- Error summary

---

# 15️⃣ FILE UPLOAD SECURITY (If applicable)

- Validate file type
- Validate file size
- Generate random file name
- Never trust client MIME type
- Store outside public root

---

# 16️⃣ ERROR HANDLING

Do not leak internal details.

Bad:

{
  "error": "pq: duplicate key value violates unique constraint"
}

Good:

{
  "error": "conflict"
}

---

# 17️⃣ DENIAL OF SERVICE PROTECTION

- Enforce max request body size
- Enforce max limit in pagination
- Prevent unbounded queries
- Avoid heavy CPU in handler
- Avoid blocking operations

---

# 18️⃣ HIGH SCALE MODE ADDITIONAL SECURITY

If high_scale=true:

- Prefer cursor pagination
- Enforce indexed filters
- Avoid OFFSET
- Monitor slow queries
- Recommend rate limit globally

---

# 19️⃣ FAILURE CONDITIONS

Reject implementation if:

- SELECT *
- Missing validation
- Dynamic SQL string building
- Missing auth middleware when required
- Password returned in API
- Missing error sanitization
- Missing sort whitelist
- No max limit cap
- Handler accesses DB directly

---

# 20️⃣ FINAL OBJECTIVE

System must be:

- SQL injection safe
- XSS resistant
- Privilege escalation safe
- Data leakage safe
- Authentication enforced
- Authorization enforced
- Rate limit aware
- Production hardened