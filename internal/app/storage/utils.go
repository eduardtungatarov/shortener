package storage

import (
	"context"
	"github.com/eduardtungatarov/shortener/internal/app/config"
)

func getUserIDOrPanic(ctx context.Context) string {
	if userID, ok := ctx.Value(config.UserIDKeyName).(string); ok {
		return userID
	}
	panic("User ID not found or not a string")
}
