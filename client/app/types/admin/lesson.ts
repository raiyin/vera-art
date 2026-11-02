// ─── Admin Lesson Types ─────────────────────────────────────────────

export interface AdminLessonItem {
    id: number
    product_id: number
    product_title_ru: string
    product_title_en: string
    product_type: string
    title_ru: string
    title_en: string
    content_type: string
    duration_minutes: number
    sort_order: number
    is_preview: boolean
    is_required: boolean
    resources_count: number
    created_at: string
    updated_at: string
}

export interface AdminListLessonsResponse {
    items: AdminLessonItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}
