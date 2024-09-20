package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/google/uuid"
)

type subreq struct {
	SourceCode string `json:"source_code"`
	LanguageID int    `json:"language_id"`
	QuestionID string `json:"question_id"`
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

	question_id, _ := uuid.Parse(req.QuestionID)

	payload, err := submission.CreateSubmission(ctx, question_id, req.LanguageID, req.SourceCode)
	if err != nil {
		logger.Errof("Error creating submission: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to create submission")
		return
	}

	subID, err := uuid.NewV7()
	if err != nil {
		logger.Errof("Error in generating uuid for submission: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Error in generating uuid for submission")
		return
	}

	judge0URL, _ := url.Parse(JUDGE0_URI + "/submissions/batch")

	params := url.Values{}
	params.Add("base64_encoded", "true")
	judge0URL.RawQuery = params.Encode()
	resp, err := http.Post(judge0URL.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Errof("Error sending request to Judge0: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to send request to Judge0")
		return
	}
	defer resp.Body.Close()

	err = submission.StoreTokens(ctx, subID, resp)
	if err != nil {
		logger.Errof("Error storing tokens for submission ID %s: %v", subID, err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Error storing tokens for the submission")
		return
	}

	userID, _ := auth.GetUserID(w, r)
	qID, _ := uuid.Parse(req.QuestionID)
	nullUserID := uuid.NullUUID{UUID: userID, Valid: true}

	qualified, err := auth.VerifyRound(ctx, userID, qID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if !qualified {
		httphelpers.WriteError(w, http.StatusForbidden, "User is not qualified for this round")
		return
	}

	err = database.Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         subID,
		UserID:     nullUserID,
		QuestionID: qID,
		LanguageID: int32(req.LanguageID),
	})
	if err != nil {
		logger.Errof("Error creating submission in database: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Error creating submission in database")
		return
	}

	type response struct {
		SubmissionID string `json:"submission_id"`
	}
	respData := response{
		SubmissionID: subID.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(respData)
}
