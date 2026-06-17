import { getHttpClient, } from '~/api/http-client';

// ─── Types matching server DTOs ───────────────────────────────────────

export interface WorkItem {
    id: number
    title: string
    description: string | null
    image_path: string
    year: number | null
    technique: string | null
    size: string | null
    status: string
    sort_order: number
    material_ids: number[] | null
    base_ids: number[] | null
    created_at: string
    updated_at: string
}

export interface SaleItem {
    id: number
    title: string
    description: string | null
    image_path: string
    price: number
    old_price: number | null
    year: number | null
    technique: string | null
    size: string | null
    status: string
    sort_order: number
    sold: boolean
    material_ids: number[] | null
    base_ids: number[] | null
    created_at: string
    updated_at: string
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Fetches a list of gallery works.
 */
export async function fetchWorks(params?: {
    status?: string
    limit?: number
    offset?: number
},): Promise<WorkItem[]> {
    try {
        const { data, } = await getHttpClient().get<WorkItem[]>('works', { params, },);
        return data;
    } catch (e) {
        console.error('fetchWorks error:', e,);
        return [];
    }
}

/**
 * Fetches a single work by ID.
 */
export async function fetchWorkById(id: number,): Promise<WorkItem | null> {
    try {
        const { data, } = await getHttpClient().get<WorkItem>(`works/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchWorkById error:', e,);
        return null;
    }
}

/**
 * Fetches a list of sales (shop items).
 */
export async function fetchSales(params?: {
    status?: string
    sold?: boolean
    limit?: number
    offset?: number
},): Promise<SaleItem[]> {
    try {
        const { data, } = await getHttpClient().get<SaleItem[]>('sales', { params, },);
        return data;
    } catch (e) {
        console.error('fetchSales error:', e,);
        return [];
    }
}

/**
 * Fetches a single sale by ID.
 */
export async function fetchSaleById(id: number,): Promise<SaleItem | null> {
    try {
        const { data, } = await getHttpClient().get<SaleItem>(`sales/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchSaleById error:', e,);
        return null;
    }
}
