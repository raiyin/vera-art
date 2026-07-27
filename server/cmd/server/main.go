package main

import (
	"database/sql"
	"log/slog"
	"os"
	"strings"
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

	// Set busy timeout to 5 seconds to avoid "database is locked" errors
	_, err = db.Exec("PRAGMA busy_timeout=5000")
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

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
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

// runMigrations applies database schema migrations.
func runMigrations(db *sql.DB) {
	// Add avatar column to users table if it doesn't exist
	_, err := db.Exec("ALTER TABLE users ADD COLUMN avatar TEXT DEFAULT ''")
	if err != nil {
		slog.Warn("Migration (add avatar column) — this is normal if column already exists", "error", err)
	}

	// Add username column to users table if it doesn't exist
	_, err = db.Exec("ALTER TABLE users ADD COLUMN username TEXT DEFAULT ''")
	if err != nil {
		slog.Warn("Migration (add username column) — this is normal if column already exists", "error", err)
	}

	// Populate username from email for existing users where username is empty
	_, err = db.Exec("UPDATE users SET username = email WHERE username IS NULL OR username = ''")
	if err != nil {
		slog.Warn("Migration (populate username)", "error", err)
	}

	// Add unique index on username
	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username)")
	if err != nil {
		slog.Warn("Migration (create username unique index)", "error", err)
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
				video_paths TEXT,
				image_paths TEXT,
				status TEXT NOT NULL DEFAULT 'draft',
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL
			)
		`)
		if err != nil {
			slog.Error("Migration (create news table)", "error", err)
		} else {
			slog.Info("Migration (create news table): applied successfully")
		}
	} else if strings.EqualFold(newsIDType, "text") {
		// Legacy table exists (id is TEXT) — rename it and create new one
		slog.Info("Migration: detected legacy news table, renaming to news_legacy and creating new schema")
		_, err = db.Exec("ALTER TABLE news RENAME TO news_legacy")
		if err != nil {
			slog.Error("Migration (rename legacy news table)", "error", err)
		} else {
			_, err = db.Exec(`
				CREATE TABLE news (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					title TEXT NOT NULL,
					description TEXT,
					content TEXT,
					image_path TEXT NOT NULL DEFAULT '',
					video_path TEXT,
					video_paths TEXT,
					image_paths TEXT,
					status TEXT NOT NULL DEFAULT 'draft',
					created_at TIMESTAMP NOT NULL,
					updated_at TIMESTAMP NOT NULL
				)
			`)
			if err != nil {
				slog.Error("Migration (create new news table)", "error", err)
			} else {
				slog.Info("Migration (news table): created new schema, migrating data...")
				// Migrate data from legacy table to new schema.
				// Mapping:
				//   title       ← COALESCE(title_ru, title_en)
				//   description ← COALESCE(NULLIF(subTitle_ru, ''), NULLIF(subTitle_en, ''), substr(text_ru, 1, 200))
				//   content     ← COALESCE(text_ru, text_en)
				//   image_path  ← dir || img_back
				//   video_path  ← dir || first video from videos list
				//   video_paths ← semicolon-separated list of dir || video filename for each video
				//   status      ← 'published'
				//   created_at  ← datetime
				//   updated_at  ← datetime
				_, err = db.Exec(`
					INSERT INTO news (title, description, content, image_path, video_path, video_paths, image_paths, status, created_at, updated_at)
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
						NULL,
						CASE
							WHEN images IS NOT NULL AND images != '' THEN
								(SELECT group_concat(nl2.dir || TRIM(value), ';') FROM json_each('["' || replace(replace(nl2.images, CHAR(13), ''), ';', '","') || '"]') WHERE TRIM(value) != '')
							ELSE NULL
						END,
						'published',
						datetime,
						datetime
					FROM news_legacy nl2
					ORDER BY rowid
				`)
				if err != nil {
					slog.Error("Migration (migrate news data)", "error", err)
				} else {
					// Get count of migrated rows
					var count int
					_ = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
					slog.Info("Migration (news table): migrated records from legacy to new schema", "count", count)
				}
			}
		}
	} else {
		slog.Info("Migration (news table): already has correct schema, skipping")
	}

	// Add video_paths column if it doesn't exist (for existing installations)
	_, err = db.Exec("ALTER TABLE news ADD COLUMN video_paths TEXT")
	if err != nil {
		slog.Warn("Migration (add video_paths column) — this is normal if column already exists", "error", err)
	}

	// Add image_paths column if it doesn't exist (for existing installations)
	_, err = db.Exec("ALTER TABLE news ADD COLUMN image_paths TEXT")
	if err != nil {
		slog.Warn("Migration (add image_paths column) — this is normal if column already exists", "error", err)
	}

	// Populate video_paths from legacy data if it's empty and news_legacy exists.
	// We match by dir prefix: news.image_path starts with news_legacy.dir.
	// Uses json_each to split semicolon-separated video names and construct full paths.
	var emptyVideoPathsCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM news WHERE (video_paths IS NULL OR video_paths = '') AND EXISTS (SELECT 1 FROM news_legacy)").Scan(&emptyVideoPathsCount)
	if emptyVideoPathsCount > 0 {
		slog.Info("Migration: populating video_paths from legacy data", "count", emptyVideoPathsCount)

		_, err = db.Exec(`
			UPDATE news SET video_paths = (
				SELECT group_concat(nl.dir || TRIM(value), ';')
				FROM news_legacy nl
				JOIN json_each('["' || replace(replace(nl.videos, CHAR(13), ''), ';', '","') || '"]')
				WHERE nl.videos IS NOT NULL AND nl.videos != ''
				  AND TRIM(value) != ''
				  AND news.image_path LIKE nl.dir || '%'
			)
			WHERE EXISTS (
				SELECT 1 FROM news_legacy nl
				WHERE nl.videos IS NOT NULL AND nl.videos != ''
				  AND news.image_path LIKE nl.dir || '%'
			)
		`)
		if err != nil {
			slog.Error("Migration (populate video_paths)", "error", err)
		} else {
			var updated int
			_ = db.QueryRow("SELECT changes()").Scan(&updated)
			slog.Info("Migration (populate video_paths): updated records", "count", updated)
		}
	}

	// Populate image_paths from legacy data if it's empty and news_legacy exists.
	var emptyImagePathsCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM news WHERE (image_paths IS NULL OR image_paths = '') AND EXISTS (SELECT 1 FROM news_legacy)").Scan(&emptyImagePathsCount)
	if emptyImagePathsCount > 0 {
		slog.Info("Migration: populating image_paths from legacy data", "count", emptyImagePathsCount)

		_, err = db.Exec(`
			UPDATE news SET image_paths = (
				SELECT group_concat(nl.dir || TRIM(value), ';')
				FROM news_legacy nl
				JOIN json_each('["' || replace(replace(nl.images, CHAR(13), ''), ';', '","') || '"]')
				WHERE nl.images IS NOT NULL AND nl.images != ''
				  AND TRIM(value) != ''
				  AND news.image_path LIKE nl.dir || '%'
			)
			WHERE EXISTS (
				SELECT 1 FROM news_legacy nl
				WHERE nl.images IS NOT NULL AND nl.images != ''
				  AND news.image_path LIKE nl.dir || '%'
			)
		`)
		if err != nil {
			slog.Error("Migration (populate image_paths)", "error", err)
		} else {
			var updated int
			_ = db.QueryRow("SELECT changes()").Scan(&updated)
			slog.Info("Migration (populate image_paths): updated records", "count", updated)
		}
	}

	// Check if legacy news data needs to be migrated (news_legacy exists but news table is empty)
	var legacyCount int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('news_legacy') WHERE name = 'id'").Scan(&legacyCount)
	if err == nil && legacyCount > 0 {
		var newsCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&newsCount)
		if newsCount == 0 {
			slog.Info("Migration: detected news_legacy table with data, migrating to new schema...")
			_, err = db.Exec(`
				INSERT INTO news (title, description, content, image_path, video_path, video_paths, image_paths, status, created_at, updated_at)
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
					NULL,
					CASE
						WHEN images IS NOT NULL AND images != '' THEN
							(SELECT group_concat(nl2.dir || TRIM(value), ';') FROM json_each('["' || replace(replace(nl2.images, CHAR(13), ''), ';', '","') || '"]') WHERE TRIM(value) != '')
						ELSE NULL
					END,
					'published',
					datetime,
					datetime
				FROM news_legacy
				ORDER BY rowid
			`)
			if err != nil {
				slog.Error("Migration (migrate legacy news data)", "error", err)
			} else {
				var count int
				_ = db.QueryRow("SELECT COUNT(*) FROM news").Scan(&count)
				slog.Info("Migration (news data): migrated records from news_legacy", "count", count)
			}
		} else {
			slog.Info("Migration (news data): news table already has records, skipping legacy migration", "count", newsCount)
		}
	}

	// Check if gallery schema migration is needed (old materials table has 'material_ru' column)
	var needsMigration int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('materials') WHERE name = 'material_ru'").Scan(&needsMigration)
	if err == nil && needsMigration > 0 {
		slog.Info("Gallery schema migration needed, applying 003_fix_gallery_schema.sql...")
		applyMigrationFile(db, "003_fix_gallery_schema.sql")
	} else {
		slog.Info("Gallery schema migration already applied, skipping")
	}

	// Check if sales data migration is needed (old sales_old table still exists with data)
	var salesOldCount int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales_old') WHERE name = 'name_ru'").Scan(&salesOldCount)
	if err == nil && salesOldCount > 0 {
		slog.Info("Sales data migration needed, applying 004_fix_sales_data.sql...")
		applyMigrationFile(db, "004_fix_sales_data.sql")
	} else {
		slog.Info("Sales data migration already applied, skipping")
	}

	// Add email_verified column to users table if it doesn't exist
	_, err = db.Exec("ALTER TABLE users ADD COLUMN email_verified INTEGER NOT NULL DEFAULT 0")
	if err != nil {
		slog.Warn("Migration (add email_verified column) — this is normal if column already exists", "error", err)
	}

	// Add verification_token column to users table if it doesn't exist
	_, err = db.Exec("ALTER TABLE users ADD COLUMN verification_token TEXT")
	if err != nil {
		slog.Warn("Migration (add verification_token column) — this is normal if column already exists", "error", err)
	}

	// Add verification_sent_at column to users table if it doesn't exist
	_, err = db.Exec("ALTER TABLE users ADD COLUMN verification_sent_at TIMESTAMP")
	if err != nil {
		slog.Warn("Migration (add verification_sent_at column) — this is normal if column already exists", "error", err)
	}

	// Create chat_threads table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_threads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			subject TEXT,
			status TEXT NOT NULL DEFAULT 'open',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		slog.Error("Migration (create chat_threads table)", "error", err)
	} else {
		slog.Info("Migration (chat_threads table): applied successfully")
	}

	// Create chat_messages table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			thread_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			is_admin INTEGER NOT NULL DEFAULT 0,
			is_read INTEGER NOT NULL DEFAULT 0,
			read_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (thread_id) REFERENCES chat_threads(id)
		)
	`)
	if err != nil {
		slog.Error("Migration (create chat_messages table)", "error", err)
	} else {
		slog.Info("Migration (chat_messages table): applied successfully")
	}

	// Create reviews table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			rating INTEGER NOT NULL DEFAULT 0,
			text TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		slog.Error("Migration (create reviews table)", "error", err)
	} else {
		slog.Info("Migration (reviews table): applied successfully")
	}

	// Create payments table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'RUB',
			status TEXT NOT NULL DEFAULT 'pending',
			payment_method TEXT,
			yookassa_id TEXT,
			promo_code TEXT,
			discount REAL,
			metadata TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		slog.Error("Migration (create payments table)", "error", err)
	} else {
		slog.Info("Migration (payments table): applied successfully")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS promo_codes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			discount_percent REAL NOT NULL,
			max_uses INTEGER NOT NULL DEFAULT 0,
			current_uses INTEGER NOT NULL DEFAULT 0,
			expires_at TIMESTAMP NOT NULL,
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		slog.Error("Migration (create promo_codes table)", "error", err)
	} else {
		slog.Info("Migration (promo_codes table): applied successfully")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name_ru TEXT NOT NULL UNIQUE,
			name_en TEXT NOT NULL UNIQUE,
			slug TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		slog.Error("Migration (create tags table)", "error", err)
	} else {
		slog.Info("Migration (tags table): applied successfully")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			UNIQUE(product_id, tag_id)
		)
	`)
	if err != nil {
		slog.Error("Migration (create product_tags table)", "error", err)
	} else {
		slog.Info("Migration (product_tags table): applied successfully")
	}

	// Add width and height columns to works table, populating from size
	_, err = db.Exec("ALTER TABLE works ADD COLUMN width INTEGER DEFAULT 0")
	if err != nil {
		slog.Warn("Migration (add width column to works) — this is normal if column already exists", "error", err)
	} else {
		_, err = db.Exec("UPDATE works SET width = CAST(SUBSTR(size, 1, INSTR(size || 'x', 'x') - 1) AS INTEGER) WHERE size != '' AND size IS NOT NULL")
		if err != nil {
			slog.Warn("Migration (populate width from size)", "error", err)
		}
	}
	_, err = db.Exec("ALTER TABLE works ADD COLUMN height INTEGER DEFAULT 0")
	if err != nil {
		slog.Warn("Migration (add height column to works) — this is normal if column already exists", "error", err)
	} else {
		_, err = db.Exec("UPDATE works SET height = CAST(SUBSTR(size, INSTR(size || 'x', 'x') + 1) AS INTEGER) WHERE size != '' AND size IS NOT NULL")
		if err != nil {
			slog.Warn("Migration (populate height from size)", "error", err)
		}
	}

	// Add width and height columns to sales table, populating from size
	_, err = db.Exec("ALTER TABLE sales ADD COLUMN width INTEGER DEFAULT 0")
	if err != nil {
		slog.Warn("Migration (add width column to sales) — this is normal if column already exists", "error", err)
	} else {
		_, err = db.Exec("UPDATE sales SET width = CAST(SUBSTR(size, 1, INSTR(size || 'x', 'x') - 1) AS INTEGER) WHERE size != '' AND size IS NOT NULL")
		if err != nil {
			slog.Warn("Migration (populate width from size)", "error", err)
		}
	}
	_, err = db.Exec("ALTER TABLE sales ADD COLUMN height INTEGER DEFAULT 0")
	if err != nil {
		slog.Warn("Migration (add height column to sales) — this is normal if column already exists", "error", err)
	} else {
		_, err = db.Exec("UPDATE sales SET height = CAST(SUBSTR(size, INSTR(size || 'x', 'x') + 1) AS INTEGER) WHERE size != '' AND size IS NOT NULL")
		if err != nil {
			slog.Warn("Migration (populate height from size)", "error", err)
		}
	}

	// Drop size column from works table (data already migrated to width/height)
	_, err = db.Exec("ALTER TABLE works DROP COLUMN size")
	if err != nil {
		slog.Warn("Migration (drop size column from works) — this is normal if column was already dropped", "error", err)
	} else {
		slog.Info("Migration (drop size column from works): applied successfully")
	}

	// Drop size column from sales table (data already migrated to width/height)
	_, err = db.Exec("ALTER TABLE sales DROP COLUMN size")
	if err != nil {
		slog.Warn("Migration (drop size column from sales) — this is normal if column was already dropped", "error", err)
	} else {
		slog.Info("Migration (drop size column from sales): applied successfully")
	}

	// Migrate material associations from old junction tables if needed
	var workMatCount int
	err = db.QueryRow("SELECT COUNT(*) FROM works_materials").Scan(&workMatCount)
	if err == nil && workMatCount == 0 {
		_, err = db.Exec(`INSERT OR IGNORE INTO works_materials (work_id, material_id)
			SELECT work_id, material_id FROM works_materials
			WHERE work_id IN (SELECT id FROM works)`)
		if err != nil {
			slog.Warn("Migration (migrate works_materials from works_materials)", "error", err)
		} else {
			slog.Info("Migration (migrate works_materials from works_materials): applied successfully")
		}
	}

	var saleMatCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sale_materials").Scan(&saleMatCount)
	if err == nil && saleMatCount == 0 {
		// Check if old sales_materials table exists
		var exists int
		if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales_materials') WHERE name = 'sale_id'").Scan(&exists); err == nil && exists > 0 {
			_, err = db.Exec(`INSERT OR IGNORE INTO sale_materials (sale_id, material_id)
				SELECT sale_id, material_id FROM sales_materials
				WHERE sale_id IN (SELECT id FROM sales)`)
			if err != nil {
				slog.Warn("Migration (migrate sale_materials from sales_materials)", "error", err)
			} else {
				slog.Info("Migration (migrate sale_materials from sales_materials): applied successfully")
			}
		}
	}

	slog.Info("Database migrations completed")
}

// applyMigrationFile reads and executes a SQL migration file from the migrations directory.
func applyMigrationFile(db *sql.DB, filename string) {
	migrationPath := "./db/migrations/" + filename
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		slog.Warn("Migration: could not read file, skipping", "filename", filename, "error", err)
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
			slog.Error("Migration: statement error", "filename", filename, "error", err, "statement", stmt[:min(len(stmt), 100)])
		}
	}
	slog.Info("Migration applied successfully", "filename", filename)
}
