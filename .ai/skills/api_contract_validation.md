# API Contract Validation Skill

Purpose:
Keep backend + React frontend consistent.

---

## Contract Rules

- DTO must not expose internal DB struct
- JSON field names must be explicit
- OpenAPI must reflect handler signature
- Breaking change must be documented

---

## Validation Steps

When changing API:

1. Does request struct change?
2. Does response struct change?
3. Does frontend DTO need update?
4. Is OpenAPI updated?
5. Is backward compatibility preserved?

---

## Breaking Change Detection

Breaking if:

- Field removed
- Field renamed
- Field type changed
- Required field added

If breaking:
- Require version bump
- Require migration note
- Require frontend update

---

## Sync Policy

After modifying handler:

- Update OpenAPI spec
- Ensure frontend service types match
- Regenerate client if applicable

---

## Forbidden

- Returning sqlc model directly
- Returning internal error details
- Exposing password hash
- Inconsistent field naming

---

## Cursor Pagination Contract

If endpoint supports pagination:

Request:
{
  "cursor": string,
  "limit": number
}

Response:
{
  "data": [],
  "next_cursor": string | null
}