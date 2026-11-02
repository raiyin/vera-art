import { getHttpClient, } from '~/api/http-client';
import type { MasterClass, } from '~/types/master-class';

export async function fetchMasterClasses(params?: {
    difficulty?: string
    featured?: string
    search?: string
    limit?: number
    offset?: number
},): Promise<MasterClass[]> {
    try {
        const { data, } = await getHttpClient().get<any>('master-classes', { params, },);
        return data.master_classes ?? [];
    } catch (e) {
        console.error('fetchMasterClasses error:', e,);
        return [];
    }
}

export async function fetchMasterClassById(id: number,): Promise<MasterClass | null> {
    try {
        const { data, } = await getHttpClient().get<MasterClass>(`master-classes/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchMasterClassById error:', e,);
        return null;
    }
}

export async function fetchMasterClassTags(): Promise<{ id: number, name_ru: string, name_en: string, slug: string }[]> {
    try {
        const { data, } = await getHttpClient().get<any>('tags',);
        return data.tags ?? [];
    } catch (e) {
        console.error('fetchMasterClassTags error:', e,);
        return [];
    }
}

export async function fetchMasterClassesByTag(tagSlug: string,): Promise<MasterClass[]> {
    try {
        const { data, } = await getHttpClient().get<any>('master-classes', { params: { q: tagSlug, }, },);
        return data.master_classes ?? [];
    } catch (e) {
        console.error('fetchMasterClassesByTag error:', e,);
        return [];
    }
}

export async function fetchMyPurchasedProductIds(): Promise<number[]> {
    try {
        const { data, } = await getHttpClient().get<any>('learning/my-courses',);
        return data.course_ids ?? [];
    } catch (e) {
        console.error('fetchMyPurchasedProductIds error:', e,);
        return [];
    }
}
