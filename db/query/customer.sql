-- name: CreateCustomer :exec
INSERT INTO customers (
  customer_id, name, email, image, dob, gender, account_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE customer_id = ?;

-- name: UpdateCustomer :exec
UPDATE customers
SET name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    image = COALESCE(sqlc.narg('image'), image),
    dob = COALESCE(sqlc.narg('dob'), dob),
    gender = COALESCE(sqlc.narg('gender'), gender),
    account_id = COALESCE(sqlc.narg('account_id'), account_id),
    update_date = NOW()
WHERE customer_id = ?;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE customer_id = ? LIMIT 1;

-- name: ListCustomers :many
SELECT * FROM customers;

-- name: ListCustomersPaged :many
SELECT * FROM customers
ORDER BY customer_id
LIMIT ? OFFSET ?;

-- name: GetCustomerByAccountID :one
SELECT * FROM customers
WHERE account_id = ? LIMIT 1;