-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted;

-- name: GetOrderItem :one
SELECT order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted
FROM order_items
WHERE order_item_id = $1 AND is_deleted = FALSE;

-- name: ListOrderItems :many
SELECT order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted
FROM order_items
WHERE is_deleted = FALSE
ORDER BY created_at DESC;

-- name: ListOrderItemsByOrder :many
SELECT order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted
FROM order_items
WHERE order_id = $1 AND is_deleted = FALSE
ORDER BY created_at ASC;

-- name: ListOrderItemsByProduct :many
SELECT order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted
FROM order_items
WHERE product_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetOrderItemsWithDetails :many
SELECT oi.order_item_id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.total_price, oi.created_at, oi.is_deleted,
       p.product_name, p.sku, p.description
FROM order_items oi
JOIN products p ON oi.product_id = p.product_id
WHERE oi.order_id = $1 AND oi.is_deleted = FALSE AND p.is_deleted = FALSE
ORDER BY oi.created_at ASC;

-- name: UpdateOrderItem :one
UPDATE order_items
SET quantity = $2, unit_price = $3, total_price = $4
WHERE order_item_id = $1 AND is_deleted = FALSE
RETURNING order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted;

-- name: UpdateOrderItemQuantity :one
UPDATE order_items
SET quantity = $2, total_price = quantity * unit_price
WHERE order_item_id = $1 AND is_deleted = FALSE
RETURNING order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted;

-- name: SoftDeleteOrderItem :exec
UPDATE order_items
SET is_deleted = TRUE
WHERE order_item_id = $1 AND is_deleted = FALSE;

-- name: SoftDeleteOrderItemsByOrder :exec
UPDATE order_items
SET is_deleted = TRUE
WHERE order_id = $1 AND is_deleted = FALSE;

-- name: CountOrderItems :one
SELECT COUNT(*) FROM order_items WHERE is_deleted = FALSE;

-- name: CountOrderItemsByOrder :one
SELECT COUNT(*) FROM order_items WHERE order_id = $1 AND is_deleted = FALSE;

-- name: GetOrderTotalFromItems :one
SELECT COALESCE(SUM(total_price), 0) as order_total
FROM order_items
WHERE order_id = $1 AND is_deleted = FALSE;