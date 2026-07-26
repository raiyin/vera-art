package port

import (
	"context"
	"io"

	"github.com/raiyin/artserver/internal/domain"
)

// FileRepository defines the interface for file storage operations.
type FileRepository interface {
	Save(ctx context.Context, path string, reader io.Reader) error
	Delete(ctx context.Context, path string) error
	GetPath(dir, filename string) string
	Exists(ctx context.Context, path string) (bool, error)
	Copy(ctx context.Context, src, dst string) error
	MkdirAll(ctx context.Context, path string) error
	RemoveDir(ctx context.Context, path string) error
	RenameDir(ctx context.Context, oldPath, newPath string) error
}

// AuthService defines the interface for authentication operations.
type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (accessToken string, refreshToken string, user *domain.User, err error)
	RefreshToken(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerification(ctx context.Context, email string) error
}

// UserService defines the interface for user profile operations.
type UserService interface {
	GetProfile(ctx context.Context, userID int64) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int64, name string) error
	UploadAvatar(ctx context.Context, userID int64, filename string, reader io.Reader) (string, error)
	DeleteAvatar(ctx context.Context, userID int64) error
	ServeAvatar(ctx context.Context, userID int64) (string, error)
	ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateUserRole(ctx context.Context, id int64, role string) error
	ToggleUserBlock(ctx context.Context, id int64) error
}

// GalleryService defines the interface for gallery (works + sales) operations.
type GalleryService interface {
	// Works
	GetWorks(ctx context.Context, filter domain.WorkFilter) ([]domain.Work, int, error)
	GetWorkByID(ctx context.Context, id int64) (*domain.Work, error)
	CreateWork(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error
	UpdateWork(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error
	DeleteWork(ctx context.Context, id int64) error
	BulkDeleteWorks(ctx context.Context, ids []int64) error

	// Sales
	GetSales(ctx context.Context, filter domain.SaleFilter) ([]domain.Sale, int, error)
	GetSaleByID(ctx context.Context, id int64) (*domain.Sale, error)
	CreateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error
	UpdateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error
	DeleteSale(ctx context.Context, id int64) error
	BulkDeleteSales(ctx context.Context, ids []int64) error
}

// NewsService defines the interface for news operations.
type NewsService interface {
	GetNews(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int, error)
	GetNewsByID(ctx context.Context, id int64) (*domain.News, error)
	CreateNews(ctx context.Context, news *domain.News, imageFile *domain.UploadedFile, videoFile *domain.UploadedFile) error
	UpdateNews(ctx context.Context, news *domain.News, imageFile *domain.UploadedFile, videoFile *domain.UploadedFile) error
	DeleteNews(ctx context.Context, id int64) error
	BulkDeleteNews(ctx context.Context, ids []int64) error
}

// ShopService defines the interface for shop (products, categories, promo codes, reviews) operations.
type ShopService interface {
	// Products
	GetProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error)
	GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error)
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id int64) error

	// Categories
	GetCategories(ctx context.Context) ([]domain.ProductCategory, error)
	CreateCategory(ctx context.Context, category *domain.ProductCategory) error
	UpdateCategory(ctx context.Context, category *domain.ProductCategory) error
	DeleteCategory(ctx context.Context, id int64) error

	// Promo Codes
	GetPromoCodes(ctx context.Context) ([]domain.PromoCode, int, error)
	CreatePromoCode(ctx context.Context, code *domain.PromoCode) error
	UpdatePromoCode(ctx context.Context, code *domain.PromoCode) error
	DeletePromoCode(ctx context.Context, id int64) error
	ValidatePromoCode(ctx context.Context, code string) (*domain.PromoCode, error)

	// Reviews
	GetReviews(ctx context.Context) ([]domain.Review, int, error)
	GetReviewsByProduct(ctx context.Context, productID int64) ([]domain.Review, error)
	CreateReview(ctx context.Context, review *domain.Review) error
	UpdateReview(ctx context.Context, review *domain.Review) error
	DeleteReview(ctx context.Context, id int64) error

	// Admin operations
	UpdateProductStatus(ctx context.Context, id int64, status string) error
	BulkDeleteProducts(ctx context.Context, ids []int64) error
	BulkApproveReviews(ctx context.Context, ids []int64) error
	BulkRejectReviews(ctx context.Context, ids []int64) error
	BulkDeleteReviews(ctx context.Context, ids []int64) error
}

