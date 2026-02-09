package service

import (
	"net/http"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (s ShortererService) Shorten(c *gin.Context) {
	link := &domain.Link{}

	if err := c.ShouldBindJSON(link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}

	oldLink, err := s.storage.FindByURL(string(link.OriginalUrl))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url not found",
		})
		return
	}
	if oldLink != nil {
		c.JSON(http.StatusCreated, oldLink.ToJson())
		return
	}

	alphabet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

	short_url, err := gonanoid.Generate(alphabet, 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "cant generate short url",
		})
		return
	}

	url, err := s.storage.FindByShortCode(short_url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "url not found",
		})
	}

	for url != nil {
		short_url, err = gonanoid.Generate(alphabet, 8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "cant generate short url",
			})
			return
		}

		url, err = s.storage.FindByShortCode(short_url)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "url not found",
			})
			return
		}
	}

	l, err := domain.NewLink(string(link.OriginalUrl), short_url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "cant create link object",
		})
		return
	}

	err = s.storage.Save(l)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "cant save to database",
		})
		return
	}

	c.JSON(http.StatusCreated, url.ToJson())
}
