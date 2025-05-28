package storage

import (
	"context"
	"database/sql"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"time"
)

type DB struct {
	SqlDB *sql.DB
	Timeout time.Duration
}

func MakeDB(cfg config.Database) (*DB, error)  {
	sqlDB, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &DB{
		SqlDB: sqlDB,
		Timeout: cfg.Timeout,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, db.Timeout)
	defer cancel()
	if err := db.SqlDB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}