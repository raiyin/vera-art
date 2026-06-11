import { useAuthStore, } from '~/stores/AuthStore';

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

// --- Gallery (Works) ---
export interface AdminWorkItem {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    width: number
    height: number
    type: number
    base_ru: string
    base_en: string
    images: string[]
    materials_ru: string[]
    materials_en: string[]
    created_at: string
}

export interface AdminListWorksResponse {
    items: AdminWorkItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

// --- Shop (Sales) ---
export interface AdminSaleItem {
    id: number
    str_id: string
    dir: string
    name_ru: string
    name_en: string
    year: number
    width: number
    height: number
    price: number
    base_ru: string
    base_en: string
    images: string[]
    materials_ru: string[]
    materials_en: string[]
    created_at: string
}

export interface AdminListSalesResponse {
    items: AdminSaleItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

// --- News ---
export interface AdminNewsItem {
    id: string
    title_ru: string
    title_en: string
    datetime: string
    dir: string
    img_back: string
    img_backfull: string
    text_ru: string
    text_en: string
    images: string[]
    videos: string[]
}

export interface AdminListNewsResponse {
    items: AdminNewsItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

function getAuthHeaders(): Record<string, string> {
    const authStore = useAuthStore();

    // Ensure auth state is initialized from localStorage
    // This is needed because the store may not have been initialized yet
    // if the API is called before the middleware or app.vue's onMounted runs
    if (!authStore.accessToken && typeof window !== 'undefined') {
        authStore.initFromLocalStorage();
    }

    const token = authStore.accessToken;
    if (!token) return {};
    return { Authorization: `Bearer ${token}`, };
}

async function tryRefreshToken(): Promise<boolean> {
    try {
        const authStore = useAuthStore();
        const refreshTokenValue = authStore.refreshToken;
        if (!refreshTokenValue) return false;

        const config = useRuntimeConfig();
        const SERVER_URL = config.public.serverUrl;
        const refreshResponse = await fetch(`${SERVER_URL}refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshTokenValue }),
        });

        if (refreshResponse.ok) {
            const data = await refreshResponse.json();
            authStore.updateAccessToken(data.access_token, data.access_expires);
            return true;
        }

        // Refresh failed, clear tokens
        authStore.clearTokens();
        return false;
    } catch {
        return false;
    }
}

async function fetchApi<T>(endpoint: string): Promise<T> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    let response = await fetch(`${SERVER_URL}${endpoint}`, {
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });

    // If 401, try to refresh the token and retry
    if (response.status === 401) {
        const refreshed = await tryRefreshToken();
        if (refreshed) {
            response = await fetch(`${SERVER_URL}${endpoint}`, {
                headers: {
                    'Content-Type': 'application/json',
                    ...getAuthHeaders(),
                },
            });
        }
    }

    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
    return response.json() as Promise<T>;
}

async function postApi<T>(endpoint: string, body: unknown): Promise<T> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    let response = await fetch(`${SERVER_URL}${endpoint}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
        body: JSON.stringify(body),
    });

    // If 401, try to refresh the token and retry
    if (response.status === 401) {
        const refreshed = await tryRefreshToken();
        if (refreshed) {
            response = await fetch(`${SERVER_URL}${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    ...getAuthHeaders(),
                },
                body: JSON.stringify(body),
            });
        }
    }

    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
    return response.json() as Promise<T>;
}

export async function fetchDashboardStats(): Promise<DashboardStats> {
    return fetchApi<DashboardStats>('api/admin/stats',);
}

export async function fetchRecentActivity(): Promise<RecentActivityItem[]> {
    return fetchApi<RecentActivityItem[]>('api/admin/recent-activity',);
}

// --- Gallery API ---
export async function fetchAdminWorks(params?: {
    page?: number
    per_page?: number
    search?: string
    type?: string
    base_id?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListWorksResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.type) query.set('type', params.type);
    if (params?.base_id) query.set('base_id', params.base_id);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListWorksResponse>(`admin/works${qs ? `?${qs}` : ''}`,);
}

export async function deleteAdminWorks(ids: number[]): Promise<void> {
    await postApi<{ message: string }>('admin/works/bulk-delete', { ids, });
}

// --- Shop API ---
export async function fetchAdminSales(params?: {
    page?: number
    per_page?: number
    search?: string
    base_id?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListSalesResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.base_id) query.set('base_id', params.base_id);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListSalesResponse>(`admin/sales${qs ? `?${qs}` : ''}`,);
}

export async function deleteAdminSales(ids: number[]): Promise<void> {
    await postApi<{ message: string }>('admin/sales/bulk-delete', { ids, });
}

// --- News API ---
export async function fetchAdminNews(params?: {
    page?: number
    per_page?: number
    search?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListNewsResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListNewsResponse>(`admin/news/list${qs ? `?${qs}` : ''}`,);
}

export async function deleteAdminNews(ids: string[]): Promise<void> {
    await postApi<{ message: string }>('admin/news/bulk-delete', { ids, });
}

// --- Products (Courses & Master-Classes) ---
export interface AdminProductItem {
    id: number
    type: string // "course" or "masterclass"
    title_ru: string
    title_en: string
    price: number // in kopecks
    status: string // "draft", "published", "archived"
    difficulty: string // "beginner", "intermediate", "advanced"
    total_lessons: number
    total_duration_minutes: number
    view_count: number
    is_featured: boolean
    language: string
    thumbnail_url: string
    category_name_ru: string
    category_name_en: string
    certificate_available: boolean
    created_at: string
    updated_at: string
}

export interface AdminListProductsResponse {
    items: AdminProductItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export async function fetchAdminProducts(params?: {
    page?: number
    per_page?: number
    search?: string
    type?: string
    status?: string
    difficulty?: string
    category_id?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListProductsResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.type) query.set('type', params.type);
    if (params?.status) query.set('status', params.status);
    if (params?.difficulty) query.set('difficulty', params.difficulty);
    if (params?.category_id) query.set('category_id', params.category_id);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListProductsResponse>(`admin/products/list${qs ? `?${qs}` : ''}`,);
}

export async function deleteAdminProducts(ids: number[]): Promise<void> {
    await postApi<{ message: string }>('admin/products/bulk-delete', { ids, });
}

export async function updateAdminProductStatus(id: number, status: string): Promise<void> {
    await postApi<{ message: string }>(`admin/products/${id}/status`, { status, });
}

// --- Lessons ---
export interface AdminLessonItem {
    id: number
    product_id: number
    product_title_ru: string
    product_title_en: string
    product_type: string
    title_ru: string
    title_en: string
    content_type: string
    duration_minutes: number
    sort_order: number
    is_preview: boolean
    is_required: boolean
    resources_count: number
    created_at: string
    updated_at: string
}

export interface AdminListLessonsResponse {
    items: AdminLessonItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export async function fetchAdminLessons(params?: {
    page?: number
    per_page?: number
    search?: string
    product_id?: number
    content_type?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListLessonsResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.product_id) query.set('product_id', String(params.product_id));
    if (params?.content_type) query.set('content_type', params.content_type);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListLessonsResponse>(`admin/lessons/list${qs ? `?${qs}` : ''}`,);
}

export async function deleteAdminLessons(ids: number[]): Promise<void> {
    await postApi<{ message: string }>('admin/lessons/bulk-delete', { ids, });
}

// --- Users ---
export interface AdminUserItem {
    id: number
    username: string
    full_name: string
    email: string
    role: string
    blocked: boolean
    purchases_count: number
    reviews_count: number
    created_at: string
    updated_at: string
}

export interface AdminListUsersResponse {
    items: AdminUserItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export interface AdminUserPurchase {
    id: number
    product_id: number
    product_title_ru: string
    product_title_en: string
    product_type: string
    price_paid: number
    status: string
    purchase_date: string
}

export interface AdminUserReview {
    id: number
    product_id: number
    product_title_ru: string
    product_title_en: string
    rating: number
    comment_ru: string
    comment_en: string
    is_approved: boolean
    created_at: string
}

export interface AdminUserDetail {
    id: number
    username: string
    full_name: string
    email: string
    role: string
    blocked: boolean
    created_at: string
    updated_at: string
    purchases: AdminUserPurchase[]
    reviews: AdminUserReview[]
}

export async function fetchAdminUsers(params?: {
    page?: number
    per_page?: number
    search?: string
    role?: string
    blocked?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListUsersResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.role) query.set('role', params.role);
    if (params?.blocked) query.set('blocked', params.blocked);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListUsersResponse>(`admin/users/list${qs ? `?${qs}` : ''}`,);
}

export async function fetchAdminUserDetail(id: number): Promise<AdminUserDetail> {
    return fetchApi<AdminUserDetail>(`admin/users/${id}`,);
}

export async function updateAdminUserRole(id: number, role: string): Promise<void> {
    await postApi<{ message: string }>(`admin/users/${id}/role`, { role, });
}

export async function toggleAdminUserBlock(id: number, blocked: boolean): Promise<void> {
    await postApi<{ message: string }>(`admin/users/${id}/block`, { blocked, });
}

// --- Reviews ---
export interface AdminReviewItem {
    id: number
    user_id: number
    product_id: number
    purchase_id: number
    rating: number
    title_ru: string
    title_en: string
    comment_ru: string
    comment_en: string
    is_approved: boolean
    is_visible: boolean
    status: string
    created_at: string
    updated_at: string
    username: string
    user_full_name: string
    user_email: string
    product_title_ru: string
    product_title_en: string
    product_type: string
}

export interface AdminReviewsStats {
    total_reviews: number
    pending_count: number
    approved_count: number
    rejected_count: number
    average_rating: number
    five_star_count: number
}

export interface AdminListReviewsResponse {
    items: AdminReviewItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
    stats: AdminReviewsStats
}

export async function fetchAdminReviews(params?: {
    page?: number
    per_page?: number
    search?: string
    status?: string
    rating?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListReviewsResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.status) query.set('status', params.status);
    if (params?.rating) query.set('rating', params.rating);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListReviewsResponse>(`admin/reviews/list${qs ? `?${qs}` : ''}`,);
}

export async function approveAdminReview(id: number): Promise<void> {
    await postApi<{ message: string }>(`admin/reviews/${id}/approve`, {});
}

export async function rejectAdminReview(id: number): Promise<void> {
    await postApi<{ message: string }>(`admin/reviews/${id}/reject`, {});
}

export async function deleteAdminReview(id: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}api/reviews/${id}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
}

export async function bulkApproveAdminReviews(ids: number[]): Promise<void> {
    await postApi<{ message: string, affected: number }>('admin/reviews/bulk-approve', { ids, });
}

export async function bulkRejectAdminReviews(ids: number[]): Promise<void> {
    await postApi<{ message: string, affected: number }>('admin/reviews/bulk-reject', { ids, });
}

export async function bulkDeleteAdminReviews(ids: number[]): Promise<void> {
    await postApi<{ message: string, affected: number }>('admin/reviews/bulk-delete', { ids, });
}

// --- Purchases ---
export interface AdminPurchaseItem {
    id: number
    user_id: number
    product_id: number
    purchase_date: string
    access_start: string
    access_end: string
    status: string
    payment_id: number | null
    promo_code_id: number | null
    price_paid: number
    created_at: string
    username: string
    user_full_name: string
    user_email: string
    product_title_ru: string
    product_title_en: string
    product_type: string
    days_remaining: number
}

export interface AdminListPurchasesResponse {
    items: AdminPurchaseItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export async function fetchAdminPurchases(params?: {
    page?: number
    per_page?: number
    search?: string
    status?: string
    product_type?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListPurchasesResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.status) query.set('status', params.status);
    if (params?.product_type) query.set('product_type', params.product_type);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListPurchasesResponse>(`admin/purchases/list${qs ? `?${qs}` : ''}`,);
}

export async function fetchAdminPurchaseDetail(id: number): Promise<AdminPurchaseItem> {
    return fetchApi<AdminPurchaseItem>(`admin/purchases/${id}`,);
}

export async function extendAdminPurchaseAccess(id: number, days: number): Promise<void> {
    await postApi<{ message: string }>(`admin/purchases/${id}/extend`, { days, });
}

export async function cancelAdminPurchase(id: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/purchases/${id}/cancel`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
}

// --- Payments ---
export interface AdminPaymentItem {
    id: number
    user_id: number
    external_id: string
    status: string
    amount: number
    currency: string
    description: string
    payment_method: string
    created_at: string
    updated_at: string
    username: string
    user_full_name: string
    user_email: string
}

export interface AdminListPaymentsResponse {
    items: AdminPaymentItem[]
    total: number
    page: number
    per_page: number
    total_pages: number
}

export async function fetchAdminPayments(params?: {
    page?: number
    per_page?: number
    search?: string
    status?: string
    sort_by?: string
    sort_dir?: string
}): Promise<AdminListPaymentsResponse> {
    const query = new URLSearchParams();
    if (params?.page) query.set('page', String(params.page));
    if (params?.per_page) query.set('per_page', String(params.per_page));
    if (params?.search) query.set('search', params.search);
    if (params?.status) query.set('status', params.status);
    if (params?.sort_by) query.set('sort_by', params.sort_by);
    if (params?.sort_dir) query.set('sort_dir', params.sort_dir);
    const qs = query.toString();
    return fetchApi<AdminListPaymentsResponse>(`admin/payments/list${qs ? `?${qs}` : ''}`,);
}

export async function fetchAdminPaymentDetail(id: number): Promise<AdminPaymentItem> {
    return fetchApi<AdminPaymentItem>(`admin/payments/${id}`,);
}

export async function refundAdminPayment(id: number): Promise<void> {
    await postApi<{ message: string }>(`admin/payments/${id}/refund`, {});
}

// --- Promo Codes ---
export interface AdminPromoCodeItem {
    id: number
    code: string
    discount_type: string
    discount_value: number
    max_uses: number | null
    used_count: number
    valid_from: string | null
    valid_until: string | null
    is_active: boolean
    created_at: string
}

export async function fetchAdminPromoCodes(): Promise<AdminPromoCodeItem[]> {
    return fetchApi<AdminPromoCodeItem[]>('admin/promo-codes',);
}

export async function createAdminPromoCode(data: {
    code: string
    discount_type: string
    discount_value: number
    max_uses?: number | null
    valid_from?: string | null
    valid_until?: string | null
    is_active: boolean
}): Promise<void> {
    await postApi<{ message: string }>('admin/promo-codes', data);
}

export async function updateAdminPromoCode(id: number, data: {
    code: string
    discount_type: string
    discount_value: number
    max_uses?: number | null
    valid_from?: string | null
    valid_until?: string | null
    is_active: boolean
}): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/promo-codes/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
}

// --- Admin Chat ---
export interface AdminChatThread {
    id: number
    purchase_id: number
    user_id: number
    admin_id: number | null
    last_message_at: string
    is_resolved: boolean
    created_at: string
    user: {
        username: string
        full_name: string
        email: string
    } | null
    purchase: {
        id: number
        product_id: number
        product: {
            title_ru: string
            title_en: string
        } | null
    } | null
}

export interface AdminChatMessage {
    id: number
    thread_id: number
    sender_id: number
    message_type: string
    content: string
    attachment_url: string | null
    attachment_size: number | null
    is_read: boolean
    read_at: string | null
    created_at: string
    sender: {
        username: string
        full_name: string
    } | null
}

export async function fetchAdminChatThreads(): Promise<AdminChatThread[]> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/chat/threads`, {
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
    return response.json();
}

export async function fetchAdminChatMessages(threadId: number): Promise<AdminChatMessage[]> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/chat/threads/${threadId}/messages`, {
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
    return response.json();
}

export async function sendAdminChatMessage(threadId: number, content: string): Promise<AdminChatMessage> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/chat/threads/${threadId}/messages`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
        body: JSON.stringify({ content }),
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
    return response.json();
}

export async function resolveAdminChatThread(threadId: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/chat/threads/${threadId}/resolve`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
}

export async function reopenAdminChatThread(threadId: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/chat/threads/${threadId}/reopen`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
}

export async function deleteAdminPromoCode(id: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const authStore = useAuthStore();
    const response = await fetch(`${SERVER_URL}admin/promo-codes/${id}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
        },
    });
    if (!response.ok) {
        throw new Error(`API error: ${response.status} ${response.statusText}`,);
    }
}

// ─── Categories ────────────────────────────────────────────────

export interface AdminCategoryItem {
    id: number
    name_ru: string
    name_en: string
    slug: string
    description_ru: string
    description_en: string
    sort_order: number
    is_active: boolean
    created_at: string
    updated_at: string
}

export async function fetchAdminCategories(): Promise<AdminCategoryItem[]> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}admin/categories`, {
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders(), },
    });
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`,);
    return response.json();
}

