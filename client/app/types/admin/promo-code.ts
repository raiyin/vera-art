// ─── Admin Promo Code Types ─────────────────────────────────────────

export interface AdminPromoCodeItem {
    id: number
    code: string
    discount_type: string
    discount_value: number
    max_uses: number | null
    used_count: number
    valid_from: string | null
    valid_until: string | null
    is_active: boolean
    created_at: string
}
