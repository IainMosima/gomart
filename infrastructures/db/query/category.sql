-- name: CreateCategory :one
INSERT INTO categories (category_name, parent_id)
VALUES ($1, $2)
RETURNING category_id, category_name, parent_id, created_at, updated_at, is_deleted;

-- name: CreateRootCategory :one
INSERT INTO categories (category_name)
VALUES ($1)
RETURNING category_id, category_name, parent_id, created_at, updated_at, is_deleted;

-- name: GetCategory :one
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted
FROM categories
WHERE category_id = $1 AND is_deleted = FALSE;

-- name: GetCategoryByName :one
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted
FROM categories
WHERE category_name = $1 AND is_deleted = FALSE;

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

-- name: GetCategoryDescendants :many
WITH RECURSIVE category_tree AS (
    SELECT c.category_id, c.category_name, c.parent_id, c.created_at, c.updated_at, c.is_deleted, 0 as level
    FROM categories c
    WHERE c.category_id = $1 AND c.is_deleted = FALSE
    
    UNION ALL
    
    SELECT c.category_id, c.category_name, c.parent_id, c.created_at, c.updated_at, c.is_deleted, ct.level + 1
    FROM categories c
    INNER JOIN category_tree ct ON c.parent_id = ct.category_id
    WHERE c.is_deleted = FALSE
)
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted, level
FROM category_tree
ORDER BY level, category_name;

-- name: GetCategoryPath :many
WITH RECURSIVE category_path AS (
    SELECT c.category_id, c.category_name, c.parent_id, c.created_at, c.updated_at, c.is_deleted, 0 as level
    FROM categories c
    WHERE c.category_id = $1 AND c.is_deleted = FALSE
    
    UNION ALL
    
    SELECT c.category_id, c.category_name, c.parent_id, c.created_at, c.updated_at, c.is_deleted, cp.level + 1
    FROM categories c
    INNER JOIN category_path cp ON cp.parent_id = c.category_id
    WHERE c.is_deleted = FALSE
)
SELECT category_id, category_name, parent_id, created_at, updated_at, is_deleted, level
FROM category_path
ORDER BY level DESC;

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

-- name: MoveCategoryToParent :one
UPDATE categories
SET parent_id = $2, updated_at = NOW()
WHERE category_id = $1 AND is_deleted = FALSE
RETURNING category_id, category_name, parent_id, created_at, updated_at, is_deleted;

-- name: SoftDeleteCategory :exec
UPDATE categories
SET is_deleted = TRUE, updated_at = NOW()
WHERE category_id = $1 AND is_deleted = FALSE;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories WHERE is_deleted = FALSE;

-- name: CountRootCategories :one
SELECT COUNT(*) FROM categories WHERE parent_id IS NULL AND is_deleted = FALSE;