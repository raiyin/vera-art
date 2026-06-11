// info for db
// table_link: 1-illustrations, 2-sales, 3-news
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/middleware"
	"github.com/raiyin/artserver/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var jwtKey []byte

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not connect to database (ping failed): %v", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
	log.Println("Database connection established successfully")

	// Run migrations
	runMigrations()
}

func runMigrations() {
	// Add avatar column to users table if it doesn't exist
	_, err := db.Exec("ALTER TABLE users ADD COLUMN avatar TEXT DEFAULT ''")
	if err != nil {
		// Ignore error if column already exists
		log.Printf("Migration (add avatar column): %v (this is normal if column already exists)", err)
	}
	log.Println("Database migrations completed")
}

func getMaterials(c *gin.Context) {
	query := "SELECT * FROM materials"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	materials := []models.Material{}
	for rows.Next() {
		p := models.Material{}
		err := rows.Scan(&p.Id, &p.MaterialRu, &p.MaterialEn)
		if err != nil {
			// Log the error but continue processing other rows
			log.Printf("Error scanning material row: %v", err)
			continue
		}
		materials = append(materials, p)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, materials)
}

func getBases(c *gin.Context) {
	query := "SELECT * FROM bases"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	bases := []models.Base{}
	for rows.Next() {
		p := models.Base{}
		err := rows.Scan(&p.Id, &p.BaseRu, &p.BaseEn)
		if err != nil {
			// Log the error but continue processing other rows
			log.Printf("Error scanning base row: %v", err)
			continue
		}
		bases = append(bases, p)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, bases)
}

// Rate limiting structures
type rateLimiter struct {
	mu          sync.RWMutex
	visits      map[string][]time.Time
	limit       int
	window      time.Duration
	cleanupTime time.Duration
}

var authRateLimiter = &rateLimiter{
	visits:      make(map[string][]time.Time),
	limit:       5,           // 5 requests
	window:      time.Minute, // per minute
	cleanupTime: time.Hour,
}

