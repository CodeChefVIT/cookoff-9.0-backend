package database

import (
	"context"
	"fmt"
	"os"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool
var Queries *db.Queries

func Init() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		logger.Errof("Database connection parameters are not set")
		return
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	logger.Infof(dsn)

	var err error
	DBPool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Errof("Unable to connect to database: %v", err)
		panic(err)
	}

	logger.Infof("Connected to the database successfully")
	Queries = db.New(DBPool) // Initialize Queries with the connection pool
	Ping(DBPool)
}

func Ping(pool *pgxpool.Pool) {
	if pool == nil {
		logger.Errof("Database connection is not initialized")
		return
	}

	ctx := context.Background()
	err := pool.Ping(ctx)
	if err != nil {
		logger.Errof("Unable to ping the database: %v", err)
		return
	}

	logger.Infof("Database ping successful")
}
