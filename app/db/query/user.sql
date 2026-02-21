-- name: CreateUser :one
INSERT INTO users (email, password_hash, role, is_active)
VALUES ($1, $2, $3, $4)
RETURNING id, email, role, is_active, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, role, is_active, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, role, is_active, created_at, updated_at
FROM users
WHERE
    ($1::text IS NULL OR email ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR role = $2)
    AND ($3::boolean IS NULL OR is_active = $3)
ORDER BY
    CASE WHEN @sort_field::text = 'email' AND @sort_order::text = 'asc' THEN email END ASC,
    CASE WHEN @sort_field::text = 'email' AND @sort_order::text = 'desc' THEN email END DESC,
    CASE WHEN @sort_field::text = 'role' AND @sort_order::text = 'asc' THEN role END ASC,
    CASE WHEN @sort_field::text = 'role' AND @sort_order::text = 'desc' THEN role END DESC,
    CASE WHEN @sort_field::text = 'created_at' AND @sort_order::text = 'asc' THEN created_at END ASC,
    CASE WHEN @sort_field::text = 'created_at' AND @sort_order::text = 'desc' THEN created_at END DESC,
    created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountUsers :one
SELECT count(*) FROM users
WHERE
    ($1::text IS NULL OR email ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR role = $2)
    AND ($3::boolean IS NULL OR is_active = $3);

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE(NULLIF($2, ''), email),
    role = COALESCE(NULLIF($3, ''), role),
    is_active = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, email, role, is_active, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
