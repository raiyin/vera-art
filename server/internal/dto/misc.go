package dto

import "time"

// AdminStatsResponse represents the dashboard statistics.
type AdminStatsResponse struct {
	GalleryWorksCount    int64                  `json:"gallery_works_count"`
	ShopItemsCount       int64                  `json:"shop_items_count"`
	NewsCount            int64                  `json:"news_count"`
	UsersCount           int64                  `json:"users_count"`
	CoursesCount         int64                  `json:"courses_count"`
	MasterClassesCount   int64                  `json:"master_classes_count"`
	ReviewsTotal         int64                  `json:"reviews_total"`
	ReviewsPending       int64                  `json:"reviews_pending"`
	PurchasesTotal       int64                  `json:"purchases_total"`
	RevenueTotal         int64                  `json:"revenue_total"`
	RevenueMonth         int64                  `json:"revenue_month"`
	ActiveChats          int64                  `json:"active_chats"`
	UsersRegisteredMonth int64                  `json:"users_registered_month"`
	SalesByMonth         []SalesByMonthEntry    `json:"sales_by_month"`
	PopularCategories    []PopularCategoryEntry `json:"popular_categories"`
}

// SalesByMonthEntry represents a sales by month entry.
type SalesByMonthEntry struct {
	Month   string `json:"month"`
	Count   int64  `json:"count"`
	Revenue int64  `json:"revenue"`
}

// PopularCategoryEntry represents a popular category entry.
type PopularCategoryEntry struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// RecentActivityItem represents a single activity entry.
type RecentActivityItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	Time      string `json:"time"`
	CreatedAt string `json:"created_at"`
}

// TagResponse represents a tag in API responses.
type TagResponse struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTagRequest represents a create tag request.
type CreateTagRequest struct {
	NameRu string `json:"name_ru" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
	Slug   string `json:"slug" binding:"required"`
}

// UpdateTagRequest represents an update tag request.
type UpdateTagRequest struct {
	NameRu string `json:"name_ru,omitempty"`
	NameEn string `json:"name_en,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

// MaterialResponse represents a material in API responses.
type MaterialResponse struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMaterialRequest represents a create material request.
type CreateMaterialRequest struct {
	NameRu string `json:"name_ru" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
}

// BaseResponse represents a base in API responses.
type BaseResponse struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateBaseRequest represents a create base request.
type CreateBaseRequest struct {
	NameRu string `json:"name_ru" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
}

// ConsentRequest represents a consent request.
type ConsentRequest struct {
	UserID           int    `json:"user_id,omitempty"`
	IPAddress        string `json:"ip_address,omitempty"`
	UserAgent        string `json:"user_agent,omitempty"`
	AnalyticsConsent bool   `json:"analytics_consent"`
	MarketingConsent bool   `json:"marketing_consent"`
}

// ConsentResponse represents a consent response.
type ConsentResponse struct {
	Success          bool      `json:"success"`
	Message          string    `json:"message,omitempty"`
	ConsentID        int       `json:"consent_id,omitempty"`
	ConsentGiven     bool      `json:"consent_given"`
	ConsentVersion   string    `json:"consent_version"`
	ConsentDate      time.Time `json:"consent_date,omitempty"`
	AnalyticsAllowed bool      `json:"analytics_allowed,omitempty"`
	MarketingAllowed bool      `json:"marketing_allowed,omitempty"`
}

// CookiePolicyVersionResponse represents the cookie policy version response.
type CookiePolicyVersionResponse struct {
	Version   string `json:"version"`
	PolicyURL string `json:"policy_url"`
}

// ConsentAuditLogResponse represents the consent audit log response.
type ConsentAuditLogResponse struct {
	Logs   []ConsentAuditLogEntry `json:"logs"`
	Count  int                    `json:"count"`
	Limit  string                 `json:"limit"`
	Offset string                 `json:"offset"`
}

// ConsentAuditLogEntry represents a single consent audit log entry.
type ConsentAuditLogEntry struct {
	ID             int                    `json:"id"`
	UserID         int                    `json:"user_id"`
	Username       string                 `json:"username,omitempty"`
	Email          string                 `json:"email,omitempty"`
	Action         string                 `json:"action"`
	ConsentType    string                 `json:"consent_type"`
	ConsentVersion string                 `json:"consent_version"`
	IPAddress      string                 `json:"ip_address"`
	UserAgent      string                 `json:"user_agent"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"created_at"`
}

// MasterClassResponse represents a master class in API responses.
type MasterClassResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	ImagePath   string    `json:"image_path"`
	VideoURL    string    `json:"video_url,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateMasterClassRequest represents a create master class request.
type CreateMasterClassRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" binding:"required"`
	ImagePath   string  `json:"image_path"`
	VideoURL    string  `json:"video_url,omitempty"`
	Status      string  `json:"status"`
}

// UpdateMasterClassRequest represents an update master class request.
type UpdateMasterClassRequest struct {
	Title       string   `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ImagePath   *string  `json:"image_path,omitempty"`
	VideoURL    *string  `json:"video_url,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

// UserListResponse represents a user in admin user list.
type UserListResponse struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
