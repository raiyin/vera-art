<template>
    <div class="product-card">
        <div class="product-card__image">
            <img
                v-if="product.thumbnail_url"
                :src="product.thumbnail_url"
                :alt="productTitle"
                class="product-card__thumbnail"
            />
            <div v-else class="product-card__placeholder">
                <span class="product-card__placeholder-text">{{ productTypeIcon }}</span>
            </div>

            <div v-if="product.discount_price" class="product-card__discount">
                -{{ discountPercentage }}%
            </div>

            <div class="product-card__badges">
                <span class="product-card__badge product-card__badge--type">
                    {{ productTypeLabel }}
                </span>
                <span
                    v-if="product.certificate_available"
                    class="product-card__badge product-card__badge--certificate"
                >
                    <svg
                        class="product-card__badge-icon"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M12 15l8-8-8 8-4-4-4 4" />
                    </svg>
                    Сертификат
                </span>
            </div>
        </div>

        <div class="product-card__content">
            <div class="product-card__category">
                {{ productCategory }}
            </div>

            <h3 class="product-card__title">
                {{ productTitle }}
            </h3>

            <p v-if="productShortDescription" class="product-card__description">
                {{ productShortDescription }}
            </p>

            <div class="product-card__meta">
                <div class="product-card__meta-item">
                    <svg
                        class="product-card__meta-icon"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <circle cx="12" cy="12" r="10" />
                        <polyline points="12 6 12 12 16 14" />
                    </svg>
                    <span v-if="product.duration_hours"
                        >{{ product.duration_hours }} ч.</span
                    >
                    <span v-else-if="product.duration_days"
                        >{{ product.duration_days }} дн.</span
                    >
                    <span v-else>—</span>
                </div>

                <div class="product-card__meta-item">
                    <svg
                        class="product-card__meta-icon"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path
                            d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"
                        />
                    </svg>
                    <span>{{ product.total_lessons }} уроков</span>
                </div>

                <div class="product-card__meta-item">
                    <svg
                        class="product-card__meta-icon"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                        <circle cx="12" cy="7" r="4" />
                    </svg>
                    <span>{{ difficultyLabel }}</span>
                </div>
            </div>

            <div class="product-card__footer">
                <div class="product-card__pricing">
                    <div
                        v-if="product.discount_price"
                        class="product-card__price-wrapper"
                    >
                        <span class="product-card__price product-card__price--old">
                            {{ formatPrice(product.price) }}
                        </span>
                        <span class="product-card__price product-card__price--current">
                            {{ formatPrice(product.discount_price) }}
                        </span>
                    </div>
                    <div v-else>
                        <span class="product-card__price">
                            {{ formatPrice(product.price) }}
                        </span>
                    </div>
                </div>

                <button class="product-card__button" @click="$emit('select', product)">
                    Подробнее
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Product } from '~/stores/ProductStore';

interface Props {
    product: Product;
}

const props = defineProps<Props>();
defineEmits<{
    select: [product: Product];
}>();

const productTitle = computed(() => {
    // Use Russian title by default, fallback to English
    return props.product.name_ru || props.product.name_en || 'Без названия';
});

const productShortDescription = computed(() => {
    return (
        props.product.short_description_ru ||
        props.product.short_description_en ||
        props.product.description_ru ||
        props.product.description_en
    );
});

const productCategory = computed(() => {
    if (props.product.category) {
        return (
            props.product.category.name_ru ||
            props.product.category.name_en ||
            'Без категории'
        );
    }
    return 'Без категории';
});

const productTypeLabel = computed(() => {
    return props.product.type === 'course' ? 'Курс' : 'Мастер-класс';
});

const productTypeIcon = computed(() => {
    return props.product.type === 'course' ? '📚' : '🎨';
});

const difficultyLabel = computed(() => {
    const difficultyMap: Record<string, string> = {
        beginner: 'Начинающий',
        intermediate: 'Средний',
        advanced: 'Продвинутый',
    };
    return difficultyMap[props.product.difficulty || 'beginner'] || 'Начинающий';
});

const discountPercentage = computed(() => {
    if (!props.product.discount_price) return 0;
    const discount =
        ((props.product.price - props.product.discount_price) / props.product.price) *
        100;
    return Math.round(discount);
});

const formatPrice = (price: number) => {
    // Price is stored in kopecks (1 RUB = 100 kopecks)
    const rubles = price / 100;
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(rubles);
};
</script>

<style scoped>
.product-card {
    background: var(--color-background);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    height: 100%;
    display: flex;
    flex-direction: column;
}

.product-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.product-card__image {
    position: relative;
    height: 180px;
    overflow: hidden;
    background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.product-card__thumbnail {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
}

.product-card:hover .product-card__thumbnail {
    transform: scale(1.05);
}

.product-card__placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 48px;
}

.product-card__placeholder-text {
    opacity: 0.3;
}

.product-card__discount {
    position: absolute;
    top: 12px;
    right: 12px;
    background: #ff4757;
    color: white;
    font-size: 14px;
    font-weight: 600;
    padding: 4px 8px;
    border-radius: 4px;
    z-index: 2;
}

.product-card__badges {
    position: absolute;
    top: 12px;
    left: 12px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    z-index: 2;
}

.product-card__badge {
    font-size: 12px;
    font-weight: 500;
    padding: 4px 8px;
    border-radius: 4px;
    display: inline-flex;
    align-items: center;
    gap: 4px;
}

.product-card__badge--type {
    background: rgba(255, 255, 255, 0.9);
    color: #2d3436;
}

.product-card__badge--certificate {
    background: rgba(46, 204, 113, 0.9);
    color: white;
}

.product-card__badge-icon {
    width: 12px;
    height: 12px;
}

.product-card__content {
    padding: 20px;
    flex-grow: 1;
    display: flex;
    flex-direction: column;
}

.product-card__category {
    font-size: 12px;
    color: #636e72;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 8px;
}

.product-card__title {
    font-size: 18px;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 12px 0;
    line-height: 1.4;
}

.product-card__description {
    font-size: 14px;
    color: #636e72;
    line-height: 1.5;
    margin: 0 0 16px 0;
    flex-grow: 1;
}

.product-card__meta {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
}

.product-card__meta-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: #636e72;
}

.product-card__meta-icon {
    width: 14px;
    height: 14px;
    opacity: 0.7;
}

.product-card__footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: auto;
}

.product-card__pricing {
    display: flex;
    flex-direction: column;
}

.product-card__price {
    font-size: 20px;
    font-weight: 700;
    color: #2d3436;
}

.product-card__price--current {
    color: #00b894;
}

.product-card__price--old {
    font-size: 16px;
    font-weight: 500;
    color: #636e72;
    text-decoration: line-through;
}

.product-card__button {
    background: var(--color-primary);
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.product-card__button:hover {
    background: var(--color-primary-dark);
}

@media (max-width: 768px) {
    .product-card__footer {
        flex-direction: column;
        gap: 12px;
        align-items: stretch;
    }

    .product-card__button {
        width: 100%;
    }
}
</style>
