-- name: CreateSupplier :exec
INSERT INTO suppliers (
  supplier_id, name, phone_number, email, address
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteSupplier :exec
DELETE FROM suppliers
WHERE supplier_id = ?;

-- name: UpdateSupplier :exec
UPDATE suppliers
SET name = COALESCE(sqlc.narg('name'), name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    email = COALESCE(sqlc.narg('email'), email),
    address = COALESCE(sqlc.narg('address'), address),
    update_date = NOW()
WHERE supplier_id = ?;

-- name: GetSupplier :one
SELECT * FROM suppliers
WHERE supplier_id = ? LIMIT 1;

-- name: ListSuppliers :many
SELECT * FROM suppliers;

-- name: ListSuppliersPaged :many
SELECT * FROM suppliers
ORDER BY supplier_id
LIMIT ? OFFSET ?;