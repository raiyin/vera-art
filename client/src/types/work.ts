import { GetBaseWork } from "./base_work";

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

// export interface GetWorkDto {
//     id: string;
//     str_id: string;
//     dir: string;
//     name_ru: string;
//     name_en: string;
//     year: string;
//     descr: string;
//     base_ru: string;
//     base_en: string;
//     width: string;
//     height: string;
//     type?: number;
//     images: string[];
//     material_ru: string;
//     material_en: string;
// }

export type TypedGetWorkDto = GetBaseWork & {
    __type: "GetWorkDto";
};
