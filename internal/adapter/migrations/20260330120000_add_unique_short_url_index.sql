-- +goose Up
-- Add uniqueness for short_url so ON CONFLICT (short_url) works correctly.
CREATE UNIQUE INDEX IF NOT EXISTS idx_links_short_url_unique
ON links(short_url);

-- +goose Down
DROP INDEX IF EXISTS idx_links_short_url_unique;
