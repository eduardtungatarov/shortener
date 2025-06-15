package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/eduardtungatarov/shortener/internal/app"
	"github.com/google/uuid"
	"net/http"
)

var key [32]byte

func Init() {
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
	userID := uuid.NewString()
	return decrypt(userID, key[:])
}

func getUserID(userIDCrypted string) (string, error) {
	return decrypt(userIDCrypted, key[:])
}

func encrypt(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	ciphertext := gcm.Seal(nonce, nonce, plainText, nil)

	return hex.EncodeToString(ciphertext), nil
}

func decrypt(cipherTextHex string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
