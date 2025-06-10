package storage

import "context"

func getUserIDOrPanic(ctx context.Context) string  {
	if userId, ok := ctx.Value("userId").(string); ok {
		return userId
	}
	panic("User ID not found or not a string")
}
