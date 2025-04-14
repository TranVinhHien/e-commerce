-- name: CreateDiscount :exec
INSERT INTO discounts (
  discount_id, discount_code, discount_value, start_date, end_date, min_order_value, amount
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteDiscount :exec
DELETE FROM discounts
WHERE discount_id = ?;

-- name: UpdateDiscount :exec
UPDATE discounts
SET discount_code = COALESCE(sqlc.narg('discount_code'), discount_code),
    discount_value = COALESCE(sqlc.narg('discount_value'), discount_value),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    min_order_value = COALESCE(sqlc.narg('min_order_value'), min_order_value),
    amount = COALESCE(sqlc.narg('amount'), amount),
    update_date = NOW()
WHERE discount_id = ?;

-- name: UpdateDiscountAmount :exec
UPDATE discounts
SET amount = amount - 1,
    update_date = NOW()
WHERE discount_id = ?;

-- name: GetDiscount :one
SELECT * FROM discounts
WHERE discount_id = ? LIMIT 1;

-- name: GetDiscountByCode :one
SELECT * FROM discounts
WHERE discount_code = ? LIMIT 1;

-- name: ListDiscounts :many
SELECT * FROM discounts;

-- name: ListDiscountsPaged :many
SELECT * FROM discounts
LIMIT ? OFFSET ?;

-- name: ListDiscountInusesPaged :many
SELECT * FROM discounts WHERE amount >0  
LIMIT ? OFFSET ?;

-- name: ListActiveDiscounts :many
SELECT * FROM discounts;

-- name: CountDisscounts :one
SELECT COUNT(*) as totalElements FROM discounts;
