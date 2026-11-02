import { getHttpClient, } from '~/api/http-client';

// ─── Types matching server DTOs ───────────────────────────────────────

export interface ProductItem {
    id: number
    title: string
    slug: string
    description: string | null
    full_description: string | null
    price: number
    old_price: number | null
    image_path: string
    category_id: number
    category_name: string | null
    status: string
    is_digital: boolean
    is_master_class: boolean
    sort_order: number
    tags: string[] | null
    created_at: string
    updated_at: string
}

export interface CategoryItem {
    id: number
    name: string
    slug: string
    sort_order: number
    created_at: string
}

export interface PromoCodeValidationResult {
    valid: boolean
    discount_percent?: number
    message?: string
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Fetches a list of products (courses).
 */
export async function fetchProducts(params?: {
    category_id?: number
    status?: string
    is_master_class?: boolean
    limit?: number
    offset?: number
},): Promise<ProductItem[]> {
    try {
        const { data, } = await getHttpClient().get<ProductItem[]>('products', { params, },);
        return data;
    } catch (e) {
        console.error('fetchProducts error:', e,);
        return [];
    }
}

/**
 * Fetches a single product by ID.
 */
export async function fetchProductById(id: number,): Promise<ProductItem | null> {
    try {
        const { data, } = await getHttpClient().get<ProductItem>(`products/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchProductById error:', e,);
        return null;
    }
}

/**
 * Fetches a single product by category slug and product slug.
 */
export async function fetchProductBySlug(categorySlug: string, productSlug: string,): Promise<ProductItem | null> {
    try {
        const { data, } = await getHttpClient().get<ProductItem>(
            `category/${categorySlug}/product/${productSlug}`,
        );
        return data;
    } catch (e) {
        console.error('fetchProductBySlug error:', e,);
        return null;
    }
}

/**
 * Fetches all product categories.
 */
export async function fetchCategories(): Promise<CategoryItem[]> {
    try {
        const { data, } = await getHttpClient().get<CategoryItem[]>('categories',);
        return data;
    } catch (e) {
        console.error('fetchCategories error:', e,);
        return [];
    }
}

/**
 * Validates a promo code.
 */
export async function validatePromoCode(code: string,): Promise<PromoCodeValidationResult> {
    try {
        const { data, } = await getHttpClient().post<PromoCodeValidationResult>('promo-codes/validate', { code, },);
        return data;
    } catch (e) {
        console.error('validatePromoCode error:', e,);
        return { valid: false, message: 'Validation failed', };
    }
}
