package controllers

import (
	"net/http"

	helpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/go-chi/jwtauth/v5"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		logger.Errof("Refresh token not found: %v", err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "refresh token not found")
		return
	}

	claims, err := jwtauth.VerifyToken(helpers.TokenAuth, cookie.Value)
	if err != nil || claims == nil {
		logger.Errof("Invalid refresh token: %v", err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	userName, ok := claims.PrivateClaims()["username"].(string)
	if !ok {
		logger.Errof("Invalid token claims, email not found")
		httphelpers.WriteError(w, http.StatusUnauthorized, "invalid token claims")
		return
	}

	user, err := database.Queries.GetUserByUsername(r.Context(), userName)
	if err != nil {
		logger.Errof("User not found: %s, err: %v", userName, err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "user not found")
		return
	}

	accessToken, accessExp, err := helpers.GenerateJWT(&user, false)
	if err != nil {
		logger.Errof("Failed to generate new access token for user: %s, err: %v", userName, err)
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

	refreshToken, refreshExp, err := helpers.GenerateJWT(&user, true)
	if err != nil {
		logger.Errof("Failed to generate new refresh token for user: %s, err: %v", userName, err)
		httphelpers.WriteError(w, http.StatusUnauthorized, "failed to generate token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshExp,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	httphelpers.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Token refreshed",
	})
}
