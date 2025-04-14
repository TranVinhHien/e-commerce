-- name: CreateProductSKU :exec
INSERT INTO product_skus (
  product_sku_id, value, sku_stock, price, sort, products_spu_id
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: DeleteProductSKU :exec
DELETE FROM product_skus
WHERE product_sku_id = ?;

-- name: UpdateProductSKU :exec
UPDATE product_skus
SET value = COALESCE(sqlc.narg('value'), value),
    sku_stock = COALESCE(sqlc.narg('sku_stock'), sku_stock),
    price = COALESCE(sqlc.narg('price'), price),
    sort = COALESCE(sqlc.narg('sort'), sort),
    products_spu_id = COALESCE(sqlc.narg('products_spu_id'), products_spu_id),
    update_date = NOW()
WHERE product_sku_id = ?;

-- name: GetProductSKU :one
SELECT * FROM product_skus
WHERE product_sku_id = ? LIMIT 1;

-- name: ListProductSKUs :many
SELECT * FROM product_skus
WHERE products_spu_id = ?;

-- name: ListProductSKUsPaged :many
SELECT * FROM product_skus
WHERE products_spu_id = ?
ORDER BY product_sku_id
LIMIT ? OFFSET ?;


-- name: GetProductsBySKU :many
SELECT * 
FROM product_skus 
WHERE product_sku_id IN (?);