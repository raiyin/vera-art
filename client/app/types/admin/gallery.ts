// ─── Admin Work Types ───────────────────────────────────────────────

export interface AdminWorkItem {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    width: number
    height: number
    type: number
    base_ru: string
    base_en: string
    images: string[]
    materials_ru: string[]
    materials_en: string[]
    created_at: string
}

export interface AdminListWorksResponse {
    items: AdminWorkItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

// ─── Admin Sale Types ───────────────────────────────────────────────

export interface AdminSaleItem {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    width: number
    height: number
    price: number
    base_ru: string
    base_en: string
    images: string[]
    materials_ru: string[]
    materials_en: string[]
    created_at: string
}

export interface AdminListSalesResponse {
    items: AdminSaleItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

// ─── Admin News Types ───────────────────────────────────────────────

export interface AdminNewsItem {
    id: string
    title_ru: string
    title_en: string
    datetime: string
    dir: string
    img_back: string
    img_backfull: string
    text_ru: string
    text_en: string
    images: string[]
    videos: string[]
}

export interface AdminListNewsResponse {
    items: AdminNewsItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}
