# Frontend Rules (Enterprise SPA Edition)

Purpose:
Ensure frontend implementation is scalable, maintainable,
validation-consistent, API-aligned, and architecture-compliant.

---

# STACK

- React 19
- Vite
- TailwindCSS
- Bun (dev)
- TypeScript (mandatory)
- Zod (validation)

Architecture:

- SPA
- Feature-based folder structure
- Separation of UI, API, and state logic

---

# DIRECTORY STRUCTURE

/ui/src/features/<feature>/

Required structure:

- api.ts
- types.ts
- schema.ts
- hooks/
- components/
- pages/

No cross-feature direct imports except shared utilities.

---

# LAYER RESPONSIBILITY

## api.ts

Responsible for:

- All HTTP calls
- Query param building
- Error normalization
- Mapping response types

Prohibited:

- Direct UI state handling
- Business logic
- DOM interaction

Must use centralized fetch wrapper.

---

## types.ts

- Must mirror backend DTO
- Must match OpenAPI spec
- No duplication of inferred types

---

## schema.ts

- Zod schema only
- Must align with backend validation rules
- Required fields must match backend required fields
- UUID validation required
- Email validation required
- Numeric range if defined in backend
- Timestamp RFC3339 string

Frontend must not invent validation rules inconsistent with backend.

---

## hooks/

- Encapsulate stateful logic
- Handle:
  - data fetching
  - pagination state
  - filter state
  - form submission
- No direct fetch inside components

---

## components/

- Pure UI
- Receive props
- No fetch
- No business logic
- No direct side effects

---

## pages/

- Compose components + hooks
- Route-level logic only
- No direct API calls

---

# API CONTRACT

All API responses must follow backend JSON contract:

Success:

{
  "data": ...,
  "error": null
}

Failure:

{
  "data": null,
  "error": {
    "code": string,
    "message": string,
    "details": optional
  }
}

Frontend must normalize error before UI render.

Never assume raw backend error shape.

---

# VALIDATION RULES

Form validation must:

- Use Zod schema
- Validate before submit
- Disable submit if invalid
- Show inline error messages
- Highlight invalid fields
- Trim string inputs

Field rules:

- Required → non-empty
- Email → valid email format
- UUID → valid UUID format
- Number → numeric validation
- Boolean → explicit toggle
- Timestamp → RFC3339 string

Server validation errors must be displayed inline.

---

# PAGINATION RULES

List pages must:

- Send page & limit
- Respect max limit cap
- Pass sort & order
- Send filters if defined
- Read meta.total
- Render pagination UI

Pagination state must be inside hook.

Never compute total manually.

---

# FILTER & SORT RULES

- Only allow sortable fields defined by backend
- Only allow filterable fields defined by backend
- Never allow free-form sort injection
- Filter state must be typed

Query param builder must:

- Omit undefined values
- Encode properly
- Prevent injection patterns

---

# PERFORMANCE RULES

- Avoid unnecessary re-renders
- Use memoization where appropriate
- Avoid inline functions inside render
- Avoid heavy computation in components
- Debounce search inputs
- Avoid duplicate API calls

For large dataset:

- Prefer cursor pagination
- Avoid client-side sorting on large result
- Avoid client-side filtering on large dataset

---

# SECURITY RULES

- Never store JWT in localStorage if avoidable (prefer httpOnly cookie)
- Never expose sensitive fields (password_hash)
- Do not trust client-side validation
- Escape dynamic rendering content
- Prevent XSS via unsafe innerHTML

---

# STATE MANAGEMENT RULES

- Local feature state only
- No global state unless necessary
- No business rules in UI state
- Derived state must be memoized

---

# ACCESS CONTROL

If route protected:

- Redirect if unauthorized
- Hide unauthorized UI actions
- Do not render restricted buttons

---

# ERROR HANDLING

Must handle:

- Network error
- 400 validation error
- 401 unauthorized
- 404 not found
- 500 internal error

Must show user-friendly error messages.

---

# FAILURE CONDITIONS

Reject implementation if:

- Fetch inside component
- No validation schema
- No pagination support
- No error handling
- Inline styles used
- Business logic in UI component
- Free-form sort injection
- Hardcoded API URL inside component

---

# HIGH SCALE MODE (Optional)

If high-scale enabled:

- Use cursor pagination only
- Infinite scroll allowed
- Debounced search mandatory
- Avoid total count dependency
- Avoid client-side data transformation

---

# FINAL OBJECTIVE

Frontend must be:

- Predictable
- Type-safe
- API-aligned
- Validation-consistent
- Scalable
- Secure
- Hook-driven
- UI-pure