<template>
    <div class="admin-page">
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Отзывы</span>
        </div>
        <div class="admin-page__header">
            <h1 class="admin-page__title">Модерация отзывов</h1>
            <div class="admin-page__header-actions">
                <UButton
                    v-if="selectedIds.length > 0"
                    color="success"
                    variant="solid"
                    size="sm"
                    @click="bulkApprove"
                >
                    <UIcon name="i-heroicons-check-circle" class="w-4 h-4" />
                    Одобрить ({{ selectedIds.length }})
                </UButton>
                <UButton
                    v-if="selectedIds.length > 0"
                    color="warning"
                    variant="solid"
                    size="sm"
                    @click="bulkReject"
                >
                    <UIcon name="i-heroicons-x-circle" class="w-4 h-4" />
                    Отклонить ({{ selectedIds.length }})
                </UButton>
                <UButton
                    v-if="selectedIds.length > 0"
                    color="error"
                    variant="solid"
                    size="sm"
                    @click="bulkDelete"
                >
                    <UIcon name="i-heroicons-trash" class="w-4 h-4" />
                    Удалить ({{ selectedIds.length }})
                </UButton>
            </div>
        </div>

        <!-- Stats Cards -->
        <div class="admin-reviews__stats">
            <UCard class="admin-reviews__stat-card">
                <div class="admin-reviews__stat-value">{{ stats.total_reviews }}</div>
                <div class="admin-reviews__stat-label">Всего отзывов</div>
            </UCard>
            <UCard class="admin-reviews__stat-card admin-reviews__stat-card--pending">
                <div class="admin-reviews__stat-value">{{ stats.pending_count }}</div>
                <div class="admin-reviews__stat-label">Ожидают</div>
            </UCard>
            <UCard class="admin-reviews__stat-card admin-reviews__stat-card--approved">
                <div class="admin-reviews__stat-value">{{ stats.approved_count }}</div>
                <div class="admin-reviews__stat-label">Одобрено</div>
            </UCard>
            <UCard class="admin-reviews__stat-card admin-reviews__stat-card--rejected">
                <div class="admin-reviews__stat-value">{{ stats.rejected_count }}</div>
                <div class="admin-reviews__stat-label">Отклонено</div>
            </UCard>
            <UCard class="admin-reviews__stat-card admin-reviews__stat-card--rating">
                <div class="admin-reviews__stat-value">{{ averageRatingFormatted }}</div>
                <div class="admin-reviews__stat-label">Средний рейтинг</div>
            </UCard>
            <UCard class="admin-reviews__stat-card admin-reviews__stat-card--stars">
                <div class="admin-reviews__stat-value">{{ stats.five_star_count }}</div>
                <div class="admin-reviews__stat-label">5 звёзд</div>
            </UCard>
        </div>

        <!-- Filters -->
        <UCard class="admin-page__table-card">
            <template #header>
                <div class="admin-reviews__filters">
                    <UInput
                        v-model="searchQuery"
                        placeholder="Поиск по тексту или пользователю..."
                        leading-icon="i-heroicons-magnifying-glass-20-solid"
                        class="admin-reviews__search"
                        size="sm"
                        @update:model-value="onSearchChange"
                    />
                    <USelect
                        v-model="filterStatus"
                        :items="statusOptions"
                        class="admin-reviews__filter"
                        size="sm"
                        @update:model-value="onFilterChange"
                    />
                    <USelect
                        v-model="filterRating"
                        :items="ratingOptions"
                        class="admin-reviews__filter"
                        size="sm"
                        @update:model-value="onFilterChange"
                    />
                    <UButton
                        color="neutral"
                        variant="outline"
                        size="sm"
                        @click="loadReviews"
                    >
                        <UIcon name="i-heroicons-arrow-path" class="w-4 h-4" />
                        Обновить
                    </UButton>
                </div>
            </template>

            <!-- Loading State -->
            <div v-if="loading" class="admin-page__placeholder">
                <UIcon
                    name="i-heroicons-arrow-path"
                    class="w-8 h-8 animate-spin mb-3 text-gray-400"
                />
                <p>Загрузка отзывов...</p>
            </div>

            <!-- Error State -->
            <div v-else-if="error" class="admin-page__placeholder">
                <UIcon
                    name="i-heroicons-exclamation-triangle"
                    class="w-10 h-10 mb-3 text-red-400"
                />
                <p class="text-red-500 font-medium">{{ error }}</p>
                <UButton
                    color="neutral"
                    variant="outline"
                    size="sm"
                    class="mt-3"
                    @click="loadReviews"
                >
                    Повторить
                </UButton>
            </div>

            <!-- Empty State -->
            <div v-else-if="reviews.length === 0" class="admin-page__placeholder">
                <UIcon
                    name="i-heroicons-chat-bubble-left-right"
                    class="w-10 h-10 mb-3 text-gray-300"
                />
                <p>Отзывы не найдены</p>
                <p class="admin-page__hint">Попробуйте изменить параметры фильтрации</p>
            </div>

            <!-- Reviews List -->
            <div v-else class="admin-reviews__list">
                <div
                    v-for="review in reviews"
                    :key="review.id"
                    class="admin-reviews__item"
                    :class="{
                        'admin-reviews__item--selected': selectedIds.includes(review.id),
                    }"
                >
                    <div class="admin-reviews__item-checkbox">
                        <UCheckbox
                            :model-value="selectedIds.includes(review.id)"
                            @update:model-value="toggleSelect(review.id)"
                        />
                    </div>

                    <div class="admin-reviews__item-main">
                        <div class="admin-reviews__item-header">
                            <div class="admin-reviews__item-user">
                                <UAvatar
                                    :alt="review.username || '?'"
                                    :text="
                                        ((review.username || '?')[0] || '?').toUpperCase()
                                    "
                                    size="sm"
                                />
                                <div class="admin-reviews__user-info">
                                    <div class="admin-reviews__user-name">
                                        {{
                                            review.user_full_name ||
                                            review.username ||
                                            'Пользователь'
                                        }}
                                    </div>
                                    <div class="admin-reviews__product">
                                        <UIcon
                                            name="i-heroicons-shopping-bag"
                                            class="w-3 h-3"
                                        />
                                        {{
                                            review.product_title_ru ||
                                            review.product_title_en ||
                                            'Продукт'
                                        }}
                                        <UBadge
                                            v-if="review.product_type"
                                            size="xs"
                                            color="neutral"
                                            variant="subtle"
                                        >
                                            {{
                                                review.product_type === 'course'
                                                    ? 'Курс'
                                                    : 'Мастер-класс'
                                            }}
                                        </UBadge>
                                    </div>
                                </div>
                            </div>
                            <div class="admin-reviews__item-meta">
                                <div class="admin-reviews__stars">
                                    <UIcon
                                        v-for="star in 5"
                                        :key="star"
                                        :name="
                                            star <= review.rating
                                                ? 'i-heroicons-star-solid'
                                                : 'i-heroicons-star'
                                        "
                                        class="w-4 h-4"
                                        :class="
                                            star <= review.rating
                                                ? 'text-yellow-400'
                                                : 'text-gray-300'
                                        "
                                    />
                                </div>
                                <div class="admin-reviews__date">
                                    {{ formatDate(review.created_at) }}
                                </div>
                                <UBadge
                                    :color="(statusBadgeColor(review.status) as 'success' | 'warning' | 'error' | 'neutral')"
                                    variant="subtle"
                                    size="sm"
                                >
                                    {{ statusLabel(review.status) }}
                                </UBadge>
                            </div>
                        </div>

                        <div class="admin-reviews__item-content">
                            <div
                                v-if="review.title_ru || review.title_en"
                                class="admin-reviews__title-text"
                            >
                                {{ review.title_ru || review.title_en }}
                            </div>
                            <div class="admin-reviews__comment">
                                {{
                                    review.comment_ru ||
                                    review.comment_en ||
                                    'Без комментария'
                                }}
                            </div>
                        </div>

                        <div class="admin-reviews__item-actions">
                            <template v-if="review.status === 'pending'">
                                <UButton
                                    color="success"
                                    variant="solid"
                                    size="xs"
                                    @click="approveReview(review)"
                                >
                                    <UIcon name="i-heroicons-check" class="w-3.5 h-3.5" />
                                    Одобрить
                                </UButton>
                                <UButton
                                    color="warning"
                                    variant="solid"
                                    size="xs"
                                    @click="rejectReview(review)"
                                >
                                    <UIcon
                                        name="i-heroicons-x-mark"
                                        class="w-3.5 h-3.5"
                                    />
                                    Отклонить
                                </UButton>
                            </template>
                            <template v-else>
                                <UButton
                                    v-if="review.status === 'rejected'"
                                    color="success"
                                    variant="outline"
                                    size="xs"
                                    @click="approveReview(review)"
                                >
                                    <UIcon name="i-heroicons-check" class="w-3.5 h-3.5" />
                                    Опубликовать
                                </UButton>
                                <UButton
                                    v-if="review.status === 'approved'"
                                    color="warning"
                                    variant="outline"
                                    size="xs"
                                    @click="rejectReview(review)"
                                >
                                    <UIcon
                                        name="i-heroicons-x-mark"
                                        class="w-3.5 h-3.5"
                                    />
                                    Снять с публикации
                                </UButton>
                                <UButton
                                    color="error"
                                    variant="outline"
                                    size="xs"
                                    @click="confirmDeleteReview(review)"
                                >
                                    <UIcon name="i-heroicons-trash" class="w-3.5 h-3.5" />
                                    Удалить
                                </UButton>
                            </template>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Pagination -->
            <template #footer>
                <div class="admin-reviews__footer">
                    <div class="admin-reviews__footer-info">
                        {{ (page - 1) * perPage + 1 }}–{{
                            Math.min(page * perPage, total)
                        }}
                        из {{ total }}
                    </div>
                    <UPagination
                        v-if="totalPages > 1"
                        v-model="page"
                        :total="total"
                        :page-count="perPage"
                        :max="5"
                        size="sm"
                        @update:model-value="onPageChange"
                    />
                </div>
            </template>
        </UCard>

        <!-- Delete Confirmation Modal -->
        <AdminConfirmDialog
            :visible="showDeleteModal"
            title="Подтверждение удаления"
            message="Вы уверены, что хотите удалить этот отзыв? Это действие нельзя отменить."
            type="danger"
            confirm-text="Удалить"
            cancel-text="Отмена"
            loading-text="Удаление..."
            @confirm="deleteReview"
            @cancel="showDeleteModal = false"
            @update:visible="showDeleteModal = $event"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import {
    fetchAdminReviews,
    approveAdminReview,
    rejectAdminReview,
    deleteAdminReview,
    bulkApproveAdminReviews,
    bulkRejectAdminReviews,
    bulkDeleteAdminReviews,
} from '~/api/admin';
import type { AdminReviewItem, AdminReviewsStats } from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

