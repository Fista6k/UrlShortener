package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
	controller "github.com/Fista6k/Url-Shorterer.git/internal/controller/http"
	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Info("Can't found .env file, using environment variables")
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("APP_PORT")

	logger := slog.Default()

	storage, err := adapter.ConnToStorage(context.WithValue(ctx, domain.LoggerKey, logger))
	if err != nil {
		logger.LogAttrs(
			context.TODO(),
			slog.LevelError,
			"something went wrong with storage init",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	service := service.NewShortenerService(context.WithValue(ctx, domain.LoggerKey, logger), storage)
	r := controller.NewRouter(context.WithValue(ctx, domain.LoggerKey, logger), service)
	addr := ":" + port

	server := &http.Server{
		Handler: r.Router,
		Addr:    addr,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogAttrs(
				context.TODO(),
				slog.LevelError,
				"filed while listening",
				slog.Any("error", err),
			)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down gracefully, press Ctrl+C to force")

	r.RateLimiter.Stop()
	_ = storage.Redis.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"Server forced to shutdown",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.Info("Server exiting")
}
