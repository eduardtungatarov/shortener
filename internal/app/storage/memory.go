package storage

import (
	"context"
)

type memoryStorage struct {
	m map[string]string
	userLinks map[string][]string
}

func MakeMemoryStorage() *memoryStorage {
	return &memoryStorage{
		m: make(map[string]string),
		userLinks: make(map[string][]string),
	}
}

func (s *memoryStorage) Load(ctx context.Context) error {
	return nil
}

func (s *memoryStorage) Set(ctx context.Context, key, value string) error {
	userId := getUserIDOrPanic(ctx)
	s.m[key] = value
	s.userLinks[userId] = append(s.userLinks[userId], key)
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

func (s *memoryStorage) Get(ctx context.Context, key string) (value string, ok bool) {
	v, ok := s.m[key]
	return v, ok
}

func (s *memoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *memoryStorage) Close() error {
	return nil
}
