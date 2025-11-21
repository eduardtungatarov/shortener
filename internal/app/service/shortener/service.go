package shortener

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/url"

	"github.com/eduardtungatarov/shortener/internal/app/storage"
)

// Storage интерфейс хранилища ссылок.
type Storage interface {
	Set(ctx context.Context, key, value string) error
}

type Service struct {
	storage Storage
	baseURL string
}

func New(
	storage Storage,
	baseURL string,
) *Service {
	return &Service{
		storage: storage,
		baseURL: baseURL,
	}
}

func (s *Service) GetShortenURL(ctx context.Context, URL string) (string, error) {
	key := s.getKey([]byte(URL))
	shortURL, err := url.JoinPath(s.baseURL, key)
	if err != nil {
		return "", err
	}

	err = s.storage.Set(ctx, key, URL)
	isConflict := errors.Is(err, storage.ErrConflict)
	if err != nil {
		if isConflict {
			return shortURL, err
		}
		return "", err
	}
	return shortURL, nil
}

func (h *Service) getKey(url []byte) string {
	hash := md5.Sum(url)
	hashStr := fmt.Sprintf("%x", hash)
	key := hashStr[:7]
	return key
}
