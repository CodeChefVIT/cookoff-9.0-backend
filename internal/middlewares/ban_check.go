package middlewares

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func BanCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			httphelpers.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userId := claims["user_id"].(string)

		userID, err := uuid.Parse(userId)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		user, err := database.Queries.GetUserById(r.Context(), userID)
		if err != nil {
			http.Error(w, "Unable to fetch user data", http.StatusInternalServerError)
			return
		}

		if user.IsBanned {
			http.Error(w, "Your account is banned", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
