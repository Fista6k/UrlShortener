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

func mainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Main page",
	})
}

func mainFormHandler(c *gin.Context) {
	url := c.PostForm("url")
	c.IndentedJSON(http.StatusOK, gin.H{
		"url": url,
	})
}
