-- name: CreateSubmission :exec
INSERT INTO submissions (id, user_id, question_id, language_id, source_code)
VALUES ($1, $2, $3, $4, $5);

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
WITH RankedSubmissions AS (
  SELECT 
    s.*,
    q.round,
    q.title,
    q.description AS question_description,
    ROW_NUMBER() OVER (
      PARTITION BY s.question_id 
      ORDER BY s.testcases_passed DESC, s.testcases_failed ASC, s.runtime ASC
    ) AS rank
  FROM submissions s
  INNER JOIN questions q ON s.question_id = q.id
  WHERE s.user_id = $1 and s.status = 'DONE'
)
SELECT 
  rs.id,
  rs.question_id,
  rs.user_id,
  rs.testcases_passed,
  rs.testcases_failed,
  rs.runtime,
  rs.submission_time,
  rs.language_id,
  rs.description AS submission_description,
  rs.memory,
  rs.status,
  q.round,
  q.title,
  q.description AS question_description,
  q.points
FROM RankedSubmissions rs
INNER JOIN questions q ON rs.question_id = q.id
WHERE rs.rank = 1;

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

-- name: UpdateScore :exec
WITH best_submissions AS (
    SELECT 
        s.user_id AS user_id,
        s.question_id,
        MAX((s.testcases_passed) * q.points / (s.testcases_passed + s.testcases_failed)::numeric) AS best_score
    FROM submissions s
    INNER JOIN questions q ON s.question_id = q.id
    INNER JOIN users u on s.user_id = u.id 
    WHERE s.user_id = (select user_id from submissions where id = $1) AND q.round = u.round_qualified
    GROUP BY s.user_id, s.question_id
)
UPDATE users
SET score = (
    SELECT SUM(best_score)
    FROM best_submissions
)
WHERE users.id = (select user_id from submissions s where s.id = $1);

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

-- name: GetSubmissionByUser :many
SELECT * FROM submissions WHERE user_id = $1;