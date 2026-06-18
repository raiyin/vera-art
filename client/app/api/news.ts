import { getHttpClient, } from '~/api/http-client';

/**
 * News item as returned by the server API.
 */
export interface NewsItem {
    id: number
    title: string
    description: string | null
    content: string | null
    image_path: string
    video_path: string | null
    video_paths: string[] | null
    image_paths: string[] | null
    status: string
    created_at: string
    updated_at: string
}

/**
 * Response shape for GET /news (list).
 */
export interface NewsListResponse {
    news: NewsItem[]
    total: number
}

/**
 * Fetches a paginated list of news entries.
 */
export async function fetchNews(params?: {
    status?: string
    page?: number
    limit?: number
},): Promise<NewsListResponse> {
    try {
        const { data, } = await getHttpClient().get<NewsListResponse>('news', { params, },);
        return data;
    } catch (e) {
        console.error('fetchNews error:', e,);
        return { news: [], total: 0, };
    }
}

/**
 * Fetches a single news entry by ID.
 */
export async function fetchNewsById(id: number,): Promise<NewsItem | null> {
    try {
        const { data, } = await getHttpClient().get<NewsItem>(`news/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchNewsById error:', e,);
        return null;
    }
}
