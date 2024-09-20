// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: submission.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createSubmission = `-- name: CreateSubmission :exec
INSERT INTO "submissions" ("id", "user_id", "question_id", "language_id")
VALUES ($1, $2, $3, $4)
`

type CreateSubmissionParams struct {
	ID         uuid.UUID
	UserID     uuid.NullUUID
	QuestionID uuid.UUID
	LanguageID int32
}

func (q *Queries) CreateSubmission(ctx context.Context, arg CreateSubmissionParams) error {
	_, err := q.db.Exec(ctx, createSubmission,
		arg.ID,
		arg.UserID,
		arg.QuestionID,
		arg.LanguageID,
	)
	return err
}

const getTestCases = `-- name: GetTestCases :many
SELECT id, expected_output, memory, input, hidden, question_id 
FROM "testcases"
WHERE "question_id" = $1
`

func (q *Queries) GetTestCases(ctx context.Context, questionID uuid.UUID) ([]Testcase, error) {
	rows, err := q.db.Query(ctx, getTestCases, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Testcase
	for rows.Next() {
		var i Testcase
		if err := rows.Scan(
			&i.ID,
			&i.ExpectedOutput,
			&i.Memory,
			&i.Input,
			&i.Hidden,
			&i.QuestionID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
