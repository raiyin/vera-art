// ─── Admin Review Types ─────────────────────────────────────────────

export interface AdminReviewItem {
    id: number
    user_id: number
    product_id: number
    purchase_id: number
    rating: number
    title_ru: string
    title_en: string
    comment_ru: string
    comment_en: string
    is_approved: boolean
    is_visible: boolean
    status: string
    created_at: string
    updated_at: string
    username: string
    user_full_name: string
    user_email: string
    product_title_ru: string
    product_title_en: string
    product_type: string
}

export interface AdminReviewsStats {
    total_reviews: number
    pending_count: number
    approved_count: number
    rejected_count: number
    average_rating: number
    five_star_count: number
}

export interface AdminListReviewsResponse {
    items: AdminReviewItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
    stats: AdminReviewsStats
}
