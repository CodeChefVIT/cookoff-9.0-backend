package server

import (
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hibiken/asynq"
)

func (s *Server) RegisterRoutes(taskClient *asynq.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByRealIP(50, time.Minute))

	r.Get("/ping", controllers.HealthCheck)
	r.Put("/callback", func(w http.ResponseWriter, r *http.Request) {
		controllers.CallbackUrl(w, r, taskClient)
	})

	r.Post("/user/signup", controllers.SignUp)

	r.Post("/login/user", controllers.LoginHandler)
	r.Post("/token/refresh", controllers.RefreshTokenHandler)
	r.Post("/logout", controllers.Logout)

	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(auth.TokenAuth))
		protected.Use(jwtauth.Authenticator(auth.TokenAuth))
		protected.Use(middlewares.BanCheckMiddleware)

		protected.Get("/me", controllers.MeHandler)
		protected.Get("/protected", controllers.ProtectedHandler)

		protected.Patch("/update/profile", controllers.UpdateUser)

		roundLocked := protected.With(middlewares.CheckRound)
		roundLocked.Get(
			"/question/round",
			controllers.GetQuestionsByRound,
		)
		roundLocked.Post(
			"/submit",
			controllers.SubmitCode,
		)
		roundLocked.Post(
			"/runcode",
			controllers.RunCode,
		)
		roundLocked.Get("/result/{submission_id}", controllers.GetResult)

		adminRoutes := protected.With(middlewares.RoleAuthorizationMiddleware("admin"))
		adminRoutes.Post("/question/create", controllers.CreateQuestion)
		adminRoutes.Get("/questions", controllers.GetAllQuestion)
		adminRoutes.Get("/question/{question_id}", controllers.GetQuestionById)
		adminRoutes.Delete("/question/{question_id}", controllers.DeleteQuestion)
		adminRoutes.Get("/submission/{user_id}", controllers.GetSubmissionByUser)
		adminRoutes.Patch("/question", controllers.UpdateQuestion)
		adminRoutes.Post("/upgrade", controllers.UpgradeUserToRound)
		adminRoutes.Post("/roast", controllers.BanUser)
		adminRoutes.Post("/unroast", controllers.UnbanUser)
		adminRoutes.Post("/round/enable", controllers.EnableRound)
		adminRoutes.Get("/users", controllers.GetAllUsers)
		adminRoutes.Get("/leaderboard", controllers.GetLeaderboard)

		adminRoutes.Post("/testcase", controllers.CreateTestCaseHandler)
		adminRoutes.Put("/testcase/{testcase_id}", controllers.UpdateTestCaseHandler)
		adminRoutes.Get("/testcase/{testcase_id}", controllers.GetTestCaseHandler)
		adminRoutes.Delete("/testcase/{testcase_id}", controllers.DeleteTestCaseHandler)
		adminRoutes.Get("/questions/{question_id}/testcases", controllers.GetTestCaseByQuestionID)
	})

	return r
}
