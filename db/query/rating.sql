-- name: CreateRating :exec
INSERT INTO ratings (
  rating_id, comment, star, account_id, products_spu_id
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteRating :exec
DELETE FROM ratings
WHERE rating_id = ?;

-- name: UpdateRating :exec
UPDATE ratings
SET comment = COALESCE(sqlc.narg('comment'), comment),
    star = COALESCE(sqlc.narg('star'), star),
    update_date = NOW()
WHERE rating_id = ?;

-- name: GetRating :one
SELECT * FROM ratings
WHERE rating_id = ? LIMIT 1;

-- name: ListRatings :many
SELECT * FROM ratings
WHERE products_spu_id = ?;

-- name: ListRatingsPaged :many
SELECT * FROM ratings
WHERE products_spu_id = ?
ORDER BY create_date DESC
LIMIT ? OFFSET ?;