-- name: CreateSubmission :exec
INSERT INTO "submissions" ("id", "user_id", "question_id", "language_id")
VALUES ($1, $2, $3, $4);

-- name: GetTestCases :many
SELECT * 
FROM "testcases"
WHERE "question_id" = $1;