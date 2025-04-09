    -- name: CreateCustomerAddress :exec
INSERT INTO customer_address (
  id_address, customer_id, address, phone_number
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteCustomerAddress :exec
DELETE FROM customer_address
WHERE id_address = ?;

-- name: UpdateCustomerAddress :exec
UPDATE customer_address
SET address = COALESCE(sqlc.narg('address'), address),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    update_date = NOW()
WHERE id_address = ?;

-- name: GetCustomerAddress :one
SELECT * FROM customer_address
WHERE id_address = ? LIMIT 1;

-- name: ListCustomerAddresses :many
SELECT * FROM customer_address
WHERE customer_id = ?;

-- name: ListCustomerAddressesPaged :many
SELECT * FROM customer_address
WHERE customer_id = ?
ORDER BY id_address
LIMIT ? OFFSET ?;