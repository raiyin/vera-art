interface NewsItemType {
    id: string;
    datetime: string;
    title_en: string;
    title_ru: string;
    subTitle_en: string;
    subTitle_ru: string;
    dir: string;
    img_back: string;
    img_backfull: string;
    imagescount: number;
    videoscount: number;
    text_en: string;
    text_ru: string;
}

interface ImageProps {
    id: string;
    dir: string;
    name_ru: string;
    name_en: string;
    base_ru: string;
    base_en: string;
    material_ru: string;
    material_en: string;
    width: string;
    height: string;
    year: string;
    price: string;
    img_count: number;
    str_id: string;
}

interface NewsDesc {
    id: string;
    datetime: string;
    title_en: string;
    title_ru: string;
    subTitle_en: string;
    subTitle_ru: string;
    dir: string;
    img_back: string;
    img_backfull: string;
    imagescount: number;
    videoscount: number;
    text_en: string;
    text_ru: string;
}

interface SortOption {
    value: string;
    name: string;
}

export type { NewsItemType };
export type { ImageProps };
export type { NewsDesc };
export type { SortOption };
