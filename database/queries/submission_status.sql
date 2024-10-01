-- name: CreateSubmissionStatus :exec
INSERT INTO submission_results (id, submission_id, testcase_id ,status ,runtime, memory, description)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetStatsForFinalSubEntry :many
SELECT 
    runtime, 
    memory,   
    status
FROM submission_results
WHERE submission_id = $1;