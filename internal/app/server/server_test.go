package server

import (
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/storage"
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

func makeMockStorage() storage.Storage {
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
		preloadedStorage storage.Storage
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
			name: "success post",
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
			name: "success get",
			input: input{
				preloadedStorage: func() storage.Storage {
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
			name: "post with empty contentType",
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
			name: "post with incorrect contentType",
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
			name: "post with empty body",
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
			name: "post with incorrect path",
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
			name: "get with empty shortUrl",
			input: input{
				preloadedStorage: func() storage.Storage {
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
			name: "get unexists shortUrl",
			input: input{
				preloadedStorage: func() storage.Storage {
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
			name: "incorrect method",
			input: input{
				preloadedStorage: func() storage.Storage {
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
			h := handlers.MakeHandler(
				tt.input.preloadedStorage,
				host,
			)
			r := getRouter(h)
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

			assert.Equal(t, tt.output.statusCode, resp.StatusCode)
			assert.Equal(t, tt.output.locationHeaderValue, resp.Header.Get("Location"))

			assert.Contains(t, string(respBody), tt.output.response)
		})
	}
}