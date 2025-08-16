-- name: CreateUser :one
INSERT INTO customers (user_id, phone_number, user_name, email)
VALUES ($1, $2, $3, $4)
RETURNING user_id, phone_number, user_name, email;

-- name: GetUser :one
SELECT user_id, phone_number, user_name, email
FROM customers WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT user_id, phone_number, user_name, email
FROM customers WHERE email = $1;