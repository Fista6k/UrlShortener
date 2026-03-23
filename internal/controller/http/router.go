package controller

import (
	"net/http"
	"sync"

	"github.com/Fista6k/Url-Shorterer.git/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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
	r.Router.Use(r.RateLimiter())
	r.Router.GET("/", service.MainPage)
	r.Router.POST("/shorten", service.Shorten)
	r.Router.GET("/:shortUrl", service.Redirect)
}

func (r Router) RateLimiter() gin.HandlerFunc {
	type client struct {
		limiter *rate.Limiter
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, ok := clients[ip]; !ok {
			clients[ip] = &client{limiter: rate.NewLimiter(10, 20)}
		}
		cl := clients[ip]
		mu.Unlock()

		if !cl.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
		}

		c.Next()
	}
}
