<template>
    <div class="review-card" :class="{ 'review-card--pending': !review.is_approved }">
        <div class="review-card__header">
            <div class="review-card__user">
                <div class="review-card__avatar">
                    {{ userInitials }}
                </div>
                <div class="review-card__user-info">
                    <div class="review-card__username">
                        {{ review.user?.full_name || review.user?.username || 'Аноним' }}
                    </div>
                    <div class="review-card__date">{{ formattedDate }}</div>
                </div>
            </div>
            <div class="review-card__rating">
                <div class="review-card__stars">
                    <span
                        v-for="star in 5"
                        :key="star"
                        class="review-card__star"
                        :class="{ 'review-card__star--active': star <= review.rating }"
                    >
                        ★
                    </span>
                </div>
                <div class="review-card__rating-value">{{ review.rating }}/5</div>
            </div>
        </div>

        <div class="review-card__body">
            <h4 v-if="review.title_ru || review.title_en" class="review-card__title">
                {{ reviewTitle }}
            </h4>
            <p class="review-card__comment">
                {{ reviewComment }}
            </p>
        </div>

        <div v-if="!review.is_approved" class="review-card__pending-badge">
            Ожидает модерации
        </div>

        <div v-if="showActions" class="review-card__actions">
            <button
                v-if="canEdit"
                class="review-card__button review-card__button--edit"
                @click="$emit('edit', review)"
            >
                Редактировать
            </button>
            <button
                v-if="canDelete"
                class="review-card__button review-card__button--delete"
                @click="$emit('delete', review)"
            >
                Удалить
            </button>
            <button
                v-if="canApprove && !review.is_approved"
                class="review-card__button review-card__button--approve"
                @click="$emit('approve', review)"
            >
                Одобрить
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import type { Review } from '../types';

interface Props {
    review: Review;
    showActions?: boolean;
    canEdit?: boolean;
    canDelete?: boolean;
    canApprove?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
    showActions: false,
    canEdit: false,
    canDelete: false,
    canApprove: false,
});

const emit = defineEmits<{
    edit: [review: Review];
    delete: [review: Review];
    approve: [review: Review];
}>();

const { locale } = useI18n();

const userInitials = computed(() => {
    const name = props.review.user?.full_name || props.review.user?.username || 'А';
    return name
        .split(' ')
        .map((part) => part[0])
        .join('')
        .toUpperCase()
        .slice(0, 2);
});

const formattedDate = computed(() => {
    const date = new Date(props.review.created_at);
    return date.toLocaleDateString(locale.value, {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    });
});

const reviewTitle = computed(() => {
    return locale.value === 'ru' ? props.review.title_ru : props.review.title_en;
});

const reviewComment = computed(() => {
    return locale.value === 'ru' ? props.review.comment_ru : props.review.comment_en;
});
</script>

<style scoped>
.review-card {
    background: white;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    border: 1px solid #eaeaea;
    position: relative;
}

.review-card--pending {
    border-left: 4px solid #ffb74d;
    background-color: #fff9e6;
}

.review-card__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
}

.review-card__user {
    display: flex;
    align-items: center;
    gap: 12px;
}

.review-card__avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 18px;
    flex-shrink: 0;
}

.review-card__user-info {
    display: flex;
    flex-direction: column;
}

.review-card__username {
    font-weight: 600;
    font-size: 16px;
    color: #333;
}

.review-card__date {
    font-size: 14px;
    color: #888;
    margin-top: 2px;
}

.review-card__rating {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
}

.review-card__stars {
    display: flex;
    gap: 2px;
}

.review-card__star {
    color: #ddd;
    font-size: 20px;
    line-height: 1;
}

.review-card__star--active {
    color: #ffc107;
}

.review-card__rating-value {
    font-size: 14px;
    color: #666;
    font-weight: 500;
}

.review-card__body {
    margin-bottom: 16px;
}

.review-card__title {
    font-size: 18px;
    font-weight: 600;
    color: #222;
    margin-bottom: 8px;
    line-height: 1.3;
}

.review-card__comment {
    font-size: 15px;
    line-height: 1.5;
    color: #444;
    white-space: pre-line;
}

.review-card__pending-badge {
    display: inline-block;
    background: #ffb74d;
    color: #5d4037;
    font-size: 12px;
    font-weight: 600;
    padding: 4px 10px;
    border-radius: 20px;
    margin-top: 12px;
}

.review-card__actions {
    display: flex;
    gap: 10px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #eee;
}

.review-card__button {
    padding: 6px 14px;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
}

.review-card__button--edit {
    background: #e3f2fd;
    color: #1976d2;
}

.review-card__button--edit:hover {
    background: #bbdefb;
}

.review-card__button--delete {
    background: #ffebee;
    color: #d32f2f;
}

.review-card__button--delete:hover {
    background: #ffcdd2;
}

.review-card__button--approve {
    background: #e8f5e9;
    color: #388e3c;
}

.review-card__button--approve:hover {
    background: #c8e6c9;
}
</style>
