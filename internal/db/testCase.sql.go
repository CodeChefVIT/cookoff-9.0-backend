package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// TestCase represents a test case with its relevant fields
type TestCase struct {
	ID             uuid.UUID `json:"id"`
	ExpectedOutput string    `json:"expected_output"`
	Memory         string    `json:"memory"`
	Input          string    `json:"input"`
	Hidden         bool      `json:"hidden"`
	Runtime        string    `json:"runtime"` // Assuming runtime is stored as a string
	QuestionID     uuid.UUID `json:"question_id"`
}

// CreateTestCase inserts a new test case into the database
func CreateTestCase(ctx context.Context, db *sql.DB, testCase TestCase) error {
	fmt.Println("creating")
	query := `
		INSERT INTO testcases (id, expected_output, memory, input, hidden, runtime, question_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.ExecContext(ctx, query, testCase.ID, testCase.ExpectedOutput, testCase.Memory, testCase.Input, testCase.Hidden, testCase.Runtime, testCase.QuestionID)
	return err
}

// GetTestCase retrieves a single test case by its ID
func GetTestCase(ctx context.Context, db *sql.DB, id uuid.UUID) (TestCase, error) {
	var testCase TestCase
	query := `SELECT id, expected_output, memory, input, hidden, runtime, question_id FROM testcases WHERE id=$1`
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(&testCase.ID, &testCase.ExpectedOutput, &testCase.Memory, &testCase.Input, &testCase.Hidden, &testCase.Runtime, &testCase.QuestionID)
	if err != nil {
		return testCase, err
	}
	return testCase, nil
}

// UpdateTestCase updates the expected output and runtime of a specific test case by its ID
func UpdateTestCase(ctx context.Context, db *sql.DB, id uuid.UUID, newExpectedOutput, newRuntime string) error {
	query := `UPDATE testcases SET expected_output=$1, runtime=$2 WHERE id=$3`
	_, err := db.ExecContext(ctx, query, newExpectedOutput, newRuntime, id)
	return err
}

// DeleteTestCase deletes a test case by its ID
func DeleteTestCase(ctx context.Context, db *sql.DB, id uuid.UUID) error {
	query := `DELETE FROM testcases WHERE id=$1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}
