package middlewares

import (
	"net/http"
	"strconv"

	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
)

func CheckRound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := auth.GetUserID(w, r)
		user, _ := database.Queries.GetUserById(r.Context(), id)

		round, _ := database.RedisClient.Get(r.Context(), "round:enabled").Result()
		roundVal, _ := strconv.Atoi(round)

		if user.RoundQualified != int32(roundVal) {
			httphelpers.WriteError(w, http.StatusForbidden, "not qualified")
			return
		}
		next.ServeHTTP(w, r)
	})
}
