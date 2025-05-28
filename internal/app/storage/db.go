package storage

import (
	"context"
	"database/sql"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"time"
)

type DB struct {
	SQLDB *sql.DB
	Timeout time.Duration
}

func MakeDB(cfg config.Database) (*DB, error)  {
	SQLDB, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &DB{
		SQLDB: SQLDB,
		Timeout: cfg.Timeout,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()
	if err := db.SQLDB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}