-- name: CreateOrder :exec
INSERT INTO orders (
  order_id, total_amount, customer_address_id, discount_id, payment_method_id, customer_id
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: DeleteOrder :exec
UPDATE orders
SET order_status = 'Đã Hủy',
    update_date = NOW()
WHERE order_id = ?;

-- name: UpdateOrder :exec
UPDATE orders
SET total_amount = COALESCE(sqlc.narg('total_amount'), total_amount),
    customer_address_id = COALESCE(sqlc.narg('customer_address_id'), customer_address_id),
    discount_id = COALESCE(sqlc.narg('discount_id'), discount_id),
    payment_method_id = COALESCE(sqlc.narg('payment_method_id'), payment_method_id),
    payment_status = COALESCE(sqlc.narg('payment_status'), payment_status),
    order_status = COALESCE(sqlc.narg('order_status'), order_status),
    customer_id = COALESCE(sqlc.narg('customer_id'), customer_id),
    update_date = NOW()
WHERE order_id = ?;

-- name: GetOrder :one
SELECT * FROM orders
WHERE order_id = ? LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders;

-- name: ListOrdersPaged :many
SELECT * FROM orders
ORDER BY create_date DESC
LIMIT ? OFFSET ?;

-- name: ListCustomerOrders :many
SELECT * FROM orders
WHERE customer_id = ?;