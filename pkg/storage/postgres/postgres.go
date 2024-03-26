package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"tanerincode/go-generic-modules/pkg/storage"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(connectionString string) (storage.Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) Disconnect() error {
	return p.db.Close()
}
