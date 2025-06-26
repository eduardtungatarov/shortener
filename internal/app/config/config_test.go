package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadFromFlag(t *testing.T) {
	type got struct {
		serverHostPortFlag  string
		baseURLFlag         string
		fileStoragePathFlag string
		serverHostPortEnv   string
		baseURLEnv          string
		fileStoragePathEnv  string
		databaseDSNFlag     string
		databaseDSNEnv      string
	}
	type want struct {
		serverHostPort  string
		baseURL         string
		fileStoragePath string
		databaseDSN     string
	}

	tests := []struct {
		name   string
		got  got
		want want
	}{
		{
			name: "not_got_flags_and_envs",
			got: got{
				serverHostPortFlag:  "",
				baseURLFlag:         "",
				fileStoragePathFlag: "",
				databaseDSNFlag:     "",
				serverHostPortEnv:   "",
				baseURLEnv:          "",
			},
			want: want{
				serverHostPort:  "localhost:8080",
				baseURL:         "http://localhost:8080",
				fileStoragePath: "/tmp/short-url-db.json",
				databaseDSN:     "",
			},
		},
		{
			name: "got_envs",
			got: got{
				serverHostPortFlag:  "blabla:80",
				baseURLFlag:         "http://blabla:80",
				fileStoragePathFlag: "/flag/path",
				databaseDSNFlag:     "host=flag user=flag pwd=flag",
				serverHostPortEnv:   "",
				baseURLEnv:          "",
			},
			want: want{
				serverHostPort:  "blabla:80",
				baseURL:         "http://blabla:80",
				fileStoragePath: "/flag/path",
				databaseDSN:     "host=flag user=flag pwd=flag",
			},
		},
		{
			name: "got_flags",
			got: got{
				serverHostPortFlag: "",
				baseURLFlag:        "",
				serverHostPortEnv:  "test:8080",
				baseURLEnv:         "http://test:8080",
				fileStoragePathEnv: "/env/path",
				databaseDSNEnv:     "host=env user=env pwd=env",
			},
			want: want{
				serverHostPort:  "test:8080",
				baseURL:         "http://test:8080",
				fileStoragePath: "/env/path",
				databaseDSN:     "host=env user=env pwd=env",
			},
		},
		{
			name: "got_server_flag_and_server_env",
			got: got{
				serverHostPortFlag: "blabla:80",
				baseURLFlag:        "",
				serverHostPortEnv:  "test:8080",
				baseURLEnv:         "",
			},
			want: want{
				serverHostPort: "test:8080",
				baseURL:        "http://localhost:8080",
				//
				fileStoragePath: "/tmp/short-url-db.json",
				databaseDSN:     "",
			},
		},
		{
			name: "got_server_flag_and_baseurl_env",
			got: got{
				serverHostPortFlag: "blabla:80",
				baseURLFlag:        "",
				serverHostPortEnv:  "",
				baseURLEnv:         "http://test:8080",
			},
			want: want{
				serverHostPort: "blabla:80",
				baseURL:        "http://test:8080",
				//
				fileStoragePath: "/tmp/short-url-db.json",
				databaseDSN:     "",
			},
		},
		{
			name: "got_baseurl_flag_and_server_env",
			got: got{
				serverHostPortFlag: "",
				baseURLFlag:        "http://blabla:80",
				serverHostPortEnv:  "test:8080",
				baseURLEnv:         "",
			},
			want: want{
				serverHostPort: "test:8080",
				baseURL:        "http://blabla:80",
				//
				fileStoragePath: "/tmp/short-url-db.json",
				databaseDSN:     "",
			},
		},
		{
			name: "got_filestoragepathflag",
			got: got{
				fileStoragePathFlag: "/path/flag",
			},
			want: want{
				fileStoragePath: "/path/flag",
				//
				baseURL:        "http://localhost:8080",
				serverHostPort: "localhost:8080",
			},
		},
		{
			name: "got_filestoragepathflag_and_env",
			got: got{
				fileStoragePathFlag: "/path/flag",
				fileStoragePathEnv:  "/path/env",
			},
			want: want{
				fileStoragePath: "/path/env",
				//
				baseURL:        "http://localhost:8080",
				serverHostPort: "localhost:8080",
			},
		},
		{
			name: "not_got_filestoragepathflag_and_env",
			got: got{
				fileStoragePathFlag: "",
				fileStoragePathEnv:  "",
			},
			want: want{
				fileStoragePath: "/tmp/short-url-db.json",
				//
				baseURL:        "http://localhost:8080",
				serverHostPort: "localhost:8080",
			},
		},
		{
			name: "got_databasedsnflag",
			got: got{
				databaseDSNFlag: "host=flag port=port user=myuser password=xxxx dbname=mydb",
			},
			want: want{
				databaseDSN: "host=flag port=port user=myuser password=xxxx dbname=mydb",
				//
				baseURL:         "http://localhost:8080",
				serverHostPort:  "localhost:8080",
				fileStoragePath: "/tmp/short-url-db.json",
			},
		},
		{
			name: "got_databasedsnflag_and_env",
			got: got{
				databaseDSNFlag: "host=flag port=port user=myuser password=xxxx dbname=mydb",
				databaseDSNEnv:  "host=env port=port user=myuser password=xxxx dbname=mydb",
			},
			want: want{
				databaseDSN: "host=env port=port user=myuser password=xxxx dbname=mydb",
				//
				baseURL:         "http://localhost:8080",
				serverHostPort:  "localhost:8080",
				fileStoragePath: "/tmp/short-url-db.json",
			},
		},
		{
			name: "not_got_databasedsnflag_and_env",
			got: got{
				databaseDSNFlag: "",
				databaseDSNEnv:  "",
			},
			want: want{
				databaseDSN: "",
				//
				baseURL:         "http://localhost:8080",
				serverHostPort:  "localhost:8080",
				fileStoragePath: "/tmp/short-url-db.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// настраиваем флаги для теста
			oldOsArgs := os.Args
			os.Args = []string{"cmd"}
			if tt.got.serverHostPortFlag != "" {
				os.Args = append(os.Args, "-a", tt.got.serverHostPortFlag)
			}
			if tt.got.baseURLFlag != "" {
				os.Args = append(os.Args, "-b", tt.got.baseURLFlag)
			}
			if tt.got.fileStoragePathFlag != "" {
				os.Args = append(os.Args, "-f", tt.got.fileStoragePathFlag)
			}
			if tt.got.databaseDSNFlag != "" {
				os.Args = append(os.Args, "-d", tt.got.databaseDSNFlag)
			}

			// настраиваем env для теста
			err := os.Unsetenv("SERVER_ADDRESS")
			assert.NoError(t, err)
			err = os.Unsetenv("BASE_URL")
			assert.NoError(t, err)
			err = os.Unsetenv("FILE_STORAGE_PATH")
			assert.NoError(t, err)
			err = os.Unsetenv("DATABASE_DSN")
			assert.NoError(t, err)
			if tt.got.serverHostPortEnv != "" {
				err := os.Setenv("SERVER_ADDRESS", tt.got.serverHostPortEnv)
				assert.NoError(t, err)
			}
			if tt.got.baseURLEnv != "" {
				err := os.Setenv("BASE_URL", tt.got.baseURLEnv)
				assert.NoError(t, err)
			}
			if tt.got.fileStoragePathEnv != "" {
				err := os.Setenv("FILE_STORAGE_PATH", tt.got.fileStoragePathEnv)
				assert.NoError(t, err)
			}
			if tt.got.databaseDSNEnv != "" {
				err := os.Setenv("DATABASE_DSN", tt.got.databaseDSNEnv)
				assert.NoError(t, err)
			}

			// проверяем
			resetCommandLineFlagSet()
			config := LoadFromFlag()
			assert.Equal(t, tt.want.serverHostPort, config.ServerHostPort, "Ожидается что хост и порт сервера = %v, по факту = %v", tt.want.serverHostPort, config.ServerHostPort)
			assert.Equal(t, tt.want.baseURL, config.BaseURL, "Ожидается что base URL = %v, по факту = %v", tt.want.baseURL, config.BaseURL)
			assert.Equal(t, tt.want.fileStoragePath, config.FileStoragePath, "Ожидается что fileStoragePath = %v, по факту = %v", tt.want.fileStoragePath, config.FileStoragePath)
			assert.Equal(t, tt.want.databaseDSN, config.Database.DSN, "Ожидается что databaseDSN = %v, по факту = %v", tt.want.databaseDSN, config.Database.DSN)

			// восстанавливаем флаги
			os.Args = oldOsArgs
		})
	}
}

func resetCommandLineFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
