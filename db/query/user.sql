-- name: CreateUser :one
INSERT INTO users (
  username, password,full_name,create_at
) VALUES (
  $1, $2, $3, $4
)RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1 ;

-- name: UpdateUser :one
UPDATE users
SET password = COALESCE(SQLC.narg(password), password),
full_name = COALESCE(SQLC.narg(full_name), full_name),
is_active = COALESCE(SQLC.narg(is_active), is_active)
WHERE username = SQLC.arg(username) 
RETURNING *;
