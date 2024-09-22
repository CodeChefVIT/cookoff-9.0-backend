package middlewares

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
)

func RoleAuthorizationMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok := auth.RoleFromToken(w, r, requiredRole)
			if !ok {
				httphelpers.WriteError(w, http.StatusUnauthorized, "Not Authorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
