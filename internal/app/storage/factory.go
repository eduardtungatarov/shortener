package storage

import (
	"context"
	"github.com/eduardtungatarov/shortener/internal/app/config"
)

type Storage interface {
	Load(ctx context.Context) error
	Set(ctx context.Context, key, value string) error
	SetBatch(ctx context.Context, keyValues map[string]string) error
	Get(ctx context.Context, key string) (value string, ok bool)
	Ping(ctx context.Context) error
	GetByUserID(ctx context.Context) ([]map[string]string, error)
	Close() error
}

func MakeStorage(cfg config.Config) (Storage, error) {
	if cfg.Database.DSN != config.DefaultDatabaseDSN {
		return MakeDBStorage(cfg.Database)
	}

	if cfg.FileStoragePath != config.DefaultFileStoragePath {
		return MakeFileStorage(cfg.FileStoragePath)
	}

	return MakeMemoryStorage(), nil
}
