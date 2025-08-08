-- name: CreateCategory :one
INSERT INTO categories (category_name)
VALUES ($1)
RETURNING category_id, category_name, created_at, updated_at, is_deleted;

-- name: GetCategory :one
SELECT category_id, category_name, created_at, updated_at, is_deleted
FROM categories
WHERE category_id = $1 AND is_deleted = FALSE;

-- name: GetCategoryByName :one
SELECT category_id, category_name, created_at, updated_at, is_deleted
FROM categories
WHERE category_name = $1 AND is_deleted = FALSE;

-- name: ListCategories :many
SELECT category_id, category_name, created_at, updated_at, is_deleted
FROM categories
WHERE is_deleted = FALSE
ORDER BY category_name ASC;

-- name: UpdateCategory :one
UPDATE categories
SET category_name = $2, updated_at = NOW()
WHERE category_id = $1 AND is_deleted = FALSE
RETURNING category_id, category_name, created_at, updated_at, is_deleted;

-- name: SoftDeleteCategory :exec
UPDATE categories
SET is_deleted = TRUE, updated_at = NOW()
WHERE category_id = $1 AND is_deleted = FALSE;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories WHERE is_deleted = FALSE;