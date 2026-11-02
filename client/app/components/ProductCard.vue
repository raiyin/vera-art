<template>
    <div class="product-card group">
        <div class="product-card__image">
            <img
                v-if="product.thumbnail_url"
                :src="product.thumbnail_url"
                :alt="productTitle"
                class="product-card__thumbnail"
            >
            <div
                v-else
                class="product-card__placeholder"
            >
                <UIcon
                    :name="product.type === 'course' ? 'i-heroicons-academic-cap' : 'i-heroicons-video-camera'"
                    class="w-10 h-10 text-gray-300 dark:text-gray-600"
                />
            </div>

            <div
                v-if="product.discount_price"
                class="product-card__discount"
            >
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
                    <UIcon
                        name="i-heroicons-check-circle"
                        class="w-3 h-3"
                    />
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

            <p
                v-if="productShortDescription"
                class="product-card__description"
            >
                {{ productShortDescription }}
            </p>

            <div class="product-card__meta">
                <div class="product-card__meta-item">
                    <UIcon
                        name="i-heroicons-clock"
                        class="product-card__meta-icon"
                    />
                    <span v-if="product.duration_hours">{{ product.duration_hours }} ч.</span>
                    <span v-else-if="product.duration_days">{{ product.duration_days }} дн.</span>
                    <span v-else>—</span>
                </div>

                <div class="product-card__meta-item">
                    <UIcon
                        name="i-heroicons-book-open"
                        class="product-card__meta-icon"
                    />
                    <span>{{ product.total_lessons }} уроков</span>
                </div>

                <div class="product-card__meta-item">
                    <UIcon
                        name="i-heroicons-user-group"
                        class="product-card__meta-icon"
                    />
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
                            {{ formatPrice(product.price,) }}
                        </span>
                        <span class="product-card__price product-card__price--current">
                            {{ formatPrice(product.discount_price,) }}
                        </span>
                    </div>
                    <div v-else>
                        <span class="product-card__price">
                            {{ formatPrice(product.price,) }}
                        </span>
                    </div>
                </div>

                <button
                    class="product-card__button"
                    @click="$emit('select', product,)"
                >
                    Подробнее
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, } from 'vue';
import type { Product, } from '~/stores/ProductStore';

interface Props {
    product: Product
}

const props = defineProps<Props>();
defineEmits<{
    select: [product: Product,]
}>();

const productTitle = computed(() => {
    // Use Russian title by default, fallback to English
    return props.product.name_ru || props.product.name_en || 'Без названия';
});

const productShortDescription = computed(() => {
    return (
        props.product.short_description_ru
            || props.product.short_description_en
        || props.product.description_ru
            || props.product.description_en
    );
});

const productCategory = computed(() => {
    if (props.product.category) {
        return (
            props.product.category.name_ru
                || props.product.category.name_en
            || 'Без категории'
        );
    }
    return 'Без категории';
});

const productTypeLabel = computed(() => {
    return props.product.type === 'course' ? 'Курс' : 'Мастер-класс';
});

const productTypeIcon = computed(() => {
    return props.product.type === 'course' ? 'i-heroicons-academic-cap' : 'i-heroicons-video-camera';
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
    const discount
            = ((props.product.price - props.product.discount_price) / props.product.price)
            * 100;
    return Math.round(discount,);
});

const formatPrice = (price: number,) => {
    // Price is stored in kopecks (1 RUB = 100 kopecks)
    const rubles = price / 100;
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(rubles,);
};
</script>

<style scoped>
.product-card {
    background: var(--color-background);
    border-radius: 16px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    height: 100%;
    display: flex;
    flex-direction: column;
    border: 1px solid var(--color-border, #e5e7eb);
}

.product-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
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
    background: linear-gradient(135deg, #ff4757, #ff6b81);
    color: white;
    font-size: 13px;
    font-weight: 700;
    padding: 4px 10px;
    border-radius: 20px;
    z-index: 2;
    box-shadow: 0 2px 8px rgba(255, 71, 87, 0.3);
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
    font-size: 11px;
    font-weight: 600;
    padding: 3px 10px;
    border-radius: 20px;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    backdrop-filter: blur(4px);
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
    background: var(--color-primary, #4B9E90);
    color: white;
    border: none;
    padding: 10px 24px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
}

.product-card__button:hover {
    background: var(--color-primary-hover, #3d8a7d);
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(75, 158, 144, 0.3);
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
