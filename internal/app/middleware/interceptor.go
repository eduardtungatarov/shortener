package middleware

import (
	"context"

	"github.com/eduardtungatarov/shortener/internal/app/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) > 0 {
		userID := getUserID(authHeaders[0])
		ctx = context.WithValue(ctx, config.UserIDKeyName, userID)
	}

	return handler(ctx, req)
}
