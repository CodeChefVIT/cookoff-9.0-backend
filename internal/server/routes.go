package server

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers" // Correct import path for your controllers package
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Health Check
	r.Get("/ping", controllers.HealthCheck)

	// POST route for code submission
	r.Post("/api/submission", controllers.SubmitCode) // Corrected reference to the exported SubmitCode function

	return r
}
