package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"github.com/eduardtungatarov/shortener/internal/app"
	"github.com/google/uuid"
	"net/http"
)

var key [32]byte

func Init()  {
	key = sha256.Sum256([]byte("fvbrqaq!33"))
}

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
			userIDCrypted, err = generateUserID()
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
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

func generateUserID() (string, error) {
	userId := []byte(uuid.NewString())

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	encrypted := aesgcm.Seal(nil, nonce, userId, nil)

	return string(encrypted), nil
}

func getUserID(userIDCrypted string) (string, error) {
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	decrypted, err := aesgcm.Open(nil, nonce, []byte(userIDCrypted), nil) // расшифровываем
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
