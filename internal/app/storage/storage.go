package storage

type storage struct {
	m map[string]string
}

func MakeStorage() *storage {
	return &storage{
		m: make(map[string]string),
	}
}

func (s *storage) Set(key, value string) {
	s.m[key] = value
}

func (s *storage) Get(key string) (value string, ok bool) {
	v, ok := s.m[key]
	return v, ok
}