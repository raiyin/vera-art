export interface MasterClassTag {
    name_ru: string;
    name_en: string;
    slug: string;
}

export interface MasterClass {
    id: number;
    title_ru: string;
    title_en: string;
    description_ru: string;
    description_en: string;
    short_description_ru: string;
    short_description_en: string;
    price: number;
    is_free: boolean;
    thumbnail_url: string;
    video_url: string;
    duration_minutes: number;
    difficulty: string;
    category_id: number;
    tags: MasterClassTag[];
    is_featured: boolean;
    view_count: number;
}
