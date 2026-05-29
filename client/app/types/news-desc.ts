export interface NewsDesc {
    id: string
    title_en: string
    title_ru: string
    img_back: string
    img_backfull: string
    images: string[]
    videos: string[]
    datetime: string
    text_en: string
    text_ru: string
    dir: string
}

export interface NewsDescDto {
    id: string
    datetime: string
    title_en: string
    title_ru: string
    dir: string
    img_back: string
    img_backfull: string
    text_en: string
    text_ru: string
    images: string[]
    videos: string[]
}