// --- State ---
const reviews = ref<AdminReviewItem[]>([]);
const stats = ref<AdminReviewsStats>({
    total_reviews: 0,
    pending_count: 0,
    approved_count: 0,
    rejected_count: 0,
    average_rating: 0,
    five_star_count: 0,
});
const loading = ref(true);
const error = ref('');
const page = ref(1);
const perPage = ref(20);
const total = ref(0);
const totalPages = ref(0);
const searchQuery = ref('');
const filterStatus = ref('all');
const filterRating = ref('all');
const selectedIds = ref<number[]>([]);
const showDeleteModal = ref(false);
const deleteTarget = ref<AdminReviewItem | null>(null);

const averageRatingFormatted = computed(() => {
    return stats.value.average_rating.toFixed(1);
});

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

// --- Options ---
const statusOptions = [
    { label: 'Все статусы', value: 'all' },
    { label: 'Ожидают модерации', value: 'pending' },
    { label: 'Одобренные', value: 'approved' },
    { label: 'Отклоненные', value: 'rejected' },
];

const ratingOptions = [
    { label: 'Любой рейтинг', value: 'all' },
    { label: '⭐ 1 звезда', value: '1' },
    { label: '⭐⭐ 2 звезды', value: '2' },
    { label: '⭐⭐⭐ 3 звезды', value: '3' },
    { label: '⭐⭐⭐⭐ 4 звезды', value: '4' },
    { label: '⭐⭐⭐⭐⭐ 5 звезд', value: '5' },
];

