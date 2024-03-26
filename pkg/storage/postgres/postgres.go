package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tanerincode/go-generic-modules/pkg/config"
	"github.com/tanerincode/go-generic-modules/pkg/storage"
	"log"
	"time"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() (storage.Storage, error) {
	// Create a new context with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Always call cancel function to release resources.
	defer cancel()

	// Get the database configuration values from the config.
	host, _ := config.GetConfig("database.host").(string)
	port, _ := config.GetConfig("database.port").(int)
	user, _ := config.GetConfig("database.user").(string)
	password, _ := config.GetConfig("database.password").(string)
	name, _ := config.GetConfig("database.name").(string)
	maxOpenConns, _ := config.GetConfig("database.max_open_conns").(int)
	maxIdleConns, _ := config.GetConfig("database.max_idle_conns").(int)
	connMaxLifetime, _ := config.GetConfig("database.conn_max_lifetime").(int)

	// Create the connection string.
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	// Open a new database connection.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		// Log and return the error if there was an issue opening the database.
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Set the maximum number of open connections to the database.
	db.SetMaxOpenConns(maxOpenConns)

	// Set the maximum number of idle connections to the database.
	db.SetMaxIdleConns(maxIdleConns)

	// Set the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)

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
