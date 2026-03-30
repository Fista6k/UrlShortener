package controller

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Router struct {
	Router      *gin.Engine
	RateLimiter *RateLimiter
	ctx         context.Context
}

type Client struct {
	limiter *rate.Limiter
}

type RateLimiter struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewRouter(ctx context.Context, service *service.ShortenerService) *Router {
	router := &Router{
		gin.Default(),
		&RateLimiter{
			clients: make(map[string]*Client),
		},
		ctx,
	}

	router.Router.LoadHTMLGlob("static/*.html")

	router.AddEndPoints(service)

	return router
}

func (r Router) AddEndPoints(service *service.ShortenerService) {
	r.Router.Use(r.RateLimiterFunc())
	r.Router.GET("/", service.MainPage)
	r.Router.POST("/shorten", service.Shorten)
	r.Router.GET("/:shortUrl", service.Redirect)
}

func (r Router) RateLimiterFunc() gin.HandlerFunc {
	logger := r.ctx.Value(domain.LoggerKey).(*slog.Logger)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		r.RateLimiter.mu.Lock()
		if _, ok := r.RateLimiter.clients[ip]; !ok {
			r.RateLimiter.clients[ip] = &Client{limiter: rate.NewLimiter(10, 20)}
		}
		cl := r.RateLimiter.clients[ip]
		r.RateLimiter.mu.Unlock()

		if !cl.limiter.Allow() {
			logger.LogAttrs(
				r.ctx,
				slog.LevelDebug,
				"rate for user limited",
				slog.String("client ip", ip),
			)

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": domain.ErrLimitExceeded,
			})
			return
		}

		logger.LogAttrs(
			r.ctx,
			slog.LevelDebug,
			"allowed request from client",
			slog.String("client ip", ip),
		)

		c.Next()
	}
}

func (rl *RateLimiter) Stop() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.clients = make(map[string]*Client)
}
