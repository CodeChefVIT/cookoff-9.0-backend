package auth

import (
	"net/http"

	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/jwtauth/v5"
)

func RoleFromToken(w http.ResponseWriter, r *http.Request, requiredRole string) bool {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		httphelpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return false
	}
	role := claims["role"].(string)

	return role == requiredRole
}
