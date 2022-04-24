package data

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Tries to initialize connection with PostgreSQL database.
// In case of failure, logs the error and panics.
// Returns connection instance on success.
func InitSQL(uri string) *pgxpool.Pool {
	conn, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		log.Panicf("Database connection error: %v\n", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Panicf("Database ping error: %v\n", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	return conn
}
