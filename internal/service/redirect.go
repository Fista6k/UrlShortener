package service

import (
	"log/slog"
	"net/http"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
)

func (s ShortererService) Redirect(c *gin.Context) {
	logger := s.ctx.Value(domain.LoggerKey).(*slog.Logger)
	shortUrl := c.Param("shortUrl")

	if shortUrl == "favicon.ico" {
		c.Status(http.StatusNotFound)
		return
	}

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"trying to find link by short code",
		slog.String("shortCode", shortUrl),
	)

	link, err := s.storage.FindByShortCode(shortUrl)

	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

			logger.LogAttrs(
				c,
				slog.LevelError,
				"link nt found by short code",
				slog.Any("error", err),
			)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       domain.ErrInternalProblems,
				"description": "can't find link by short code for a while",
			})

			logger.LogAttrs(
				c,
				slog.LevelError,
				"something went wrong while executing query (finding link by short code)",
				slog.Any("error", err),
			)
		}
		return
	}

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"Successfully find link by shortCode",
		slog.String("url", link),
		slog.String("shortCode", shortUrl),
	)

	logger.LogAttrs(
		c,
		slog.LevelInfo,
		"Redirecting to the link",
		slog.String("url", link),
	)

	c.Redirect(http.StatusPermanentRedirect, link)
}
