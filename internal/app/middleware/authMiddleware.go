package middleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

const userIdNameCookie = "userId"

func (m *Middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var userId string
		var userIdCrypted string
		var setCookie bool

		c, err := req.Cookie(userIdNameCookie)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			userIdCrypted = generateUserId()
			setCookie = true
		} else {
			userIdCrypted = c.Value
		}

		userId, err = getUserId(userIdCrypted)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := req.Context()
		newCtx := context.WithValue(ctx, userIdNameCookie, userId)
		req = req.WithContext(newCtx)

		if setCookie {
			cookie := &http.Cookie{
				Name:    userIdNameCookie,
				Value:   userId,
			}
			http.SetCookie(res, cookie)
		}

		next.ServeHTTP(res, req)
	})
}

func generateUserId() string {
	return uuid.NewString()
}

func getUserId(userIdCrypted string) (string, error) {
	return userIdCrypted, nil
}