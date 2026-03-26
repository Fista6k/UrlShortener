package service

import (
	"context"

	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
)

type ShortererService struct {
	storage adapter.Storage
	ctx     context.Context
}

func NewShortererService(ctx context.Context, storage adapter.Storage) *ShortererService {
	return &ShortererService{
		storage: storage,
		ctx:     ctx,
	}
}
