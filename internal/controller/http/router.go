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

	router.Router.AddEndPoints(service)
}

func (r Router) AddEndPoints(service *service.ShortererService) {
	r.Router.POST("/", service.Shorten)
}
