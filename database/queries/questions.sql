-- name: GetQuestion :one
SELECT * FROM questions
WHERE id = $1 LIMIT 1;

-- name: GetQuestions :many
SELECT id, description, title, input_format, points, round, constraints, output_format
FROM questions;

-- name: CreateQuestion :one
INSERT INTO questions (id, description, title, "input_format", points, round, constraints, output_format)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions 
WHERE id = $1;

-- name: UpdateQuestion :exec
UPDATE questions
SET description = $1, title = $2, input_format = $3, points = $4, round = $5, constraints = $6, output_format = $7
WHERE id = $8;

-- name: GetQuestionByRound :many
SELECT * FROM questions
WHERE round = $1;
