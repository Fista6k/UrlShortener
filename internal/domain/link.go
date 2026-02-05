package domain

import (
	"time"

	"github.com/google/uuid"
)

type OriginalUrl string

type ShortUrl string

type Link struct {
	ID          uuid.UUID   `json: id`
	OriginalUrl OriginalUrl `json: url`
	ShortUrl    ShortUrl    `json: short_url`
	CreatedAt   time.Time   `json: created_at`
}

