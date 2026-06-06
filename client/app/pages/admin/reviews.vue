<template>
    <div class="admin-reviews">
        <div class="admin-reviews__header">
            <h1 class="admin-reviews__title">Модерация отзывов</h1>
            <p class="admin-reviews__subtitle">Управление отзывами пользователей</p>
        </div>

        <div class="admin-reviews__filters">
            <div class="admin-reviews__filter-group">
                <label class="admin-reviews__filter-label">Статус</label>
                <select v-model="filterStatus" class="admin-reviews__filter-select">
                    <option value="all">Все</option>
                    <option value="pending">Ожидают модерации</option>
                    <option value="approved">Одобренные</option>
                    <option value="rejected">Отклоненные</option>
                </select>
            </div>
            <div class="admin-reviews__filter-group">
                <label class="admin-reviews__filter-label">Рейтинг</label>
                <select v-model="filterRating" class="admin-reviews__filter-select">
                    <option value="all">Все</option>
                    <option value="1">1 звезда</option>
                    <option value="2">2 звезды</option>
                    <option value="3">3 звезды</option>
                    <option value="4">4 звезды</option>
                    <option value="5">5 звезд</option>
                </select>
            </div>
            <button class="admin-reviews__refresh-button" @click="loadReviews">
                Обновить
            </button>
        </div>

        <div v-if="loading" class="admin-reviews__loading">Загрузка...</div>
        <div v-else-if="error" class="admin-reviews__error">{{ error }}</div>
        <div v-else-if="reviews.length === 0" class="admin-reviews__empty">
            Отзывы не найдены
        </div>
        <div v-else class="admin-reviews__list">
            <div
                v-for="review in filteredReviews"
                :key="review.id"
                class="admin-reviews__item"
            >
                <div class="admin-reviews__item-header">
                    <div class="admin-reviews__item-user">
                        <div class="admin-reviews__user-avatar">
                            {{ review.user?.username?.[0] || 'П' }}
                        </div>
                        <div class="admin-reviews__user-info">
                            <div class="admin-reviews__user-name">
                                {{
                                    review.user?.full_name ||
                                    review.user?.username ||
                                    'Пользователь'
                                }}
                            </div>
                            <div class="admin-reviews__product">
                                {{ review.product?.title_ru || 'Продукт' }}
                            </div>
                        </div>
                    </div>
                    <div class="admin-reviews__item-meta">
                        <div class="admin-reviews__rating">
                            <span class="admin-reviews__stars">
                                <span
                                    v-for="star in 5"
                                    :key="star"
                                    class="admin-reviews__star"
                                    :class="{
                                        'admin-reviews__star--active':
                                            star <= review.rating,
                                    }"
                                >
                                    ★
                                </span>
                            </span>
                            <span class="admin-reviews__rating-value"
                                >{{ review.rating }}/5</span
                            >
                        </div>
                        <div class="admin-reviews__date">
                            {{ formatDate(review.created_at) }}
                        </div>
                        <div
                            class="admin-reviews__status"
                            :class="`admin-reviews__status--${review.status}`"
                        >
                            {{ statusLabel(review.status) }}
                        </div>
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
                        {{ review.comment_ru || review.comment_en || 'Без комментария' }}
                    </div>
                </div>

                <div class="admin-reviews__item-actions">
                    <template v-if="review.status === 'pending'">
                        <button
                            class="admin-reviews__button admin-reviews__button--approve"
                            @click="approveReview(review)"
                        >
                            Одобрить
                        </button>
                        <button
                            class="admin-reviews__button admin-reviews__button--reject"
                            @click="rejectReview(review)"
                        >
                            Отклонить
                        </button>
                    </template>
                    <template v-else>
                        <button
                            class="admin-reviews__button admin-reviews__button--edit"
                            @click="editReview(review)"
                        >
                            Редактировать
                        </button>
                        <button
                            class="admin-reviews__button admin-reviews__button--delete"
                            @click="deleteReview(review)"
                        >
                            Удалить
                        </button>
                        <button
                            v-if="review.status === 'approved'"
                            class="admin-reviews__button admin-reviews__button--reject"
                            @click="rejectReview(review)"
                        >
                            Снять с публикации
                        </button>
                        <button
                            v-if="review.status === 'rejected'"
                            class="admin-reviews__button admin-reviews__button--approve"
                            @click="approveReview(review)"
                        >
                            Опубликовать
                        </button>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../../stores/AuthStore';
import type { Review } from '../../types';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const authStore = useAuthStore();

const reviews = ref<Review[]>([]);
const loading = ref(true);
const error = ref('');
const filterStatus = ref('all');
const filterRating = ref('all');

function computeReviewStatus(review: Review): string {
    if (review.is_approved) {
        return 'approved';
    }
    // If not approved and not visible, consider rejected
    if (!review.is_visible) {
        return 'rejected';
    }
    // Default pending
    return 'pending';
}

onMounted(async () => {
    await loadReviews();
});

async function loadReviews() {
    loading.value = true;
    error.value = '';
    try {
        const serverUrl = 'http://localhost:8000/';
        const response = await fetch(`${serverUrl}api/reviews/admin`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const data = (await response.json()) as Review[];
        // Ensure each review has a status field
        reviews.value = data.map((review) => ({
            ...review,
            status: review.status ?? computeReviewStatus(review),
        }));
    } catch (err) {
        error.value = 'Не удалось загрузить отзывы';
        console.error(err);
    } finally {
        loading.value = false;
    }
}

const filteredReviews = computed(() => {
    return reviews.value.filter((review) => {
        if (filterStatus.value !== 'all' && review.status !== filterStatus.value) {
            return false;
        }
        if (
            filterRating.value !== 'all' &&
            review.rating !== parseInt(filterRating.value)
        ) {
            return false;
        }
        return true;
    });
});

function statusLabel(status: string | undefined): string {
    const labels: Record<string, string> = {
        pending: 'Ожидает',
        approved: 'Одобрен',
        rejected: 'Отклонен',
    };
    if (!status) return 'Неизвестно';
    return labels[status] || status;
}

function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
    });
}

