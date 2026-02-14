package service

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itchyny/base58-go"
)

func (s ShortererService) Shorten(c *gin.Context) {

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

func GenerateShortLink(originalUrl, userId string) string {
	urlHashBytes := hashing(originalUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	result := encoding([]byte(fmt.Sprintf("%d", generatedNumber)))
	return result[:8]
}
