package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var urls = map[string]string{}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", mainPageHandler)
	r.POST("/submit", mainFormHandler)

	r.Run()
}