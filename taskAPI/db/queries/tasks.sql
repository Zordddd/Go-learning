-- name: GetTasks :many
SELECT * FROM tasks;

-- name: CreateTask :one
INSERT INTO tasks (title, text, completed, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks SET title=$2, text=$3, completed=$4, updated_at=$5 WHERE id=$1 RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;