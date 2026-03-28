-- +goose Up
CREATE TABLE IF NOT EXISTS links(
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_links_short_url ON links(short_url);

-- +goose Down
DROP TABLE IF EXISTS links;
