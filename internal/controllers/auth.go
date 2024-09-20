package controllers

import (
	"net/http"

	helpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := httphelpers.ParseJSON(r, &req); err != nil {
		logger.Errof("Invalid request payload: %v", err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "invalid request payload")
		return
	}

	user, err := database.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		logger.Errof("User not found or invalid credentials for email: %s, err: %v", req.Email, err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if user.Password != req.Password {
		logger.Errof("Invalid password for user: %s", user.Email)
		httphelpers.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	accessToken, accessExp, err := helpers.GenerateJWT(&user, false)
	if err != nil {
		logger.Errof("Failed to generate access token for user: %s, err: %v", user.Email, err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "failed to generate token")
		return
	}

	refreshToken, refreshExp, err := helpers.GenerateJWT(&user, true)
	if err != nil {
		logger.Errof("Failed to generate refresh token for user: %s, err: %v", user.Email, err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "failed to generate token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    accessToken,
		Expires:  accessExp,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExp,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	data := map[string]any{
		"username": user.Name,
		"round":    user.RoundQualified,
		"score":    user.Score,
	}

	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Login successful",
		"user":    data,
	})
}
