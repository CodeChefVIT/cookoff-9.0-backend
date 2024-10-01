package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type createTestCaseRequest struct {
	ExpectedOutput string         `json:"expected_output" validate:"required,omitempty"`
	Memory         pgtype.Numeric `json:"memory"          validate:"required,omitempty"`
	Input          string         `json:"input"           validate:"required,omitempty"`
	Hidden         *bool          `json:"hidden"          validate:"required,omitempty"`
	QuestionID     string         `json:"question_id"     validate:"required,omitempty"`
	Runtime        pgtype.Numeric `json:"runtime"         validate:"required,omitempty"`
}

func CreateTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	var testCase createTestCaseRequest

	if err := httphelpers.ParseJSON(r, &testCase); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	questionID, _ := uuid.Parse(testCase.QuestionID)

	newTestCase := db.CreateTestCaseParams{
		ExpectedOutput: testCase.ExpectedOutput,
		Memory:         testCase.Memory,
		Input:          testCase.Input,
		Hidden:         *testCase.Hidden,
		QuestionID:     questionID,
		Runtime:        testCase.Runtime,
	}
	newTestCase.ID, _ = uuid.NewV7()

	err := database.Queries.CreateTestCase(context.Background(), newTestCase)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to create test case")
		return
	}

	httphelpers.WriteJSON(w, http.StatusCreated, testCase)
}

func GetTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid test case ID")
		return
	}

	testCase, err := database.Queries.GetTestCase(context.Background(), testcaseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httphelpers.WriteError(w, http.StatusNotFound, "test case does not exist")
			return
		}
		logger.Warnf("%s", err.Error())
		httphelpers.WriteError(w, http.StatusInternalServerError, "some error occured")
		return
	}

	httphelpers.WriteJSON(w, http.StatusOK, testCase)
}

func GetAllTestCasesHandler(w http.ResponseWriter, r *http.Request) {
	testCases, err := database.Queries.GetAllTestCases(context.Background())
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to fetch test cases")
		return
	}

	httphelpers.WriteJSON(w, http.StatusOK, testCases)
}

func UpdateTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid test case ID")
		return
	}

	var payload createTestCaseRequest
	if err := httphelpers.ParseJSON(r, &payload); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	originalTestCase, err := database.Queries.GetTestCase(r.Context(), testcaseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httphelpers.WriteError(w, http.StatusNotFound, "test case does not exist")
			return
		}
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updateData := db.UpdateTestCaseParams{
		ExpectedOutput: originalTestCase.ExpectedOutput,
		Memory:         originalTestCase.Memory,
		Input:          originalTestCase.Input,
		Hidden:         originalTestCase.Hidden,
		Runtime:        originalTestCase.Runtime,
		ID:             testcaseID,
	}

	if payload.Hidden != nil {
		updateData.Hidden = *payload.Hidden
	}
	if payload.ExpectedOutput != "" {
		updateData.ExpectedOutput = payload.ExpectedOutput
	}
	if payload.Memory.Valid {
		updateData.Memory = payload.Memory
	}
	if payload.Input != "" {
		updateData.Input = payload.Input
	}
	if payload.Runtime.Valid {
		updateData.Runtime = payload.Runtime
	}

	updateData.ID = testcaseID

	err = database.Queries.UpdateTestCase(context.Background(), updateData)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to update test case")
		return
	}

	httphelpers.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"message": "Test case updated successfully"},
	)
}

func DeleteTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	testcaseID, err := uuid.Parse(chi.URLParam(r, "testcase_id"))
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid test case ID")
		return
	}

	err = database.Queries.DeleteTestCase(context.Background(), testcaseID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to delete test case")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetTestCaseByQuestionID(w http.ResponseWriter, r *http.Request) {
	questionID, err := uuid.Parse(chi.URLParam(r, "question_id"))
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	testCases, err := database.Queries.GetTestCasesByQuestion(r.Context(), questionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httphelpers.WriteError(w, http.StatusNotFound, "test case not found")
			return
		}
		logger.Warnf("%s", err.Error())
		httphelpers.WriteError(w, http.StatusInternalServerError, "some error occured")
		return
	}

	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "fetched testcases",
		"data":    testCases,
	})
}
