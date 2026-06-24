import { getHttpClient, } from '~/api/http-client';
import type { DashboardStats, RecentActivityItem, } from '~/types/dashboard';
import type {
    AdminListWorksResponse,
    AdminWorkItem,
    AdminListSalesResponse,
    AdminSaleItem,
    AdminListNewsResponse,
    AdminListProductsResponse,
    AdminListLessonsResponse,
    AdminListUsersResponse,
    AdminUserDetail,
    AdminListReviewsResponse,
    AdminPurchaseItem,
    AdminListPurchasesResponse,
    AdminPaymentItem,
    AdminListPaymentsResponse,
    AdminPromoCodeItem,
    AdminChatThread,
    AdminChatMessage,
    AdminCategoryItem,
    AdminTagItem,
} from '~/types';

// ─── Dashboard ──────────────────────────────────────────────────────

export async function fetchDashboardStats(): Promise<DashboardStats> {
    const { data, } = await getHttpClient().get<DashboardStats>('api/admin/stats',);
    return data;
}

export async function fetchRecentActivity(): Promise<RecentActivityItem[]> {
    const { data, } = await getHttpClient().get<RecentActivityItem[]>('api/admin/recent-activity',);
    return data;
}

// ─── Gallery (Works) ────────────────────────────────────────────────

