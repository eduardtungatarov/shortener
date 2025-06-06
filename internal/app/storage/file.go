package storage

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"os"
)

type storageString struct {
	UUID string `json:"uuid"`
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type fileStorage struct {
	m map[string]string
	file *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func MakeFileStorage(filename string) (*fileStorage, error) {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_RDWR |os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileStorage{
		m: make(map[string]string),
		file: file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (s *fileStorage) Load(ctx context.Context) error {
	v := storageString{}

	for  {
		err := s.decoder.Decode(&v)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		s.m[v.ShortURL] = v.OriginalURL
	}

	return nil
}

func (s *fileStorage) Set(ctx context.Context, key, value string) error {
	v := storageString{
		UUID: uuid.New().String(),
		ShortURL: key,
		OriginalURL: value,
	}
	err := s.encoder.Encode(v)
	if err != nil {
		return err
	}

	s.m[key] = value
	return nil
}

func (s *fileStorage) SetBatch(ctx context.Context, keyValues map[string]string) error {
	return nil
}

func (s *fileStorage) Get(ctx context.Context, key string) (value string, ok bool) {
	v, ok := s.m[key]
	return v, ok
}

func (s *fileStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *fileStorage) Close() error {
	return nil
}