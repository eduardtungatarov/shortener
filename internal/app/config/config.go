package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerHostPort string
	BaseURL string
}

func LoadFromFlag() Config  {
	a := flag.String("a", "localhost:8080", "отвечает за адрес запуска HTTP-сервера")
	b := flag.String("b", "http://localhost:8080", "отвечает за базовый адрес результирующего сокращённого URL")
	flag.Parse()

	aEnv, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		*a = aEnv
	}

	bEnv, ok := os.LookupEnv("BASE_URL")
	if ok {
		*b = bEnv
	}

	return Config{
		ServerHostPort: *a,
		BaseURL: *b,
	}
}