-- name: DeleteCustomerAddress :exec
DELETE FROM customer_address
WHERE id_address = ?;

-- name: UpdateCustomerAddress :exec
UPDATE customer_address
SET address = COALESCE(sqlc.narg('address'), address),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    update_date = NOW()
WHERE id_address = ?;

-- name: GetCustomerAddressByAddressAndCustomer :one
SELECT customer_address.*,customers.* FROM customer_address join  customers ON  customers.customer_id = customer_address.customer_id
WHERE customer_address.id_address = ? AND customer_address.customer_id = ? LIMIT 1;
-- name: GetCustomerAddress :one
SELECT * FROM customer_address
WHERE id_address =  ? LIMIT 1;

-- name: ListCustomerAddresses :many
SELECT * FROM customer_address
WHERE customer_id = ?;

-- name: ListCustomerAddressesPaged :many
SELECT * FROM customer_address
WHERE customer_id = ?
ORDER BY id_address
LIMIT ? OFFSET ?;

-- name: CreateCustomerAddress :exec
INSERT INTO customer_address (
  id_address, customer_id, address, phone_number
) VALUES (
  ?, ?, ?, ?
);