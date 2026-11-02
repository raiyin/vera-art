// ─── Admin User Types ───────────────────────────────────────────────

export interface AdminUserItem {
    id: number
    username: string
    full_name: string
    email: string
    role: string
    blocked: boolean
    purchases_count: number
    reviews_count: number
    created_at: string
    updated_at: string
}

export interface AdminListUsersResponse {
    items: AdminUserItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export interface AdminUserPurchase {
    id: number
    product_id: number
    product_title_ru: string
    product_title_en: string
    product_type: string
    price_paid: number
    status: string
    purchase_date: string
}

export interface AdminUserReview {
    id: number
    product_id: number
    product_title_ru: string
    product_title_en: string
    rating: number
    comment_ru: string
    comment_en: string
    is_approved: boolean
    created_at: string
}

export interface AdminUserDetail {
    id: number
    username: string
    full_name: string
    email: string
    role: string
    blocked: boolean
    created_at: string
    updated_at: string
    purchases: AdminUserPurchase[]
    reviews: AdminUserReview[]
}