// Cleanup old entries periodically
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.cleanupTime)
	for ip, times := range rl.visits {
		// Filter out old timestamps
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

// Check if IP is allowed
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing visits
	visits := rl.visits[ip]

	// Filter out visits outside the window
	var filtered []time.Time
	for _, t := range visits {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	// Check if under limit
	if len(filtered) >= rl.limit {
		return false
	}

	// Add current visit
	filtered = append(filtered, now)
	rl.visits[ip] = filtered
	return true
}

// RateLimitMiddleware for authentication endpoints
func RateLimitMiddleware() gin.HandlerFunc {
	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(authRateLimiter.cleanupTime)
		defer ticker.Stop()
		for range ticker.C {
			authRateLimiter.cleanup()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !authRateLimiter.allow(ip) {
			// Log rate limit hit
			log.Printf("[SECURITY] type=rate_limit ip=%s path=%s", ip, c.Request.URL.Path)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("dir is ", dir)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assuming environment variables are set externally")
	}

	// Set JWT key from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mydevsecret" // fallback for development
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	jwtKey = []byte(secret)

	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Access configuration values
	appConfig := config.AppConfigInstance
	log.Printf("Starting %s on port %d", appConfig.App.Name, appConfig.App.Port)
	log.Printf("dir is %s and %s", appConfig.Directories.RelWorksDir, appConfig.Directories.AbsWorksDir)

	r_gin := gin.Default()
	// r.Run(fmt.Sprintf(":%d", appConfig.App.Port))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true // Allow sending cookies/auth headers
	corsConfig.AllowOrigins = appConfig.CORS.AllowedOrigins
	if appConfig.CORS.Debug {
		log.Printf("CORS configured with allowed origins: %v", appConfig.CORS.AllowedOrigins)
	}
	r_gin.Use(cors.New(corsConfig))

	// Auth routes with rate limiting
	r_gin.POST("/register", RateLimitMiddleware(), Register)
	r_gin.POST("/login", RateLimitMiddleware(), Login)
	r_gin.POST("/refresh", RateLimitMiddleware(), Refresh)

	// User profile routes
	r_gin.GET("/profile", AuthMiddleware(), GetProfile)
	r_gin.PUT("/profile", AuthMiddleware(), UpdateProfile)

	// Avatar routes
	r_gin.POST("/profile/avatar", AuthMiddleware(), UploadAvatar)
	r_gin.DELETE("/profile/avatar", AuthMiddleware(), DeleteAvatar)
	r_gin.GET("/profile/avatar/:id", ServeAvatar)

	r_gin.GET("/sales", GetSales)
	r_gin.GET("/sales/:id", GetSaleById)
	r_gin.GET("/sales/:id/edit", GetSaleByIdForEdit)
	r_gin.POST("/sales", CreateSale)
	r_gin.PUT("/sales/:id", AuthMiddleware(), UpdateSale)
	r_gin.DELETE("/sales/:id", AuthMiddleware(), DeleteSale)

	r_gin.GET("/works", GetWorks)
	r_gin.GET("/works/:id", GetWorkById)
	r_gin.GET("/works/:id/edit", GetWorkByIdForEdit)
	r_gin.POST("/works", AddWork)
	r_gin.PUT("/works/:id", AuthMiddleware(), UpdateWork)
	r_gin.DELETE("/works/:id", AuthMiddleware(), DeleteWork)

	r_gin.GET("/news", GetNews)
	r_gin.GET("/news/:id", GetNewsById)
	r_gin.POST("/news", AuthMiddleware(), AddNews)
	r_gin.DELETE("/news/:id", AuthMiddleware(), DeleteNews)
	r_gin.PUT("/news/:id", AuthMiddleware(), UpdateNews)

	r_gin.GET("/materials", getMaterials)
	r_gin.GET("/bases", getBases)

	// Категории продуктов (публичные)
	r_gin.GET("/categories", GetCategories)
	r_gin.GET("/categories/:slug", GetCategoryBySlug)

	// Продукты по slug категории и slug продукта
	r_gin.GET("/category/:category_slug/product/:product_slug", GetProductBySlug)

	// Админские маршруты для категорий
	r_gin.GET("/admin/categories", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetCategories)
	r_gin.POST("/admin/categories", AuthMiddleware(), middleware.AdminMiddleware(), AdminCreateCategory)
	r_gin.PUT("/admin/categories/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdateCategory)
	r_gin.DELETE("/admin/categories/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteCategory)

	// Админские маршруты для тегов
	r_gin.GET("/admin/tags", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetTags)
	r_gin.POST("/admin/tags", AuthMiddleware(), middleware.AdminMiddleware(), AdminCreateTag)
	r_gin.PUT("/admin/tags/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdateTag)
	r_gin.DELETE("/admin/tags/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteTag)

	// Продукты (публичные)
	r_gin.GET("/products", GetProducts)
	r_gin.GET("/products/:id", GetProductById)

	// Уроки для продукта
	r_gin.GET("/product/:product_id/lessons", GetLessonsByProduct)

	// Админские маршруты для продуктов
	r_gin.GET("/admin/products", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetProducts)
	r_gin.POST("/admin/products", AuthMiddleware(), middleware.AdminMiddleware(), AdminCreateProduct)
	r_gin.PUT("/admin/products/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdateProduct)
	r_gin.DELETE("/admin/products/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteProduct)

	// Уроки (публичные с проверкой доступа)
	r_gin.GET("/lessons/:id", GetLessonById)

	// Админские маршруты для уроков
	r_gin.GET("/admin/lessons", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetLessons)
	r_gin.POST("/admin/lessons", AuthMiddleware(), middleware.AdminMiddleware(), AdminCreateLesson)
	r_gin.PUT("/admin/lessons/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdateLesson)
	r_gin.DELETE("/admin/lessons/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteLesson)

	// Платежи и покупки
	r_gin.POST("/payments/create", AuthMiddleware(), CreatePayment)
	r_gin.GET("/payments/:id", AuthMiddleware(), GetPaymentStatus)
	r_gin.POST("/payments/webhook", WebhookYooKassa)
	r_gin.POST("/purchases", AuthMiddleware(), CreatePurchase)

	// Личный кабинет обучения
	r_gin.GET("/learning/my-courses", AuthMiddleware(), GetMyCourses)
	r_gin.GET("/learning/progress/:purchase_id", AuthMiddleware(), GetLearningProgress)
	r_gin.POST("/learning/progress/:lesson_id", AuthMiddleware(), UpdateLessonProgress)

	// Промокоды
	r_gin.POST("/promo-codes/validate", ValidatePromoCode)
	r_gin.GET("/admin/promo-codes", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetPromoCodes)
	r_gin.POST("/admin/promo-codes", AuthMiddleware(), middleware.AdminMiddleware(), AdminCreatePromoCode)
	r_gin.PUT("/admin/promo-codes/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdatePromoCode)
	r_gin.DELETE("/admin/promo-codes/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeletePromoCode)

	// Отзывы
	r_gin.GET("/products/:id/reviews", GetProductReviews)
	r_gin.POST("/products/:id/reviews", AuthMiddleware(), CreateReview)
	r_gin.PUT("/reviews/:id", AuthMiddleware(), UpdateReview)
	r_gin.DELETE("/reviews/:id", AuthMiddleware(), DeleteReview)
	r_gin.GET("/admin/reviews/pending", AuthMiddleware(), middleware.AdminMiddleware(), GetPendingReviews)
	r_gin.PUT("/admin/reviews/:id/approve", AuthMiddleware(), middleware.AdminMiddleware(), ApproveReview)
	r_gin.PUT("/admin/reviews/:id/reject", AuthMiddleware(), middleware.AdminMiddleware(), AdminRejectReview)

	// Admin reviews management (paginated list with filters)
	r_gin.GET("/admin/reviews/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetReviewsList)
	r_gin.POST("/admin/reviews/bulk-approve", AuthMiddleware(), middleware.AdminMiddleware(), AdminBulkApproveReviews)
	r_gin.POST("/admin/reviews/bulk-reject", AuthMiddleware(), middleware.AdminMiddleware(), AdminBulkRejectReviews)
	r_gin.POST("/admin/reviews/bulk-delete", AuthMiddleware(), middleware.AdminMiddleware(), AdminBulkDeleteReviews)
	// Согласие на использование cookie и обработку персональных данных
	r_gin.POST("/consent", AuthMiddleware(), SaveConsent)
	r_gin.GET("/consent/status", AuthMiddleware(), GetConsentStatus)
	r_gin.POST("/consent/withdraw", AuthMiddleware(), WithdrawConsent)
	r_gin.GET("/consent/policy-version", GetCookiePolicyVersion)
	r_gin.GET("/admin/consent/audit", AuthMiddleware(), middleware.AdminMiddleware(), GetConsentAuditLog)

	// Мастер-классы (публичные)
	r_gin.GET("/master-classes", GetMasterClasses)
	r_gin.GET("/master-classes/tags", GetMasterClassTags)
	r_gin.GET("/master-classes/tag/:tag_slug", GetMasterClassesByTag)
	r_gin.GET("/master-classes/:id", GetMasterClassByID)

	// Видео и обложки мастер-классов (с проверкой доступа)
	r_gin.GET("/master-classes/:id/video", AuthMiddlewareOptional(), ServeMasterClassVideo)
	r_gin.GET("/master-classes/:id/thumbnail", ServeMasterClassThumbnail)

	// Чат
	r_gin.GET("/chat/threads", AuthMiddleware(), GetChatThreads)
	r_gin.POST("/chat/threads", AuthMiddleware(), CreateChatThread)
	r_gin.GET("/chat/threads/:id/messages", AuthMiddleware(), GetChatMessages)
	r_gin.POST("/chat/threads/:id/messages", AuthMiddleware(), SendChatMessage)
	r_gin.POST("/chat/messages/:id/read", AuthMiddleware(), MarkMessageAsRead)
	r_gin.GET("/chat/poll", AuthMiddleware(), PollChatMessages)
	r_gin.GET("/admin/chat/threads", AuthMiddleware(), middleware.AdminMiddleware(), GetAdminChatThreads)
	r_gin.GET("/admin/chat/threads/:id/messages", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetChatMessages)
	r_gin.POST("/admin/chat/threads/:id/messages", AuthMiddleware(), middleware.AdminMiddleware(), AdminSendChatMessage)
	r_gin.PUT("/admin/chat/threads/:id/resolve", AuthMiddleware(), middleware.AdminMiddleware(), ResolveChatThread)
	r_gin.PUT("/admin/chat/threads/:id/reopen", AuthMiddleware(), middleware.AdminMiddleware(), AdminReopenChatThread)

	// Admin dashboard
	r_gin.GET("/api/admin/stats", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetStats)
	r_gin.GET("/api/admin/recent-activity", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetRecentActivity)

	// Admin gallery management
	r_gin.GET("/admin/works", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetWorks)
	r_gin.POST("/admin/works/bulk-delete", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteWorks)

	// Admin shop management
	r_gin.GET("/admin/sales", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetSales)
	r_gin.POST("/admin/sales/bulk-delete", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteSales)
	// Admin news management
	r_gin.GET("/admin/news/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetNews)
	r_gin.POST("/admin/news/bulk-delete", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteNews)

	// Admin products management (courses & master-classes)
	r_gin.GET("/admin/products/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetProductsList)
	r_gin.POST("/admin/products/bulk-delete", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteProducts)
	r_gin.PUT("/admin/products/:id/status", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdateProductStatus)

	// Admin lessons management (paginated list)
	r_gin.GET("/admin/lessons/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetLessonsList)
	r_gin.POST("/admin/lessons/bulk-delete", AuthMiddleware(), middleware.AdminMiddleware(), AdminDeleteLessons)

	// Admin users management
	r_gin.GET("/admin/users/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetUsersList)
	r_gin.GET("/admin/users/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetUserDetail)
	r_gin.PUT("/admin/users/:id/role", AuthMiddleware(), middleware.AdminMiddleware(), AdminUpdateUserRole)
	r_gin.PUT("/admin/users/:id/block", AuthMiddleware(), middleware.AdminMiddleware(), AdminToggleUserBlock)

	// Admin purchases management
	r_gin.GET("/admin/purchases/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetPurchasesList)
	r_gin.GET("/admin/purchases/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetPurchaseDetail)
	r_gin.PUT("/admin/purchases/:id/extend", AuthMiddleware(), middleware.AdminMiddleware(), AdminExtendPurchaseAccess)
	r_gin.PUT("/admin/purchases/:id/cancel", AuthMiddleware(), middleware.AdminMiddleware(), AdminCancelPurchase)

	// Admin payments management
	r_gin.GET("/admin/payments/list", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetPaymentsList)
	r_gin.GET("/admin/payments/:id", AuthMiddleware(), middleware.AdminMiddleware(), AdminGetPaymentDetail)
	r_gin.POST("/admin/payments/:id/refund", AuthMiddleware(), middleware.AdminMiddleware(), AdminRefundPayment)

	defer db.Close()
	if err := r_gin.Run("localhost:8000"); err != nil {
		log.Fatal(err)
	}
}
