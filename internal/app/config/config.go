// Package config для работы с настройками приложения.
// Вызов LoadFromFlag() заполнит настройки из env, flags, config file, если не задано - дефолтные настройки.
// Настройки вернутся в виде структуры Config.
//
//	cfg := config.LoadFromFlag()
package config

import (
	"encoding/json"
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
	// DefaultConfigPath дефолтный путь к файлу конфигурации.
	DefaultConfigPath = ""
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

// JSONConfig структура для парсинга JSON конфигурации.
type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// LoadFromFlag инициализация конфига приложения.
// Приоритет настроек: env, flags, config file, default.
func LoadFromFlag() Config {
	flagServer := flag.String("a", DefaultServerHostPort, "отвечает за адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", DefaultBaseURL, "отвечает за базовый адрес результирующего сокращённого URL")
	flagFileStoragePath := flag.String("f", DefaultFileStoragePath, "путь до файла, куда сохраняются все сокращенные URL")
	databaseDSN := flag.String("d", DefaultDatabaseDSN, "строка с адресом подключения к БД")
	enableHTTPS := flag.Bool("s", DefaultEnableHTTPS, "запуск сервера по защищенному протоколу HTTPS")
	configPath := flag.String("c", DefaultConfigPath, "путь к файлу конфигурации в формате JSON")
	flag.Parse()

	jsonConfig := loadJSONConfig(*configPath)

	if envVal, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		*flagServer = envVal
	} else if jsonConfig.ServerAddress != "" {
		*flagServer = jsonConfig.ServerAddress
	}

	if envVal, ok := os.LookupEnv("BASE_URL"); ok {
		*flagBaseURL = envVal
	} else if jsonConfig.BaseURL != "" {
		*flagBaseURL = jsonConfig.BaseURL
	}

	if envVal, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		*flagFileStoragePath = envVal
	} else if jsonConfig.FileStoragePath != "" {
		*flagFileStoragePath = jsonConfig.FileStoragePath
	}

	if envVal, ok := os.LookupEnv("DATABASE_DSN"); ok {
		*databaseDSN = envVal
	} else if jsonConfig.DatabaseDSN != "" {
		*databaseDSN = jsonConfig.DatabaseDSN
	}

	if envVal, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		*enableHTTPS = envVal == "true"
	} else if jsonConfig.EnableHTTPS {
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

// loadJSONConfig загружает конфигурацию из JSON файла.
func loadJSONConfig(configPath string) JSONConfig {
	if configPath == "" {
		if envVal, ok := os.LookupEnv("CONFIG"); ok {
			configPath = envVal
		} else {
			return JSONConfig{}
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return JSONConfig{}
	}

	var jsonConfig JSONConfig
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return JSONConfig{}
	}

	return jsonConfig
}
