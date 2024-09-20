package controllers

import (
	"net/http"

	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"

	"github.com/go-chi/jwtauth/v5"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Welcome " + claims["username"].(string),
	})
}