// --- Lifecycle ---
onMounted(async () => {
    await loadReviews();
});

// --- Data Loading ---
async function loadReviews() {
    loading.value = true;
    error.value = '';
    try {
        const response = await fetchAdminReviews({
            page: page.value,
            per_page: perPage.value,
            search: searchQuery.value || undefined,
            status: filterStatus.value !== 'all' ? filterStatus.value : undefined,
            rating: filterRating.value !== 'all' ? filterRating.value : undefined,
            sort_by: 'created_at',
            sort_dir: 'desc',
        });
        reviews.value = response.items;
        total.value = response.total;
        totalPages.value = response.total_pages;
        stats.value = response.stats;
    } catch (err) {
        error.value = 'Не удалось загрузить отзывы';
        console.error(err);
    } finally {
        loading.value = false;
    }
}

// --- Search with debounce ---
function onSearchChange() {
    if (searchTimeout) clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
        page.value = 1;
        loadReviews();
    }, 400);
}

function onFilterChange() {
    page.value = 1;
    loadReviews();
}

function onPageChange(newPage: number) {
    page.value = newPage;
    loadReviews();
}

// --- Selection ---
function toggleSelect(id: number) {
    const idx = selectedIds.value.indexOf(id);
    if (idx === -1) {
        selectedIds.value.push(id);
    } else {
        selectedIds.value.splice(idx, 1);
    }
}

