import { GetBaseWork } from "./base_work";

export interface CreateSaleDto {
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

export interface UpdateSaleResponse {
    id: number;
    str_id: string;
    dir: string;
    name_ru: string;
    name_en: string;
    base_id: number;
    year: number;
    descr: string;
    width: number;
    height: number;
    price: number;
    images: string[];
    materials_ids: number[];
}

export interface UpdateSaleRequest {
    id: number;
    name_ru: string;
    name_en: string;
    base_id: number;
    year: number;
    descr: string;
    width: number;
    height: number;
    price: number;
    images: string[];
    materials_ids: number[];
}

export type GetSaleDto = GetBaseWork & {
    price: number;
};

export type TypedGetSaleDto = GetSaleDto & {
    __type: "GetSaleDto";
};
