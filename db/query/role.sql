-- name: CreateRole :exec
INSERT INTO roles (
  role_id, name, description
) VALUES (
  ?, ?, ?
);

-- name: DeleteRole :exec
DELETE FROM roles
WHERE role_id = ?;

-- name: UpdateRole :exec
UPDATE roles
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description)
WHERE role_id = ?;

-- name: GetRole :one
SELECT * FROM roles
WHERE role_id = ? LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles;

-- name: ListRolesPaged :many
SELECT * FROM roles
ORDER BY role_id
LIMIT ? OFFSET ?;