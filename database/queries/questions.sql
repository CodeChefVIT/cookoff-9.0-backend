-- name: GetQuestion :one
SELECT * FROM questions
WHERE id = $1 LIMIT 1;

-- name: GetQuestions :many
SELECT * FROM questions;

-- name: CreateQuestion :one
INSERT INTO questions (id, description, title, "input_format", points, round, constraints, output_format, sample_test_input, sample_test_output, explanation)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions 
WHERE id = $1;

-- name: UpdateQuestion :exec
UPDATE questions
SET description = $1, title = $2, input_format = $3, points = $4, round = $5, constraints = $6, output_format = $7, sample_test_input = $8, sample_test_output = $9, explanation = $10
WHERE id = $11;

-- name: GetQuestionByRound :many
SELECT * FROM questions
WHERE round = $1;
