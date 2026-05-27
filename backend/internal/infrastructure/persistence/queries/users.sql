-- name: CreateUser :one
INSERT INTO users (email, email_hash, password_hash, full_name, phone, role, super_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmailHash :one
SELECT * FROM users WHERE email_hash = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUser :one
UPDATE users
SET email = COALESCE($2, email),
    email_hash = COALESCE($3, email_hash),
    full_name = COALESCE($4, full_name),
    phone = COALESCE($5, phone),
    role = COALESCE($6, role),
    active = COALESCE($7, active),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: SetUserActive :one
UPDATE users SET active = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
