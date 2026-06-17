// ─── Admin Category Types ───────────────────────────────────────────

export interface AdminCategoryItem {
    id: number
    name_ru: string
    name_en: string
    slug: string
    description_ru: string
    description_en: string
    sort_order: number
    is_active: boolean
    created_at: string
    updated_at: string
}
