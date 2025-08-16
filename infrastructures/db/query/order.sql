-- name: CreateOrder :one
INSERT INTO orders (customer_id, order_number, status, total_amount)
VALUES ($1, $2, $3, $4)
RETURNING order_id, customer_id, order_number, status, total_amount, created_at, updated_at, is_deleted;

-- name: GetOrder :one
SELECT order_id, customer_id, order_number, status, total_amount, created_at, updated_at, is_deleted
FROM orders
WHERE order_id = $1 AND is_deleted = FALSE;