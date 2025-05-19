package storage

import (
	"encoding/json"
	"github.com/google/uuid"
	"os"
)

type storageString struct {
	UUID string `json:"uuid"`
	ShortUrl string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type storage struct {
	m map[string]string
	file *os.File
	encoder *json.Encoder
}

func MakeStorage(filename string) (*storage, error) {
	file, err := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY |os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &storage{
		m: make(map[string]string),
		file: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (s *storage) Set(key, value string) error {
	v := storageString{
		UUID: uuid.New().String(),
		ShortUrl: key,
		OriginalUrl: value,
	}
	err := s.encoder.Encode(v)
	if err != nil {
		return err
	}

	s.m[key] = value
	return nil
}

func (s *storage) Get(key string) (value string, ok bool) {
	v, ok := s.m[key]
	return v, ok
}