-- name: CreateDescriptionAttr :exec
INSERT INTO description_attr (
  description_attr_id, name, value, products_spu_id
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteDescriptionAttr :exec
DELETE FROM description_attr
WHERE description_attr_id = ?;

-- name: UpdateDescriptionAttr :exec
UPDATE description_attr
SET name = COALESCE(sqlc.narg('name'), name),
    value = COALESCE(sqlc.narg('value'), value),
    products_spu_id = COALESCE(sqlc.narg('products_spu_id'), products_spu_id)
WHERE description_attr_id = ?;

-- name: GetDescriptionAttr :one
SELECT * FROM description_attr
WHERE description_attr_id = ? LIMIT 1;

-- name: ListDescriptionAttrs :many
SELECT * FROM description_attr
WHERE products_spu_id = ?;