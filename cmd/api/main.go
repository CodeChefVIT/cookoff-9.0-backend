package main

import (
	"fmt"
	"os"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	database "github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/queue"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/validator"
	"github.com/CodeChefVIT/cookoff-backend/internal/server"
	"github.com/CodeChefVIT/cookoff-backend/internal/worker"
	"github.com/hibiken/asynq"
)

func main() {
	// Initialize services
	logger.InitLogger()
	database.Init()
	database.InitCache()
	auth.InitJWT()
	submission.Init(database.RedisClient)
	validator.InitValidator()

	redisURI := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), "6379")
	taskServer, taskClient := queue.InitQueue(redisURI, 2)

	go func() {
		mux := asynq.NewServeMux()
		mux.HandleFunc("submission:process", worker.ProcessSubmissionTask)
		queue.StartQueueServer(taskServer, mux)
	}()

	server := server.NewServer(taskClient)
	logger.Infof("Server started at port 8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
