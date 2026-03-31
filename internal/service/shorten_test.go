package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
)

type fakeStorage struct {
	st map[string]string
}

func (f *fakeStorage) SaveOrGet(link *domain.Link) (string, error) {
	if f.st == nil {
		return "", nil
	}

	return f.st[link.ShortUrl], nil
}

func (f *fakeStorage) FindByShortCode(code string) (string, error) {
	if f.st == nil {
		return "", nil
	}

	if v, ok := f.st[code]; ok {
		return v, nil
	}

	return "", domain.ErrNotFound
}

type testWriter struct {
}

func (testWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestGenerateShortLink_DeterministicAndLength(t *testing.T) {
	url := "https://example.com/page"

	s1, err := generateShortLink(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s2, err := generateShortLink(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s1 != s2 {
		t.Fatalf("expected deterministic short url, got %q and %q", s1, s2)
	}

	if len(s1) != 8 {
		t.Fatalf("expected length=8, got %d (%q)", len(s1), s1)
	}
}

func TestGenerateShortLink_NoCollision_ReturnsBaseShortLink(t *testing.T) {
	url := "https://example.com/page"

	baseShort, err := generateShortLink(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fakeStorage := &fakeStorage{
		st: map[string]string{
			baseShort: url,
		},
	}

	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	got, err := testService.CreateShortLink(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != baseShort {
		t.Fatalf("expected %q, got %q", baseShort, got)
	}
}

func TestGenerateShortLink_CollisionResolveWithSaltedAttempt(t *testing.T) {
	url := "https://example.com/page"

	baseShort, _ := generateShortLink(url)
	saltedShort1, _ := generateSalted(url, 1)

	otherUrl := "https://example.com/other"

	fakeStorage := &fakeStorage{
		st: map[string]string{
			baseShort:    otherUrl,
			saltedShort1: url,
		},
	}

	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	got, err := testService.CreateShortLink(url)
	if err != nil {
		t.Fatalf("unexpected error^ %v", err)
	}

	if got != saltedShort1 {
		t.Fatalf("expected %q, got %q", saltedShort1, got)
	}
}

func generateSalted(url string, attempt int) (string, error) {
	salted := fmt.Sprintf("%s|%d", url, attempt)
	short, _ := generateShortLink(salted)
	return short, nil
}

func TestGenerateShortLink_ExceedsAttempts_ReturnsMaxAttemptsError(t *testing.T) {
	url := "https://example.com/page"

	fakeStorage := &fakeStorage{
		st: map[string]string{},
	}

	otherUrl := "https://example.com/other"

	baseShort, _ := generateShortLink(url)
	fakeStorage.st[baseShort] = otherUrl

	for attempt := 1; attempt <= 5; attempt += 1 {
		saltedShort, _ := generateSalted(url, attempt)
		fakeStorage.st[saltedShort] = otherUrl
	}

	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	_, err := testService.CreateShortLink(url)
	if !errors.Is(err, domain.ErrMaxAttemptsToGenerateShortUrl) {
		t.Fatalf("expected max attempts error, got %v", err)
	}
}

func TestCreateShortLink_StatusCreated_Returns201AndJson(t *testing.T) {
	gin.SetMode(gin.TestMode)

	link := "https://example.com/page"
	shortUrl, _ := generateShortLink(link)

	fakeStorage := &fakeStorage{st: map[string]string{
		shortUrl: link,
	}}
	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	testService := NewShortenerService(context.WithValue(context.Background(), domain.LoggerKey, logger), fakeStorage)

	router := gin.New()
	router.LoadHTMLGlob("../../static/*.html")
	router.POST("/shorten", testService.Shorten)

	values := url.Values{}
	values.Set("url", link)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, shortUrl) {
		t.Fatalf("expected body to contain short url %q, got %q", shortUrl, body)
	}
}
