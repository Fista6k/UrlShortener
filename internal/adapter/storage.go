package adapter

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

type storage struct {
	db    *sql.DB
	Redis *redis.Client
	ctx   context.Context
}

type Storage interface {
	Save(*domain.Link) error
	FindByShortCode(string) (string, error)
}

var embedMigrations embed.FS

func ConnToStorage(ctx context.Context) (*storage, error) {
	logger := ctx.Value("logger").(*slog.Logger)

	connStr := makeConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"can't open database with connect string",
			slog.Any("error", err),
		)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.LogAttrs(ctx,
			slog.LevelError,
			"error while ping data base, data base don't answer",
			slog.Any("error", err),
		)
		return nil, err
	}

	goose.SetBaseFS(embedMigrations)

	if err = goose.Up(db, "migrations"); err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"can't apply migrations",
			slog.Any("error", err),
		)
	}

	redis, err := ConnToRedis()
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"Can't connect to redis db",
			slog.Any("error", err),
		)
		return nil, err
	}

	return &storage{
		db:    db,
		ctx:   ctx,
		Redis: redis,
	}, nil
}

func ConnToRedis() (*redis.Client, error) {
	connString := makeConnStringForRedis()

	opt, err := redis.ParseURL(connString)

	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}

func makeConnStr() string {
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
}

func makeConnStringForRedis() string {
	db := os.Getenv("DB_NAME_REDIS")
	password := os.Getenv("REDIS_PASSWORD")
	user := os.Getenv("REDIS_USER")
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", user, password, host, port, db)
}
