-- name: CreateOrderDetail :exec
INSERT INTO order_detail (
  order_detail_id, quantity, unit_price, product_sku_id, order_id
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteOrderDetail :exec
DELETE FROM order_detail
WHERE order_detail_id = ?;

-- name: UpdateOrderDetail :exec
UPDATE order_detail
SET quantity = COALESCE(sqlc.narg('quantity'), quantity),
    unit_price = COALESCE(sqlc.narg('unit_price'), unit_price),
    product_sku_id = COALESCE(sqlc.narg('product_sku_id'), product_sku_id),
    order_id = COALESCE(sqlc.narg('order_id'), order_id)
WHERE order_detail_id = ?;

-- name: GetOrderDetail :one
SELECT * FROM order_detail
WHERE order_detail_id = ? LIMIT 1;

-- name: ListOrderDetails :many
SELECT * FROM order_detail
WHERE order_id = ?;