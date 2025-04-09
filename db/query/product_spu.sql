-- name: CreateProductSPU :exec
INSERT INTO products_spu (
  products_spu_id, name, brand_id, description, short_description, image, media, `key`, category_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteProductSPU :exec
UPDATE products_spu
SET delete_status = 'Deleted',
    update_date = NOW()
WHERE products_spu_id = ?;

-- name: UpdateProductSPU :exec
UPDATE products_spu
SET name = COALESCE(sqlc.narg('name'), name),
    brand_id = COALESCE(sqlc.narg('brand_id'), brand_id),
    description = COALESCE(sqlc.narg('description'), description),
    short_description = COALESCE(sqlc.narg('short_description'), short_description),
    stock_status = COALESCE(sqlc.narg('stock_status'), stock_status),
    sort = COALESCE(sqlc.narg('sort'), sort),
    image = COALESCE(sqlc.narg('image'), image),
    media = COALESCE(sqlc.narg('media'), media),
    `key` = COALESCE(sqlc.narg('key'), `key`),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    update_date = NOW()
WHERE products_spu_id = ?;

-- name: GetProductSPU :one
SELECT * FROM products_spu
WHERE products_spu_id = ? LIMIT 1;

-- name: ListProductSPUs :many
SELECT * FROM products_spu
WHERE delete_status = 'Active';

-- name: ListProductSPUsPaged :many 
SELECT * FROM products_spu
WHERE delete_status = 'Active'
ORDER BY products_spu_id
LIMIT ? OFFSET ?;

-- name: ListProductSPUsByCategory :many
SELECT * FROM products_spu
WHERE category_id = ? AND delete_status = 'Active';