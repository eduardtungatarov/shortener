package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

var ErrConflict = errors.New("data conflict")

type dbStorage struct {
	sqlDB   *sql.DB
	timeout time.Duration
}

func MakeDBStorage(cfg config.Database) (*dbStorage, error) {
	sqlDB, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &dbStorage{
		sqlDB:   sqlDB,
		timeout: cfg.Timeout,
	}, nil
}

func (s *dbStorage) Load(ctx context.Context) error {
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS urls (
            uuid UUID PRIMARY KEY,
            short_url VARCHAR(255) NOT NULL UNIQUE,
            original_url TEXT NOT NULL
        );
    `
	createShortURLIndexSQL := `
        CREATE INDEX IF NOT EXISTS idx_short_url ON urls (short_url);
    `

	addUserUUIDColumn := `
    DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name='urls' 
            AND column_name='user_uuid'
        ) THEN
            ALTER TABLE urls ADD COLUMN user_uuid UUID;
        END IF;
    END $$;`

	createUserUUIDIndex := `
		CREATE INDEX IF NOT EXISTS idx_user_uuid ON urls (user_uuid);
	`

	_, err := s.sqlDB.ExecContext(ctx, createTableSQL)
	if err != nil {
		return err
	}

	_, err = s.sqlDB.ExecContext(ctx, createShortURLIndexSQL)
	if err != nil {
		return err
	}

	_, err = s.sqlDB.ExecContext(ctx, addUserUUIDColumn)
	if err != nil {
		return err
	}

	_, err = s.sqlDB.ExecContext(ctx, createUserUUIDIndex)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbStorage) Set(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.sqlDB.ExecContext(ctx, `INSERT INTO urls (uuid, short_url, original_url, user_uuid)
		VALUES ($1, $2, $3, $4)
	`, uuid.NewString(), key, value, getUserIDOrPanic(ctx))

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
	}

	return err
}

func (s *dbStorage) SetBatch(ctx context.Context, keyValues map[string]string) error {
	tx, err := s.sqlDB.Begin()
	if err != nil {
		return err
	}

	for key, value := range keyValues {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO urls (uuid, short_url, original_url, user_uuid)"+
				" VALUES($1, $2, $3, $4)", uuid.NewString(), key, value, getUserIDOrPanic(ctx))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *dbStorage) Get(ctx context.Context, key string) (value string, ok bool) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	row := s.sqlDB.QueryRowContext(ctx, `SELECT original_url FROM urls WHERE short_url = $1`, key)

	var originalURL string
	err := row.Scan(&originalURL)
	if err != nil {
		return "", false
	}
	return originalURL, true
}

func (s *dbStorage) GetByUserID(ctx context.Context) ([]map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	rows, err := s.sqlDB.QueryContext(ctx, `SELECT original_url, short_url FROM urls WHERE user_uuid = $1`, getUserIDOrPanic(ctx))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []map[string]string
	for rows.Next() {
		var v UserURL
		err = rows.Scan(&v.OriginalURL, &v.ShortURL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, map[string]string{
			"short_url":    v.ShortURL,
			"original_url": v.OriginalURL,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return urls, nil
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
