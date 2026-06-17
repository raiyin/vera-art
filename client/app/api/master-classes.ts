import { getHttpClient, } from '~/api/http-client';
import type { MasterClass, } from '~/types/master-class';

/**
 * Возвращает полный URL для доступа к видео мастер-класса.
 * Для бесплатных мастер-классов видео доступно без авторизации.
 * Для платных — требуется токен авторизации (будет добавлен автоматически).
 */
export function getMasterClassVideoUrl(masterClassId: number,): string {
    const config = useRuntimeConfig();
    return `${config.public.serverUrl}master-classes/${masterClassId}/video`;
}

/**
 * Возвращает полный URL для доступа к обложке мастер-класса.
 */
export function getMasterClassThumbnailUrl(masterClassId: number,): string {
    const config = useRuntimeConfig();
    return `${config.public.serverUrl}master-classes/${masterClassId}/thumbnail`;
}

/**
 * Fetches all published master-classes with optional filters.
 */
export async function fetchMasterClasses(params?: {
    difficulty?: string
    featured?: string
    search?: string
    limit?: number
    offset?: number
},): Promise<MasterClass[]> {
    try {
        const { data, } = await getHttpClient().get<MasterClass[]>('master-classes', { params, },);
        return data;
    } catch (e) {
        console.error('fetchMasterClasses error:', e,);
        return [];
    }
}

/**
 * Fetches a single master-class by ID.
 */
export async function fetchMasterClassById(id: number,): Promise<MasterClass | null> {
    try {
        const { data, } = await getHttpClient().get<MasterClass>(`master-classes/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchMasterClassById error:', e,);
        return null;
    }
}

/**
 * Fetches all tags used in master-classes.
 */
export async function fetchMasterClassTags(): Promise<{ id: number, name_ru: string, name_en: string, slug: string }[]> {
    try {
        const { data, } = await getHttpClient().get<{ id: number, name_ru: string, name_en: string, slug: string }[]>('master-classes/tags',);
        return data;
    } catch (e) {
        console.error('fetchMasterClassTags error:', e,);
        return [];
    }
}

/**
 * Fetches master-classes by tag slug.
 */
export async function fetchMasterClassesByTag(tagSlug: string,): Promise<MasterClass[]> {
    try {
        const { data, } = await getHttpClient().get<MasterClass[]>(`master-classes/tag/${tagSlug}`,);
        return data;
    } catch (e) {
        console.error('fetchMasterClassesByTag error:', e,);
        return [];
    }
}

/**
 * Fetches IDs of all products (courses and master-classes) purchased by the current user.
 * Requires authentication. Returns an array of product IDs.
 */
export async function fetchMyPurchasedProductIds(): Promise<number[]> {
    try {
        const { data, } = await getHttpClient().get<Array<{ product: { id: number, type: string } }>>('learning/my-courses',);
        // data is an array of { product: { id, type, ... }, ... }
        return data
            .filter(item => item?.product?.type === 'masterclass',)
            .map(item => item.product.id,);
    } catch (e) {
        console.error('fetchMyPurchasedProductIds error:', e,);
        return [];
    }
}
