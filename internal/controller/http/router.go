package controller

import (
	"log"
	"net/http"
	"sync"

	"github.com/Fista6k/Url-Shorterer.git/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Router struct {
	Router      *gin.Engine
	RateLimiter *RateLimiter
}

type Client struct {
	limiter *rate.Limiter
}

type RateLimiter struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewRouter(service *service.ShortererService) *Router {
	router := &Router{
		gin.Default(),
		&RateLimiter{
			clients: make(map[string]*Client),
		},
	}

	router.Router.LoadHTMLGlob("static/*.html")

	router.AddEndPoints(service)

	return router
}

func (r Router) AddEndPoints(service *service.ShortererService) {
	r.Router.Use(r.RateLimiterFunc())
	r.Router.GET("/", service.MainPage)
	r.Router.POST("/shorten", service.Shorten)
	r.Router.GET("/:shortUrl", service.Redirect)
}

func (r Router) RateLimiterFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		r.RateLimiter.mu.Lock()
		if _, ok := r.RateLimiter.clients[ip]; !ok {
			r.RateLimiter.clients[ip] = &Client{limiter: rate.NewLimiter(10, 20)}
		}
		cl := r.RateLimiter.clients[ip]
		r.RateLimiter.mu.Unlock()

		if !cl.limiter.Allow() {
			log.Println("rate limited", ip)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		log.Println("allowed request from", ip)

		c.Next()
	}
}

func (rl *RateLimiter) Stop() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.clients = make(map[string]*Client)
}
