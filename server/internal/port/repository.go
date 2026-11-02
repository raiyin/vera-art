package port

import (
	"context"

	"github.com/raiyin/artserver/internal/domain"
)

// TransactionManager defines the interface for database transactions.
type TransactionManager interface {
	BeginTx(ctx context.Context) (Tx, error)
}

// Tx represents a database transaction.
type Tx interface {
	Commit() error
	Rollback() error
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByVerificationToken(ctx context.Context, token string) (*domain.User, error)
	List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}

// WorkRepository defines the interface for work (gallery) data access.
type WorkRepository interface {
	Create(ctx context.Context, work *domain.Work) error
	GetByID(ctx context.Context, id int64) (*domain.Work, error)
	List(ctx context.Context, filter domain.WorkFilter) ([]domain.Work, int, error)
	Update(ctx context.Context, work *domain.Work) error
	Delete(ctx context.Context, id int64) error
	SetMaterials(ctx context.Context, workID int64, materialIDs []int64) error
}

// SaleRepository defines the interface for sale data access.
type SaleRepository interface {
	Create(ctx context.Context, sale *domain.Sale) error
	GetByID(ctx context.Context, id int64) (*domain.Sale, error)
	List(ctx context.Context, filter domain.SaleFilter) ([]domain.Sale, int, error)
	Update(ctx context.Context, sale *domain.Sale) error
	Delete(ctx context.Context, id int64) error
	SetMaterials(ctx context.Context, saleID int64, materialIDs []int64) error
	SetBases(ctx context.Context, saleID int64, baseIDs []int64) error
	GetMaterialIDs(ctx context.Context, saleID int64) ([]int64, error)
	GetBaseIDs(ctx context.Context, saleID int64) ([]int64, error)
}

// NewsRepository defines the interface for news data access.
type NewsRepository interface {
	Create(ctx context.Context, news *domain.News) error
	GetByID(ctx context.Context, id string) (*domain.News, error)
	List(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int, error)
	Update(ctx context.Context, news *domain.News) error
	Delete(ctx context.Context, id string) error
}

// ProductRepository defines the interface for product data access.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Product, error)
	List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
}

// ProductCategoryRepository defines the interface for product category data access.
type ProductCategoryRepository interface {
	Create(ctx context.Context, category *domain.ProductCategory) error
	GetByID(ctx context.Context, id int64) (*domain.ProductCategory, error)
	List(ctx context.Context) ([]domain.ProductCategory, error)
	Update(ctx context.Context, category *domain.ProductCategory) error
	Delete(ctx context.Context, id int64) error
}

// LessonRepository defines the interface for lesson data access.
type LessonRepository interface {
	Create(ctx context.Context, lesson *domain.Lesson) error
	GetByID(ctx context.Context, id int64) (*domain.Lesson, error)
	ListByProduct(ctx context.Context, productID int64, onlyPublic bool) ([]domain.Lesson, error)
	List(ctx context.Context, filter domain.LessonFilter) ([]domain.Lesson, int, error)
	Update(ctx context.Context, lesson *domain.Lesson) error
	Delete(ctx context.Context, id int64) error
}

// LearningProgressRepository defines the interface for learning progress data access.
type LearningProgressRepository interface {
	Upsert(ctx context.Context, progress *domain.LearningProgress) error
	GetByUserAndLesson(ctx context.Context, userID, lessonID int64) (*domain.LearningProgress, error)
	ListByUserAndProduct(ctx context.Context, userID, productID int64) ([]domain.LearningProgress, error)
	ListCoursesByUser(ctx context.Context, userID int64) ([]int64, error)
}

// PaymentRepository defines the interface for payment data access.
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id int64) (*domain.Payment, error)
	GetByYooKassaID(ctx context.Context, yookassaID string) (*domain.Payment, error)
	List(ctx context.Context) ([]domain.Payment, int, error)
	Update(ctx context.Context, payment *domain.Payment) error
}

// PurchaseRepository defines the interface for purchase data access.
type PurchaseRepository interface {
	Create(ctx context.Context, purchase *domain.Purchase) error
	GetByID(ctx context.Context, id int64) (*domain.Purchase, error)
	List(ctx context.Context) ([]domain.Purchase, int, error)
	ListByUser(ctx context.Context, userID int64) ([]domain.Purchase, error)
	GetByUserAndProduct(ctx context.Context, userID, productID int64) (*domain.Purchase, error)
}

// PromoCodeRepository defines the interface for promo code data access.
type PromoCodeRepository interface {
	Create(ctx context.Context, code *domain.PromoCode) error
	GetByCode(ctx context.Context, code string) (*domain.PromoCode, error)
	List(ctx context.Context) ([]domain.PromoCode, int, error)
	Update(ctx context.Context, code *domain.PromoCode) error
	Delete(ctx context.Context, id int64) error
}

// ReviewRepository defines the interface for review data access.
type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id int64) (*domain.Review, error)
	ListByProduct(ctx context.Context, productID int64) ([]domain.Review, error)
	List(ctx context.Context) ([]domain.Review, int, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id int64) error
}

// ChatThreadRepository defines the interface for chat thread data access.
type ChatThreadRepository interface {
	Create(ctx context.Context, thread *domain.ChatThread) error
	GetByID(ctx context.Context, id int64) (*domain.ChatThread, error)
	ListByUser(ctx context.Context, userID int64) ([]domain.ChatThread, error)
	List(ctx context.Context) ([]domain.ChatThread, error)
	Update(ctx context.Context, thread *domain.ChatThread) error
}

// ChatMessageRepository defines the interface for chat message data access.
type ChatMessageRepository interface {
	Create(ctx context.Context, message *domain.ChatMessage) error
	ListByThread(ctx context.Context, threadID int64) ([]domain.ChatMessage, error)
	MarkAsRead(ctx context.Context, messageID int64) error
	GetUnreadCount(ctx context.Context, threadID int64) (int, error)
	PollNewMessages(ctx context.Context, threadID int64, lastMessageID int64) ([]domain.ChatMessage, bool, error)
}

// TagRepository defines the interface for tag data access.
type TagRepository interface {
	Create(ctx context.Context, tag *domain.Tag) error
	Update(ctx context.Context, tag *domain.Tag) error
	List(ctx context.Context) ([]domain.Tag, error)
	Delete(ctx context.Context, id int64) error
}

// MaterialRepository defines the interface for material data access.
type MaterialRepository interface {
	Create(ctx context.Context, material *domain.Material) error
	List(ctx context.Context) ([]domain.Material, error)
	Delete(ctx context.Context, id int64) error
}

// BaseRepository defines the interface for base data access.
type BaseRepository interface {
	Create(ctx context.Context, base *domain.Base) error
	List(ctx context.Context) ([]domain.Base, error)
	Delete(ctx context.Context, id int64) error
}

// ConsentRepository defines the interface for user consent data access.
type ConsentRepository interface {
	Create(ctx context.Context, consent *domain.UserConsent) error
	ListByUser(ctx context.Context, userID int64) ([]domain.UserConsent, error)
}

// MasterClassRepository defines the interface for master class data access.
type MasterClassRepository interface {
	Create(ctx context.Context, mc *domain.MasterClass) error
	GetByID(ctx context.Context, id int64) (*domain.MasterClass, error)
	List(ctx context.Context) ([]domain.MasterClass, int, error)
	Update(ctx context.Context, mc *domain.MasterClass) error
	Delete(ctx context.Context, id int64) error
}
