-- name: CreateProduct :one
INSERT INTO products (product_name, description, price, sku, stock_quantity, category_id, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted;

-- name: GetProduct :one
SELECT product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted
FROM products
WHERE product_id = $1 AND is_deleted = FALSE;

-- name: GetProductBySKU :one
SELECT product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted
FROM products
WHERE sku = $1 AND is_deleted = FALSE;

-- name: ListProducts :many
SELECT product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted
FROM products
WHERE is_deleted = FALSE
ORDER BY created_at DESC;

-- name: ListActiveProducts :many
SELECT product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted
FROM products
WHERE is_active = true AND is_deleted = FALSE
ORDER BY product_name ASC;

-- name: ListProductsByCategory :many
SELECT product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted
FROM products
WHERE category_id = $1 AND is_deleted = FALSE
ORDER BY product_name ASC;

-- name: ListProductsInStock :many
SELECT product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted
FROM products
WHERE stock_quantity > 0 AND is_deleted = FALSE
ORDER BY product_name ASC;

-- name: UpdateProduct :one
UPDATE products
SET product_name = $2, description = $3, price = $4, stock_quantity = $5, category_id = $6, is_active = $7, updated_at = NOW()
WHERE product_id = $1 AND is_deleted = FALSE
RETURNING product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted;

-- name: UpdateProductStock :one
UPDATE products
SET stock_quantity = $2, updated_at = NOW()
WHERE product_id = $1 AND is_deleted = FALSE
RETURNING product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted;

-- name: UpdateProductStatus :one
UPDATE products
SET is_active = $2, updated_at = NOW()
WHERE product_id = $1 AND is_deleted = FALSE
RETURNING product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted;

-- name: SoftDeleteProduct :exec
UPDATE products
SET is_deleted = TRUE, updated_at = NOW()
WHERE product_id = $1 AND is_deleted = FALSE;

-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE is_deleted = FALSE;

-- name: CountActiveProducts :one
SELECT COUNT(*) FROM products WHERE is_active = true AND is_deleted = FALSE;