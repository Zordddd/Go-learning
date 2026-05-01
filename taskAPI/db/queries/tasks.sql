-- name: GetTasks :many
SELECT *
FROM tasks
ORDER BY created_at DESC;

-- name: GetTask :one
SELECT *
FROM tasks
WHERE id = $1
LIMIT 1;

-- name: CreateTask :one
INSERT INTO tasks (title, text, completed)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET title      = COALESCE($2, title),
    text       = COALESCE($3, text),
    completed  = COALESCE($4, completed),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: DeleteTask :exec
DELETE
FROM tasks
WHERE id = $1;