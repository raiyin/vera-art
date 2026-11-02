<template>
    <div class="product-list">
        <div
            v-if="isLoading"
            class="product-list__loading"
        >
            <div class="product-list__loading-spinner" />
            <p>Загрузка курсов...</p>
        </div>

        <div
            v-else-if="error"
            class="product-list__error"
        >
            <p class="product-list__error-message">
                {{ error }}
            </p>
            <button
                class="product-list__error-retry"
                @click="$emit('retry',)"
            >
                Попробовать снова
            </button>
        </div>

        <div
            v-else-if="products.length === 0"
            class="product-list__empty"
        >
            <p class="product-list__empty-message">
                Курсы не найдены
            </p>
        </div>

        <div
            v-else
            class="product-list__grid"
        >
            <ProductCard
                v-for="product in products"
                :key="product.id"
                :product="product"
                class="product-list__item"
                @select="$emit('select', product,)"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
    import { computed, } from 'vue';
    import ProductCard from './ProductCard.vue';
    import type { Product, } from '~/stores/ProductStore';

    interface Props {
        products: Product[]
        isLoading?: boolean
        error?: string | null
        columns?: number
    }

    const props = withDefaults(defineProps<Props>(), {
        isLoading: false,
        error: null,
        columns: 3,
    });

    defineEmits<{
        select: [product: Product,]
        retry: []
    }>();

    const gridColumns = computed(() => {
        return `repeat(${props.columns}, 1fr)`;
    });
</script>

<style scoped>
.product-list {
    width: 100%;
}

.product-list__loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    gap: 20px;
    color: var(--color-on-surface, #636e72);
}

.product-list__loading-spinner {
    width: 48px;
    height: 48px;
    border: 4px solid var(--color-border, #e5e7eb);
    border-top: 4px solid var(--color-primary, #4B9E90);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}

.product-list__error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
    gap: 20px;
    background: #fef2f2;
    border-radius: 12px;
    border: 1px solid #fecaca;
}

.body_theme_dark .product-list__error {
    background: rgba(239, 68, 68, 0.1);
    border-color: rgba(239, 68, 68, 0.2);
}

.product-list__error-message {
    color: #dc2626;
    text-align: center;
    margin: 0;
    font-weight: 500;
}

.body_theme_dark .product-list__error-message {
    color: #fca5a5;
}

.product-list__error-retry {
    background: #dc2626;
    color: white;
    border: none;
    padding: 10px 24px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
}

.product-list__error-retry:hover {
    background: #b91c1c;
    transform: translateY(-1px);
}

.product-list__empty {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    background: #f8f9fa;
    border-radius: 12px;
    border: 2px dashed #d1d5db;
}

.body_theme_dark .product-list__empty {
    background: rgba(30, 30, 30, 0.5);
    border-color: #374151;
}

.product-list__empty-message {
    color: #6c757d;
    font-size: 16px;
    margin: 0;
}

.product-list__grid {
    display: grid;
    grid-template-columns: v-bind(gridColumns);
    gap: 24px;
}

.product-list__item {
    height: 100%;
}

@media (max-width: 1200px) {
    .product-list__grid {
        grid-template-columns: repeat(2, 1fr);
        gap: 20px;
    }
}

@media (max-width: 768px) {
    .product-list__grid {
        grid-template-columns: 1fr;
        gap: 16px;
    }
}
</style>
