package adapter

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/redis/go-redis/v9"
)

func (s storage) SaveOrGet(link *domain.Link) (string, error) {
	logger := s.ctx.Value(domain.LoggerKey).(*slog.Logger)

	query := `
		INSERT INTO links (original_url, short_url, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
		ON CONFLICT (short_url) DO NOTHING;
	`

	err := s.db.QueryRow(query, link.OriginalUrl, link.ShortUrl, link.CreatedAt).Scan(&link.ID)
	if err != nil && err != sql.ErrNoRows {
		logger.LogAttrs(
			s.ctx,
			slog.LevelError,
			"Can't save link in postgres",
			slog.Any("error", err),
		)
		return "", err
	}

	existingUrl, err := s.FindByShortCode(link.ShortUrl)
	if err != nil {
		return "", err
	}

	err = s.Redis.Set(s.ctx, link.ShortUrl, existingUrl, 24*time.Hour).Err()
	if err != nil {
		logger.LogAttrs(
			s.ctx,
			slog.LevelError,
			"Can't save link in redis",
			slog.Any("error", err),
		)
		return "", err
	}

	return existingUrl, nil
}

func (s storage) FindByShortCode(code string) (string, error) {
	logger := s.ctx.Value(domain.LoggerKey).(*slog.Logger)

	originalUrl, err := s.Redis.Get(s.ctx, code).Result()
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

	var link domain.Link
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

func (s storage) FindByURL(url string) (*domain.Link, error) {
	query := `
		SELECT id, original_url, short_url, created_at
		FROM links
		WHERE original_url = $1
	`

	var link domain.Link
	err := s.db.QueryRow(query, url).Scan(&link.ID, &link.OriginalUrl, &link.ShortUrl, &link.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}

	return &link, nil
}