export async function createAdminCategory(data: {
    name_ru: string
    name_en: string
    slug: string
    description_ru?: string
    description_en?: string
    sort_order?: number
    is_active?: boolean
}): Promise<{ id: number, message: string }> {
    return postApi('admin/categories', data);
}

export async function updateAdminCategory(id: number, data: {
    name_ru?: string
    name_en?: string
    slug?: string
    description_ru?: string
    description_en?: string
    sort_order?: number
    is_active?: boolean
}): Promise<{ message: string }> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}admin/categories/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders(), },
        body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`,);
    return response.json();
}

export async function deleteAdminCategory(id: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}admin/categories/${id}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders(), },
    });
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`,);
}

// ─── Tags ──────────────────────────────────────────────────────

export interface AdminTagItem {
    id: number
    name_ru: string
    name_en: string
    slug: string
    created_at: string
}

export async function fetchAdminTags(): Promise<AdminTagItem[]> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}admin/tags`, {
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders(), },
    });
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`,);
    return response.json();
}

export async function createAdminTag(data: {
    name_ru: string
    name_en: string
    slug: string
}): Promise<{ id: number, message: string }> {
    return postApi('admin/tags', data);
}

export async function updateAdminTag(id: number, data: {
    name_ru?: string
    name_en?: string
    slug?: string
}): Promise<{ message: string }> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}admin/tags/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders(), },
        body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`,);
    return response.json();
}

export async function deleteAdminTag(id: number): Promise<void> {
    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;
    const response = await fetch(`${SERVER_URL}admin/tags/${id}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json', ...getAuthHeaders(), },
    });
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`,);
}
