package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// DB instance (assumed to be initialized in your server setup)
var dbInstance *sql.DB

// AdminOnly middleware that checks user role from the database
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assuming user ID is passed through the request context or JWT token
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid user token", http.StatusUnauthorized)
			return
		}

		// Query the database to get the user's role
		role, err := getUserRoleFromDB(context.Background(), userID)
		if err != nil {
			http.Error(w, "Unauthorized: Could not fetch user role", http.StatusInternalServerError)
			return
		}

		// Check if the role is "admin"
		if strings.ToLower(role) != "admin" {
			http.Error(w, "Forbidden: You do not have access to this resource", http.StatusForbidden)
			return
		}

		// If the role is admin, continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// Helper function to extract user ID from the request (e.g., from a JWT or session)
func getUserIDFromRequest(r *http.Request) (uuid.UUID, error) {
	// Example of extracting a user ID from context or JWT (you'll need to adjust this)
	// Assuming the user ID is available as a UUID string in the request context or token
	userIDStr := r.Header.Get("UserID") // Adjust this based on your auth setup
	return uuid.Parse(userIDStr)
}

// Helper function to get the user's role from the database by user ID
func getUserRoleFromDB(ctx context.Context, userID uuid.UUID) (string, error) {
	var role string
	query := "SELECT role FROM users WHERE id = $1" // Replace "users" with your actual user table
	err := dbInstance.QueryRowContext(ctx, query, userID).Scan(&role)
	return role, err
}

// CreateTestCase handles the creation of a new test case
func CreateTestCase(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateTestCase handler invoked")
	var testCase db.TestCase

	// Parse request body into the testCase struct
	if err := json.NewDecoder(r.Body).Decode(&testCase); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a new UUID for the test case
	testCase.ID = uuid.New()

	// Insert the test case into the database
	if err := db.CreateTestCase(context.Background(), dbInstance, testCase); err != nil {
		http.Error(w, "Failed to create test case", http.StatusInternalServerError)
		return
	}

	// Return the created test case as a JSON response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testCase)
}

// GetTestCase retrieves a test case by its ID
func GetTestCase(w http.ResponseWriter, r *http.Request) {
	// Extract testcase_id from the URL
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	// Fetch the test case from the database
	testCase, err := db.GetTestCase(context.Background(), dbInstance, testcaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Test case not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch test case", http.StatusInternalServerError)
		}
		return
	}

	// Return the test case as a JSON response
	json.NewEncoder(w).Encode(testCase)
}

// UpdateTestCase handles updating a test case by its ID
func UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	// Extract testcase_id from the URL
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	// Parse request body for updated values (e.g., ExpectedOutput and Runtime)
	var updateData struct {
		ExpectedOutput string `json:"expected_output"`
		Runtime        string `json:"runtime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the test case in the database
	if err := db.UpdateTestCase(context.Background(), dbInstance, testcaseID, updateData.ExpectedOutput, updateData.Runtime); err != nil {
		http.Error(w, "Failed to update test case", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test case updated successfully"))
}

// DeleteTestCase handles deleting a test case by its ID
func DeleteTestCase(w http.ResponseWriter, r *http.Request) {
	// Extract testcase_id from the URL
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		http.Error(w, "Invalid test case ID", http.StatusBadRequest)
		return
	}

	// Delete the test case from the database
	if err := db.DeleteTestCase(context.Background(), dbInstance, testcaseID); err != nil {
		http.Error(w, "Failed to delete test case", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusNoContent)
}
