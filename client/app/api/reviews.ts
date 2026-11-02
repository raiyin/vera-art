import { getHttpClient, } from '~/api/http-client';

// ─── Types matching server DTOs ───────────────────────────────────────

export interface ReviewItem {
    id: number
    user_id: number
    product_id: number
    rating: number
    text: string | null
    status: string
    created_at: string
    updated_at: string
}

export interface CreateReviewPayload {
    product_id: number
    rating: number
    text?: string
}

export interface UpdateReviewPayload {
    rating?: number
    text?: string
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Fetches reviews for a product.
 */
export async function fetchReviewsByProduct(productId: number,): Promise<ReviewItem[]> {
    try {
        const { data, } = await getHttpClient().get<ReviewItem[]>(`products/${productId}/reviews`,);
        return data;
    } catch (e) {
        console.error('fetchReviewsByProduct error:', e,);
        return [];
    }
}

/**
 * Creates a new review for a product.
 */
export async function createReview(payload: CreateReviewPayload,): Promise<ReviewItem | null> {
    try {
        const { data, } = await getHttpClient().post<ReviewItem>(
            `products/${payload.product_id}/reviews`,
            payload,
        );
        return data;
    } catch (e) {
        console.error('createReview error:', e,);
        return null;
    }
}

/**
 * Updates an existing review.
 */
export async function updateReview(id: number, payload: UpdateReviewPayload,): Promise<ReviewItem | null> {
    try {
        const { data, } = await getHttpClient().put<ReviewItem>(`reviews/${id}`, payload,);
        return data;
    } catch (e) {
        console.error('updateReview error:', e,);
        return null;
    }
}

/**
 * Deletes a review.
 */
export async function deleteReview(id: number,): Promise<boolean> {
    try {
        await getHttpClient().delete(`reviews/${id}`,);
        return true;
    } catch (e) {
        console.error('deleteReview error:', e,);
        return false;
    }
}
