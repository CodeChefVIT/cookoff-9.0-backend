package server

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Health check
	r.Get("/ping", controllers.HealthCheck)

	// Code submission
	r.Post("/submit", controllers.SubmitCode)

	// Questions
	r.Post("/question/create", controllers.CreateQuestion)
	r.Get("/question", controllers.GetAllQuestion)
	r.Get("/question/{question_id}", controllers.GetQuestionById)
	r.Delete("/question/{question_id}", controllers.DeleteQuestion)
	r.Patch("/question/{question_id}", controllers.UpdateQuestion)

	// Apply admin middleware to test case routes
	r.Route("/testcase", func(r chi.Router) {
		// Apply admin middleware to all routes except the GetTestCaseHandler
		r.Group(func(r chi.Router) {
			// r.Use(controllers.AdminOnly) // Apply admin middleware here

			// Test cases routes, restricted to admins
			r.Post("/", controllers.CreateQuestion)                       // Create test case
			r.Patch("/{testcase_id}", controllers.UpdateTestCaseHandler)  // Update test case by ID
			r.Delete("/{testcase_id}", controllers.DeleteTestCaseHandler) // Delete test case by ID
			r.Get("/{testcase_id}", controllers.GetTestCaseHandler)       // Get test case by ID
		})

		r.Get("/", controllers.GetAllTestCasesHandler)

	})

	return r
}
