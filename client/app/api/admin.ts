import { useAuthStore, } from '~/stores/AuthStore'

export interface DashboardStats {
    gallery_works_count: number
    shop_items_count: number
    news_count: number
    users_count: number
    courses_count: number
    master_classes_count: number
    reviews_total: number
    reviews_pending: number
    purchases_total: number
    revenue_total: number
    revenue_month: number
    active_chats: number
    users_registered_month: number
    sales_by_month: SalesByMonthEntry[]
    popular_categories: PopularCategoryEntry[]
}

export interface SalesByMonthEntry {
    month: string
    count: number
    revenue: number
}

export interface PopularCategoryEntry {
    name: string
    count: number
}

export interface RecentActivityItem {
    id: string
    type: 'gallery' | 'shop' | 'news' | 'review' | 'user' | 'purchase'
    text: string
    time: string
    created_at: string
}

function getAuthHeaders(): Record<string, string> {
    const authStore = useAuthStore()
    const token = authStore.accessToken
    if (!token) return {}
    return { Authorization: `Bearer ${token}`, }
}

async function fetchApi<T>(endpoint: string): Promise<T> {
    const config = useRuntimeConfig()
    const SERVER_URL = config.public.serverUrl
    const response = await fetch(`${SERVER_URL}${endpoint}`, {
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    })
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,)
    }
    return response.json() as Promise<T>
}

export async function fetchDashboardStats(): Promise<DashboardStats> {
    return fetchApi<DashboardStats>('api/admin/stats',)
}

export async function fetchRecentActivity(): Promise<RecentActivityItem[]> {
    return fetchApi<RecentActivityItem[]>('api/admin/recent-activity',)
}
