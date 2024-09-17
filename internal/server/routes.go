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
	r.Post("/submit", controllers.SubmitCode)

	// Questions
	r.Post("/question/create", controllers.CreateQuestion)
	r.Get("/question", controllers.GetAllQuestion)
	r.Get("/question/{question_id}", controllers.GetQuestionById)
	r.Delete("/question/{question_id}", controllers.DeleteQuestion)
	r.Patch("/question/{question_id}", controllers.UpdateQuestion)

	return r
}
