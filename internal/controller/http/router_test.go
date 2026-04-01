package controller

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
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

func TestRateLimiter_AllowsAndThenLimitsRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.WithValue(context.Background(), domain.LoggerKey, logger)
	router := &Router{
		gin.Default(),
		&RateLimiter{
			clients: make(map[string]*Client),
		},
		ctx,
	}
	router.Router.Use(router.RateLimiterFunc())
	router.Router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	ok, limited := 0, 0

	for i := 0; i < 100; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.RemoteAddr = "1.2.3.4:12345"

		router.Router.ServeHTTP(rec, req)

		if rec.Code == http.StatusTooManyRequests {
			limited += 1
		} else {
			ok += 1
		}
	}

	if ok == 0 {
		t.Fatal("expected some non-429 responses, got 0")
	}

	if limited == 0 {
		t.Fatal("expected some 429 responses, got 0")
	}
}

func TestRateLimiter_CallStop_NothingRemains(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.WithValue(context.Background(), domain.LoggerKey, logger)
	router := &Router{
		gin.Default(),
		&RateLimiter{
			clients: make(map[string]*Client),
		},
		ctx,
	}

	router.Router.Use(router.RateLimiterFunc())
	router.Router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 10; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.RemoteAddr = "1.2.3.4:12345"

		router.Router.ServeHTTP(rec, req)
	}

	if len(router.RateLimiter.clients) == 0 {
		t.Fatal("expected clients map to be non-empty before Stop")
	}

	router.RateLimiter.Stop()

	if len(router.RateLimiter.clients) != 0 {
		t.Fatal("expected empty clients map after Stop")
	}
}

func TestRouter_SuccessfullyCreated_ReturnsAnyResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	fakeStorage := &fakeStorage{
		st: map[string]string{
			"short": "long",
		},
	}
	ctx := context.WithValue(context.Background(), domain.LoggerKey, logger)
	testService := service.NewShortenerService(ctx, fakeStorage)

	router := gin.New()
	router.LoadHTMLGlob("../../../static/*.html")
	router.GET("/", testService.MainPage)
	router.POST("/shorten", testService.Shorten)
	router.GET("/:shortUrl", testService.Redirect)

	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(rec1, req1)
	if rec1.Code == http.StatusNotFound {
		t.Fatal("GET / route is not registered")
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/shorten", nil)
	router.ServeHTTP(rec2, req2)
	if rec2.Code == http.StatusNotFound {
		t.Fatal("POST /shorten route is not registered")
	}

	rec3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/short", nil)
	router.ServeHTTP(rec3, req3)
	if rec3.Code == http.StatusNotFound {
		t.Fatal("GET /:shortUrl route is not registered")
	}
}
