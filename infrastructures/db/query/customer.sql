-- name: CreateCustomer :one
INSERT INTO customers (email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted;

-- name: GetCustomer :one
SELECT customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted
FROM customers
WHERE customer_id = $1 AND is_deleted = FALSE;

-- name: GetCustomerByEmail :one
SELECT customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted
FROM customers
WHERE email = $1 AND is_deleted = FALSE;

-- name: GetCustomerByOpenIDSub :one
SELECT customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted
FROM customers
WHERE openid_sub = $1 AND is_deleted = FALSE;

-- name: ListCustomers :many
SELECT customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted
FROM customers
WHERE is_deleted = FALSE
ORDER BY created_at DESC;

-- name: ListCustomersBySetupStatus :many
SELECT customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted
FROM customers
WHERE setup_completed = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: UpdateCustomer :one
UPDATE customers
SET email = $2, first_name = $3, last_name = $4, phone = $5, address = $6, city = $7, postal_code = $8, updated_at = NOW()
WHERE customer_id = $1 AND is_deleted = FALSE
RETURNING customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted;

-- name: UpdateCustomerSetupStatus :one
UPDATE customers
SET setup_completed = $2, updated_at = NOW()
WHERE customer_id = $1 AND is_deleted = FALSE
RETURNING customer_id, email, first_name, last_name, phone, address, city, postal_code, openid_sub, setup_completed, created_at, updated_at, is_deleted;

-- name: SoftDeleteCustomer :exec
UPDATE customers
SET is_deleted = TRUE, updated_at = NOW()
WHERE customer_id = $1 AND is_deleted = FALSE;

-- name: CountCustomers :one
SELECT COUNT(*) FROM customers WHERE is_deleted = FALSE;