package controllers

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/google/uuid"
)

type resultreq struct {
	SubID string `json:"submission_id" validate:"required"`
}

func GetResult(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*2)
	defer cancel()

	var req resultreq

	req.SubID = chi.URLParam(r, "submission_id")
	subid, err := uuid.Parse(req.SubID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	processed, err := submission.CheckStatus(ctx, subid)
	if err != nil {
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			"Internal server error while getting status",
		)
		return
	}

	if processed {
		result, err := submission.GetSubResult(ctx, subid)
		if err != nil {
			httphelpers.WriteError(
				w,
				http.StatusInternalServerError,
				"Internal server error while getting submission result",
			)
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
				httphelpers.WriteError(
					w,
					http.StatusInternalServerError,
					"Internal server error while getting status",
				)
				return
			}

			if processed {
				result, err := submission.GetSubResult(ctx, subid)
				if err != nil {
					httphelpers.WriteError(
						w,
						http.StatusInternalServerError,
						"Internal server error while getting submission result",
					)
					return
				}
				httphelpers.WriteJSON(w, http.StatusOK, result)
				return
			}
		}
	}
}

type GetStatus struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
}

type GetSub struct {
	CompileOutput *string   `json:"compile_output"`
	Memory        *int      `json:"memory"`
	Message       *string   `json:"message"`
	Status        GetStatus `json:"status"`
	Stderr        *string   `json:"stderr"`
	Stdout        *string   `json:"stdout"`
	Time          *string   `json:"time"`
	Token         string    `json:"token"`
}

type Response struct {
	Submissions []GetSub `json:"submissions"`
}


func parseTime(timeStr string) (float64, error) {
	if timeStr == "" {
		log.Println("Time value is empty, setting time to 0 for this submission.")
		return 0, nil
	}

	timeValue, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return 0, err
	}
	return timeValue, nil
}

func HandleCompilationError(
	ctx context.Context,
	idUUID uuid.UUID,
	data GetSub,
	time int,
	testcase uuid.UUID,
	status string,
) error {
	subID, err := uuid.NewV7()
	if err != nil {
		log.Println("Error updating submission for compilation error: ", err)
		return err
	}

	err = database.Queries.CreateSubmissionStatus(ctx, db.CreateSubmissionStatusParams{
		ID:           subID,
		SubmissionID: idUUID,
		TestcaseID:   uuid.NullUUID{UUID: testcase, Valid: true},
		Runtime:      pgtype.Numeric{Int: big.NewInt(int64(time)), Valid: true},
		Memory:       pgtype.Numeric{Int: big.NewInt(int64(*data.Memory)), Valid: true},
		Description:  &data.Status.Description,
		Status:       status,
	})
	if err != nil {
		log.Println("Error creating submission status error: ", err)
		return err
	}
	return nil
}
