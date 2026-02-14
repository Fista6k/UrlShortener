package service

import (
	"net/http"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
)

func (s ShortererService) Redirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	link, err := s.storage.FindByShortCode(shortUrl)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.Redirect(http.StatusPermanentRedirect, link.OriginalUrl)
}
