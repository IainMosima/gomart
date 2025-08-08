-- name: CreateOrder :one
INSERT INTO orders (customer_id, order_number, status, total_amount, shipping_address, billing_address, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted;

-- name: GetOrder :one
SELECT order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted
FROM orders
WHERE order_id = $1 AND is_deleted = FALSE;

-- name: GetOrderByNumber :one
SELECT order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted
FROM orders
WHERE order_number = $1 AND is_deleted = FALSE;

-- name: ListOrders :many
SELECT order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted
FROM orders
WHERE is_deleted = FALSE
ORDER BY created_at DESC;

-- name: ListOrdersByCustomer :many
SELECT order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted
FROM orders
WHERE customer_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: ListOrdersByStatus :many
SELECT order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted
FROM orders
WHERE status = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: ListOrdersByDateRange :many
SELECT order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted
FROM orders
WHERE created_at >= $1 AND created_at <= $2 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: UpdateOrder :one
UPDATE orders
SET status = $2, total_amount = $3, shipping_address = $4, billing_address = $5, notes = $6, updated_at = NOW()
WHERE order_id = $1 AND is_deleted = FALSE
RETURNING order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE order_id = $1 AND is_deleted = FALSE
RETURNING order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted;

-- name: UpdateOrderTotal :one
UPDATE orders
SET total_amount = $2, updated_at = NOW()
WHERE order_id = $1 AND is_deleted = FALSE
RETURNING order_id, customer_id, order_number, status, total_amount, shipping_address, billing_address, notes, created_at, updated_at, is_deleted;

-- name: SoftDeleteOrder :exec
UPDATE orders
SET is_deleted = TRUE, updated_at = NOW()
WHERE order_id = $1 AND is_deleted = FALSE;

-- name: CountOrders :one
SELECT COUNT(*) FROM orders WHERE is_deleted = FALSE;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) FROM orders WHERE status = $1 AND is_deleted = FALSE;