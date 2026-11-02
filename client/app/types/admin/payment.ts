// ─── Admin Payment Types ────────────────────────────────────────────

export interface AdminPurchaseItem {
    id: number
    user_id: number
    product_id: number
    purchase_date: string
    access_start: string
    access_end: string
    status: string
    payment_id: number | null
    promo_code_id: number | null
    price_paid: number
    created_at: string
    username: string
    user_full_name: string
    user_email: string
    product_title_ru: string
    product_title_en: string
    product_type: string
    days_remaining: number
}

export interface AdminListPurchasesResponse {
    items: AdminPurchaseItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export interface AdminPaymentItem {
    id: number
    user_id: number
    external_id: string
    status: string
    amount: number
    currency: string
    description: string
    payment_method: string
    created_at: string
    updated_at: string
    username: string
    user_full_name: string
    user_email: string
}

export interface AdminListPaymentsResponse {
    items: AdminPaymentItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}
