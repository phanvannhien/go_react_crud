-- name: CreateOrder :one
INSERT INTO orders (user_id, status, total_amount)
VALUES ($1, $2, $3)
RETURNING id, user_id, status, total_amount, created_at, updated_at;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING id, order_id, product_id, quantity, price;

-- name: GetOrderByID :one
SELECT id, user_id, status, total_amount, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, price
FROM order_items
WHERE order_id = $1;

-- name: ListOrdersByUserID :many
SELECT id, user_id, status, total_amount, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountOrdersByUserID :one
SELECT COUNT(*) FROM orders WHERE user_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, user_id, status, total_amount, created_at, updated_at;

-- name: UpdateOrderTotal :one
UPDATE orders
SET total_amount = $2, updated_at = now()
WHERE id = $1
RETURNING id, user_id, status, total_amount, created_at, updated_at;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items WHERE order_id = $1;
