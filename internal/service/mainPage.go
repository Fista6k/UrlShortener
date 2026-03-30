package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *ShortenerService) MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
