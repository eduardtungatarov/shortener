package handlers

import (
	"context"

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
	if err != nil {
		return nil, err
	}
	return &v1.URLShortenResponse{
		Result: shortURL,
	}, nil
}

func (h *grpcHandler) ExpandURL(ctx context.Context, req *v1.URLExpandRequest) (*v1.URLExpandResponse, error) {
	URL, err := h.shortenerService.GetFullURL(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.URLExpandResponse{
		Result: URL,
	}, nil
}

func (h *grpcHandler) ListUserURLs(ctx context.Context, req *emptypb.Empty) (*v1.UserURLsResponse, error) {
	return &v1.UserURLsResponse{}, nil
}
