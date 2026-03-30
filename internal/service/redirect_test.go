package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
)

func TestRedirect_NotFound_Returns404AndJson(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeStorage := &fakeStorage{
		st: map[string]string{},
	}

	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/abc", nil)
	c.Params = gin.Params{gin.Param{Key: "shortUrl", Value: "/abc123"}}

	testService.Redirect(c)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v, body=%q", err, rec.Body.String())
	}

	if body["error"] != domain.ErrNotFound.Error() {
		t.Fatalf("expected error=%q, got %v", domain.ErrNotFound.Error(), body["error"])
	}
}

func TestRedirect_Found_Returns308AndLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	url := "https://example.com/page"
	short := "abs123"

	fakeStorage := &fakeStorage{
		st: map[string]string{
			short: url,
		},
	}

	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/"+short, nil)
	c.Params = gin.Params{gin.Param{Key: "shortUrl", Value: short}}

	testService.Redirect(c)

	if rec.Code != http.StatusPermanentRedirect {
		t.Fatalf("expected 308, got %d", rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != url {
		t.Fatalf("expected location=%q, got %q", url, loc)
	}
}

func TestRedirect_Favicon_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	fakeStorage := &fakeStorage{st: map[string]string{}}
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	router := gin.New()
	router.GET("/:shortUrl", testService.Redirect)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
