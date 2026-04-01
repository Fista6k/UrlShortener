//go:build integration

package adapter

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
)

func newTestCtx() context.Context {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return context.WithValue(context.Background(), domain.LoggerKey, logger)
}

func clearData(t *testing.T, st *storage) {
	t.Helper()

	_, err := st.db.Exec(`DELETE FROM links`)
	if err != nil {
		t.Fatalf("failed to clean links table: %v", err)
	}

	if err := st.Redis.FlushDB(st.ctx).Err(); err != nil {
		t.Fatalf("failed to flush redis: %v", err)
	}
}

func TestStorage_SaveOrGet_And_FindByShortCode(t *testing.T) {
	ctx := newTestCtx()

	st, err := ConnToStorage(ctx)
	if err != nil {
		t.Fatalf("failed to connect to storage: %v", err)
	}
	defer st.Redis.Close()
	defer st.db.Close()

	clearData(t, st)

	url := "https://example.com/page"
	short := "abc12345"

	got, err := st.SaveOrGet(&domain.Link{
		OriginalUrl: url,
		ShortUrl:    short,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		t.Fatalf("SaveOrGet failed: %v", err)
	}

	if got != url {
		t.Fatalf("expected %q, got %q", url, got)
	}

	found, err := st.FindByShortCode(short)
	if err != nil {
		t.Fatalf("FindByShortCode failed: %v", err)
	}

	if found != url {
		t.Fatalf("expected %q, got %q", url, found)
	}
}

func TestStorage_SaveOrGet_Collision_ReturnsExistingOriginalUrl(t *testing.T) {
	ctx := newTestCtx()

	st, err := ConnToStorage(ctx)
	if err != nil {
		t.Fatalf("failed to connect storage: %v", err)
	}
	defer st.Redis.Close()
	defer st.db.Close()

	clearData(t, st)

	short := "short"
	url1 := "https://example.com/first"
	url2 := "https://example.com/second"

	_, err = st.SaveOrGet(&domain.Link{
		OriginalUrl: url1,
		ShortUrl:    short,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("SaveOrGet failed: %v", err)
	}

	got, err := st.SaveOrGet(&domain.Link{
		OriginalUrl: url2,
		ShortUrl:    short,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("SaveOrGet failed: %v", err)
	}

	if got != url1 {
		t.Fatalf("expected %q, got %q", url1, got)
	}
}
