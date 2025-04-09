-- name: CreateAccount :exec
INSERT INTO accounts (
  account_id, username, password
) VALUES (
  ?, ?, ?
);

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_id = ?;

-- name: UpdateAccount :exec
UPDATE accounts
SET password = COALESCE(sqlc.narg('password'), password),
    active_status = COALESCE(sqlc.narg('active_status'), active_status),
    update_date = NOW()
WHERE username = ?;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE account_id = ? LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts;

-- name: ListAccountsPaged :many
SELECT * FROM accounts
ORDER BY account_id
LIMIT ? OFFSET ?;

-- name: GetAccountByUsername :one
SELECT * FROM accounts
WHERE username = ? LIMIT 1;

-- name: Login :one
SELECT accounts.*,role_account.role_id  FROM accounts join role_account 
ON accounts.account_id = role_account.account_id 
WHERE username = ? LIMIT 1;