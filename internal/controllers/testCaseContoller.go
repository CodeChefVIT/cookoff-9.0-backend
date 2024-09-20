package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestCase represents a test case with its relevant fields
type TestCase struct {
	ID             uuid.UUID `json:"id"`
	ExpectedOutput string    `json:"expected_output"`
	Memory         string    `json:"memory"`
	Input          string    `json:"input"`
	Hidden         bool      `json:"hidden"`
	QuestionID     uuid.UUID `json:"question_id"`
}

// RoleKey is the key used to store and retrieve the user's role in the context
type RoleKey string

// AdminOnly is a middleware that ensures only users with the "admin" role can access the route
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user's role from the context (modify this based on your auth system)
		role := r.Context().Value(RoleKey("role"))

		// Check if the role is admin
		if role != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		// Call the next handler in the chain if the user is an admin
		next.ServeHTTP(w, r)
	})
}

// CreateTestCase inserts a new test case into the database
func CreateTestCase(ctx context.Context, db *pgxpool.Pool, testCase TestCase) error {
	query := `
		INSERT INTO testcases (id, expected_output, memory, input, hidden, question_id) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(ctx, query, testCase.ID, testCase.ExpectedOutput, testCase.Memory, testCase.Input, testCase.Hidden, testCase.QuestionID)
	if err != nil {
		log.Printf("Error creating test case: %v\nQuery: %s\n", err, query)
		return err
	}
	return nil
}

// GetTestCase retrieves a single test case by its ID
func GetTestCase(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (TestCase, error) {
	var testCase TestCase
	query := `SELECT id, expected_output, memory, input, hidden, question_id FROM testcases WHERE id=$1`
	err := db.QueryRow(ctx, query, id).Scan(
		&testCase.ID,
		&testCase.ExpectedOutput,
		&testCase.Memory,
		&testCase.Input,
		&testCase.Hidden,
		&testCase.QuestionID,
	)
	if err != nil {
		log.Printf("Error fetching test case: %v\nQuery: %s\n", err, query)
		return testCase, err
	}
	return testCase, nil
}

// UpdateTestCase updates the expected output of a specific test case by its ID
func UpdateTestCase(ctx context.Context, db *pgxpool.Pool, id uuid.UUID, newExpectedOutput string) error {
	query := `UPDATE testcases SET expected_output=$1 WHERE id=$2`
	_, err := db.Exec(ctx, query, newExpectedOutput, id)
	if err != nil {
		log.Printf("Error updating test case: %v\nQuery: %s\n", err, query)
		return err
	}
	return nil
}

// DeleteTestCase deletes a test case by its ID
func DeleteTestCase(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) error {
	query := `DELETE FROM testcases WHERE id=$1`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting test case: %v\nQuery: %s\n", err, query)
		return err
	}
	return nil
}

// CreateTestCaseHandler handles creating a new test case
func CreateTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateTestCase handler invoked")

	var testCase TestCase
	if err := json.NewDecoder(r.Body).Decode(&testCase); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	testCase.ID = uuid.New() // Generate a new UUID for the test case

	// Call the CreateTestCase function
	err := CreateTestCase(context.Background(), database.DBPool, testCase)
	if err != nil {
		http.Error(w, "Failed to create test case", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testCase) // Respond with the created test case
}

// GetTestCaseHandler retrieves a test case by its ID
func GetTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	// Call the GetTestCase function
	testCase, err := GetTestCase(context.Background(), database.DBPool, testcaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Test case not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch test case", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testCase) // Respond with test case data
}

// GetAllTestCases retrieves all test cases from the database
func GetAllTestCases(ctx context.Context, db *pgxpool.Pool) ([]TestCase, error) {
	var testCases []TestCase
	query := `SELECT id, expected_output, memory, input, hidden, question_id FROM testcases`

	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Printf("Error fetching all test cases: %v\nQuery: %s\n", err, query)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and scan the data into the TestCase struct
	for rows.Next() {
		var testCase TestCase
		err := rows.Scan(
			&testCase.ID,
			&testCase.ExpectedOutput,
			&testCase.Memory,
			&testCase.Input,
			&testCase.Hidden,
			&testCase.QuestionID,
		)
		if err != nil {
			log.Printf("Error scanning test case: %v\n", err)
			return nil, err
		}
		testCases = append(testCases, testCase)
	}

	return testCases, nil
}

// GetAllTestCasesHandler handles retrieving all test cases
func GetAllTestCasesHandler(w http.ResponseWriter, r *http.Request) {
	// Call the GetAllTestCases function
	testCases, err := GetAllTestCases(context.Background(), database.DBPool)
	if err != nil {
		http.Error(w, "Failed to fetch test cases", http.StatusInternalServerError)
		return
	}

	// Respond with the retrieved test cases as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testCases)
}

// UpdateTestCaseHandler handles updating a test case by its ID
func UpdateTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	var updateData struct {
		ExpectedOutput string `json:"expected_output"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the UpdateTestCase function
	err = UpdateTestCase(context.Background(), database.DBPool, testcaseID, updateData.ExpectedOutput)
	if err != nil {
		http.Error(w, "Failed to update test case", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Test case updated successfully"})
}

// DeleteTestCaseHandler handles deleting a test case by its ID
func DeleteTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteTestCase function
	err = DeleteTestCase(context.Background(), database.DBPool, testcaseID)
	if err != nil {
		http.Error(w, "Failed to delete test case", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
