package controller

import (
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router *gin.Engine
}

func NewRouter(service *service.ShortererService) *Router {
	router := &Router{
		gin.Default(),
	}

	router.Router.LoadHTMLGlob("static/*.html")

	router.AddEndPoints(service)

	return router
}

func (r Router) AddEndPoints(service *service.ShortererService) {
	r.Router.GET("/", service.MainPage)
	r.Router.POST("/shorten", service.Shorten)
	r.Router.GET("/:shortUrl", service.Redirect)
}
