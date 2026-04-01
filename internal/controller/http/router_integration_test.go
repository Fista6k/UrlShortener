//go:build integration

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
)

func testCtx() context.Context {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return context.WithValue(context.Background(), domain.LoggerKey, logger)
}

func TestE2E_ShortenThenRedirect(t *testing.T) {
	ctx := testCtx()

	st, err := adapter.ConnToStorage(ctx)
	if err != nil {
		t.Fatalf("failed to connect storage: %v", err)
	}
	defer st.Redis.Close()
	defer st.Db.Close()

	service := service.NewShortenerService(ctx, st)
	router := NewRouter(ctx, service)

	link := "https://example.com/page"
	form := url.Values{}
	form.Set("url", link)

	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(form.Encode()))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.Router.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec1.Code)
	}

	shortUrl, err := extractShortUrlFromHTML(rec1.Body.String())
	if err != nil {
		t.Fatalf("failed to parse short url from html: %v", err)
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/"+shortUrl, nil)
	router.Router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusPermanentRedirect {
		t.Fatalf("expected 308, got %d", rec2.Code)
	}

	if loc := rec2.Header().Get("Location"); loc != "https://example.com/page" {
		t.Fatalf("expected redirect to original url, got %q", loc)
	}
}

var ShortUrlRe = regexp.MustCompile(`href="/([A-Za-z0-9]+)"`)

func extractShortUrlFromHTML(body string) (string, error) {
	m := ShortUrlRe.FindStringSubmatch(body)
	if len(m) < 2 {
		return "", fmt.Errorf("short url not found in html body: %q", body)
	}
	return m[1], nil
}

func TestE2E_RedirectUnknown_Returns404(t *testing.T) {
	ctx := testCtx()

	st, err := adapter.ConnToStorage(ctx)
	if err != nil {
		t.Fatalf("failed to connect storage: %v", err)
	}
	defer st.Redis.Close()
	defer st.Db.Close()

	_, err = st.Db.Exec(`DELETE FROM links`)
	if err != nil {
		t.Fatalf("failed to clear table links: %v", err)
	}

	service := service.NewShortenerService(ctx, st)
	router := NewRouter(ctx, service)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/abs123", nil)
	router.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	var body map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if body["error"] != domain.ErrNotFound.Error() {
		t.Fatalf("expected %q, got %q", domain.ErrNotFound.Error(), body["error"])
	}
}
