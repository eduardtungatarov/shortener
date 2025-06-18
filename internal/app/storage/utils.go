package storage

import (
	"context"
	"github.com/eduardtungatarov/shortener/internal/app"
)

func getUserIDOrPanic(ctx context.Context) string {
	if userID, ok := ctx.Value(app.UserIDKeyName).(string); ok {
		return userID
	}
	panic("User ID not found or not a string")
}
