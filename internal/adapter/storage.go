package adapter

import "github.com/Fista6k/Url-Shorterer.git/internal/domain"

type storage []*domain.Link

type Storage interface {
	Save(*domain.Link) error
	FindByShortCode(string) (*domain.Link, error)
	FindByURL(string) (*domain.Link, error)
}

func makeStorage() (storage, error) {
	storage := storage{}

	return storage, nil
}
