package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/handler"
	"github.com/raiyin/artserver/internal/repository/file"
	"github.com/raiyin/artserver/internal/repository/sqlite"
	"github.com/raiyin/artserver/internal/router"
	"github.com/raiyin/artserver/internal/service"
	"github.com/raiyin/artserver/pkg/email"
	"github.com/raiyin/artserver/pkg/jwt"
)

// setupServer builds a fully wired HTTP router backed by a temporary SQLite
// database and a temporary filesystem directory. Each call is isolated.
func setupServer(t *testing.T) (*gin.Engine, *sql.DB, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	baseDir := t.TempDir()
	db := openTestDB(t)

	userRepo := sqlite.NewUserRepository(db)
	workRepo := sqlite.NewWorkRepository(db)
	saleRepo := sqlite.NewSaleRepository(db)
	newsRepo := sqlite.NewNewsRepository(db)
	productRepo := sqlite.NewProductRepository(db)
	categoryRepo := sqlite.NewProductCategoryRepository(db)
	promoRepo := sqlite.NewPromoCodeRepository(db)
	reviewRepo := sqlite.NewReviewRepository(db)
	lessonRepo := sqlite.NewLessonRepository(db)
	progressRepo := sqlite.NewLearningProgressRepository(db)
	purchaseRepo := sqlite.NewPurchaseRepository(db)
	paymentRepo := sqlite.NewPaymentRepository(db)
	threadRepo := sqlite.NewChatThreadRepository(db)
	messageRepo := sqlite.NewChatMessageRepository(db)
	tagRepo := sqlite.NewTagRepository(db)
	materialRepo := sqlite.NewMaterialRepository(db)
	baseRepo := sqlite.NewBaseRepository(db)
	consentRepo := sqlite.NewConsentRepository(db)
	mcRepo := sqlite.NewMasterClassRepository(db)

	fileRepo := file.NewRepository(baseDir)

	worksDir := filepath.Join(baseDir, "content", "works")
	salesDir := filepath.Join(baseDir, "content", "sales")
	newsDir := filepath.Join(baseDir, "content", "news")

	jwtManager := jwt.NewManager([]byte("integration-test-secret"), 15*time.Minute, 24*time.Hour)
	emailSender := email.NewSender("", 0, "", "", "")

	authService := service.NewAuthService(userRepo, jwtManager, emailSender)
	userService := service.NewUserService(userRepo, fileRepo, filepath.Join(baseDir, "avatars"))
	galleryService := service.NewGalleryService(workRepo, saleRepo, fileRepo, worksDir, "/content/works/")
	newsService := service.NewNewsService(newsRepo, fileRepo, newsDir)
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

	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(userService)
	galleryHandler := handler.NewGalleryHandler(galleryService, worksDir, "/content/works/", salesDir, "/content/sales/")
	newsHandler := handler.NewNewsHandler(newsService)
	shopHandler := handler.NewShopHandler(shopService)
	learningHandler := handler.NewLearningHandler(learningService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	chatHandler := handler.NewChatHandler(chatService)
	miscHandler := handler.NewMiscHandler(adminService, tagService, materialService, baseService, consentService, masterClassService, userService)

	r := router.NewRouter(
		authHandler, profileHandler, galleryHandler, newsHandler, shopHandler,
		learningHandler, paymentHandler, chatHandler, miscHandler,
		jwtManager, config.FeaturesConfig{RegistrationEnabled: true},
	)

	return r, db, baseDir
}

// openTestDB opens a fresh SQLite database with the schema used in production
// (tables referenced by the exercised flows) and registers its cleanup.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })

	schema := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			email TEXT,
			name TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			avatar_path TEXT DEFAULT '',
			email_verified INTEGER NOT NULL DEFAULT 0,
			verification_token TEXT,
			verification_sent_at TIMESTAMP,
			avatar TEXT DEFAULT ''
		)`,
		`CREATE TABLE works (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			str_id TEXT NOT NULL DEFAULT '',
			width INTEGER NOT NULL DEFAULT 0,
			height INTEGER NOT NULL DEFAULT 0,
			year INTEGER NOT NULL DEFAULT 0,
			name_ru TEXT NOT NULL DEFAULT '',
			name_en TEXT NOT NULL DEFAULT '',
			base_id INTEGER NOT NULL DEFAULT 1,
			descr_ru TEXT DEFAULT '',
			descr_en TEXT DEFAULT '',
			work_path TEXT DEFAULT '',
			images TEXT DEFAULT ''
		)`,
		`CREATE TABLE works_materials (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			work_id INTEGER NOT NULL,
			material_id INTEGER NOT NULL
		)`,
		`CREATE TABLE sales (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			descr_ru TEXT DEFAULT '',
			descr_en TEXT DEFAULT '',
			image_path TEXT DEFAULT '',
			price REAL DEFAULT 0,
			year INTEGER,
			technique TEXT DEFAULT '',
			status TEXT DEFAULT 'published',
			sort_order INTEGER DEFAULT 0,
			sold INTEGER DEFAULT 0,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			width INTEGER DEFAULT 0,
			height INTEGER DEFAULT 0,
			name_ru TEXT NOT NULL DEFAULT '',
			name_en TEXT NOT NULL DEFAULT '',
			sale_path TEXT DEFAULT ''
		)`,
		`CREATE TABLE sales_materials (
			sale_id INTEGER NOT NULL,
			material_id INTEGER NOT NULL,
			PRIMARY KEY (sale_id, material_id)
		)`,
		`CREATE TABLE sales_bases (
			sale_id INTEGER NOT NULL,
			base_id INTEGER NOT NULL,
			PRIMARY KEY (sale_id, base_id)
		)`,
		`CREATE TABLE news (
			id TEXT NOT NULL PRIMARY KEY,
			datetime TEXT NOT NULL,
			title_ru TEXT NOT NULL,
			title_en TEXT NOT NULL,
			dir TEXT NOT NULL,
			main_image TEXT NOT NULL,
			text_ru TEXT NOT NULL,
			text_en TEXT NOT NULL,
			images TEXT,
			videos TEXT
		)`,
	}

	for _, stmt := range schema {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("exec schema: %v\n%s", err, stmt)
		}
	}
	return db
}

// seedAdminUser inserts an administrator user (role "admin") so the test can
// log in the same way a user does through the /login endpoint.
func seedAdminUser(t *testing.T, db *sql.DB) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash admin password: %v", err)
	}

	user := &domain.User{
		Username:     "admin",
		Email:        "admin@test.local",
		PasswordHash: string(hash),
		Name:         "Admin",
		Role:         "admin",
	}
	if err := sqlite.NewUserRepository(db).Create(context.Background(), user); err != nil {
		t.Fatalf("create admin user: %v", err)
	}
}

const adminUsername = "admin"
const adminPassword = "Secret123!"

// login performs a real POST /login and returns the access token.
func login(t *testing.T, r http.Handler) string {
	t.Helper()
	body := strings.NewReader(fmt.Sprintf(`{"username":%q,"password":%q}`, adminUsername, adminPassword))

	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login: status %d, body %s", w.Code, w.Body.String())
	}

	var resp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatalf("login returned empty access token")
	}
	return resp.AccessToken
}

// doRequest performs an HTTP request against the router and returns the recorder.
func doRequest(t *testing.T, r http.Handler, method, path, token string, body io.Reader, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// formFile is a single file part to attach to a multipart request.
type formFile struct {
	filename string
	content  []byte
}

// buildMultipart builds a multipart/form-data request body.
// fields is a map of form field name to its repeated values;
// files is a map of form field name to the files to attach.
func buildMultipart(t *testing.T, fields map[string][]string, files map[string][]formFile) (io.Reader, string) {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	for name, values := range fields {
		for _, v := range values {
			if err := mw.WriteField(name, v); err != nil {
				t.Fatalf("write field %s: %v", name, err)
			}
		}
	}
	for field, ffs := range files {
		for _, ff := range ffs {
			fh, err := mw.CreateFormFile(field, ff.filename)
			if err != nil {
				t.Fatalf("create form file %s: %v", field, err)
			}
			if _, err := fh.Write(ff.content); err != nil {
				t.Fatalf("write form file %s: %v", field, err)
			}
		}
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	return &buf, mw.FormDataContentType()
}

// testFile returns the content of a fixture. It prefers the real files under
// <repo>/test_data (as in the actual user upload flow) and falls back to a
// minimal valid PNG so the tests still run on machines without fixtures.
func testFile(t *testing.T, name string) []byte {
	t.Helper()
	if data, ok := readFromTestData(name); ok {
		return data
	}
	return minimalPNG
}

// readFromTestData locates <repo>/test_data/<name> by walking up from the
// package working directory.
func readFromTestData(name string) ([]byte, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, false
	}
	for {
		p := filepath.Join(dir, "test_data", name)
		if data, err := os.ReadFile(p); err == nil {
			return data, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, false
		}
		dir = parent
	}
}

var minimalPNG = func() []byte {
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==")
	if err != nil {
		panic(err)
	}
	return data
}()

// countRows returns the number of rows in a table matching an optional filter.
func countRows(t *testing.T, db *sql.DB, query string, args ...interface{}) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("count query %q: %v", query, err)
	}
	return n
}

// mustStatus asserts the HTTP status code of a recorded response.
func mustStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Fatalf("status: got %d, want %d, body %s", w.Code, want, w.Body.String())
	}
}

// mustEqual asserts that got equals want.
func mustEqual(t *testing.T, name string, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: got %#v, want %#v", name, got, want)
	}
}

func int64Str(n int64) string {
	return strconv.FormatInt(n, 10)
}