// LearningService defines the interface for learning (lessons, progress) operations.
type LearningService interface {
	// Lessons
	GetLessonsByProduct(ctx context.Context, productID int64, userID int64) ([]domain.Lesson, error)
	GetLessonByID(ctx context.Context, id int64, userID int64) (*domain.Lesson, error)
	AdminGetLessons(ctx context.Context, filter domain.LessonFilter) ([]domain.Lesson, int, error)
	CreateLesson(ctx context.Context, lesson *domain.Lesson) error
	UpdateLesson(ctx context.Context, lesson *domain.Lesson) error
	DeleteLesson(ctx context.Context, id int64) error
	BulkDeleteLessons(ctx context.Context, ids []int64) error

	// Learning Progress
	GetMyCourses(ctx context.Context, userID int64) ([]int64, error)
	GetLearningProgress(ctx context.Context, userID, productID int64) ([]domain.LearningProgress, error)
	UpdateLessonProgress(ctx context.Context, userID, lessonID int64, completed bool) error
}

// PaymentService defines the interface for payment operations.
type PaymentService interface {
	CreatePayment(ctx context.Context, userID, productID int64, promoCode string) (*domain.Payment, string, error)
	GetPaymentStatus(ctx context.Context, paymentID int64) (*domain.Payment, error)
	HandleWebhook(ctx context.Context, payload []byte) error
	GetPayments(ctx context.Context) ([]domain.Payment, int, error)
	GetPurchases(ctx context.Context) ([]domain.Purchase, int, error)
	GetUserPurchases(ctx context.Context, userID int64) ([]domain.Purchase, error)
	HasUserPurchasedProduct(ctx context.Context, userID, productID int64) (bool, error)
	GetPaymentByID(ctx context.Context, id int64) (*domain.Payment, error)
	GetPurchaseByID(ctx context.Context, id int64) (*domain.Purchase, error)
	ExtendPurchaseAccess(ctx context.Context, id int64, days int) error
	CancelPurchase(ctx context.Context, id int64) error
	RefundPayment(ctx context.Context, id int64) error
}

// ChatService defines the interface for chat operations.
type ChatService interface {
	GetThreads(ctx context.Context, userID int64) ([]domain.ChatThread, error)
	GetAllThreads(ctx context.Context) ([]domain.ChatThread, error)
	CreateThread(ctx context.Context, userID int64, subject string) (*domain.ChatThread, error)
	GetMessages(ctx context.Context, threadID int64, userID int64) ([]domain.ChatMessage, error)
	AdminGetMessages(ctx context.Context, threadID int64) ([]domain.ChatMessage, error)
	SendMessage(ctx context.Context, threadID, userID int64, content string) (*domain.ChatMessage, error)
	AdminSendMessage(ctx context.Context, threadID int64, content string) (*domain.ChatMessage, error)
	MarkMessageAsRead(ctx context.Context, messageID int64) error
	PollMessages(ctx context.Context, threadID int64, lastMessageID int64, userID int64) ([]domain.ChatMessage, bool, error)
	ResolveThread(ctx context.Context, id int64) error
	ReopenThread(ctx context.Context, id int64) error
}

// AdminService defines the interface for admin operations.
type AdminService interface {
	GetDashboardStats(ctx context.Context) (*domain.AdminDashboardStats, error)
	ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error)
}

// ConsentService defines the interface for consent operations.
type ConsentService interface {
	RecordConsent(ctx context.Context, userID int64, consentType string, granted bool, ipAddress string) error
	GetUserConsents(ctx context.Context, userID int64) ([]domain.UserConsent, error)
}

// MasterClassService defines the interface for master class operations.
type MasterClassService interface {
	GetMasterClasses(ctx context.Context) ([]domain.MasterClass, int, error)
	GetMasterClassByID(ctx context.Context, id int64) (*domain.MasterClass, error)
	CreateMasterClass(ctx context.Context, mc *domain.MasterClass) error
	UpdateMasterClass(ctx context.Context, mc *domain.MasterClass) error
	DeleteMasterClass(ctx context.Context, id int64) error
}

// TagService defines the interface for tag operations.
type TagService interface {
	GetTags(ctx context.Context) ([]domain.Tag, error)
	CreateTag(ctx context.Context, tag *domain.Tag) error
	UpdateTag(ctx context.Context, tag *domain.Tag) error
	DeleteTag(ctx context.Context, id int64) error
}

// MaterialService defines the interface for material operations.
type MaterialService interface {
	GetMaterials(ctx context.Context) ([]domain.Material, error)
	CreateMaterial(ctx context.Context, material *domain.Material) error
	DeleteMaterial(ctx context.Context, id int64) error
}

// BaseService defines the interface for base operations.
type BaseService interface {
	GetBases(ctx context.Context) ([]domain.Base, error)
	CreateBase(ctx context.Context, base *domain.Base) error
	DeleteBase(ctx context.Context, id int64) error
}