// --- Single Actions ---
async function approveReview(review: AdminReviewItem) {
    try {
        await approveAdminReview(review.id);
        review.is_approved = true;
        review.is_visible = true;
        review.status = 'approved';
        await loadReviews();
    } catch (err) {
        console.error(err);
    }
}

async function rejectReview(review: AdminReviewItem) {
    try {
        await rejectAdminReview(review.id);
        review.is_approved = false;
        review.is_visible = false;
        review.status = 'rejected';
        await loadReviews();
    } catch (err) {
        console.error(err);
    }
}

function confirmDeleteReview(review: AdminReviewItem) {
    deleteTarget.value = review;
    showDeleteModal.value = true;
}

async function deleteReview() {
    if (!deleteTarget.value) return;
    try {
        await deleteAdminReview(deleteTarget.value.id);
        showDeleteModal.value = false;
        deleteTarget.value = null;
        await loadReviews();
    } catch (err) {
        console.error(err);
    }
}

// --- Bulk Actions ---
async function bulkApprove() {
    if (selectedIds.value.length === 0) return;
    try {
        await bulkApproveAdminReviews([...selectedIds.value]);
        selectedIds.value = [];
        await loadReviews();
    } catch (err) {
        console.error(err);
    }
}

async function bulkReject() {
    if (selectedIds.value.length === 0) return;
    try {
        await bulkRejectAdminReviews([...selectedIds.value]);
        selectedIds.value = [];
        await loadReviews();
    } catch (err) {
        console.error(err);
    }
}

async function bulkDelete() {
    if (selectedIds.value.length === 0) return;
    try {
        await bulkDeleteAdminReviews([...selectedIds.value]);
        selectedIds.value = [];
        await loadReviews();
    } catch (err) {
        console.error(err);
    }
}

// --- Helpers ---
function statusLabel(status: string): string {
    const labels: Record<string, string> = {
        pending: 'Ожидает',
        approved: 'Одобрен',
        rejected: 'Отклонен',
    };
    return labels[status] || status;
}

