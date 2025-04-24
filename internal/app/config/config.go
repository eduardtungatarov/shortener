package config

import "flag"

type Config struct {
	ServerHostPort string
	BaseURL string
}

func LoadFromFlag() Config  {
	a := flag.String("a", "localhost:8080", "отвечает за адрес запуска HTTP-сервера")
	b := flag.String("b", "http://localhost:8080", "отвечает за базовый адрес результирующего сокращённого URL")
	flag.Parse()

	return Config{
		ServerHostPort: *a,
		BaseURL: *b,
	}
}