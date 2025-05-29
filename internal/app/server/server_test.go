package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/eduardtungatarov/shortener/internal/app/handlers"
	"github.com/eduardtungatarov/shortener/internal/app/logger"
	"github.com/eduardtungatarov/shortener/internal/app/middleware"
	"github.com/eduardtungatarov/shortener/internal/app/mocks"
	"github.com/golang/mock/gomock"
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

func (s *mockStorage) Set(ctx context.Context, key, value string) error {
	s.m[key] = value
	return nil
}

func (s *mockStorage) Get(ctx context.Context, key string) (value string, ok bool) {
	v, ok := s.m[key]
	return v, ok
}

func (s *mockStorage) Ping(ctx context.Context) error {
	return nil
}

func TestServer(t *testing.T) {
	type input struct {
		preloadedStorage handlers.Storage
		httpMethod string
		requestURI string
		contentType string
		acceptEncoding string
		contentEncoding string
		body string
	}
	type output struct {
		statusCode int
		locationHeaderValue string
		contentTypeHeaderValue string
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
			name: "success_post_with_gzip",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/",
				contentType: "text/plain",
				contentEncoding: "gzip",
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
					s.Set(context.Background(),"0dd1981", "https://practicum.yandex.ru/")
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
					s.Set(context.Background(),"0dd1981", "https://practicum.yandex.ru/")
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
					s.Set(context.Background(), "0dd1981", "https://practicum.yandex.ru/")
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
					s.Set(context.Background(), "0dd1981", "https://practicum.yandex.ru/")
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
		{
			name: "success_post_api_shorten",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/api/shorten",
				contentType: "application/json",
				body: `{"url":"https://practicum.yandex.ru"}`,
			},
			output: output{
				statusCode: 201,
				contentTypeHeaderValue: "application/json",
				response: `{"result":"http://localhost:8080/`,
			},
		},
		{
			name: "success_post_api_shorten_with_accept_encoding",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/api/shorten",
				contentType: "application/json",
				acceptEncoding: "gzip",
				body: `{"url":"https://practicum.yandex.ru"}`,
			},
			output: output{
				statusCode: 201,
				contentTypeHeaderValue: "application/json",
				response: `{"result":"http://localhost:8080/`,
			},
		},
		{
			name: "success_post_api_shorten_with_content_encoding",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/api/shorten",
				contentType: "application/json",
				contentEncoding: "gzip",
				acceptEncoding: "gzip",
				body: `{"url":"https://practicum.yandex.ru"}`,
			},
			output: output{
				statusCode: 201,
				contentTypeHeaderValue: "application/json",
				response: `{"result":"http://localhost:8080/`,
			},
		},
		{
			name: "post_api_shorten_without_application_json_header",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/api/shorten",
				contentType: "",
				body: `{"url":"https://practicum.yandex.ru"}`,
			},
			output: output{
				statusCode: 400,
				contentTypeHeaderValue: "",
				response: ``,
			},
		},
		{
			name: "post_api_shorten_another_method",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "GET",
				requestURI: "/api/shorten",
				contentType: "application/json",
				body: `{"url":"https://practicum.yandex.ru"}`,
			},
			output: output{
				statusCode: 405,
				contentTypeHeaderValue: "",
				response: ``,
			},
		},
		{
			name: "post_api_shorten_not_json_body",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/api/shorten",
				contentType: "application/json",
				body: ``,
			},
			output: output{
				statusCode: 500,
				contentTypeHeaderValue: "",
				response: ``,
			},
		},
		{
			name: "post_api_shorten_body_without_json_url",
			input: input{
				preloadedStorage: makeMockStorage(),
				httpMethod: "POST",
				requestURI: "/api/shorten",
				contentType: "application/json",
				body: `{"test": "bla"}`,
			},
			output: output{
				statusCode: 400,
				contentTypeHeaderValue: "",
				response: ``,
			},
		},
		{
			name: "get_ping_success",
			input: input{
				preloadedStorage: func() handlers.Storage {
					ctrl := gomock.NewController(t)
					m := mocks.NewMockStorage(ctrl)
					m.EXPECT().Ping(gomock.Any()).Return(nil)
					return m
				}(),
				httpMethod: "GET",
				requestURI: "/ping",
				contentType: "",
				body: "",
			},
			output: output{
				statusCode: 200,
				contentTypeHeaderValue: "",
				response: ``,
			},
		},
		{
			name: "get_ping_error",
			input: input{
				preloadedStorage: func() handlers.Storage {
					ctrl := gomock.NewController(t)
					m := mocks.NewMockStorage(ctrl)
					m.EXPECT().Ping(gomock.Any()).Return(errors.New("ping failed"))
					return m
				}(),
				httpMethod: "GET",
				requestURI: "/ping",
				contentType: "",
				body: "",
			},
			output: output{
				statusCode: 500,
				contentTypeHeaderValue: "",
				response: ``,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// заводим сервер
			log, err := logger.MakeNop()
			if err != nil {
				panic(err)
			}

			m := middleware.MakeMiddleware(log)
			h := handlers.MakeHandler(
				context.Background(),
				tt.input.preloadedStorage,
				"http://localhost:8080",
			)

			r := getRouter(h, m)
			ts := httptest.NewServer(r)
			defer ts.Close()

			// подгатавливаем реквест
			var reqReader io.Reader
			reqReader = strings.NewReader(tt.input.body)
			if tt.input.contentEncoding == "gzip" {
				var b bytes.Buffer
				w := gzip.NewWriter(&b)
				_, err := w.Write([]byte(tt.input.body))
				require.NoError(t, err)
				err = w.Close()
				require.NoError(t, err)
				reqReader = bytes.NewReader(b.Bytes())
			}

			req, err := http.NewRequest(tt.input.httpMethod, ts.URL+tt.input.requestURI, reqReader)
			require.NoError(t, err)
			req.Header.Set("Content-Type", tt.input.contentType)
			req.Header.Set("Accept-Encoding", tt.input.acceptEncoding)
			req.Header.Set("Content-Encoding", tt.input.contentEncoding)

			//шлем запрос на сервер
			client := ts.Client()
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.output.statusCode, resp.StatusCode, "Ожидался http статус ответа %v, а не %v", tt.output.statusCode, resp.StatusCode)
			assert.Equal(t, tt.output.locationHeaderValue, resp.Header.Get("Location"), "Ожидался редирект на url: %v, по факту редирект на: %v", tt.output.locationHeaderValue, resp.Header.Get("Location"))
			if tt.output.contentTypeHeaderValue != "" {
				assert.Equal(t, tt.output.contentTypeHeaderValue, resp.Header.Get("Content-Type"), "Ожидался Content-Type в ответе: %v, по факту: %v", tt.output.contentTypeHeaderValue, resp.Header.Get("Content-Type"))
			}

			if resp.StatusCode == http.StatusCreated {
				body := resp.Body
				if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
					gzipR, err := gzip.NewReader(body)
					require.NoError(t, err)
					defer gzipR.Close()
					body = gzipR
				}

				// обрабатываем ответ
				respBody, err := io.ReadAll(body)
				require.NoError(t, err)

				assert.Contains(t, string(respBody), tt.output.response, "Ссылка в ответе должна начинаться с %v, получено = %v", tt.output.response, string(respBody))
			}
		})
	}
}