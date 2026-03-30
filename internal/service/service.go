package service

import (
	"context"

	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
)

type ShortenerService struct {
	storage adapter.Storage
	ctx     context.Context
}

func NewShortenerService(ctx context.Context, storage adapter.Storage) *ShortenerService {
	return &ShortenerService{
		storage: storage,
		ctx:     ctx,
	}
}
