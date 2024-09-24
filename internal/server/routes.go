package server

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/hibiken/asynq"
)

func (s *Server) RegisterRoutes(taskClient *asynq.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", controllers.HealthCheck)
	r.Put("/callback", func(w http.ResponseWriter, r *http.Request) {
		controllers.CallbackUrl(w, r, taskClient)
	})
	r.Get("/testcase", controllers.GetAllTestCasesHandler)

	r.Post("/user/signup", controllers.SignUp)

	r.Post("/login/user", controllers.LoginHandler)
	r.Post("/token/refresh", controllers.RefreshTokenHandler)
	r.Post("/logout", controllers.Logout)

	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(auth.TokenAuth))
		protected.Use(jwtauth.Authenticator(auth.TokenAuth))

		protected.Get("/me", controllers.MeHandler)
		protected.Get("/protected", controllers.ProtectedHandler)
		protected.Post("/submit", controllers.SubmitCode)
		protected.Post("/runcode", controllers.RunCode)
		protected.Get("/result", controllers.GetResult)
		protected.Get("/question/round", controllers.GetQuestionsByRound)

		adminRoutes := protected.With(middlewares.RoleAuthorizationMiddleware("admin"))

		adminRoutes.Post("/question/create", controllers.CreateQuestion)
		adminRoutes.Get("/questions", controllers.GetAllQuestion)
		adminRoutes.Get("/question/{question_id}", controllers.GetQuestionById)
		adminRoutes.Delete("/question/{question_id}", controllers.DeleteQuestion)
		adminRoutes.Patch("/question", controllers.UpdateQuestion)
	})

	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(auth.TokenAuth))
		protected.Use(jwtauth.Authenticator(auth.TokenAuth))
		protected.Use(middlewares.RoleAuthorizationMiddleware("admin"))

		protected.Post("/testcase", controllers.CreateTestCaseHandler)
		protected.Put("/testcase/{testcase_id}", controllers.UpdateTestCaseHandler)
		protected.Delete("/testcase/{testcase_id}", controllers.DeleteTestCaseHandler)
		protected.Get("/testcase/{testcase_id}", controllers.GetTestCaseHandler)
	})

	return r
}
