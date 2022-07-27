package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

// Database - It's the struct that represents our Database
type Database struct {
	Client *sqlx.DB
}

// NewDatabase - Instantiates our database
func NewDatabase() (*Database, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s passwowrd=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	dbConn, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{Client: dbConn}, nil
}

// Ping - Ping our database to check its health.
func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}