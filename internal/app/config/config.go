// Package config для работы с настройками приложения.
// Вызов LoadFromFlag() заполнит настройки из env, flags, если не задано - дефолтные настройки.
// Настройки вернутся в виде структуры Config.
//
//	cfg := config.LoadFromFlag()
package config

import (
	"flag"
	"os"
	"time"
)

// UserIDKey тип ключа в контексте для поиска userID авторизованного пользователя приложения.
type UserIDKey string

// Дефолтные настройки.
const (
	// DefaultServerHostPort дефолтный адрес запуска HTTP сервера.
	DefaultServerHostPort = "localhost:8080"
	// DefaultBaseURL дефолтный базовый адрес результирующего сокращённого URL.
	DefaultBaseURL = "http://localhost:8080"
	// DefaultFileStoragePath путь до файла, куда сохраняются все сокращенные URL.
	DefaultFileStoragePath = "/tmp/short-url-db.json"
	// DefaultDatabaseDSN строка с адресом подключения к БД.
	DefaultDatabaseDSN = ""
	// DefaultEnableHTTPS по дефолту HTTPS отключен.
	DefaultEnableHTTPS = false
)

// UserIDKeyName имя ключа для поиска в контексте userID авторизованного пользователя сервиса.
const UserIDKeyName UserIDKey = "userId"

// Config настройки сервиса.
type Config struct {
	ServerHostPort  string // адрес запуска HTTP сервера
	EnableHTTPS     bool   // включить https
	BaseURL         string // базовый адрес результирующего сокращённого URL
	FileStoragePath string // путь до файла, куда сохраняются все сокращенные URL
	Database               // настройки бд
}

// Database настройки БД хранения сокращенных ссылок.
type Database struct {
	DSN     string        // строка с адресом подключения к БД
	Timeout time.Duration // таймаут
}

// LoadFromFlag инициализация конфига приложения.
// Приоритет настроек: env, flags, default.
func LoadFromFlag() Config {
	flagServer := flag.String("a", DefaultServerHostPort, "отвечает за адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", DefaultBaseURL, "отвечает за базовый адрес результирующего сокращённого URL")
	flagFileStoragePath := flag.String("f", DefaultFileStoragePath, "путь до файла, куда сохраняются все сокращенные URL")
	databaseDSN := flag.String("d", DefaultDatabaseDSN, "строка с адресом подключения к БД")
	enableHTTPS := flag.Bool("s", DefaultEnableHTTPS, "запуск сервера по защищенному протоколу HTTPS")
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

	sEnv, ok := os.LookupEnv("ENABLE_HTTPS")
	if ok && sEnv == "true" {
		*enableHTTPS = true
	}

	return Config{
		ServerHostPort:  *flagServer,
		EnableHTTPS:     *enableHTTPS,
		BaseURL:         *flagBaseURL,
		FileStoragePath: *flagFileStoragePath,
		Database: Database{
			DSN:     *databaseDSN,
			Timeout: time.Second * 1,
		},
	}
}
