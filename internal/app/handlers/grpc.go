package handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eduardtungatarov/shortener/internal/app/storage"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/eduardtungatarov/shortener/internal/contracts/shortener/v1"
)

type grpcHandler struct {
	v1.UnimplementedShortenerServiceServer
	shortenerService ShortenerService
}

func NewGrpcHandler(shortenerService ShortenerService) *grpcHandler {
	return &grpcHandler{
		shortenerService: shortenerService,
	}
}

func (h *grpcHandler) ShortenURL(ctx context.Context, req *v1.URLShortenRequest) (*v1.URLShortenResponse, error) {
	shortURL, err := h.shortenerService.GetShortenURL(ctx, req.Url)
	isConflict := errors.Is(err, storage.ErrConflict)
	if err != nil {
		if isConflict {
			return nil, status.Error(codes.AlreadyExists, "ссылка уже есть такая в бд")
		}
		return nil, err
	}
	return &v1.URLShortenResponse{
		Result: shortURL,
	}, nil
}

func (h *grpcHandler) ExpandURL(ctx context.Context, req *v1.URLExpandRequest) (*v1.URLExpandResponse, error) {
	return &v1.URLExpandResponse{}, nil
}

func (h *grpcHandler) ListUserURLs(ctx context.Context, req *emptypb.Empty) (*v1.UserURLsResponse, error) {
	return &v1.UserURLsResponse{}, nil
}
