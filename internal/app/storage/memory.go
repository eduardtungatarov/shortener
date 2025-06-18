package storage

import (
	"context"
	"errors"
)

type memoryStorage struct {
	m         map[string]string
	userLinks map[string][]string
}

func MakeMemoryStorage() *memoryStorage {
	return &memoryStorage{
		m:         make(map[string]string),
		userLinks: make(map[string][]string),
	}
}

func (s *memoryStorage) Load(ctx context.Context) error {
	return nil
}

func (s *memoryStorage) Set(ctx context.Context, key, value string) error {
	userID := getUserIDOrPanic(ctx)
	s.m[key] = value
	s.userLinks[userID] = append(s.userLinks[userID], key)
	return nil
}

func (s *memoryStorage) SetBatch(ctx context.Context, keyValues map[string]string) error {
	for key, originalURL := range keyValues {
		err := s.Set(ctx, key, originalURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *memoryStorage) Get(ctx context.Context, key string) (string, error) {
	v, ok := s.m[key]
	if !ok {
		return "", errors.New("not found")
	}

	return v, nil
}

func (s *memoryStorage) GetByUserID(ctx context.Context) ([]map[string]string, error) {
	var urls []map[string]string

	userLinks := s.userLinks[getUserIDOrPanic(ctx)]

	for _, v := range userLinks {
		urls = append(urls, map[string]string{
			"short_url":    v,
			"original_url": s.m[v],
		})
	}

	return urls, nil
}

func (s *memoryStorage) DeleteBatch(ctx context.Context, keys []string, userID string) error {
	return nil
}

func (s *memoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *memoryStorage) Close() error {
	return nil
}
