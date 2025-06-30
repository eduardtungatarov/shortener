package config

import (
	"flag"
	"os"
	"time"
)

type UserIDKey string

const (
	DefaultServerHostPort  = "localhost:8080"
	DefaultBaseURL         = "http://localhost:8080"
	DefaultFileStoragePath = "/tmp/short-url-db.json"
	DefaultDatabaseDSN     = ""

 	UserIDKeyName UserIDKey = "userId"
)

type Config struct {
	ServerHostPort  string
	BaseURL         string
	FileStoragePath string
	Database
}

type Database struct {
	DSN     string
	Timeout time.Duration
}

func LoadFromFlag() Config {
	flagServer := flag.String("a", DefaultServerHostPort, "отвечает за адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", DefaultBaseURL, "отвечает за базовый адрес результирующего сокращённого URL")
	flagFileStoragePath := flag.String("f", DefaultFileStoragePath, "путь до файла, куда сохраняются все сокращенные URL")
	databaseDSN := flag.String("d", DefaultDatabaseDSN, "строка с адресом подключения к БД")
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

	dEnv, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		*databaseDSN = dEnv
	}

	return Config{
		ServerHostPort:  *flagServer,
		BaseURL:         *flagBaseURL,
		FileStoragePath: *flagFileStoragePath,
		Database: Database{
			DSN:     *databaseDSN,
			Timeout: time.Second * 1,
		},
	}
}
