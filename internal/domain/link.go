package domain

import (
	"time"
)

type Link struct {
	ID          int       `json: id`
	OriginalUrl string    `json: url`
	ShortUrl    string    `json: short_url`
	CreatedAt   time.Time `json: created_at`
}
