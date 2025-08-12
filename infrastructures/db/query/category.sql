-- name: CreateCategory :one
INSERT INTO categories (category_name, parent_id)
VALUES ($1, $2)
RETURNING category_id, category_name, parent_id, created_at, updated_at, is_deleted;

-- name: GetCategory :one
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted
FROM categories
WHERE category_id = $1 AND is_deleted = FALSE;

-- name: GetRootCategories :many
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted
FROM categories
WHERE parent_id IS NULL AND is_deleted = FALSE
ORDER BY category_name ASC;

-- name: GetCategoryChildren :many
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted
FROM categories
WHERE parent_id = $1 AND is_deleted = FALSE
ORDER BY category_name ASC;


-- name: ListCategories :many
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted
FROM categories
WHERE is_deleted = FALSE
ORDER BY category_name ASC;

-- name: UpdateCategory :one
UPDATE categories
SET category_name = $2, updated_at = NOW()
WHERE category_id = $1 AND is_deleted = FALSE
RETURNING category_id, category_name, parent_id, created_at, updated_at, is_deleted;


-- name: SoftDeleteCategory :exec
UPDATE categories
SET is_deleted = TRUE, updated_at = NOW()
WHERE category_id = $1 AND is_deleted = FALSE;


-- name: GetCategoryAverageProductPrice :one
SELECT COALESCE(AVG(p.price), 0.00)::DECIMAL(10,2) as average_price
FROM categories c
LEFT JOIN products p ON c.category_id = p.category_id 
WHERE c.category_id = $1 AND c.is_deleted = FALSE 
AND (p.is_deleted = FALSE OR p.is_deleted IS NULL) 
AND (p.is_active = TRUE OR p.is_active IS NULL);