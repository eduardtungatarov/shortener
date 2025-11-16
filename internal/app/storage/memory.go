// Package storage варианты хранения ссылок.
package storage

import (
	"context"
	"errors"
)

type memoryStorage struct {
	m         map[string]string
	userLinks map[string][]string
}

// MakeMemoryStorage конструктор стораджа в памяти.
func MakeMemoryStorage() *memoryStorage {
	return &memoryStorage{
		m:         make(map[string]string),
		userLinks: make(map[string][]string),
	}
}

// Load загрузить сторадж. Инициализация.
func (s *memoryStorage) Load(ctx context.Context) error {
	return nil
}

// Set установить.
func (s *memoryStorage) Set(ctx context.Context, key, value string) error {
	userID, err := getUserIDOrPanic(ctx)
	if err != nil {
		return err
	}

	s.m[key] = value
	s.userLinks[userID] = append(s.userLinks[userID], key)
	return nil
}

// SetBatch установить пачкой.
func (s *memoryStorage) SetBatch(ctx context.Context, keyValues map[string]string) error {
	for key, originalURL := range keyValues {
		err := s.Set(ctx, key, originalURL)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get получить по ключу.
func (s *memoryStorage) Get(ctx context.Context, key string) (string, error) {
	v, ok := s.m[key]
	if !ok {
		return "", errors.New("not found")
	}

	return v, nil
}

// GetByUserID получить по юзеру.
func (s *memoryStorage) GetByUserID(ctx context.Context) ([]map[string]string, error) {
	var urls []map[string]string

	userID, err := getUserIDOrPanic(ctx)
	if err != nil {
		return nil, err
	}

	userLinks := s.userLinks[userID]

	for _, v := range userLinks {
		urls = append(urls, map[string]string{
			"short_url":    v,
			"original_url": s.m[v],
		})
	}

	return urls, nil
}

// DeleteBatch удалить пачкой.
func (s *memoryStorage) DeleteBatch(ctx context.Context, keys []string, userID string) error {
	return nil
}

// Ping статуса.
func (s *memoryStorage) Ping(ctx context.Context) error {
	return nil
}

// Close закроем.
func (s *memoryStorage) Close() error {
	return nil
}

func (s *memoryStorage) GetStats(ctx context.Context) (map[string]int, error) {
	return map[string]int{
		"urls":  len(s.m),
		"users": len(s.userLinks),
	}, nil
}