function statusBadgeColor(status: string): string {
    const colors: Record<string, string> = {
        pending: 'warning',
        approved: 'success',
        rejected: 'error',
    };
    return colors[status] || 'neutral';
}

function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
}
</script>

<style>
@import url('./_shared.css');

.admin-page__header-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
}

/* Stats */
.admin-reviews__stats {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
    margin-bottom: 20px;
}

.admin-reviews__stat-card {
    text-align: center;
}

.admin-reviews__stat-value {
    font-size: 28px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.2;
}

.admin-reviews__stat-label {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    margin-top: 4px;
}

.admin-reviews__stat-card--pending .admin-reviews__stat-value {
    color: #e17055;
}

.admin-reviews__stat-card--approved .admin-reviews__stat-value {
    color: #00b894;
}

.admin-reviews__stat-card--rejected .admin-reviews__stat-value {
    color: #d63031;
}

.admin-reviews__stat-card--rating .admin-reviews__stat-value {
    color: #fdcb6e;
}

.admin-reviews__stat-card--stars .admin-reviews__stat-value {
    color: #e17055;
}

/* Filters */
.admin-reviews__filters {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    align-items: center;
}

.admin-reviews__search {
    flex: 1;
    min-width: 200px;
}

.admin-reviews__filter {
    width: 180px;
}

/* Reviews List */
.admin-reviews__list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.admin-reviews__item {
    display: flex;
    gap: 12px;
    padding: 16px;
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 8px;
    transition: border-color 0.2s, box-shadow 0.2s;
}

.admin-reviews__item:hover {
    border-color: var(--admin-primary, #6c5ce7);
    box-shadow: 0 1px 4px rgba(108, 92, 231, 0.1);
}

.admin-reviews__item--selected {
    border-color: var(--admin-primary, #6c5ce7);
    background: rgba(108, 92, 231, 0.03);
}

.admin-reviews__item-checkbox {
    padding-top: 4px;
}

.admin-reviews__item-main {
    flex: 1;
    min-width: 0;
}

.admin-reviews__item-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16px;
    margin-bottom: 12px;
}

.admin-reviews__item-user {
    display: flex;
    gap: 10px;
    align-items: center;
    min-width: 0;
}

.admin-reviews__user-info {
    min-width: 0;
}

.admin-reviews__user-name {
    font-weight: 600;
    font-size: 14px;
    color: var(--admin-text-primary, #2d3436);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.admin-reviews__product {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    display: flex;
    align-items: center;
    gap: 4px;
    margin-top: 2px;
}

.admin-reviews__item-meta {
    display: flex;
    flex-direction: column;
    gap: 6px;
    align-items: flex-end;
    flex-shrink: 0;
}

.admin-reviews__stars {
    display: flex;
    gap: 1px;
}

.admin-reviews__date {
    font-size: 11px;
    color: var(--admin-text-secondary, #636e72);
    white-space: nowrap;
}

.admin-reviews__item-content {
    margin-bottom: 12px;
}

.admin-reviews__title-text {
    font-size: 15px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin-bottom: 6px;
}

.admin-reviews__comment {
    font-size: 13px;
    line-height: 1.6;
    color: var(--admin-text-secondary, #636e72);
    word-break: break-word;
}

.admin-reviews__item-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
}

/* Footer */
.admin-reviews__footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
}

.admin-reviews__footer-info {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

@media (max-width: 640px) {
    .admin-reviews__stats {
        grid-template-columns: repeat(2, 1fr);
    }

    .admin-reviews__filters {
        flex-direction: column;
    }

    .admin-reviews__filter {
        width: 100%;
    }

    .admin-reviews__item-header {
        flex-direction: column;
    }

    .admin-reviews__item-meta {
        align-items: flex-start;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 8px;
    }

    .admin-reviews__footer {
        flex-direction: column;
        align-items: stretch;
        text-align: center;
    }
}
</style>
