// ─── Dashboard Stats ────────────────────────────────────────────────

export interface DashboardStats {
    gallery_works_count: number
    shop_items_count: number
    news_count: number
    users_count: number
    courses_count: number
    master_classes_count: number
    reviews_total: number
    reviews_pending: number
    purchases_total: number
    revenue_total: number
    revenue_month: number
    active_chats: number
    users_registered_month: number
    sales_by_month: SalesByMonthEntry[]
    popular_categories: PopularCategoryEntry[]
}

export interface SalesByMonthEntry {
    month: string
    count: number
    revenue: number
}

export interface PopularCategoryEntry {
    name: string
    count: number
}

export interface RecentActivityItem {
    id: string
    type: 'gallery' | 'shop' | 'news' | 'review' | 'user' | 'purchase'
    text: string
    time: string
    created_at: string
}
