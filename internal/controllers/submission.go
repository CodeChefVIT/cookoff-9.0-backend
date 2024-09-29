package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/validator"
	"github.com/google/uuid"
)

type subreq struct {
	SourceCode string `json:"source_code" validate:"required"`
	QuestionID string `json:"question_id" validate:"required"`
	LanguageID int    `json:"language_id" validate:"required"`
}

var JUDGE0_URI = os.Getenv("JUDGE0_URI")

func SubmitCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req subreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.ValidatePayload(w, req); err != nil {
		httphelpers.WriteError(
			w,
			http.StatusNotAcceptable,
			"Please provide values for all required fields.",
		)
		return
	}

	userID, _ := auth.GetUserID(w, r)
	nullUserID := uuid.NullUUID{UUID: userID, Valid: true}

	question_id, err := uuid.Parse(req.QuestionID)
	if err != nil {
		httphelpers.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid question id, unable to parse it to uuid",
		)
		return
	}

	qualified, err := auth.VerifyRound(ctx, userID, question_id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if !qualified {
		httphelpers.WriteError(w, http.StatusForbidden, "User is not qualified for this round")
		return
	}

	payload, testcase_id, err := submission.CreateSubmission(ctx, question_id, req.LanguageID, req.SourceCode)
	if err != nil {
		if errors.Is(err, submission.ErrNoTestCases) {
			httphelpers.WriteError(
				w,
				http.StatusNotFound,
				err.Error(),
			)
			return
		}
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to create submission payload: %v", err),
		)
		return
	}

	subID, err := uuid.NewV7()
	if err != nil {
		logger.Errof("Error in generating uuid for submission: %v", err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error in generating uuid for submission: %v", err),
		)
		return
	}

	judge0URL, err := url.Parse(JUDGE0_URI + "/submissions/batch")
	if err != nil {
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error parsing Judge0 URL: %v", err),
		)
		return
	}

	params := url.Values{}
	params.Add("base64_encoded", "true")

	resp, err := submission.SendToJudge0(judge0URL, params, payload)
	if err != nil {
		logger.Errof("Error sending request to Judge0: %v", err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error sending request to Judge0: %v", err),
		)
		return
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errof("Error reading response body from Judge0: %v", err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error reading response body from Judge0: %v", err),
		)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		logger.Errof(
			"Unexpected status code from Judge0: %d, error: %v",
			resp.StatusCode,
			string(respBytes),
		)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Internal server error!")
		return
	}
	err = submission.StoreTokens(ctx, subID, respBytes, testcase_id)
	if err != nil {
		logger.Errof("Error storing tokens for submission ID %s: %v", subID, err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error storing tokens for submission ID %s: %v", subID, err),
		)
		return
	}

	err = database.Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         subID,
		UserID:     nullUserID,
		QuestionID: question_id,
		LanguageID: int32(req.LanguageID),
		SourceCode: req.SourceCode,
	})
	if err != nil {
		logger.Errof("Error creating submission in database: %v", err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error creating submission in database: %v", err),
		)
		return
	}

	type response struct {
		SubmissionID string `json:"submission_id"`
	}
	respData := response{
		SubmissionID: subID.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respData); err != nil {
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error encoding response: %v", err),
		)
		return
	}
}
