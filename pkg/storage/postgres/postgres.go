package postgres

import (
	"context"                                    // Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.
	"database/sql"                               // Package sql provides a generic interface around SQL (or SQL-like) databases.
	_ "github.com/lib/pq"                        // pq is a pure Go Postgres driver for the database/sql package.
	"tanerincode/go-generic-modules/pkg/storage" // Importing the storage package which defines the Storage interface.
	"time"                                       // Package time provides functionality for measuring and displaying time.
)

// Postgres is a struct that holds a sql.DB pointer which represents a pool of zero or more underlying connections.
type Postgres struct {
	db *sql.DB
}

// NewPostgres is a function that takes a connection string and returns an implementation of the Storage interface and an error.
// It opens a new database connection and pings the database to ensure that the connection is live.
func NewPostgres(connectionString string) (storage.Storage, error) {
	// Create a new context with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// The cancel function must be called to release resources even if the context is already done.
	defer cancel()

	// sql.Open opens a database specified by its database driver name and a driver-specific data source name.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		// Return nil and the error if there was an issue opening the database.
		return nil, err
	}

	// PingContext verifies a connection to the database is still alive, establishing a connection if necessary.
	err = db.PingContext(ctx)
	if err != nil {
		// Return nil and the error if there was an issue pinging the database.
		return nil, err
	}

	// Return a new Postgres struct with the db and nil for the error.
	return &Postgres{db: db}, nil
}

// Disconnect is a method on the Postgres struct that closes the database, releasing any open resources.
func (p *Postgres) Disconnect() error {
	return p.db.Close()
}
