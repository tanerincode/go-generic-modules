package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/tanerincode/go-generic-modules/pkg/config"
	"github.com/tanerincode/go-generic-modules/pkg/storage"
	"log"
	"strconv"
	"time"
)

type Postgres struct {
	Db *pgxpool.Pool
}

func NewPostgres() (storage.Storage, error) {
	// Create a new context with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Always call cancel function to release resources.
	defer cancel()

	// Get the database configuration values from the config.
	host, _ := config.GetConfig("database.host").(string)
	portStr, _ := config.GetConfig("database.port").(string)
	user, _ := config.GetConfig("database.user").(string)
	password, _ := config.GetConfig("database.password").(string)
	name, _ := config.GetConfig("database.name").(string)
	sslMode, _ := config.GetConfig("database.ssl_mode").(string)

	maxOpenConnsString, _ := config.GetConfig("database.max_open_conns").(string)
	maxIdleConnsString, _ := config.GetConfig("database.max_idle_conns").(string)
	connMaxLifetimeString, _ := config.GetConfig("database.conn_max_lifetime").(string)

	port, _ := strconv.Atoi(portStr)
	maxOpenConns, _ := strconv.Atoi(maxOpenConnsString)
	maxIdleConns, _ := strconv.Atoi(maxIdleConnsString)
	connMaxLifetime, _ := strconv.Atoi(connMaxLifetimeString)

	// Create the connection string.
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, name, sslMode)

	// Create a new config for the connection pool.
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Printf("Error parsing config: %v", err)
		return nil, err
	}

	// Set the maximum number of open connections to the database.
	config.MaxConns = int32(maxOpenConns)

	// Set the maximum number of idle connections to the database.
	config.MinConns = int32(maxIdleConns)

	// Set the maximum amount of time a connection may be reused.
	config.MaxConnLifetime = time.Duration(connMaxLifetime) * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		// Log and return the error if there was an issue opening the database.
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Ping the database to ensure the connection is live.
	err = db.Ping(ctx)
	if err != nil {
		// Log and return the error if there was an issue pinging the database.
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	// Return a new Postgres struct with the db and nil for the error.
	return &Postgres{Db: db}, nil
}

func (p *Postgres) Disconnect() error {
	p.Db.Close()
	return nil
}
