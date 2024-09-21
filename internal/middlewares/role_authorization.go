package middlewares

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
)

func RoleAuthorizationMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok := controllers.RoleFromToken(w, r, requiredRole)
			if !ok {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
