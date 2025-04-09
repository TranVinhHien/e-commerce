-- name: CreateCategory :exec
INSERT INTO categorys (
  category_id, name, `key`, `path`, parent
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteCategory :exec
DELETE FROM categorys
WHERE category_id = ?;

-- name: UpdateCategory :exec
UPDATE categorys
SET name = COALESCE(sqlc.narg('name'), name),
    `key` = COALESCE(sqlc.narg('key'), `key`),
    `path` = COALESCE(sqlc.narg('path'), `path`),
    parent = COALESCE(sqlc.narg('parent'), parent)
WHERE category_id = ?;

-- name: GetCategory :one
SELECT * FROM categorys
WHERE category_id = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categorys;

-- name: ListCategoriesPaged :many
SELECT * FROM categorys
ORDER BY category_id
LIMIT ? OFFSET ?;