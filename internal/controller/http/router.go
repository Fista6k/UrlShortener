package http

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

	router.AddEndPoints(service)

	return router
}

func (r Router) AddEndPoints(service *service.ShortererService) {
	r.Router.POST("/create-shortUrl", service.Shorten)
	r.Router.GET("/:shortUrl", service.Redirect)
}
