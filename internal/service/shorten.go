package service

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/itchyny/base58-go"
)

func (s ShortenerService) Shorten(c *gin.Context) {
	logger := s.ctx.Value(domain.LoggerKey).(*slog.Logger)
	original_url := c.PostForm("url")

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"generate short link by original link",
		slog.String("url", original_url),
	)

	shortUrl, err := s.CreateShortLink(original_url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       domain.ErrInternalProblems,
			"description": "can't generate short link for a while",
		})

		logger.LogAttrs(
			c,
			slog.LevelError,
			"Error with generating short link",
			slog.Any("error", err),
		)
		return
	}

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"short link was successfully generated",
		slog.String("url", original_url),
		slog.String("shortUrl", shortUrl),
	)

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"trying to save link in db",
		slog.String("url", original_url),
		slog.String("shortUrl", shortUrl),
	)

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"link was successfully saved in db",
		slog.String("url", original_url),
		slog.String("shortUrl", shortUrl),
	)

	c.HTML(http.StatusCreated, "index.html", gin.H{
		"ShortUrl": shortUrl,
	})
}

func hashing(input string) []byte {
	s := sha256.New()
	s.Write([]byte(input))
	return s.Sum(nil)
}

func encoding(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func generateShortLink(originalUrl string) (string, error) {
	urlHashBytes := hashing(originalUrl)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	result, err := encoding([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return "", err
	}
	return result[:8], nil
}

func (s *ShortenerService) CreateShortLink(original_url string) (string, error) {

	shortUrl, err := generateShortLink(original_url)

	if err != nil {
		return "", err
	}

	existingUrl, err := s.storage.SaveOrGet(&domain.Link{
		OriginalUrl: original_url,
		ShortUrl:    shortUrl,
		CreatedAt:   time.Now(),
	})

	if existingUrl == original_url {
		return shortUrl, nil
	}

	return s.resolveCollision(original_url, 1)
}

func (s *ShortenerService) resolveCollision(original_url string, attempt int) (string, error) {
	if attempt > 5 {
		return "", domain.ErrMaxAttemptsToGenerateShortUrl
	}

	salted := fmt.Sprintf("%s|%d", original_url, attempt)
	shortLink, err := generateShortLink(salted)
	if err != nil {
		return "", err
	}

	existingUrl, err := s.storage.SaveOrGet(&domain.Link{
		OriginalUrl: original_url,
		ShortUrl:    shortLink,
		CreatedAt:   time.Now(),
	})

	if existingUrl == original_url {
		return shortLink, nil
	}

	return s.resolveCollision(original_url, attempt+1)
}
