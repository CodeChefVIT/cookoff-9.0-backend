-- name: CreateTestCase :exec
INSERT INTO testcases (
    id, 
    expected_output, 
    memory, 
    input, 
    hidden, 
    question_id, 
    runtime
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdateTestCase :exec
UPDATE testcases
SET 
    expected_output = $1,
    memory = $2,
    input = $3,
    hidden = $4,
    runtime = $5
WHERE id = $6;

-- name: GetTestCase :one
SELECT 
    id, 
    expected_output, 
    memory, 
    input, 
    hidden, 
    question_id, 
    runtime
FROM testcases
WHERE id = $1;

-- name: DeleteTestCase :exec
DELETE FROM testcases
WHERE id = $1;

-- name: GetAllTestCases :many
SELECT 
    *
FROM testcases
ORDER BY id ASC;

-- name: GetTestCasesByQuestion :many
SELECT 
    id, 
    expected_output, 
    memory, 
    input, 
    hidden, 
    question_id, 
    runtime
FROM testcases
WHERE question_id = $1;

-- name: GetPublicTestCasesByQuestion :many
SELECT * FROM testcases
WHERE question_id = $1 AND hidden = false;
