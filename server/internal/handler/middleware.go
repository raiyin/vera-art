package handler

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/pkg/apperror"
	"github.com/raiyin/artserver/pkg/jwt"
)

// RecoveryMiddleware returns a Gin middleware that recovers from panics and logs them.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				slog.Error("PANIC recovered",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"client_ip", c.ClientIP(),
					"panic", r,
					"stack", stack,
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, apperror.APIError{
					Status:  http.StatusInternalServerError,
					Code:    "INTERNAL_ERROR",
					Message: "internal server error",
				})
			}
		}()
		c.Next()
	}
}

// RequestLoggerMiddleware returns a Gin middleware that logs all HTTP requests.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
			latency := time.Since(start)
			status := c.Writer.Status()
			clientIP := c.ClientIP()
			method := c.Request.Method

			// Get user ID if available
			userID, _ := c.Get("user_id")
			userIDVal := int64(0)
			if uid, ok := userID.(int64); ok {
				userIDVal = uid
			}

			// Build log attributes
			logArgs := []any{
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("client_ip", clientIP),
				slog.Int64("user_id", userIDVal),
			}
			if query != "" {
				logArgs = append(logArgs, slog.String("query", query))
			}

			if status >= 500 {
				slog.Warn("HTTP request completed with server error", logArgs...)
			} else if status >= 400 {
				slog.Info("HTTP request completed with client error", logArgs...)
			} else if latency > time.Second {
				slog.Info("HTTP request completed (slow)", logArgs...)
			} else {
				slog.Debug("HTTP request completed", logArgs...)
			}
	}
}

// AuthMiddleware returns a Gin middleware that requires a valid JWT access token.
func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			slog.Warn("Auth failed: missing authorization header",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
			)
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
			slog.Warn("Auth failed: invalid authorization header format",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
			)
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
			slog.Warn("Auth failed: invalid or expired token",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
				"error", err,
			)
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
		c.Set("username", claims.Username)
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
		c.Set("username", claims.Username)
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
			userID, _ := c.Get("user_id")
			slog.Warn("Admin access denied",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
				"user_id", userID,
				"role", role,
			)
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
			slog.Warn("Rate limit exceeded",
				"path", c.Request.URL.Path,
				"client_ip", ip,
			)
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
