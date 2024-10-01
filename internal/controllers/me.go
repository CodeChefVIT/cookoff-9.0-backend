package controllers

import (
	"fmt"
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	helpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DashboardSubmission struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`
	MaxScore    int     `json:"max_score"`
}

type UpdateUserReq struct {
	ID    uuid.UUID `json:"id"`
	RegNo string    `json:"reg_no"`
	Name  string    `json:"name"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetUserID(w, r)
	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := database.Queries.GetUserById(r.Context(), id)
	if err != nil {
		logger.Errof("Failed to get user: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	submissions, err := database.Queries.GetSubmissionsWithRoundByUserId(
		r.Context(),
		uuid.NullUUID{UUID: id, Valid: true},
	)
	if err != nil {
		logger.Errof("Failed to get submissions: %v", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to get submissions")
		return
	}

	submissionsByRound := make(map[string][]DashboardSubmission)

	for _, submission := range submissions {
		round := fmt.Sprint(submission.Round)

		pointsPerSubmission := float64(
			submission.Points,
		) / float64(
			submission.TestcasesPassed.Int32+submission.TestcasesFailed.Int32,
		)
		submissionsByRound[round] = append(submissionsByRound[round], DashboardSubmission{
			Title:       submission.Title,
			Description: submission.QuestionDescription,
			Score:       float64(submission.TestcasesPassed.Int32) * pointsPerSubmission,
			MaxScore:    int(submission.Points),
		})
	}

	data := map[string]any{
		"username":    user.Name,
		"round":       user.RoundQualified,
		"score":       user.Score,
		"submissions": submissionsByRound,
	}
	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "User details fetched successfully",
		"data":    data,
	})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updateUser UpdateUserReq
	if err := httphelpers.ParseJSON(r, &updateUser); err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "could not parse body")
		return
	}
	id, _ := helpers.GetUserID(w, r)

	user, err := database.Queries.GetUserById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			httphelpers.WriteError(w, http.StatusNotFound, err.Error())
			return
		} else {
			logger.Infof("received error from database %v", err.Error())
			httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if updateUser.RegNo != "" {
		user.RegNo = updateUser.RegNo
	}
	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}

	err = database.Queries.UpdateProfile(ctx, db.UpdateProfileParams{
		ID:    id,
		RegNo: user.RegNo,
		Name:  user.Name,
	})

	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "profile updated",
		"data": map[string]any{
			"reg_no": user.RegNo,
			"name":   user.Name,
		},
	})

}
