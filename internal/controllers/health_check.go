package controllers

import (
	"net/http"

	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httphelpers.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "pong",
	})
}