export async function fetchAdminWorks(params?: {
    page?: number
    per_page?: number
    search?: string
    type?: string
    base_id?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListWorksResponse> {
    const backendParams: Record<string, string | number | undefined> = {
        page: params?.page,
        limit: params?.per_page,
        q: params?.search,
        base_id: params?.base_id,
        sort_by: params?.sort_by,
        sort_order: params?.sort_dir,
        status: params?.type !== undefined ? String(params.type) : undefined,
    };
    const { data, } = await getHttpClient().get<any>('admin/works', { params: backendParams, },);

    const items: AdminWorkItem[] = (data.works || []).map((w: any,) => ({
        id: w.id || 0,
        str_id: String(w.id || ''),
        dir: w.image_path || '',
        name_ru: w.title || '',
        name_en: '',
        year: w.year || 0,
        width: w.width || 0,
        height: w.height || 0,
        type: 0,
        base_ru: '',
        base_en: '',
        images: w.images || [],
        materials_ru: [],
        materials_en: [],
        created_at: w.created_at || '',
        material_ids: w.material_ids || [],
        base_ids: w.base_ids || [],
    }),);

    return {
        items,
        total: data.total || items.length,
        page: params?.page || 1,
        per_page: params?.per_page || items.length,
        total_pages: Math.ceil((data.total || items.length) / (params?.per_page || 20)),
    };
}

export async function deleteAdminWorks(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/works/bulk-delete', { ids, },);
}

// ─── Shop (Sales) ───────────────────────────────────────────────────

export async function fetchAdminSales(params?: {
    page?: number
    per_page?: number
    search?: string
    base_id?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListSalesResponse> {
    const backendParams: Record<string, string | number | undefined> = {
        page: params?.page,
        limit: params?.per_page,
        q: params?.search,
        base_id: params?.base_id,
        sort_by: params?.sort_by,
        sort_order: params?.sort_dir,
    };
    const { data, } = await getHttpClient().get<any>('admin/sales', { params: backendParams, },);

    const items: AdminSaleItem[] = (data.sales || []).map((s: any,) => ({
        id: s.id || 0,
        str_id: String(s.id || ''),
        dir: s.image_path || '',
        name_ru: s.title || '',
        name_en: '',
        year: s.year || 0,
        width: s.width || 0,
        height: s.height || 0,
        price: s.price || 0,
        base_ru: '',
        base_en: '',
        images: s.images || [],
        materials_ru: [],
        materials_en: [],
        created_at: s.created_at || '',
        material_ids: s.material_ids || [],
        base_ids: s.base_ids || [],
    }),);

    return {
        items,
        total: data.total || items.length,
        page: params?.page || 1,
        per_page: params?.per_page || items.length,
        total_pages: Math.ceil((data.total || items.length) / (params?.per_page || 20)),
    };
}

export async function deleteAdminSales(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/sales/bulk-delete', { ids, },);
}

// ─── News ───────────────────────────────────────────────────────────

export async function fetchAdminNews(params?: {
    page?: number
    per_page?: number
    search?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListNewsResponse> {
    const { data, } = await getHttpClient().get<AdminListNewsResponse>('admin/news/list', { params, },);
    return data;
}

export async function deleteAdminNews(ids: string[],): Promise<void> {
    await getHttpClient().post('admin/news/bulk-delete', { ids, },);
}

// ─── Products (Courses & Master-Classes) ────────────────────────────

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
},): Promise<AdminListProductsResponse> {
    const { data, } = await getHttpClient().get<AdminListProductsResponse>('admin/products/list', { params, },);
    return data;
}

export async function deleteAdminProducts(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/products/bulk-delete', { ids, },);
}

export async function updateAdminProductStatus(id: number, status: string,): Promise<void> {
    await getHttpClient().post(`admin/products/${id}/status`, { status, },);
}

// ─── Lessons ────────────────────────────────────────────────────────

export async function fetchAdminLessons(params?: {
    page?: number
    per_page?: number
    search?: string
    product_id?: number
    content_type?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListLessonsResponse> {
    const { data, } = await getHttpClient().get<AdminListLessonsResponse>('admin/lessons/list', { params, },);
    return data;
}

export async function deleteAdminLessons(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/lessons/bulk-delete', { ids, },);
}

// ─── Users ──────────────────────────────────────────────────────────

export async function fetchAdminUsers(params?: {
    page?: number
    per_page?: number
    search?: string
    role?: string
    blocked?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListUsersResponse> {
    const { data, } = await getHttpClient().get<AdminListUsersResponse>('admin/users/list', { params, },);
    return data;
}

export async function fetchAdminUserDetail(id: number,): Promise<AdminUserDetail> {
    const { data, } = await getHttpClient().get<AdminUserDetail>(`admin/users/${id}`,);
    return data;
}

export async function updateAdminUserRole(id: number, role: string,): Promise<void> {
    await getHttpClient().post(`admin/users/${id}/role`, { role, },);
}

export async function toggleAdminUserBlock(id: number, blocked: boolean,): Promise<void> {
    await getHttpClient().post(`admin/users/${id}/block`, { blocked, },);
}

// ─── Reviews ────────────────────────────────────────────────────────

export async function fetchAdminReviews(params?: {
    page?: number
    per_page?: number
    search?: string
    status?: string
    rating?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListReviewsResponse> {
    const { data, } = await getHttpClient().get<AdminListReviewsResponse>('admin/reviews/list', { params, },);
    return data;
}

export async function approveAdminReview(id: number,): Promise<void> {
    await getHttpClient().post(`admin/reviews/${id}/approve`, {},);
}

export async function rejectAdminReview(id: number,): Promise<void> {
    await getHttpClient().post(`admin/reviews/${id}/reject`, {},);
}

export async function deleteAdminReview(id: number,): Promise<void> {
    await getHttpClient().delete(`api/reviews/${id}`,);
}

export async function bulkApproveAdminReviews(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/reviews/bulk-approve', { ids, },);
}

export async function bulkRejectAdminReviews(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/reviews/bulk-reject', { ids, },);
}

export async function bulkDeleteAdminReviews(ids: number[],): Promise<void> {
    await getHttpClient().post('admin/reviews/bulk-delete', { ids, },);
}

// ─── Purchases ──────────────────────────────────────────────────────

export async function fetchAdminPurchases(params?: {
    page?: number
    per_page?: number
    search?: string
    status?: string
    product_type?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListPurchasesResponse> {
    const { data, } = await getHttpClient().get<AdminListPurchasesResponse>('admin/purchases/list', { params, },);
    return data;
}

export async function fetchAdminPurchaseDetail(id: number,): Promise<AdminPurchaseItem> {
    const { data, } = await getHttpClient().get<AdminPurchaseItem>(`admin/purchases/${id}`,);
    return data;
}

export async function extendAdminPurchaseAccess(id: number, days: number,): Promise<void> {
    await getHttpClient().post(`admin/purchases/${id}/extend`, { days, },);
}

export async function cancelAdminPurchase(id: number,): Promise<void> {
    await getHttpClient().put(`admin/purchases/${id}/cancel`,);
}

// ─── Payments ───────────────────────────────────────────────────────

export async function fetchAdminPayments(params?: {
    page?: number
    per_page?: number
    search?: string
    status?: string
    sort_by?: string
    sort_dir?: string
},): Promise<AdminListPaymentsResponse> {
    const { data, } = await getHttpClient().get<AdminListPaymentsResponse>('admin/payments/list', { params, },);
    return data;
}

export async function fetchAdminPaymentDetail(id: number,): Promise<AdminPaymentItem> {
    const { data, } = await getHttpClient().get<AdminPaymentItem>(`admin/payments/${id}`,);
    return data;
}

export async function refundAdminPayment(id: number,): Promise<void> {
    await getHttpClient().post(`admin/payments/${id}/refund`, {},);
}

// ─── Promo Codes ────────────────────────────────────────────────────

export async function fetchAdminPromoCodes(): Promise<AdminPromoCodeItem[]> {
    const { data, } = await getHttpClient().get<AdminPromoCodeItem[]>('admin/promo-codes',);
    return data;
}

export async function createAdminPromoCode(data: {
    code: string
    discount_type: string
    discount_value: number
    max_uses?: number | null
    valid_from?: string | null
    valid_until?: string | null
    is_active: boolean
},): Promise<void> {
    await getHttpClient().post('admin/promo-codes', data,);
}

export async function updateAdminPromoCode(id: number, data: {
    code: string
    discount_type: string
    discount_value: number
    max_uses?: number | null
    valid_from?: string | null
    valid_until?: string | null
    is_active: boolean
},): Promise<void> {
    await getHttpClient().put(`admin/promo-codes/${id}`, data,);
}

export async function deleteAdminPromoCode(id: number,): Promise<void> {
    await getHttpClient().delete(`admin/promo-codes/${id}`,);
}

// ─── Chat ───────────────────────────────────────────────────────────

export async function fetchAdminChatThreads(): Promise<AdminChatThread[]> {
    const { data, } = await getHttpClient().get<AdminChatThread[]>('admin/chat/threads',);
    return data;
}

export async function fetchAdminChatMessages(threadId: number,): Promise<AdminChatMessage[]> {
    const { data, } = await getHttpClient().get<AdminChatMessage[]>(`admin/chat/threads/${threadId}/messages`,);
    return data;
}

export async function sendAdminChatMessage(threadId: number, content: string,): Promise<AdminChatMessage> {
    const { data, } = await getHttpClient().post<AdminChatMessage>(`admin/chat/threads/${threadId}/messages`, { content, },);
    return data;
}

export async function resolveAdminChatThread(threadId: number,): Promise<void> {
    await getHttpClient().put(`admin/chat/threads/${threadId}/resolve`,);
}

export async function reopenAdminChatThread(threadId: number,): Promise<void> {
    await getHttpClient().put(`admin/chat/threads/${threadId}/reopen`,);
}

// ─── Categories ─────────────────────────────────────────────────────

export async function fetchAdminCategories(): Promise<AdminCategoryItem[]> {
    const { data, } = await getHttpClient().get<AdminCategoryItem[]>('admin/categories',);
    return data;
}

export async function createAdminCategory(data: {
    name_ru: string
    name_en: string
    slug: string
    description_ru?: string
    description_en?: string
    sort_order?: number
    is_active?: boolean
},): Promise<{ id: number, message: string }> {
    const { data: responseData, } = await getHttpClient().post<{ id: number, message: string }>('admin/categories', data,);
    return responseData;
}

export async function updateAdminCategory(id: number, data: {
    name_ru?: string
    name_en?: string
    slug?: string
    description_ru?: string
    description_en?: string
    sort_order?: number
    is_active?: boolean
},): Promise<{ message: string }> {
    const { data: responseData, } = await getHttpClient().put<{ message: string }>(`admin/categories/${id}`, data,);
    return responseData;
}

export async function deleteAdminCategory(id: number,): Promise<void> {
    await getHttpClient().delete(`admin/categories/${id}`,);
}

// ─── Tags ───────────────────────────────────────────────────────────

export async function fetchAdminTags(): Promise<AdminTagItem[]> {
    const { data, } = await getHttpClient().get<AdminTagItem[]>('admin/tags',);
    return data;
}

export async function createAdminTag(data: {
    name_ru: string
    name_en: string
    slug: string
},): Promise<{ id: number, message: string }> {
    const { data: responseData, } = await getHttpClient().post<{ id: number, message: string }>('admin/tags', data,);
    return responseData;
}

export async function updateAdminTag(id: number, data: {
    name_ru?: string
    name_en?: string
    slug?: string
},): Promise<{ message: string }> {
    const { data: responseData, } = await getHttpClient().put<{ message: string }>(`admin/tags/${id}`, data,);
    return responseData;
}

export async function deleteAdminTag(id: number,): Promise<void> {
    await getHttpClient().delete(`admin/tags/${id}`,);
}
