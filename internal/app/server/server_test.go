package server

import (
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/logger"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockStorage struct {
	m map[string]string
}

func makeMockStorage() *mockStorage{
	return &mockStorage{
		m: make(map[string]string),
	}
}

func (s *mockStorage) Set(key, value string) {
	s.m[key] = value
}

func (s *mockStorage) Get(key string) (value string, ok bool) {
	v, ok := s.m[key]
	return v, ok
}

func TestServer(t *testing.T) {
	type input struct {
		preloadedStorage handlers.Storage
		httpMethod string
		requestURI string
		contentType string
		body string
	}
	type output struct {
		statusCode int
		locationHeaderValue string
		response string
	}
	tests := []struct {
		name   string
		input input
		output output
	}{
		{
			name: "success_post",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/",
				contentType: "text/plain",
				body: "https://practicum.yandex.ru/",
			},
			output: output{
				statusCode: 201,
				locationHeaderValue: "",
				response: "http://localhost:8080/",
			},
		},
		{
			name: "success_get",
			input: input{
				preloadedStorage: func() handlers.Storage {
					s := makeMockStorage()
					s.Set("0dd1981", "https://practicum.yandex.ru/")
					return s
				}(),
				httpMethod: "GET",
				requestURI: "/0dd1981",
				contentType: "",
				body: "",
			},
			output: output{
				statusCode: 307,
				locationHeaderValue: "https://practicum.yandex.ru/",
				response: "",
			},
		},
		{
			name: "post_with_empty_contentType",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/",
				contentType: "",
				body: "https://practicum.yandex.ru/",
			},
			output: output{
				statusCode: 400,
				locationHeaderValue: "",
				response: "",
			},
		},
		{
			name: "post_with_incorrect_contentType",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/",
				contentType: "application/json",
				body: "https://practicum.yandex.ru/",
			},
			output: output{
				statusCode: 400,
				locationHeaderValue: "",
				response: "",
			},
		},
		{
			name: "post_with_empty_body",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/",
				contentType: "text/plain",
				body: "",
			},
			output: output{
				statusCode: 400,
				locationHeaderValue: "",
				response: "",
			},
		},
		{
			name: "post_with_incorrect_path",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/bla",
				contentType: "text/plain",
				body: "https://practicum.yandex.ru/",
			},
			output: output{
				statusCode: 405,
				locationHeaderValue: "",
				response: "",
			},
		},
		{
			name: "get_with_empty_shortUrl",
			input: input{
				preloadedStorage: func() handlers.Storage {
					s := makeMockStorage()
					s.Set("0dd1981", "https://practicum.yandex.ru/")
					return s
				}(),
				httpMethod: "GET",
				requestURI: "/",
				contentType: "",
				body: "",
			},
			output: output{
				statusCode: 405,
				locationHeaderValue: "",
				response: "",
			},
		},
		{
			name: "get_unexists_shortUrl",
			input: input{
				preloadedStorage: func() handlers.Storage {
					s := makeMockStorage()
					s.Set("0dd1981", "https://practicum.yandex.ru/")
					return s
				}(),
				httpMethod: "GET",
				requestURI: "/blabla",
				contentType: "",
				body: "",
			},
			output: output{
				statusCode: 400,
				locationHeaderValue: "",
				response: "",
			},
		},
		{
			name: "incorrect_method",
			input: input{
				preloadedStorage: func() handlers.Storage {
					s := makeMockStorage()
					s.Set("0dd1981", "https://practicum.yandex.ru/")
					return s
				}(),
				httpMethod: "PATCH",
				requestURI: "/",
				contentType: "",
				body: "",
			},
			output: output{
				statusCode: 405,
				locationHeaderValue: "",
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// заводим сервер
			log, err := logger.MakeLogger()
			if err != nil {
				panic(err)
			}

			m := middleware.MakeMiddleware(log)
			h := handlers.MakeHandler(
				tt.input.preloadedStorage,
				"http://localhost:8080",
			)

			r := getRouter(h, m)
			ts := httptest.NewServer(r)
			defer ts.Close()

			// подгатавливаем реквест
			req, err := http.NewRequest(tt.input.httpMethod, ts.URL+tt.input.requestURI, strings.NewReader(tt.input.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", tt.input.contentType)

			//шлем запрос на сервер
			client := ts.Client()
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// обрабатываем ответ
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.output.statusCode, resp.StatusCode, "Ожидался http статус ответа %v, а не %v", tt.output.statusCode, resp.StatusCode)
			assert.Equal(t, tt.output.locationHeaderValue, resp.Header.Get("Location"), "Ожидался редирект на url: %v, по факту редирект на: %v", tt.output.locationHeaderValue, resp.Header.Get("Location"))

			assert.Contains(t, string(respBody), tt.output.response, "Ссылка в ответе должна начинаться с %v, получено = %v", tt.output.response, string(respBody))
		})
	}
}