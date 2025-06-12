package storage

import "context"

func getUserIDOrPanic(ctx context.Context) string  {
	if userID, ok := ctx.Value("userId").(string); ok {
		return userID
	}
	panic("User ID not found or not a string")
}
