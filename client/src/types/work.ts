export interface CreateWorkDto {
    width: number;
    height: number;
    year: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    materials_ids: number[];
    descr: string;
    type: number;
    images: string[];
    removed_indices?: number[];
}
