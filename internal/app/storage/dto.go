package storage

// UserURL полный и сокращенный URL юзера.
type UserURL struct {
	ShortURL    string
	OriginalURL string
}

// GetShortURL получить сокращенный вариант.
func (s *UserURL) GetShortURL() string {
	return s.ShortURL
}

// GetOriginalURL получить полный вариант.
func (s *UserURL) GetOriginalURL() string {
	return s.OriginalURL
}
