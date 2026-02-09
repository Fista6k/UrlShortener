package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s ShortererService) Redirect(c *gin.Context) {
	var shortCode string
	if err := c.ShouldBindUri(&shortCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	link, err := s.storage.FindByShortCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url not found",
		})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, string(link.OriginalUrl))
	return
}
