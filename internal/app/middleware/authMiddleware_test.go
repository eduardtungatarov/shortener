package middleware

import (
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWithAuth проверяет Middleware WithAuth
func TestWithAuth(t *testing.T) {
	m := &Middleware{}

	// Тест для случая, когда куки отсутствуют
	t.Run("when_cookie_is_missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NotNil(t, r.Context().Value(config.UserIDKeyName))
		})

		m.WithAuth(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.NotEmpty(t, res.Header().Get("Set-Cookie"))
	})

	// Тест для случая, когда куки есть, но токен недействительный
	t.Run("when_cookie_is_invalid", func(t *testing.T) {
		tokenString := "invalid.token.string"
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  string(config.UserIDKeyName),
			Value: tokenString,
		})
		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "next handler should not be called with invalid token")
		})

		m.WithAuth(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	// Тест для случая, когда куки валидный токен
	t.Run("when_cookie_is_valid", func(t *testing.T) {
		token, err := buildJWTString()
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  string(config.UserIDKeyName),
			Value: token,
		})
		res := httptest.NewRecorder()

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(config.UserIDKeyName).(string)
			assert.NotEmpty(t, userID)
		})

		m.WithAuth(nextHandler).ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}