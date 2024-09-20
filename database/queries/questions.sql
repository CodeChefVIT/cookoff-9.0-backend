-- name: GetQuestion :one
SELECT * FROM "questions" 
WHERE "id" = $1 LIMIT 1;

-- name: GetQuestions :many
SELECT id, description, title, "inputFormat", points, round, constraints, "outputFormat" 
FROM "questions";

-- name: CreateQuestion :one
INSERT INTO "questions" (id, description, title, "inputFormat", points, round, constraints,"outputFormat")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM "questions" 
WHERE "id" = $1;

-- name: UpdateQuestion :exec
UPDATE "questions" 
SET description = $1, title = $2, "inputFormat" = $3, points = $4, round = $5, constraints = $6, "outputFormat" = $7
WHERE id = $8;

-- name: GetQuestionByRound :many
SELECT * FROM "questions"
WHERE round = $1;