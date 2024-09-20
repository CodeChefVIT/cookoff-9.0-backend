package controllers

import (
	"fmt"
	"net/http"

	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/jwtauth/v5"
)

func RoleFromToken (w http.ResponseWriter,r *http.Request, user string) bool {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	role, ok := claims["role"].(string)
	if !ok {
		http.Error(w, "Role not found in token", http.StatusUnauthorized)
		return ok
	}
	if role != user {
		msg := fmt.Sprintf("Access Denied: %s only", role)
		http.Error(w, msg, http.StatusForbidden)
		return false
	}
	return true
}