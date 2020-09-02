package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func main() {
	logger := log15adapter.NewLogger(log.New("module", "pgx"))
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	poolConfig.ConnConfig.Logger = logger

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
}
