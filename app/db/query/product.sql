-- name: CreateProduct :one
INSERT INTO products (name, description, price, category_id, stock, is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, description, price, category_id, stock, is_active, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, price, category_id, stock, is_active, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, category_id, stock, is_active, created_at, updated_at
FROM products
WHERE
    ($1::text IS NULL OR name ILIKE '%' || $1 || '%')
    AND ($2::uuid IS NULL OR category_id = $2)
    AND ($3::boolean IS NULL OR is_active = $3)
    AND ($4::integer IS NULL OR stock >= $4)
    AND ($5::numeric IS NULL OR price >= $5)
    AND ($6::numeric IS NULL OR price <= $6)
    AND ($7::timestamptz IS NULL OR created_at < $7)
ORDER BY created_at DESC
LIMIT $8;

-- name: UpdateProduct :one
UPDATE products
SET
    name = COALESCE(NULLIF($2, ''), name),
    description = $3,
    price = $4,
    category_id = $5,
    stock = $6,
    is_active = $7,
    updated_at = now()
WHERE id = $1
RETURNING id, name, description, price, category_id, stock, is_active, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
