-- name: CreateEmployee :exec
INSERT INTO employees (
  employee_id, gender, dob, name, email, phone_number, address, salary, account_id
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE employee_id = ?;

-- name: UpdateEmployee :exec
UPDATE employees
SET gender = COALESCE(sqlc.narg('gender'), gender),
    dob = COALESCE(sqlc.narg('dob'), dob),
    name = COALESCE(sqlc.narg('name'), name),
    email = COALESCE(sqlc.narg('email'), email),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    address = COALESCE(sqlc.narg('address'), address),
    salary = COALESCE(sqlc.narg('salary'), salary),
    account_id = COALESCE(sqlc.narg('account_id'), account_id),
    update_date = NOW()
WHERE employee_id = ?;

-- name: GetEmployee :one
SELECT * FROM employees
WHERE employee_id = ? LIMIT 1;

-- name: ListEmployees :many
SELECT * FROM employees;

-- name: ListEmployeesPaged :many
SELECT * FROM employees
ORDER BY employee_id
LIMIT ? OFFSET ?;

-- name: GetEmployeeByAccountID :one
SELECT * FROM employees
WHERE account_id = ? LIMIT 1;