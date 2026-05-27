-- name: CreateCategory :one
INSERT INTO categories (user_id, name, type, icon, color)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: ListCategoriesByUserID :many
SELECT * FROM categories WHERE user_id = $1 ORDER BY name;

-- name: ListCategoriesByType :many
SELECT * FROM categories WHERE user_id = $1 AND type = $2 ORDER BY name;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;
