package auth

import (
	"context"
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func GetUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.UUID{}, err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		httphelpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.UUID{}, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return uuid.UUID{}, err
	}

	return userUUID, nil
}

func VerifyRound(ctx context.Context, userID uuid.UUID, questionID uuid.UUID) (bool, error) {
	user, err := database.Queries.GetUserById(ctx, userID)
	if err != nil {
		return false, err
	}

	question, err := database.Queries.GetQuestion(ctx, questionID)
	if err != nil {
		return false, err
	}

	return user.RoundQualified == question.Round, nil
}
