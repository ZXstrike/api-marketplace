package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type CacheRouteInfo struct {
	TargetUrl string `json:"target_url"`
	ApiID     string `json:"api_id"`
}

func ProxyHandler(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	// This is a placeholder for the actual proxy logic.
	// You would typically forward the request to another service here.

	return func(c *gin.Context) {
		// Example: Forward the request to another service
		// You can use a library like "net/http" or "github.com/gin-gonic/gin" to forward the request.

		targetUrl := ResolveTargetUrlFromCache(c, redisClient)
		if targetUrl == nil {
			c.JSON(500, gin.H{
				"error": "Failed to resolve target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			req.URL.Scheme = targetUrl.Scheme
			req.URL.Host = targetUrl.Host

			req.Host = targetUrl.Host
			req.RequestURI = req.URL.RequestURI()

		}

		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

func ResolveTargetUrlFromCache(c *gin.Context, redisClient *redis.Client) *url.URL {
	// This function should resolve the target URL based on the incoming request.
	// For now, we'll just return a placeholder URL.

	return &url.URL{
		Scheme: "http", Host: "example.com",
	}
}

func ResolveTargetUrlFromDBandCache(c *gin.Context, db *gorm.DB, redisClient *redis.Client) *url.URL {
	// This function should resolve the target URL based on the incoming request.
	// For now, we'll just return a placeholder URL.

	return &url.URL{
		Scheme: "http", Host: "example.com",
	}
}
