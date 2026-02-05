package service

import (
	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
)

type ShortererService struct {
	storage domain.LinkRepository
}

func NewShortererService(storage domain.LinkRepository) *ShortererService {
	return &ShortererService{
		storage: storage,
	}
}
