package config

import (
	"flag"
	"os"
)

const (
	DefaultServerHostPort = "localhost:8080"
	DefaultBaseURL        = "http://localhost:8080"
)

type Config struct {
	ServerHostPort string
	BaseURL string
}

func LoadFromFlag() Config  {
	flagServer := flag.String("a", DefaultServerHostPort, "отвечает за адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", DefaultBaseURL, "отвечает за базовый адрес результирующего сокращённого URL")
	flag.Parse()

	aEnv, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		*flagServer = aEnv
	}

	bEnv, ok := os.LookupEnv("BASE_URL")
	if ok {
		*flagBaseURL = bEnv
	}

	return Config{
		ServerHostPort: *flagServer,
		BaseURL: *flagBaseURL,
	}
}
