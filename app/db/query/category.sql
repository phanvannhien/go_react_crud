-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: GetCategoryByID :one
SELECT id, name, created_at, updated_at
FROM categories
WHERE id = $1;

-- name: ListCategories :many
SELECT id, name, created_at, updated_at
FROM categories
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: CountCategories :one
SELECT count(*) FROM categories;
