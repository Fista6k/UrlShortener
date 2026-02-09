package main

import (
	"github.com/gin-gonic/gin"
)

var urls = map[string]string{}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Run()
}
