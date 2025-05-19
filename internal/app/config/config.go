package config

import (
	"flag"
	"os"
)

const (
	DefaultServerHostPort = "localhost:8080"
	DefaultBaseURL        = "http://localhost:8080"
	FileStoragePath       = "./storage.file"
)

type Config struct {
	ServerHostPort string
	BaseURL string
	FileStoragePath string
}

func LoadFromFlag() Config  {
	flagServer := flag.String("a", DefaultServerHostPort, "отвечает за адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", DefaultBaseURL, "отвечает за базовый адрес результирующего сокращённого URL")
	flagFileStoragePath := flag.String("f", FileStoragePath, "путь до файла, куда сохраняются все сокращенные URL")
	flag.Parse()

	aEnv, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		*flagServer = aEnv
	}

	bEnv, ok := os.LookupEnv("BASE_URL")
	if ok {
		*flagBaseURL = bEnv
	}

	fEnv, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		*flagFileStoragePath = fEnv
	}

	return Config{
		ServerHostPort: *flagServer,
		BaseURL: *flagBaseURL,
		FileStoragePath: *flagFileStoragePath,
	}
}
