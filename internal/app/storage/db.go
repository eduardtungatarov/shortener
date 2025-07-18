package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"time"
)

var ErrConflict = errors.New("data conflict")
var ErrDeleted = errors.New("url deleted")

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

	addDeletedFlagColumn := `
    DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name='urls' 
            AND column_name='deleted_flag'
        ) THEN
            ALTER TABLE urls ADD COLUMN deleted_flag integer default 0;
        END IF;
    END $$;`

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

	_, err = s.sqlDB.ExecContext(ctx, addDeletedFlagColumn)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbStorage) Set(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	userID, err := getUserIDOrPanic(ctx)
	if err != nil {
		return err
	}

	_, err = s.sqlDB.ExecContext(ctx, `INSERT INTO urls (uuid, short_url, original_url, user_uuid)
		VALUES ($1, $2, $3, $4)
	`, uuid.NewString(), key, value, userID)

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

	userID, err := getUserIDOrPanic(ctx)
	if err != nil {
		return err
	}

	for key, value := range keyValues {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO urls (uuid, short_url, original_url, user_uuid)"+
				" VALUES($1, $2, $3, $4)", uuid.NewString(), key, value, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *dbStorage) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	row := s.sqlDB.QueryRowContext(ctx, `SELECT original_url, deleted_flag FROM urls WHERE short_url = $1`, key)

	var originalURL string
	var deletedFlag bool
	err := row.Scan(&originalURL, &deletedFlag)
	if err != nil {
		return "", err
	}

	if deletedFlag {
		return "", ErrDeleted
	}

	return originalURL, nil
}

func (s *dbStorage) GetByUserID(ctx context.Context) ([]map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	userID, err := getUserIDOrPanic(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.sqlDB.QueryContext(ctx, `SELECT original_url, short_url FROM urls WHERE user_uuid = $1`, userID)
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

func (s *dbStorage) DeleteBatch(ctx context.Context, keys []string, userID string) error {
	placeholders := make([]string, len(keys))
	for i := range keys {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("UPDATE urls SET deleted_flag = 1 WHERE short_url IN (%s) AND user_uuid = $%d;",
		strings.Join(placeholders, ", "), len(keys)+1)

	args := make([]interface{}, len(keys)+1)
	for i, key := range keys {
		args[i] = key
	}
	args[len(keys)] = userID

	_, err := s.sqlDB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
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
