export interface GetBaseWork {
    id: string
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: string
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
