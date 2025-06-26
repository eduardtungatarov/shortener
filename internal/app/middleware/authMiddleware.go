package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/eduardtungatarov/shortener/internal/app/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
)

func (m *Middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var userID string
		var token string
		var setCookie bool

		c, err := req.Cookie(string(config.UserIDKeyName))
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			token, err = buildJWTString()
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			setCookie = true
		} else {
			token = c.Value
		}

		userID = getUserID(token)
		if userID == "" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := req.Context()
		newCtx := context.WithValue(ctx, config.UserIDKeyName, userID)
		req = req.WithContext(newCtx)

		if setCookie {
			cookie := &http.Cookie{
				Name:  string(config.UserIDKeyName),
				Value: token,
			}
			http.SetCookie(res, cookie)
		}

		next.ServeHTTP(res, req)
	})
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const secretKey = "supersecretkey"

func buildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims {
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID: uuid.NewString(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}

	return claims.UserID
}