async function approveReview(review: Review) {
    try {
        const serverUrl = 'http://localhost:8000/';
        const response = await fetch(`${serverUrl}api/reviews/${review.id}/approve`, {
            method: 'POST',
            headers: { Authorization: `Bearer ${authStore.token}` },
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        // Update local state
        review.status = 'approved';
    } catch (err) {
        error.value = 'Не удалось одобрить отзыв';
        console.error(err);
    }
}

async function rejectReview(review: Review) {
    try {
        const serverUrl = 'http://localhost:8000/';
        const response = await fetch(`${serverUrl}api/reviews/${review.id}/reject`, {
            method: 'POST',
            headers: { Authorization: `Bearer ${authStore.token}` },
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        // Update local state
        review.status = 'rejected';
    } catch (err) {
        error.value = 'Не удалось отклонить отзыв';
        console.error(err);
    }
}

async function deleteReview(review: Review) {
    if (!confirm('Вы уверены, что хотите удалить этот отзыв?')) {
        return;
    }

    try {
        const serverUrl = 'http://localhost:8000/';
        const response = await fetch(`${serverUrl}api/reviews/${review.id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${authStore.token}` },
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        // Remove from local list
        const index = reviews.value.findIndex((r) => r.id === review.id);
        if (index !== -1) {
            reviews.value.splice(index, 1);
        }
    } catch (err) {
        error.value = 'Не удалось удалить отзыв';
        console.error(err);
    }
}

function editReview(review: Review) {
    // Navigate to edit page or show edit modal
    console.log('Edit review:', review);
}
</script>

<style scoped>
.admin-reviews {
    max-width: 1200px;
    margin: 0 auto;
    padding: 40px 20px;
}

.admin-reviews__header {
    text-align: center;
    margin-bottom: 40px;
}

.admin-reviews__title {
    font-size: 36px;
    font-weight: 700;
    color: #222;
    margin-bottom: 8px;
}

.admin-reviews__subtitle {
    font-size: 16px;
    color: #666;
    line-height: 1.5;
}

.admin-reviews__filters {
    display: flex;
    gap: 20px;
    margin-bottom: 30px;
    padding: 20px;
    background: #f9f9f9;
    border-radius: 12px;
    align-items: flex-end;
}

.admin-reviews__filter-group {
    flex: 1;
}

.admin-reviews__filter-label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: #333;
}

.admin-reviews__filter-select {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    background: white;
    font-size: 14px;
}

.admin-reviews__refresh-button {
    padding: 10px 20px;
    background: #6a11cb;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.2s;
}

.admin-reviews__refresh-button:hover {
    background: #5a0db3;
}

.admin-reviews__loading,
.admin-reviews__error,
.admin-reviews__empty {
    text-align: center;
    padding: 40px;
    background: #f9f9f9;
    border-radius: 12px;
    color: #666;
}

.admin-reviews__error {
    color: #d32f2f;
    background: #ffebee;
}

.admin-reviews__list {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.admin-reviews__item {
    background: white;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    border: 1px solid #eaeaea;
    padding: 24px;
}

.admin-reviews__item-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;
    padding-bottom: 20px;
    border-bottom: 1px solid #eee;
}

.admin-reviews__item-user {
    display: flex;
    gap: 16px;
    align-items: center;
}

.admin-reviews__user-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    font-weight: bold;
}

.admin-reviews__user-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.admin-reviews__user-name {
    font-weight: 600;
    color: #222;
}

.admin-reviews__product {
    font-size: 14px;
    color: #666;
}

.admin-reviews__item-meta {
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: flex-end;
}

.admin-reviews__rating {
    display: flex;
    align-items: center;
    gap: 8px;
}

.admin-reviews__stars {
    display: flex;
    gap: 2px;
}

.admin-reviews__star {
    color: #ddd;
    font-size: 18px;
}

.admin-reviews__star--active {
    color: #ffc107;
}

.admin-reviews__rating-value {
    font-weight: 600;
    color: #333;
}

.admin-reviews__date {
    font-size: 12px;
    color: #999;
}

.admin-reviews__status {
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
}

.admin-reviews__status--pending {
    background: #fff3cd;
    color: #856404;
}

.admin-reviews__status--approved {
    background: #d4edda;
    color: #155724;
}

.admin-reviews__status--rejected {
    background: #f8d7da;
    color: #721c24;
}

.admin-reviews__item-content {
    margin-bottom: 20px;
}

.admin-reviews__title-text {
    font-size: 18px;
    font-weight: 600;
    color: #222;
    margin-bottom: 12px;
}

.admin-reviews__comment {
    line-height: 1.6;
    color: #444;
}

.admin-reviews__item-actions {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}

.admin-reviews__button {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
    font-size: 14px;
    transition: all 0.2s;
}

.admin-reviews__button--approve {
    background: #28a745;
    color: white;
}

.admin-reviews__button--approve:hover {
    background: #218838;
}

.admin-reviews__button--reject {
    background: #dc3545;
    color: white;
}

.admin-reviews__button--reject:hover {
    background: #c82333;
}

.admin-reviews__button--edit {
    background: #17a2b8;
    color: white;
}

.admin-reviews__button--edit:hover {
    background: #138496;
}

.admin-reviews__button--delete {
    background: #6c757d;
    color: white;
}

.admin-reviews__button--delete:hover {
    background: #5a6268;
}
</style>
