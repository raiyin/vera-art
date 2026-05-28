import type { GetBaseWork, } from './base_work';

export interface CreateWorkDto {
    width: number
    height: number
    year: number
    name_ru: string
    name_en: string
    base_id: number
    materials_ids: number[]
    descr: string
    type: number | null
    images: string[]
}

export interface UpdateWorkResponse {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    base_id: number
    year: number
    descr: string
    width: number
    height: number
    type: number
    images: string[]
    materials_ids: number[]
}

export interface UpdateWorkRequest {
    id: number
    name_ru: string
    name_en: string
    base_id: number
    year: number
    descr: string
    width: number
    height: number
    type: number
    images: string[]
    materials_ids: number[]
}

export type TypedGetWorkDto = GetBaseWork & {
    __type: 'GetWorkDto'
};
