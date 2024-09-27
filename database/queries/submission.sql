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
SET 
    runtime = $1, 
    memory = $2, 
    status = $3,
    testcases_passed = $4,
    testcases_failed = $5
WHERE id = $6;

-- name: UpdateSubmissionStatus :exec
UPDATE submissions
SET status = $1
WHERE id = $2;

-- name: UpdateDescriptionStatus :exec
UPDATE submissions
SET description = $1
WHERE id = $2;

-- name: GetSubmission :one
SELECT 
    testcases_passed, 
    testcases_failed 
FROM 
    submissions 
WHERE 
    id = $1;

-- name: GetSubmissionsWithRoundByUserId :many
SELECT q.round, q.title, q.description, s.*
FROM submissions s
INNER JOIN questions q ON s.question_id = q.id
WHERE s.user_id = $1;

-- name: GetSubmissionByID :one
SELECT
    id,
    question_id,
    testcases_passed,
    testcases_failed,
    runtime,
    memory,
    submission_time,
    description,
    user_id
FROM submissions
WHERE id = $1;

-- name: GetSubmissionStatusByID :one
SELECT
    status
FROM submissions
WHERE id = $1;

-- name: GetSubmissionResultsBySubmissionID :many
SELECT 
    id,
    testcase_id,
    submission_id,
    runtime,
    memory,
    status,
    description
FROM 
    submission_results
WHERE 
    submission_id = $1;