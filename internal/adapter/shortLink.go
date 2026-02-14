package adapter

import (
	"database/sql"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
)

func (s storage) Save(link *domain.Link) {
	query := `
		INSERT INTO links original_url, short_url, created_at
		VALUES ($1, $2, $3)
		returning id
	`

	s.db.QueryRow(query, link.OriginalUrl, link.ShortUrl, link.CreatedAt).Scan(&link.ID)
}

func (s storage) FindByShortCode(code string) (*domain.Link, error) {
	query := `
		SELECT id, original_url, short_url, created_at
		FROM links
		WHERE short_url = $1
	`

	var link *domain.Link
	err := s.db.QueryRow(query, code).Scan(&link.ID, &link.OriginalUrl, &link.ShortUrl, &link.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}

	return link, nil
}

func (s storage) FindByURL(url string) (*domain.Link, error) {
	query := `
		SELECT id, original_url, short_url, created_at
		FROM links
		WHERE original_url = $1
	`

	var link *domain.Link
	err := s.db.QueryRow(query, url).Scan(&link.ID, &link.OriginalUrl, &link.ShortUrl, &link.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}

	return link, nil
}
