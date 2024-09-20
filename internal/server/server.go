package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

func NewServer(taskClient *asynq.Client) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	newServer := &Server{
		port: port,
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(taskClient),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}
