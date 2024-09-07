package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/google/uuid"
)

type subreq struct {
	SourceCode string `json:"source_code"`
	LanguageID int    `json:"language_id"`
	QuestionID string `json:"question_id"`
}

func SubmitCode(w http.ResponseWriter, r *http.Request) {
	JUDGE0_URI := os.Getenv("JUDGE0_URI")
	ctx := r.Context()

	var req subreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	question_id, _ := uuid.Parse(req.QuestionID)

	payload, err := submission.CreateSubmission(ctx, question_id, req.LanguageID, req.SourceCode)
	if err != nil {
		logger.Errof("Error creating submission: %v", err)
		http.Error(w, "Failed to create submission", http.StatusInternalServerError)
		return
	}

	subID, err := uuid.NewV7()
	if err != nil {
		logger.Errof("Error in generating uuid for submission: %v", err)
		http.Error(w, "Error in generating uuid for submission", http.StatusInternalServerError)
		return
	}

	judge0URL := JUDGE0_URI + "/submissions/batch?base64_encoded=true"
	resp, err := http.Post(judge0URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Errof("Error sending request to Judge0: %v", err)
		http.Error(w, "Failed to send request to Judge0", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	err = submission.StoreTokens(ctx, subID, resp)
	if err != nil {
		logger.Errof("Error storing tokens for submission ID %s: %v", subID, err)
		http.Error(w, "Error storing tokens for the submission", http.StatusInternalServerError)
	}

	user_id, _ := r.Context().Value("user_id").(string)
	userID, _ := uuid.Parse(user_id)
	qID, _ := uuid.Parse(req.QuestionID)
	nullUserID := uuid.NullUUID{UUID: userID, Valid: true}

	err = database.Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         subID,
		UserID:     nullUserID,
		QuestionID: qID,
		LanguageID: int32(req.LanguageID),
	})
	if err != nil {
		logger.Errof("Error creating submission in database: %v", err)
		http.Error(w, "Error creating submission in database", http.StatusInternalServerError)
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
