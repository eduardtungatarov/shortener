package storage

type UserURL struct {
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *UserURL) GetShortURL() string {
	return s.ShortURL
}

func (s *UserURL) GetOriginalURL() string {
	return s.OriginalURL
}