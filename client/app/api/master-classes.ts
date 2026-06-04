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
 * Fetches all published master-classes with optional filters.
 */
export async function fetchMasterClasses(params?: {
    difficulty?: string;
    featured?: string;
    search?: string;
    limit?: number;
    offset?: number;
}): Promise<MasterClass[]> {
    try {
        const query = new URLSearchParams();
        if (params?.difficulty) query.set('difficulty', params.difficulty);
        if (params?.featured) query.set('featured', params.featured);
        if (params?.search) query.set('search', params.search);
        if (params?.limit) query.set('limit', String(params.limit));
        if (params?.offset) query.set('offset', String(params.offset));

        const url = `${getServerUrl()}master-classes${query.toString() ? '?' + query.toString() : ''}`;
        const response = await fetch(url);
        if (!response.ok) throw new Error('Failed to fetch master classes');
        return await response.json();
    } catch (e) {
        console.error('fetchMasterClasses error:', e);
        return [];
    }
}

/**
 * Fetches a single master-class by ID.
 */
export async function fetchMasterClassById(id: number): Promise<MasterClass | null> {
    try {
        const response = await fetch(`${getServerUrl()}master-classes/${id}`);
        if (!response.ok) throw new Error('Master class not found');
        return await response.json();
    } catch (e) {
        console.error('fetchMasterClassById error:', e);
        return null;
    }
}

/**
 * Fetches all tags used in master-classes.
 */
export async function fetchMasterClassTags(): Promise<{ id: number; name_ru: string; name_en: string; slug: string; }[]> {
    try {
        const response = await fetch(`${getServerUrl()}master-classes/tags`);
        if (!response.ok) throw new Error('Failed to fetch tags');
        return await response.json();
    } catch (e) {
        console.error('fetchMasterClassTags error:', e);
        return [];
    }
}

/**
 * Fetches master-classes by tag slug.
 */
export async function fetchMasterClassesByTag(tagSlug: string): Promise<MasterClass[]> {
    try {
        const response = await fetch(`${getServerUrl()}master-classes/tag/${tagSlug}`);
        if (!response.ok) throw new Error('Failed to fetch master classes by tag');
        return await response.json();
    } catch (e) {
        console.error('fetchMasterClassesByTag error:', e);
        return [];
    }
}
