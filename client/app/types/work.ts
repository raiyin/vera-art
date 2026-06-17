// ─── Public Work Types ──────────────────────────────────────────────

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

export interface TypedGetWorkDto {
    __type: 'GetWorkDto'
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
