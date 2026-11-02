import { getHttpClient, } from '~/api/http-client';

export interface NewsItem {
    id: string
    datetime: string
    title_ru: string
    title_en: string
    dir: string
    main_image: string
    text_ru: string
    text_en: string
    images: string[]
    videos: string[]
}

export interface NewsListResponse {
    news: NewsItem[]
    total: number
}

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

export async function fetchNewsById(id: string,): Promise<NewsItem | null> {
    try {
        const { data, } = await getHttpClient().get<NewsItem>(`news/${id}`,);
        return data;
    } catch (e) {
        console.error('fetchNewsById error:', e,);
        return null;
    }
}
