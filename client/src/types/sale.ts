export interface Sale {
    width: number;
    height: number;
    year: number;
    price: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    materials_ids: number[];
    img_count: number;
    descr: string;
    images: string[];
    removed_indices?: number[];
}

export interface AddSaleDto {
    width: number;
    height: number;
    year: number;
    price: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    materials_ids: number[];
    img_count: number;
    descr: string;
    images: string[];
    removed_indices?: number[];
}
