package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
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

	// Create news table with the new schema.
	// The old legacy table (id TEXT, title_ru, title_en, datetime, dir, img_back, etc.)
	// may already exist, so we check the schema first.
	var newsIDType string
	err = db.QueryRow("SELECT type FROM pragma_table_info('news') WHERE name = 'id'").Scan(&newsIDType)
	if err != nil {
		// Table doesn't exist at all — create it fresh
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS news (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT NOT NULL,
				description TEXT,
				content TEXT,
				image_path TEXT NOT NULL DEFAULT '',
				video_path TEXT,
				status TEXT NOT NULL DEFAULT 'draft',
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL
			)
		`)
		if err != nil {
			log.Printf("Migration (create news table): %v", err)
		} else {
			log.Println("Migration (create news table): applied successfully")
		}
	} else if strings.EqualFold(newsIDType, "text") {
		// Legacy table exists (id is TEXT) — rename it and create new one
		log.Println("Migration: detected legacy news table, renaming to news_legacy and creating new schema")
		_, err = db.Exec("ALTER TABLE news RENAME TO news_legacy")
		if err != nil {
			log.Printf("Migration (rename legacy news table): %v", err)
		} else {
			_, err = db.Exec(`
				CREATE TABLE news (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					title TEXT NOT NULL,
					description TEXT,
					content TEXT,
					image_path TEXT NOT NULL DEFAULT '',
					video_path TEXT,
					status TEXT NOT NULL DEFAULT 'draft',
					created_at TIMESTAMP NOT NULL,
					updated_at TIMESTAMP NOT NULL
				)
			`)
			if err != nil {
				log.Printf("Migration (create new news table): %v", err)
			} else {
				log.Println("Migration (news table): created new schema, migrating data...")
				// Migrate data from legacy table to new schema.
				// Mapping:
				//   title       ← COALESCE(title_ru, title_en)
				//   description ← COALESCE(NULLIF(subTitle_ru, ''), NULLIF(subTitle_en, ''), substr(text_ru, 1, 200))
				//   content     ← COALESCE(text_ru, text_en)
				//   image_path  ← dir || img_back
				//   video_path  ← dir || first video from videos list
				//   status      ← 'published'
				//   created_at  ← datetime
				//   updated_at  ← datetime
				_, err = db.Exec(`
					INSERT INTO news (title, description, content, image_path, video_path, status, created_at, updated_at)
					SELECT
						COALESCE(NULLIF(title_ru, ''), title_en, ''),
						COALESCE(NULLIF(subTitle_ru, ''), NULLIF(subTitle_en, ''), SUBSTR(COALESCE(NULLIF(text_ru, ''), text_en, ''), 1, 200), ''),
						COALESCE(NULLIF(text_ru, ''), text_en, ''),
						CASE
							WHEN img_back != '' THEN dir || img_back
							ELSE ''
						END,
						CASE
							WHEN videos IS NOT NULL AND videos != '' THEN dir || SUBSTR(videos, 1, INSTR(videos || ';', ';') - 1)
							ELSE NULL
						END,
						'published',
						datetime,
						datetime
					FROM news_legacy
					ORDER BY rowid
				`)
				if err != nil {
					log.Printf("Migration (migrate news data): %v", err)
				} else {
					// Get count of migrated rows
					var count int
					_ = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
					log.Printf("Migration (news table): migrated %d records from legacy → new schema successfully", count)
				}
			}
		}
	} else {
		log.Println("Migration (news table): already has correct schema, skipping")
	}

	// Check if legacy news data needs to be migrated (news_legacy exists but news table is empty)
	var legacyCount int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('news_legacy') WHERE name = 'id'").Scan(&legacyCount)
	if err == nil && legacyCount > 0 {
		var newsCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&newsCount)
		if newsCount == 0 {
			log.Println("Migration: detected news_legacy table with data, migrating to new schema...")
			_, err = db.Exec(`
				INSERT INTO news (title, description, content, image_path, video_path, status, created_at, updated_at)
				SELECT
					COALESCE(NULLIF(title_ru, ''), title_en, ''),
					COALESCE(NULLIF(subTitle_ru, ''), NULLIF(subTitle_en, ''), SUBSTR(COALESCE(NULLIF(text_ru, ''), text_en, ''), 1, 200), ''),
					COALESCE(NULLIF(text_ru, ''), text_en, ''),
					CASE
						WHEN img_back != '' THEN dir || img_back
						ELSE ''
					END,
					CASE
						WHEN videos IS NOT NULL AND videos != '' THEN dir || SUBSTR(videos, 1, INSTR(videos || ';', ';') - 1)
						ELSE NULL
					END,
					'published',
					datetime,
					datetime
				FROM news_legacy
				ORDER BY rowid
			`)
			if err != nil {
				log.Printf("Migration (migrate legacy news data): %v", err)
			} else {
				var count int
				_ = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
				log.Printf("Migration (news data): migrated %d records from news_legacy successfully", count)
			}
		} else {
			log.Printf("Migration (news data): news table already has %d records, skipping legacy migration", newsCount)
		}
	}

	// Check if gallery schema migration is needed (old materials table has 'material_ru' column)
	var needsMigration int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('materials') WHERE name = 'material_ru'").Scan(&needsMigration)
	if err == nil && needsMigration > 0 {
		log.Println("Gallery schema migration needed, applying 003_fix_gallery_schema.sql...")
		applyMigrationFile(db, "003_fix_gallery_schema.sql")
	} else {
		log.Println("Gallery schema migration already applied, skipping")
	}

	// Check if sales data migration is needed (old sales_old table still exists with data)
	var salesOldCount int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales_old') WHERE name = 'name_ru'").Scan(&salesOldCount)
	if err == nil && salesOldCount > 0 {
		log.Println("Sales data migration needed, applying 004_fix_sales_data.sql...")
		applyMigrationFile(db, "004_fix_sales_data.sql")
	} else {
		log.Println("Sales data migration already applied, skipping")
	}

	log.Println("Database migrations completed")
}

// applyMigrationFile reads and executes a SQL migration file from the migrations directory.
func applyMigrationFile(db *sql.DB, filename string) {
	migrationPath := "./db/migrations/" + filename
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Printf("Migration %s: could not read file: %v (skipping)", filename, err)
		return
	}

	// Split by semicolons and execute each statement
	statements := strings.Split(string(migrationSQL), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			log.Printf("Migration %s: statement error: %v\nStatement: %.100s", filename, err, stmt)
		}
	}
	log.Printf("Migration %s applied successfully", filename)
}
