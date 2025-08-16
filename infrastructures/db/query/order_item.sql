-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at, is_deleted;