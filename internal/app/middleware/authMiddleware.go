package middleware

import (
	"context"
	"errors"
	"github.com/eduardtungatarov/shortener/internal/app"
	"github.com/google/uuid"
	"net/http"
)

func (m *Middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var userID string
		var userIDCrypted string
		var setCookie bool

		c, err := req.Cookie(string(app.UserIDKeyName))
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			userIDCrypted = generateUserID()
			setCookie = true
		} else {
			userIDCrypted = c.Value
		}

		userID, err = getUserID(userIDCrypted)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := req.Context()
		newCtx := context.WithValue(ctx, app.UserIDKeyName, userID)
		req = req.WithContext(newCtx)

		if setCookie {
			cookie := &http.Cookie{
				Name:  string(app.UserIDKeyName),
				Value: userID,
			}
			http.SetCookie(res, cookie)
		}

		next.ServeHTTP(res, req)
	})
}

func generateUserID() string {
	return uuid.NewString()
}

func getUserID(userIDCrypted string) (string, error) {
	return userIDCrypted, nil
}
