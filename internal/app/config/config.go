package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerHostPort string
	BaseURL string
}

var flagServer *string
var flagBaseURL *string

func init() {
	flagServer = flag.String("a", "", "отвечает за адрес запуска HTTP-сервера")
	flagBaseURL = flag.String("b", "", "отвечает за базовый адрес результирующего сокращённого URL")
}

func setDefaultValuesToFlags() {
	*flagServer = "localhost:8080"
	*flagBaseURL = "http://localhost:8080"
}

func LoadFromFlag() Config  {
	setDefaultValuesToFlags()
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