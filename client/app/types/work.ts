// ─── WorkBase — shared fields for both Work and Sale ────────────────

export interface WorkBase {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    descr: string
    base_ru: string
    base_en: string
    width: number
    height: number
    type?: number
    images: string[]
    materials_ids: number[]
    materials_en?: string[]
    materials_ru?: string[]
}

// ─── Work (gallery) ─────────────────────────────────────────────────

export interface Work extends WorkBase {
    __type: 'GetWorkDto'
}

// ─── Sale (art-store) ───────────────────────────────────────────────

export interface Sale extends WorkBase {
    price: number
    __type: 'GetSaleDto'
}

// ─── Discriminated union ────────────────────────────────────────────

export type CommonWork = Work | Sale;

// ─── Create/Update DTOs ─────────────────────────────────────────────

export interface CreateWorkDto {
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    descr: string
    base_ru: string
    base_en: string
    width: number
    height: number
    type: number
    images: string[]
    materials_ids: number[]
}

export interface UpdateWorkResponse {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    descr: string
    base_ru: string
    base_en: string
    width: number
    height: number
    type: number
    images: string[]
    materials_ids: number[]
}

export interface UpdateWorkRequest {
    str_id?: string
    dir?: string
    name_ru?: string
    name_en?: string
    year?: number
    descr?: string
    base_ru?: string
    base_en?: string
    width?: number
    height?: number
    type?: number
    images?: string[]
    materials_ids?: number[]
}

// ─── Sale DTOs ──────────────────────────────────────────────────────

export interface CreateSaleDto {
    width: number
    height: number
    year: number
    price: number
    name_ru: string
    name_en: string
    base_id: number
    materials_ids: number[]
    descr_ru: string
    descr_en: string
    images: string[]
}

export interface UpdateSaleResponse {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    base_id: number
    year: number
    descr_ru: string
    descr_en: string
    width: number
    height: number
    price: number
    images: string[]
    materials_ids: number[]
}

export interface UpdateSaleRequest {
    id: number
    name_ru: string
    name_en: string
    base_id: number
    year: number
    descr_ru: string
    descr_en: string
    width: number
    height: number
    price: number
    images: string[]
    materials_ids: number[]
}
