package controllers

import (
	"errors"
	"net/http"
	"os"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	helpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	RegNo string `json:"reg_no"`
	Key   string `json:"fuck_you"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignupRequest

	if err := httphelpers.ParseJSON(r, &payload); err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Key != os.Getenv("SECRET_KEY_FUCKERS") {
		httphelpers.WriteError(w, http.StatusUnauthorized, "I WILL POP YOUR CHERRY BRO")
		return
	}

	passwd := auth.PasswordGenerator(6)
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd), 10)
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = database.Queries.CreateUser(r.Context(), db.CreateUserParams{
		ID:             id,
		Email:          payload.Email,
		RegNo:          payload.RegNo,
		Password:       string(hashed),
		Role:           "user",
		RoundQualified: 0,
		Score:          pgtype.Int4{},
		Name:           payload.Name,
	})
	if err != nil {
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message":  "user added",
		"email":    payload.Email,
		"password": passwd,
	})
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			httphelpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"message": "invalid password",
			})
			return
		}
		httphelpers.WriteError(w, http.StatusInternalServerError, err.Error())
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
