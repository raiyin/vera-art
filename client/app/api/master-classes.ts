import type { MasterClass, } from '../types/master-class';

/**
 * Lazily resolves the API base URL inside a Nuxt-compatible context.
 * Must be called within a Vue setup function, plugin, or Nuxt hook.
 */
function getServerUrl(): string {
    const config = useRuntimeConfig();
    return config.public.serverUrl;
}

/**
 * Возвращает полный URL для доступа к видео мастер-класса.
 * Для бесплатных мастер-классов видео доступно без авторизации.
 * Для платных — требуется токен авторизации (будет добавлен автоматически).
 */
export function getMasterClassVideoUrl(masterClassId: number,): string {
    return `${getServerUrl()}master-classes/${masterClassId}/video`;
}

/**
 * Возвращает полный URL для доступа к обложке мастер-класса.
 */
export function getMasterClassThumbnailUrl(masterClassId: number,): string {
    return `${getServerUrl()}master-classes/${masterClassId}/thumbnail`;
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
        const query = new URLSearchParams();
        if (params?.difficulty) query.set('difficulty', params.difficulty,);
        if (params?.featured) query.set('featured', params.featured,);
        if (params?.search) query.set('search', params.search,);
        if (params?.limit) query.set('limit', String(params.limit,),);
        if (params?.offset) query.set('offset', String(params.offset,),);

        const url = `${getServerUrl()}master-classes${query.toString() ? '?' + query.toString() : ''}`;
        const response = await fetch(url,);
        if (!response.ok) throw new Error('Failed to fetch master classes',);
        return await response.json();
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
        const response = await fetch(`${getServerUrl()}master-classes/${id}`,);
        if (!response.ok) throw new Error('Master class not found',);
        return await response.json();
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
        const response = await fetch(`${getServerUrl()}master-classes/tags`,);
        if (!response.ok) throw new Error('Failed to fetch tags',);
        return await response.json();
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
        const response = await fetch(`${getServerUrl()}master-classes/tag/${tagSlug}`,);
        if (!response.ok) throw new Error('Failed to fetch master classes by tag',);
        return await response.json();
    } catch (e) {
        console.error('fetchMasterClassesByTag error:', e,);
        return [];
    }
}

/**
 * Fetches IDs of all products (courses and master-classes) purchased by the current user.
 * Requires authentication. Returns an array of product IDs.
 */
export async function fetchMyPurchasedProductIds(token: string,): Promise<number[]> {
    try {
        const response = await fetch(`${getServerUrl()}learning/my-courses`, {
            headers: { Authorization: `Bearer ${token}`, },
        },);
        if (!response.ok) return [];
        const data: Array<{ product: { id: number, type: string } }> = await response.json();
        // data is an array of { product: { id, type, ... }, ... }
        return data
            .filter(item => item?.product?.type === 'masterclass',)
            .map(item => item.product.id,);
    } catch (e) {
        console.error('fetchMyPurchasedProductIds error:', e,);
        return [];
    }
}
