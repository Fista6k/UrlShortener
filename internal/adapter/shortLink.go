package adapter

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/redis/go-redis/v9"
)

func (s storage) Save(link *domain.Link) error {
	query := `
		INSERT INTO links (original_url, short_url, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := s.db.QueryRow(query, link.OriginalUrl, link.ShortUrl, link.CreatedAt).Scan(&link.ID)
	if err != nil {
		return err
	}

	err = s.redis.Set(s.ctx, link.ShortUrl, link.OriginalUrl, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s storage) FindByShortCode(code string) (string, error) {
	logger := s.ctx.Value("logger").(*slog.Logger)

	originalUrl, err := s.redis.Get(s.ctx, code).Result()
	if err == nil {
		return originalUrl, nil
	}

	if err != redis.Nil {
		logger.LogAttrs(
			s.ctx,
			slog.LevelError,
			"Can't take value from redis db",
			slog.Any("error", err),
			slog.String("shortCode", code),
		)
	}

	query := `
		SELECT id, original_url, short_url, created_at
		FROM links
		WHERE short_url = $1
	`

	var link *domain.Link
	err = s.db.QueryRow(query, code).Scan(&link.ID, &link.OriginalUrl, &link.ShortUrl, &link.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", domain.ErrNotFound
		} else {
			return "", err
		}
	}

	return link.OriginalUrl, nil
}
