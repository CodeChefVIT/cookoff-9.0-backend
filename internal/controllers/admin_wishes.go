package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := database.Queries.GetAllUsers(ctx)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Unable to fetch users")
		return
	}
	httphelpers.WriteJSON(w, http.StatusOK, users)
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	users, err := database.Queries.GetLeaderboard(r.Context())
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httphelpers.WriteJSON(w, http.StatusAccepted, users)
}

func UpgradeUserToRound(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserIDs []string `json:"user_ids"`
		Round   float64  `json:"round"`
	}

	if err := httphelpers.ParseJSON(r, &requestBody); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(requestBody.UserIDs) == 0 {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid user_ids format")
		return
	}

	var upgradeParams db.UpgradeUsersToRoundParams
	upgradeParams.Column1 = requestBody.UserIDs
	upgradeParams.RoundQualified = int32(requestBody.Round)
	ctx := r.Context()
	err := database.Queries.UpgradeUsersToRound(ctx, upgradeParams)
	if err != nil {
		log.Println("Error in upgrading the user :- ", err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			"Unable to upgrade users to round",
		)
		return
	}

	httphelpers.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"message": "Users upgraded successfully"},
	)
}

func BanUser(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := httphelpers.ParseJSON(r, &requestBody); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userIDStr, ok := requestBody["user_id"]
	if !ok {
		httphelpers.WriteError(w, http.StatusBadRequest, "user_id must be a string")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	ctx := r.Context()
	err = database.Queries.BanUser(ctx, userID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Unable to ban user")
		return
	}

	httphelpers.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"message": "User banned successfully"},
	)
}

func UnbanUser(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := httphelpers.ParseJSON(r, &requestBody); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userIDStr, ok := requestBody["user_id"]
	if !ok {
		httphelpers.WriteError(w, http.StatusBadRequest, "user_id must be a string")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	ctx := r.Context()
	err = database.Queries.UnbanUser(ctx, userID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, "Unable to unban user")
		return
	}

	httphelpers.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"message": "User unbanned successfully"},
	)
}

func GetSubmissionByUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	uid, err := uuid.Parse(id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	sub, err := database.Queries.GetSubmissionByUser(r.Context(), uuid.NullUUID{UUID: uid, Valid: true})
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(sub)
	httphelpers.WriteJSON(w, http.StatusAccepted, map[string]interface{}{"data": sub})
}

type RoundRequest struct {
	RoundID int `json:"round_id"`
}

func EnableRound(w http.ResponseWriter, r *http.Request) {
	var reqBody RoundRequest
	if err := httphelpers.ParseJSON(r, &reqBody); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	ctx := r.Context()
	redisKey := "round:enabled"
	roundIDStr := fmt.Sprintf("%d", reqBody.RoundID)
	err := database.RedisClient.Set(ctx, redisKey, roundIDStr, 0).Err()
	if err != nil {
		log.Printf("Failed to enable round: %v\n", err)
		httphelpers.WriteError(w, http.StatusInternalServerError, "Failed to enable round")
		return
	}

	httphelpers.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"message": "Round enabled successfully"},
	)
}
