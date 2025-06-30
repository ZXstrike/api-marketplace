package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/ZXstrike/shared/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
			targetUrl = ResolveTargetUrlFromDBandCacheIt(c, db, redisClient)
			if targetUrl == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "No target URL found for this request"})
				return
			}
		}

		proxy := httputil.NewSingleHostReverseProxy(targetUrl)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			// Overwrite the path to use the one from our resolved URL
			req.URL.Path = targetUrl.Path

			// Set the scheme and host for the request
			req.URL.Scheme = targetUrl.Scheme
			req.URL.Host = targetUrl.Host

			// Set the Host header and clear the RequestURI
			req.Host = targetUrl.Host
			req.RequestURI = ""
		}

		proxy.ServeHTTP(c.Writer, c.Request)

	}
}
func ResolveTargetUrlFromCache(c *gin.Context, redisClient *redis.Client) *url.URL {
	tokenHeader := c.GetHeader("Token")
	if tokenHeader == "" {
		return nil
	}

	val, err := redisClient.Get(c.Request.Context(), tokenHeader).Result()
	if err != nil {
		return nil
	}

	var cacheInfo CacheRouteInfo
	err = json.Unmarshal([]byte(val), &cacheInfo)
	if err != nil {
		return nil
	}

	targetUrl, err := url.Parse(cacheInfo.TargetUrl)
	if err != nil {
		return nil
	}

	return targetUrl
}

func ResolveTargetUrlFromDBandCacheIt(c *gin.Context, db *gorm.DB, redisClient *redis.Client) *url.URL {
	// Get the host from the request header, not from URL
	host := c.Request.Host
	if host == "" {
		return nil
	}

	// Extract subdomain from host (e.g., "john.api.zxsttm" -> "john")
	hostParts := strings.Split(host, ".")
	if len(hostParts) < 3 {
		return nil
	}

	subdomain := hostParts[0] // This is your provider_username

	// Get the path to extract api_slug
	path := c.Request.URL.Path
	pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(pathParts) < 1 || pathParts[0] == "" {
		return nil
	}

	apiSlug := pathParts[0]

	targetPath := strings.Join(pathParts[1:], "/")

	var user models.User
	if err := db.Where("username = ?", subdomain).First(&user).Error; err != nil {
		// User not found or error occurred
		return nil
	}

	// Fetch the API by slug and user ID
	var api models.API
	if err := db.Where("slug = ? AND provider_id = ?", apiSlug, user.ID).First(&api).Error; err != nil {
		// API not found or error occurred
		return nil
	}

	// Chek if the api mathc with the one resolved from the auth middleware
	if api.ID != c.GetString("api_id") {
		// API ID does not match, return nil
		return nil
	}

	targetUrlString := api.BaseURL + targetPath

	targetUrl, err := url.Parse(targetUrlString)
	if err != nil {
		// Invalid URL format
		return nil
	}

	// For now, returning a placeholder
	return targetUrl
}
