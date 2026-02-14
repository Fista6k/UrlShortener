package domain

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID `json: id`
	OriginalUrl string    `json: url`
	ShortUrl    string    `json: short_url`
	CreatedAt   time.Time `json: created_at`
}

func NewLink(url, short_url string) (*Link, error) {
	l := &Link{
		ID:          uuid.New(),
		OriginalUrl: url,
		ShortUrl:    short_url,
		CreatedAt:   time.Now(),
	}

	return l, nil
}

func (l *Link) ToJson() map[string]interface {
}
