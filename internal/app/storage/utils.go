package storage

import (
	"context"
	"errors"
	"github.com/eduardtungatarov/shortener/internal/app/config"
)

func getUserIDOrPanic(ctx context.Context) (string, error) {
	if userID, ok := ctx.Value(config.UserIDKeyName).(string); ok {
		return userID, nil
	}
	return "", errors.New("userID not found or not a string")
}
