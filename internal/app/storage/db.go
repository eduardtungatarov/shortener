package storage

import (
	"context"
	"database/sql"
	"github.com/eduardtungatarov/shortener/internal/app/config"
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

func (s *dbStorage) Set(key, value string) error {
	return nil
}

func (s *dbStorage) Get(key string) (value string, ok bool) {
	return "", false
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