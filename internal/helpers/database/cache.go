package database

import (
	"context"
	"fmt"
	"os"
	"time"

	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitCache() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	pswd := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		DB:           0,
		Password:     pswd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  2 * time.Second,
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	logger.Infof("Connected to DragonflyDB")
}
