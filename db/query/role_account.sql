-- name: CreateRoleAccount :exec
INSERT INTO role_account (
  role_account_id, account_id, role_id
) VALUES (
  ?, ?, ?
);

-- name: DeleteRoleAccount :exec
DELETE FROM role_account
WHERE role_account_id = ?;

-- name: UpdateRoleAccount :exec
UPDATE role_account
SET account_id = COALESCE(sqlc.narg('account_id'), account_id),
    role_id = COALESCE(sqlc.narg('role_id'), role_id)
WHERE role_account_id = ?;

-- name: GetRoleAccount :one
SELECT * FROM role_account
WHERE role_account_id = ? LIMIT 1;

-- name: GetAccountRoles :many
SELECT r.* FROM roles r
JOIN role_account ra ON r.role_id = ra.role_id
WHERE ra.account_id = ?;

-- name: ListRoleAccounts :many
SELECT * FROM role_account;

-- name: ListRoleAccountsPaged :many
SELECT * FROM role_account
ORDER BY role_account_id
LIMIT ? OFFSET ?;