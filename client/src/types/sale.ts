import { GetBaseWork } from "./base_work";

export interface Sale {
    width: number;
    height: number;
    year: number;
    price: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    materials_ids: number[];
    descr: string;
    images: string[];
}

// export interface GetSaleDto {
//     id: number;
//     dir: string;
//     width: number;
//     height: number;
//     year: number;
//     price: number;
//     name_ru: string;
//     name_en: string;
//     str_id: string;
//     descr: string;
//     images: string[];
//     base_ru: string;
//     base_en: string;
//     materials_ru: string[];
//     materials_en: string[];
// }

export type GetSaleDto = GetBaseWork & {
    price: number;
};

export type TypedGetSaleDto = GetSaleDto & {
    __type: "GetSaleDto";
};

export interface AddSaleDto {
    width: number;
    height: number;
    year: number;
    price: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    materials_ids: number[];
    descr: string;
    images: string[];
}
