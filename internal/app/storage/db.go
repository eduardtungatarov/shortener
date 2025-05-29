package storage

import (
	"context"
	"database/sql"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/google/uuid"
	"time"
)

type dbStorage struct {
	sqlDB *sql.DB
	timeout time.Duration
}

func MakeDBStorage(cfg config.Database) (*dbStorage, error)  {
	SQLDB, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &dbStorage{
		sqlDB: SQLDB,
		timeout: cfg.Timeout,
	}, nil
}

func (s *dbStorage) Load(ctx context.Context) error {
	_, err := s.sqlDB.ExecContext(ctx, `CREATE TABLE urls (
        uuid UUID PRIMARY KEY,
        short_url VARCHAR(255) NOT NULL,
        original_url TEXT NOT NULL
    );
    CREATE INDEX idx_short_url ON urls (short_url);`)
	return err
}

func (s *dbStorage) Set(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.sqlDB.ExecContext(ctx, `INSERT INTO urls (uuid, short_url, original_url)
		VALUES ($1, $2, $3)
	`, uuid.NewString(), key, value)
	return err
}

func (s *dbStorage) Get(ctx context.Context, key string) (value string, ok bool) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	row := s.sqlDB.QueryRowContext(ctx, `SELECT FROM urls WHERE short_url = $1`, key)

	var originalUrl string
	err := row.Scan(&originalUrl)
	if err != nil {
		return "", false
	}
	return originalUrl, true
}

func (s *dbStorage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := s.sqlDB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *dbStorage) Close() error {
	return s.sqlDB.Close()
}