package storage

type UserURL struct {
	ShortURL    string
	OriginalURL string
}

func (s *UserURL) GetShortURL() string {
	return s.ShortURL
}

func (s *UserURL) GetOriginalURL() string {
	return s.OriginalURL
}
