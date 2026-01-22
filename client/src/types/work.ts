export interface CreateWorkDto {
    width: number;
    height: number;
    year: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    materials_ids: number[];
    img_count: number;
    descr: string;
    type: number;
    removed_indices?: number[];
}
