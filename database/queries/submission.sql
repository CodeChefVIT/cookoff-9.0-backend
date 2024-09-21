-- name: CreateSubmission :exec
INSERT INTO submissions (id, user_id, question_id, language_id)
VALUES ($1, $2, $3, $4);

-- name: GetTestCases :many
SELECT * 
FROM testcases
WHERE question_id = $1
  AND (CASE WHEN $2 = TRUE THEN hidden = FALSE ELSE TRUE END);

-- name: UpdateSubmission :exec
UPDATE submissions
SET testcases_passed = $1, testcases_failed = $2, runtime = $3, memory = $4
WHERE id = $5;

-- name: GetSubmission :one
SELECT 
    testcases_passed, 
    testcases_failed 
FROM 
    submissions 
WHERE 
    id = $1;
