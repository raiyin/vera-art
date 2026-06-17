package handler

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/pkg/apperror"
	"github.com/raiyin/artserver/pkg/jwt"
)

// AuthMiddleware returns a Gin middleware that requires a valid JWT access token.
func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, apperror.APIError{
				Status:  http.StatusUnauthorized,
				Code:    "UNAUTHORIZED",
				Message: "authorization header required",
			})
			c.Abort()
			return
		}

		// Check for Bearer prefix
		const bearerPrefix = "Bearer "
		if len(tokenString) <= len(bearerPrefix) || !strings.HasPrefix(tokenString, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, apperror.APIError{
				Status:  http.StatusUnauthorized,
				Code:    "INVALID_TOKEN",
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}
		tokenString = tokenString[len(bearerPrefix):]

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apperror.APIError{
				Status:  http.StatusUnauthorized,
				Code:    "INVALID_TOKEN",
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// AuthMiddlewareOptional returns a Gin middleware that optionally authenticates.
// If a valid token is provided, user info is set in context.
// If no token or invalid token, the request still proceeds.
func AuthMiddlewareOptional(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Next()
			return
		}

		const bearerPrefix = "Bearer "
		if len(tokenString) <= len(bearerPrefix) || !strings.HasPrefix(tokenString, bearerPrefix) {
			c.Next()
			return
		}
		tokenString = tokenString[len(bearerPrefix):]

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// AdminMiddleware returns a Gin middleware that requires admin role.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			c.JSON(http.StatusForbidden, apperror.APIError{
				Status:  http.StatusForbidden,
				Code:    "FORBIDDEN",
				Message: "admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// rateLimiter implements a simple IP-based rate limiter.
type rateLimiter struct {
	mu          sync.RWMutex
	visits      map[string][]time.Time
	limit       int
	window      time.Duration
	cleanupTime time.Duration
}

// newRateLimiter creates a new rate limiter.
func newRateLimiter(limit int, window, cleanupTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visits:      make(map[string][]time.Time),
		limit:       limit,
		window:      window,
		cleanupTime: cleanupTime,
	}

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(cleanupTime)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// cleanup removes old entries.
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.cleanupTime)
	for ip, times := range rl.visits {
		var filtered []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			delete(rl.visits, ip)
		} else {
			rl.visits[ip] = filtered
		}
	}
}

// allow checks if an IP is allowed to make a request.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	visits := rl.visits[ip]

	var filtered []time.Time
	for _, t := range visits {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= rl.limit {
		return false
	}

	filtered = append(filtered, now)
	rl.visits[ip] = filtered
	return true
}

// RateLimitMiddleware returns a Gin middleware that rate limits requests by IP.
func RateLimitMiddleware() gin.HandlerFunc {
	rl := newRateLimiter(5, time.Minute, time.Hour)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, apperror.APIError{
				Status:  http.StatusTooManyRequests,
				Code:    "RATE_LIMITED",
				Message: "too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORSMiddleware returns a Gin middleware that handles CORS.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Range")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
