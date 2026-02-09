package service

import (
	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
)

type ShortererService struct {
	storage adapter.Storage
}

func NewShortererService(storage adapter.Storage) *ShortererService {
	return &ShortererService{
		storage: storage,
	}
}
