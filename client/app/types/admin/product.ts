// ─── Admin Product Types ────────────────────────────────────────────

export interface AdminProductItem {
    id: number
    type: string // "course" or "masterclass"
    title_ru: string
    title_en: string
    price: number // in kopecks
    status: string // "draft", "published", "archived"
    difficulty: string // "beginner", "intermediate", "advanced"
    total_lessons: number
    total_duration_minutes: number
    view_count: number
    is_featured: boolean
    language: string
    thumbnail_url: string
    category_name_ru: string
    category_name_en: string
    certificate_available: boolean
    created_at: string
    updated_at: string
}

export interface AdminListProductsResponse {
    items: AdminProductItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}
