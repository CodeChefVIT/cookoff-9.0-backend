package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/validator"
	"github.com/google/uuid"
)

type resultreq struct {
	SubID string `json:"submission_id" validate:"required"`
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*2)
	defer cancel()

	var req resultreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.ValidatePayload(w, req); err != nil {
		httphelpers.WriteError(w, http.StatusNotAcceptable, "Please provide values for all required fields.")
		return
	}

	subid, err := uuid.Parse(req.SubID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	processed, err := submission.CheckStatus(ctx, subid)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Internal server error while getting status")
		return
	}

	if processed {
		result, err := submission.GetSubResult(ctx, subid)
		if err != nil {
			httphelpers.WriteError(w, http.StatusInternalServerError, "Internal server error while getting submission result")
			return
		}
		httphelpers.WriteJSON(w, http.StatusOK, result)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			httphelpers.WriteError(w, http.StatusRequestTimeout, "Submission not processed")
			return

		case <-ticker.C:
			processed, err := submission.CheckStatus(ctx, subid)
			if err != nil {
				httphelpers.WriteError(w, http.StatusInternalServerError, "Internal server error while getting status")
				return
			}

			if processed {
				result, err := submission.GetSubResult(ctx, subid)
				if err != nil {
					httphelpers.WriteError(w, http.StatusInternalServerError, "Internal server error while getting submission result")
					return
				}
				httphelpers.WriteJSON(w, http.StatusOK, result)
				return
			}
		}
	}

}
