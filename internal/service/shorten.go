package service

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/itchyny/base58-go"
)

func (s ShortererService) Shorten(c *gin.Context) {
	original_url := c.PostForm("url")

	shortUrl := GenerateShortLink(original_url)
	err := s.storage.Save(&domain.Link{
		OriginalUrl: original_url,
		ShortUrl:    shortUrl,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusCreated, "index.html", gin.H{
		"ShortUrl": shortUrl,
	})
}

func hashing(input string) []byte {
	s := sha256.New()
	s.Write([]byte(input))
	return s.Sum(nil)
}

func encoding(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

func GenerateShortLink(originalUrl string) string {
	urlHashBytes := hashing(originalUrl)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	result := encoding([]byte(fmt.Sprintf("%d", generatedNumber)))
	return result[:8]
}
