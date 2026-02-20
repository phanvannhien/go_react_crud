# Skill: Filtering & Pagination Contract

Purpose:
Ensure all list endpoints are production-ready.

Requirements:

- Default limit = 20
- Max limit = 100
- Default order = created_at DESC
- Sort whitelist only
- No free-text sort injection

SQL Pattern:

SELECT columns
FROM table
WHERE
  ($1::text IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY
  CASE WHEN $2 = 'name' THEN name END,
  created_at DESC
LIMIT $3 OFFSET $4;

Agent must enforce:

- No SELECT *
- No unbounded query
- No unsafe string concatenation