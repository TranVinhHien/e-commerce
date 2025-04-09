-- name: CreateProductSKUAttr :exec
INSERT INTO product_sku_attrs (
  product_sku_attr_id, name, value, image, products_spu_id
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteProductSKUAttr :exec
DELETE FROM product_sku_attrs
WHERE product_sku_attr_id = ?;

-- name: UpdateProductSKUAttr :exec
UPDATE product_sku_attrs
SET name = COALESCE(sqlc.narg('name'), name),
    value = COALESCE(sqlc.narg('value'), value),
    image = COALESCE(sqlc.narg('image'), image),
    products_spu_id = COALESCE(sqlc.narg('products_spu_id'), products_spu_id)
WHERE product_sku_attr_id = ?;

-- name: GetProductSKUAttr :one
SELECT * FROM product_sku_attrs
WHERE product_sku_attr_id = ? LIMIT 1;

-- name: ListProductSKUAttrs :many
SELECT * FROM product_sku_attrs
WHERE products_spu_id = ?;