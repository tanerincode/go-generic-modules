package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/tanerincode/go-generic-modules/pkg/storage"
	"log"
	"os"
	"strconv"
	"time"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connectionString string) (storage.Storage, error) {
	// Create a new context with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Always call cancel function to release resources.
	defer cancel()

	// Open a new database connection.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		// Log and return the error if there was an issue opening the database.
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Get the maximum number of open connections from the environment variable.
	maxOpenConns, err := strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))
	if err != nil {
		// Log and return the error if there was an issue converting the environment variable to an integer.
		log.Printf("Error converting MAX_OPEN_CONNS to integer: %v", err)
		return nil, err
	}
	// Set the maximum number of open connections to the database.
	db.SetMaxOpenConns(maxOpenConns)

	// Get the maximum number of idle connections from the environment variable.
	maxIdleConns, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	if err != nil {
		// Log and return the error if there was an issue converting the environment variable to an integer.
		log.Printf("Error converting MAX_IDLE_CONNS to integer: %v", err)
		return nil, err
	}
	// Set the maximum number of idle connections to the database.
	db.SetMaxIdleConns(maxIdleConns)

	// Set the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Minute * 30)

	// Ping the database to ensure the connection is live.
	err = db.PingContext(ctx)
	if err != nil {
		// Log and return the error if there was an issue pinging the database.
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	// Return a new Postgres struct with the db and nil for the error.
	return &Postgres{db: db}, nil
}

func (p *Postgres) Disconnect() error {
	return p.db.Close()
}
