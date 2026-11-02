package domain

import "time"

type Tag struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type Material struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

type Base struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

type UserConsent struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ConsentType string    `json:"consent_type"`
	Granted     bool      `json:"granted"`
	IPAddress   string    `json:"ip_address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// AdminDashboardStats represents the admin dashboard statistics.
type AdminDashboardStats struct {
	GalleryWorksCount    int                     `json:"gallery_works_count"`
	ShopItemsCount       int                     `json:"shop_items_count"`
	NewsCount            int                     `json:"news_count"`
	UsersCount           int                     `json:"users_count"`
	CoursesCount         int                     `json:"courses_count"`
	MasterClassesCount   int                     `json:"master_classes_count"`
	ReviewsTotal         int                     `json:"reviews_total"`
	ReviewsPending       int                     `json:"reviews_pending"`
	PurchasesTotal       int                     `json:"purchases_total"`
	RevenueTotal         int64                   `json:"revenue_total"`
	RevenueMonth         int64                   `json:"revenue_month"`
	ActiveChats          int                     `json:"active_chats"`
	UsersRegisteredMonth int                     `json:"users_registered_month"`
	SalesByMonth         []SalesByMonthEntry     `json:"sales_by_month"`
	PopularCategories    []PopularCategoryEntry  `json:"popular_categories"`
}

type SalesByMonthEntry struct {
	Month   string `json:"month"`
	Count   int    `json:"count"`
	Revenue int64  `json:"revenue"`
}

type PopularCategoryEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
