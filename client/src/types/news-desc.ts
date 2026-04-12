export interface NewsDesc {
    id: string;
    datetime: string;
    title_en: string;
    title_ru: string;
    subTitle_en: string;
    subTitle_ru: string;
    dir: string;
    img_back: string;
    img_backfull: string;
    text_en: string;
    text_ru: string;
    images: string[];
    videos: string[];
}

export interface NewsDescDto {
    id: string;
    datetime: string;
    title_en: string;
    title_ru: string;
    subTitle_en: string;
    subTitle_ru: string;
    dir: string;
    img_back: string;
    img_backfull: string;
    text_en: string;
    text_ru: string;
    images: string[];
    videos: string[];
}
