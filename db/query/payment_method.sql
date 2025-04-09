-- name: CreatePaymentMethod :exec
INSERT INTO payment_methods (
  payment_method_id, name, description
) VALUES (
  ?, ?, ?
);

-- name: DeletePaymentMethod :exec
DELETE FROM payment_methods
WHERE payment_method_id = ?;

-- name: UpdatePaymentMethod :exec
UPDATE payment_methods
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description)
WHERE payment_method_id = ?;

-- name: GetPaymentMethod :one
SELECT * FROM payment_methods
WHERE payment_method_id = ? LIMIT 1;

-- name: ListPaymentMethods :many
SELECT * FROM payment_methods;

-- name: ListPaymentMethodsPaged :many
SELECT * FROM payment_methods
ORDER BY payment_method_id
LIMIT ? OFFSET ?;