package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadFromFlag(t *testing.T) {
	type input struct {
		settedServerHostPortFlag string
		settedBaseURLFlag string
		settedServerHostPortEnv string
		settedBaseURLEnv string
	}
	type output struct {
		serverHostPort string
		baseURL string
	}

	tests := []struct {
		name   string
		input input
		output output
	}{
		{
			name: "not_setted_flags_and_envs",
			input: input{
				settedServerHostPortFlag: "",
				settedBaseURLFlag: "",
				settedServerHostPortEnv: "",
				settedBaseURLEnv: "",
			},
			output: output{
				serverHostPort: "localhost:8080",
				baseURL: "http://localhost:8080",
			},
		},
		{
			name: "setted_envs",
			input: input{
				settedServerHostPortFlag: "blabla:80",
				settedBaseURLFlag: "http://blabla:80",
				settedServerHostPortEnv: "",
				settedBaseURLEnv: "",
			},
			output: output{
				serverHostPort: "blabla:80",
				baseURL: "http://blabla:80",
			},
		},
		{
			name: "setted_flags",
			input: input{
				settedServerHostPortFlag: "",
				settedBaseURLFlag: "",
				settedServerHostPortEnv: "test:8080",
				settedBaseURLEnv: "http://test:8080",
			},
			output: output{
				serverHostPort: "test:8080",
				baseURL: "http://test:8080",
			},
		},
		{
			name: "setted_server_flag_and_server_env",
			input: input{
				settedServerHostPortFlag: "blabla:80",
				settedBaseURLFlag: "",
				settedServerHostPortEnv: "test:8080",
				settedBaseURLEnv: "",
			},
			output: output{
				serverHostPort: "test:8080",
				baseURL: "http://localhost:8080",
			},
		},
		{
			name: "setted_server_flag_and_baseurl_env",
			input: input{
				settedServerHostPortFlag: "blabla:80",
				settedBaseURLFlag: "",
				settedServerHostPortEnv: "",
				settedBaseURLEnv: "http://test:8080",
			},
			output: output{
				serverHostPort: "blabla:80",
				baseURL: "http://test:8080",
			},
		},
		{
			name: "setted_baseurl_flag_and_server_env",
			input: input{
				settedServerHostPortFlag: "",
				settedBaseURLFlag: "http://blabla:80",
				settedServerHostPortEnv: "test:8080",
				settedBaseURLEnv: "",
			},
			output: output{
				serverHostPort: "test:8080",
				baseURL: "http://blabla:80",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// настраиваем флаги для теста
			oldOsArgs := os.Args
			os.Args = []string{"cmd"}
			if tt.input.settedServerHostPortFlag != "" {
				os.Args = append(os.Args, "-a", tt.input.settedServerHostPortFlag)
			}
			if tt.input.settedBaseURLFlag != "" {
				os.Args = append(os.Args, "-b", tt.input.settedBaseURLFlag)
			}

			// настраиваем env для теста
			err := os.Unsetenv("SERVER_ADDRESS")
			assert.NoError(t, err)
			err = os.Unsetenv("BASE_URL")
			assert.NoError(t, err)
			if tt.input.settedServerHostPortEnv != "" {
				err := os.Setenv("SERVER_ADDRESS", tt.input.settedServerHostPortEnv)
				assert.NoError(t, err)
			}
			if tt.input.settedBaseURLEnv != "" {
				err := os.Setenv("BASE_URL", tt.input.settedBaseURLEnv)
				assert.NoError(t, err)
			}

			// проверяем
			config := LoadFromFlag()
			assert.Equal(t, tt.output.serverHostPort, config.ServerHostPort, "Ожидается что хост и порт сервера = %v, по факту = %v", tt.output.serverHostPort, config.ServerHostPort)
			assert.Equal(t, tt.output.baseURL, config.BaseURL, "Ожидается что base URL = %v, по факту = %v", tt.output.baseURL, config.BaseURL)

			// восстанавливаем флаги
			os.Args = oldOsArgs
		})
	}
}
