package router

import (
	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/handler"
	"github.com/raiyin/artserver/pkg/jwt"
)

// NewRouter creates and configures a new Gin router with all routes.
func NewRouter(
	authHandler *handler.AuthHandler,
	profileHandler *handler.ProfileHandler,
	galleryHandler *handler.GalleryHandler,
	newsHandler *handler.NewsHandler,
	shopHandler *handler.ShopHandler,
	learningHandler *handler.LearningHandler,
	paymentHandler *handler.PaymentHandler,
	chatHandler *handler.ChatHandler,
	miscHandler *handler.MiscHandler,
	jwtManager *jwt.Manager,
) *gin.Engine {
	r := gin.New()

	// Global middlewares
	r.Use(handler.RecoveryMiddleware())
	r.Use(handler.RequestLoggerMiddleware())
	r.Use(handler.CORSMiddleware())

	// Pre-create middleware instances
	authMw := handler.AuthMiddleware(jwtManager)
	authOptionalMw := handler.AuthMiddlewareOptional(jwtManager)
	adminMw := handler.AdminMiddleware()
	rateLimitMw := handler.RateLimitMiddleware()

	// =========================================================================
	// Auth routes (with rate limiting)
	// =========================================================================
	auth := r.Group("/")
	auth.Use(rateLimitMw)
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/resend-verification", authHandler.ResendVerification)
	}
	r.GET("/verify-email", authHandler.VerifyEmail)

	// =========================================================================
	// User profile routes
	// =========================================================================
	r.GET("/profile", authMw, profileHandler.GetProfile)
	r.PUT("/profile", authMw, profileHandler.UpdateProfile)

	// Avatar routes
	r.POST("/profile/avatar", authMw, profileHandler.UploadAvatar)
	r.DELETE("/profile/avatar", authMw, profileHandler.DeleteAvatar)
	r.GET("/profile/avatar/:id", profileHandler.ServeAvatar)

	// =========================================================================
	// Gallery: Sales
	// =========================================================================
	r.GET("/sales", galleryHandler.GetSales)
	r.GET("/sales/:id", galleryHandler.GetSaleByID)
	r.GET("/sales/:id/edit", galleryHandler.GetSaleByID)
	r.POST("/sales", galleryHandler.CreateSale)
	r.PUT("/sales/:id", authMw, galleryHandler.UpdateSale)
	r.DELETE("/sales/:id", authMw, galleryHandler.DeleteSale)

	// =========================================================================
	// Gallery: Works
	// =========================================================================
	r.GET("/works", galleryHandler.GetWorks)
	r.GET("/works/:id", galleryHandler.GetWorkByID)
	r.GET("/works/:id/edit", galleryHandler.GetWorkByID)
	r.POST("/works", galleryHandler.CreateWork)
	r.PUT("/works/:id", authMw, galleryHandler.UpdateWork)
	r.DELETE("/works/:id", authMw, galleryHandler.DeleteWork)

	// =========================================================================
	// News
	// =========================================================================
	r.GET("/news", newsHandler.GetNews)
	r.GET("/news/:id", newsHandler.GetNewsByID)
	r.POST("/news", authMw, newsHandler.CreateNews)
	r.PUT("/news/:id", authMw, newsHandler.UpdateNews)
	r.DELETE("/news/:id", authMw, newsHandler.DeleteNews)

	// =========================================================================
	// Materials & Bases (public)
	// =========================================================================
	r.GET("/materials", miscHandler.GetMaterials)
	r.GET("/bases", miscHandler.GetBases)

	// =========================================================================
	// Categories (public)
	// =========================================================================
	r.GET("/categories", shopHandler.GetCategories)
	r.GET("/categories/:slug", shopHandler.GetCategories)

	// =========================================================================
	// Products (public)
	// =========================================================================
	r.GET("/products", shopHandler.GetProducts)
	r.GET("/products/:id", shopHandler.GetProductByID)
	r.GET("/category/:category_slug/product/:product_slug", shopHandler.GetProductBySlug)

	// =========================================================================
	// Lessons (public with access check)
	// =========================================================================
	r.GET("/product/:product_id/lessons", authOptionalMw, learningHandler.GetLessonsByProduct)
	r.GET("/lessons/:id", authOptionalMw, learningHandler.GetLessonByID)

	// =========================================================================
	// Payments & Purchases
	// =========================================================================
	r.POST("/payments/create", authMw, paymentHandler.CreatePayment)
	r.GET("/payments/:id", authMw, paymentHandler.GetPaymentStatus)
	r.POST("/payments/webhook", paymentHandler.HandleWebhook)
	r.POST("/purchases", authMw, paymentHandler.CreatePurchase)

	// =========================================================================
	// Learning (personal cabinet)
	// =========================================================================
	r.GET("/learning/my-courses", authMw, learningHandler.GetMyCourses)
	r.GET("/learning/progress/:product_id", authMw, learningHandler.GetLearningProgress)
	r.POST("/learning/progress/:id", authMw, learningHandler.UpdateLessonProgress)

	// =========================================================================
	// Promo codes
	// =========================================================================
	r.POST("/promo-codes/validate", shopHandler.ValidatePromoCode)

	// =========================================================================
	// Reviews
	// =========================================================================
	r.GET("/products/:id/reviews", shopHandler.GetReviewsByProduct)
	r.POST("/products/:id/reviews", authMw, shopHandler.CreateReview)
	r.PUT("/reviews/:id", authMw, shopHandler.UpdateReview)
	r.DELETE("/reviews/:id", authMw, shopHandler.DeleteReview)

	// =========================================================================
	// Consent
	// =========================================================================
	r.POST("/consent", authMw, miscHandler.SaveConsent)
	r.GET("/consent/status", authMw, miscHandler.GetConsentStatus)

	// =========================================================================
	// Master classes (public)
	// =========================================================================
	r.GET("/master-classes", miscHandler.GetMasterClasses)
	r.GET("/master-classes/:id", miscHandler.GetMasterClassByID)

	// =========================================================================
	// Chat
	// =========================================================================
	r.GET("/chat/threads", authMw, chatHandler.GetThreads)
	r.POST("/chat/threads", authMw, chatHandler.CreateThread)
	r.GET("/chat/threads/:thread_id/messages", authMw, chatHandler.GetMessages)
	r.POST("/chat/threads/:thread_id/messages", authMw, chatHandler.SendMessage)
	r.POST("/chat/messages/:message_id/read", authMw, chatHandler.MarkMessageAsRead)
	r.GET("/chat/poll", authMw, chatHandler.PollMessages)

	// =========================================================================
	// Admin routes
	// =========================================================================
	admin := r.Group("/admin")
	admin.Use(authMw)
	admin.Use(adminMw)
	{
		// Categories
		admin.GET("/categories", shopHandler.GetCategories)
		admin.POST("/categories", shopHandler.CreateCategory)
		admin.PUT("/categories/:id", shopHandler.UpdateCategory)
		admin.DELETE("/categories/:id", shopHandler.DeleteCategory)

		// Tags
		admin.GET("/tags", miscHandler.GetTags)
		admin.POST("/tags", miscHandler.CreateTag)
		admin.PUT("/tags/:id", miscHandler.UpdateTag)
		admin.DELETE("/tags/:id", miscHandler.DeleteTag)

		// Products
		admin.GET("/products", shopHandler.GetProducts)
		admin.POST("/products", shopHandler.CreateProduct)
		admin.PUT("/products/:id", shopHandler.UpdateProduct)
		admin.DELETE("/products/:id", shopHandler.DeleteProduct)

		// Lessons
		admin.GET("/lessons", learningHandler.AdminGetLessons)
		admin.POST("/lessons", learningHandler.CreateLesson)
		admin.PUT("/lessons/:id", learningHandler.UpdateLesson)
		admin.DELETE("/lessons/:id", learningHandler.DeleteLesson)

		// Promo codes
		admin.GET("/promo-codes", shopHandler.GetPromoCodes)
		admin.POST("/promo-codes", shopHandler.CreatePromoCode)
		admin.PUT("/promo-codes/:id", shopHandler.UpdatePromoCode)
		admin.DELETE("/promo-codes/:id", shopHandler.DeletePromoCode)

		// Reviews
		admin.GET("/reviews/list", shopHandler.GetReviews)
		admin.PUT("/reviews/:id", shopHandler.UpdateReview)
		admin.DELETE("/reviews/:id", shopHandler.DeleteReview)

		// Chat
		admin.GET("/chat/threads", chatHandler.GetThreads)
		admin.GET("/chat/threads/:thread_id/messages", chatHandler.AdminGetMessages)
		admin.POST("/chat/threads/:thread_id/messages", chatHandler.AdminSendMessage)

		// Dashboard
		admin.GET("/stats", miscHandler.GetDashboardStats)

		// Gallery management
		admin.GET("/works", galleryHandler.GetWorks)
		admin.GET("/sales", galleryHandler.GetSales)

		// News management
		admin.GET("/news/list", newsHandler.GetNews)

		// Products list
		admin.GET("/products/list", shopHandler.GetProducts)

		// Lessons list
		admin.GET("/lessons/list", learningHandler.AdminGetLessons)

		// Payments
		admin.GET("/payments/list", paymentHandler.GetPayments)

		// Purchases
		admin.GET("/purchases/list", paymentHandler.GetPurchases)

		// Materials
		admin.POST("/materials", miscHandler.CreateMaterial)
		admin.DELETE("/materials/:id", miscHandler.DeleteMaterial)

		// Bases
		admin.POST("/bases", miscHandler.CreateBase)
		admin.DELETE("/bases/:id", miscHandler.DeleteBase)

		// Master classes
		admin.POST("/master-classes", miscHandler.CreateMasterClass)
		admin.PUT("/master-classes/:id", miscHandler.UpdateMasterClass)
		admin.DELETE("/master-classes/:id", miscHandler.DeleteMasterClass)

		// Consent audit
		admin.GET("/consent/audit", miscHandler.GetConsentStatus)
	}

	// =========================================================================
	// API admin routes (separate prefix)
	// =========================================================================
	apiAdmin := r.Group("/api/admin")
	apiAdmin.Use(authMw)
	apiAdmin.Use(adminMw)
	{
		apiAdmin.GET("/stats", miscHandler.GetDashboardStats)
		apiAdmin.GET("/recent-activity", miscHandler.GetRecentActivity)
	}

	return r
}
