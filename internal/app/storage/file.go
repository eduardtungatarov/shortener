package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"os"
)

type storageString struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserUUID    string `json:"user_uuid"`
}

type fileStorage struct {
	m         map[string]string
	userLinks map[string][]string
	file      *os.File
	encoder   *json.Encoder
	decoder   *json.Decoder
}

func MakeFileStorage(filename string) (*fileStorage, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileStorage{
		m:         make(map[string]string),
		userLinks: make(map[string][]string),
		file:      file,
		encoder:   json.NewEncoder(file),
		decoder:   json.NewDecoder(file),
	}, nil
}

func (s *fileStorage) Load(ctx context.Context) error {
	v := storageString{}

	for {
		err := s.decoder.Decode(&v)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		s.m[v.ShortURL] = v.OriginalURL
		s.userLinks[v.UserUUID] = append(s.userLinks[v.UserUUID], v.ShortURL)
	}

	return nil
}

func (s *fileStorage) Set(ctx context.Context, key, value string) error {
	userID := getUserIDOrPanic(ctx)

	v := storageString{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: value,
		UserUUID:    userID,
	}
	err := s.encoder.Encode(v)
	if err != nil {
		return err
	}

	s.m[key] = value
	s.userLinks[userID] = append(s.userLinks[userID], key)
	return nil
}

func (s *fileStorage) SetBatch(ctx context.Context, keyValues map[string]string) error {
	for key, originalURL := range keyValues {
		err := s.Set(ctx, key, originalURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *fileStorage) Get(ctx context.Context, key string) (string, error) {
	v, ok := s.m[key]
	if !ok {
		return "", errors.New("not found")
	}
	
	return v, nil
}

func (s *fileStorage) GetByUserID(ctx context.Context) ([]map[string]string, error) {
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

func (s *fileStorage) DeleteBatch(ctx context.Context, keys []string, userID string) error {
	return nil
}

func (s *fileStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *fileStorage) Close() error {
	return nil
}
