package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/internal/handler"
	"github.com/raiyin/artserver/internal/logger"
	"github.com/raiyin/artserver/internal/repository/file"
	"github.com/raiyin/artserver/internal/repository/sqlite"
	"github.com/raiyin/artserver/internal/router"
	"github.com/raiyin/artserver/internal/service"
	"github.com/raiyin/artserver/pkg/email"
	"github.com/raiyin/artserver/pkg/jwt"
)

func main() {
	// Setup structured logger with colored output
	slog.SetDefault(slog.New(logger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file, assuming environment variables are set externally", "error", err)
	}

	// Load config
	if err := config.LoadConfig("."); err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}
	appConfig := config.AppConfigInstance
	slog.Info("Starting server", "name", appConfig.App.Name, "port", appConfig.App.Port)
	if !appConfig.Features.RegistrationEnabled {
		slog.Warn("Registration is DISABLED by feature flag")
	}

	// -------------------------------------------------------------------------
	// Database
	// -------------------------------------------------------------------------
	db, err := sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		slog.Error("Could not open database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Enable WAL mode for better concurrent read/write performance
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		slog.Warn("Could not set WAL mode", "error", err)
	}

	// Set busy timeout to 10 seconds to avoid "database is locked" errors
	_, err = db.Exec("PRAGMA busy_timeout=10000")
	if err != nil {
		slog.Warn("Could not set busy_timeout", "error", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		slog.Warn("Could not enable foreign keys", "error", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		slog.Error("Could not connect to database (ping failed)", "error", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(3)
	slog.Info("Database connection established successfully")

	// Run migrations
	runMigrations(db)

	// -------------------------------------------------------------------------
	// JWT Manager
	// -------------------------------------------------------------------------
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mydevsecret" // fallback for development
		slog.Warn("Using default JWT secret. Set JWT_SECRET environment variable in production.")
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
	galleryService := service.NewGalleryService(workRepo, saleRepo, fileRepo, appConfig.Directories.AbsWorksDir, appConfig.Directories.RelWorksDir)
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
	adminService := service.NewAdminService(userRepo, productRepo, purchaseRepo, workRepo, newsRepo, reviewRepo, threadRepo, mcRepo)

	// -------------------------------------------------------------------------
	// Handlers
	// -------------------------------------------------------------------------
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(userService)
	galleryHandler := handler.NewGalleryHandler(
		galleryService,
		appConfig.Directories.AbsWorksDir,
		appConfig.Directories.RelWorksDir,
		appConfig.Directories.AbsSalesDir,
		appConfig.Directories.RelSalesDir,
	)
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
		userService,
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
		appConfig.Features,
	)

	// -------------------------------------------------------------------------
	// Start server
	// -------------------------------------------------------------------------
	addr := "localhost:8000"
	slog.Info("Server starting", "address", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func runMigrations(db *sql.DB) {
	// Align the sales table with the schema expected by the code.
	// Older databases used `title` (and `old_price`); the code now expects
	// `name_ru`, `name_en` and `sale_path`.
	var nameRuExists int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales') WHERE name = 'name_ru'").Scan(&nameRuExists); err == nil && nameRuExists == 0 {
		if _, err := db.Exec("ALTER TABLE sales ADD COLUMN name_ru TEXT NOT NULL DEFAULT ''"); err != nil {
			slog.Warn("Migration (add name_ru to sales) failed", "error", err)
		} else {
			slog.Info("Migration (add name_ru to sales): column added successfully")
			if _, err := db.Exec("UPDATE sales SET name_ru = title WHERE title IS NOT NULL AND title != ''"); err != nil {
				slog.Warn("Migration (populate sales name_ru from title) failed", "error", err)
			}
		}
	}

	var nameEnExists int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales') WHERE name = 'name_en'").Scan(&nameEnExists); err == nil && nameEnExists == 0 {
		if _, err := db.Exec("ALTER TABLE sales ADD COLUMN name_en TEXT NOT NULL DEFAULT ''"); err != nil {
			slog.Warn("Migration (add name_en to sales) failed", "error", err)
		} else {
			slog.Info("Migration (add name_en to sales): column added successfully")
		}
	}

	var salePathExists int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales') WHERE name = 'sale_path'").Scan(&salePathExists); err == nil && salePathExists == 0 {
		if _, err := db.Exec("ALTER TABLE sales ADD COLUMN sale_path TEXT DEFAULT ''"); err != nil {
			slog.Warn("Migration (add sale_path to sales) failed", "error", err)
		} else {
			slog.Info("Migration (add sale_path to sales): column added successfully")
		}
	}

	// Drop legacy sales columns no longer used by the code.
	var titleExists int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales') WHERE name = 'title'").Scan(&titleExists); err == nil && titleExists > 0 {
		if _, err := db.Exec("ALTER TABLE sales DROP COLUMN title"); err != nil {
			slog.Warn("Migration (drop title from sales) failed", "error", err)
		} else {
			slog.Info("Migration (drop title from sales): applied successfully")
		}
	}
	var oldPriceExists int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales') WHERE name = 'old_price'").Scan(&oldPriceExists); err == nil && oldPriceExists > 0 {
		if _, err := db.Exec("ALTER TABLE sales DROP COLUMN old_price"); err != nil {
			slog.Warn("Migration (drop old_price from sales) failed", "error", err)
		} else {
			slog.Info("Migration (drop old_price from sales): applied successfully")
		}
	}

	slog.Info("Database migrations completed")
}
