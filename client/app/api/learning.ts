import { getHttpClient, } from '~/api/http-client';

// ─── Types matching server DTOs ───────────────────────────────────────

export interface LessonItem {
    id: number
    product_id: number
    title: string
    description: string | null
    content: string | null
    video_url: string | null
    resources: string[] | null
    duration_minutes: number | null
    sort_order: number
    status: string
    created_at: string
    updated_at: string
}

export interface LearningProgress {
    id: number
    user_id: number
    lesson_id: number
    product_id: number
    completed: boolean
    created_at: string
    updated_at: string
}

export interface MyCourseItem {
    product: {
        id: number
        type: string
    }
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Fetches lessons for a product (course).
 */
export async function fetchLessonsByProduct(productId: number,): Promise<LessonItem[]> {
    try {
        const { data, } = await getHttpClient().get<LessonItem[]>(
            `product/${productId}/lessons`,
        );
        return data;
    } catch (e) {
        console.error('fetchLessonsByProduct error:', e,);
        return [];
    }
}

/**
 * Fetches a single lesson by ID.
 */
export async function fetchLessonById(id: number,): Promise<LessonItem | null> {
    try {
        const { data, } = await getHttpClient().get<LessonItem>(`lessons/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchLessonById error:', e,);
        return null;
    }
}

/**
 * Fetches all courses purchased by the current user.
 */
export async function fetchMyCourses(): Promise<MyCourseItem[]> {
    try {
        const { data, } = await getHttpClient().get<MyCourseItem[]>('learning/my-courses',);
        return data;
    } catch (e) {
        console.error('fetchMyCourses error:', e,);
        return [];
    }
}

/**
 * Fetches learning progress for a specific product.
 */
export async function fetchLearningProgress(productId: number,): Promise<LearningProgress[]> {
    try {
        const { data, } = await getHttpClient().get<LearningProgress[]>(
            `learning/progress/${productId}`,
        );
        return data;
    } catch (e) {
        console.error('fetchLearningProgress error:', e,);
        return [];
    }
}

/**
 * Updates lesson progress (mark as completed/incomplete).
 */
export async function updateLessonProgress(progressId: number, completed: boolean,): Promise<LearningProgress | null> {
    try {
        const { data, } = await getHttpClient().post<LearningProgress>(
            `learning/progress/${progressId}`,
            { completed, },
        );
        return data;
    } catch (e) {
        console.error('updateLessonProgress error:', e,);
        return null;
    }
}
