package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/internal/handler"
	"github.com/raiyin/artserver/internal/repository/file"
	"github.com/raiyin/artserver/internal/repository/sqlite"
	"github.com/raiyin/artserver/internal/router"
	"github.com/raiyin/artserver/internal/service"
	"github.com/raiyin/artserver/pkg/email"
	"github.com/raiyin/artserver/pkg/jwt"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assuming environment variables are set externally")
	}

	// Load config
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	appConfig := config.AppConfigInstance
	log.Printf("Starting %s on port %d", appConfig.App.Name, appConfig.App.Port)

	// -------------------------------------------------------------------------
	// Database
	// -------------------------------------------------------------------------
	db, err := sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err)
	}
	defer db.Close()

	// Enable WAL mode for better concurrent read/write performance
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		log.Printf("Warning: could not set WAL mode: %v", err)
	}

	// Set busy timeout to 5 seconds to avoid "database is locked" errors
	_, err = db.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		log.Printf("Warning: could not set busy_timeout: %v", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		log.Printf("Warning: could not enable foreign keys: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not connect to database (ping failed): %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	log.Println("Database connection established successfully")

	// Run migrations
	runMigrations(db)

	// -------------------------------------------------------------------------
	// JWT Manager
	// -------------------------------------------------------------------------
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mydevsecret" // fallback for development
		log.Println("WARNING: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	jwtManager := jwt.NewManager(
		[]byte(secret),
		15*time.Minute, // access token TTL
		7*24*time.Hour, // refresh token TTL
	)

	// -------------------------------------------------------------------------
	// Email Sender
	// -------------------------------------------------------------------------
	emailSender := email.NewSender(
		appConfig.SMTP.Host,
		appConfig.SMTP.Port,
		appConfig.SMTP.Username,
		appConfig.SMTP.Password,
		appConfig.SMTP.From,
	)

	// -------------------------------------------------------------------------
	// Repositories
	// -------------------------------------------------------------------------
	userRepo := sqlite.NewUserRepository(db)
	workRepo := sqlite.NewWorkRepository(db)
	saleRepo := sqlite.NewSaleRepository(db)
	newsRepo := sqlite.NewNewsRepository(db)
	productRepo := sqlite.NewProductRepository(db)
	categoryRepo := sqlite.NewProductCategoryRepository(db)
	lessonRepo := sqlite.NewLessonRepository(db)
	progressRepo := sqlite.NewLearningProgressRepository(db)
	paymentRepo := sqlite.NewPaymentRepository(db)
	purchaseRepo := sqlite.NewPurchaseRepository(db)
	promoRepo := sqlite.NewPromoCodeRepository(db)
	reviewRepo := sqlite.NewReviewRepository(db)
	threadRepo := sqlite.NewChatThreadRepository(db)
	messageRepo := sqlite.NewChatMessageRepository(db)
	mcRepo := sqlite.NewMasterClassRepository(db)
	tagRepo := sqlite.NewTagRepository(db)
	materialRepo := sqlite.NewMaterialRepository(db)
	baseRepo := sqlite.NewBaseRepository(db)
	consentRepo := sqlite.NewConsentRepository(db)

	// File repository (for images, avatars, etc.)
	fileRepo := file.NewRepository(".")

	// -------------------------------------------------------------------------
	// Services
	// -------------------------------------------------------------------------
	authService := service.NewAuthService(userRepo, jwtManager, emailSender)
	userService := service.NewUserService(userRepo, fileRepo, appConfig.Directories.AbsAvatarsDir)
	galleryService := service.NewGalleryService(workRepo, saleRepo, fileRepo, appConfig.Directories.AbsWorksDir)
	newsService := service.NewNewsService(newsRepo, fileRepo, appConfig.Directories.AbsNewsDir)
	shopService := service.NewShopService(productRepo, categoryRepo, promoRepo, reviewRepo)
	learningService := service.NewLearningService(lessonRepo, progressRepo, purchaseRepo)
	paymentService := service.NewPaymentService(paymentRepo, purchaseRepo, productRepo, shopService)
	chatService := service.NewChatService(threadRepo, messageRepo)
	tagService := service.NewTagService(tagRepo)
	materialService := service.NewMaterialService(materialRepo)
	baseService := service.NewBaseService(baseRepo)
	consentService := service.NewConsentService(consentRepo)
	masterClassService := service.NewMasterClassService(mcRepo)
	adminService := service.NewAdminService(userRepo, productRepo, paymentRepo, purchaseRepo)

	// -------------------------------------------------------------------------
	// Handlers
	// -------------------------------------------------------------------------
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(userService)
	galleryHandler := handler.NewGalleryHandler(galleryService)
	newsHandler := handler.NewNewsHandler(newsService)
	shopHandler := handler.NewShopHandler(shopService)
	learningHandler := handler.NewLearningHandler(learningService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	chatHandler := handler.NewChatHandler(chatService)
	miscHandler := handler.NewMiscHandler(
		adminService,
		tagService,
		materialService,
		baseService,
		consentService,
		masterClassService,
	)

	// -------------------------------------------------------------------------
	// Router
	// -------------------------------------------------------------------------
	r := router.NewRouter(
		authHandler,
		profileHandler,
		galleryHandler,
		newsHandler,
		shopHandler,
		learningHandler,
		paymentHandler,
		chatHandler,
		miscHandler,
		jwtManager,
	)

	// -------------------------------------------------------------------------
	// Start server
	// -------------------------------------------------------------------------
	addr := "localhost:8000"
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

// runMigrations applies database schema migrations.
func runMigrations(db *sql.DB) {
	// Add avatar column to users table if it doesn't exist
	_, err := db.Exec("ALTER TABLE users ADD COLUMN avatar TEXT DEFAULT ''")
	if err != nil {
		log.Printf("Migration (add avatar column): %v (this is normal if column already exists)", err)
	}
	log.Println("Database migrations completed")
}
