<template>
    <div class="product-list">
        <div v-if="isLoading" class="product-list__loading">
            <div class="product-list__loading-spinner"></div>
            <p>Загрузка курсов...</p>
        </div>

        <div v-else-if="error" class="product-list__error">
            <p class="product-list__error-message">{{ error }}</p>
            <button class="product-list__error-retry" @click="$emit('retry')">
                Попробовать снова
            </button>
        </div>

        <div v-else-if="products.length === 0" class="product-list__empty">
            <p class="product-list__empty-message">Курсы не найдены</p>
        </div>

        <div v-else class="product-list__grid">
            <ProductCard
                v-for="product in products"
                :key="product.id"
                :product="product"
                @select="$emit('select', product)"
                class="product-list__item"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import ProductCard from './ProductCard.vue';
import type { Product } from '~/stores/ProductStore';

interface Props {
    products: Product[];
    isLoading?: boolean;
    error?: string | null;
    columns?: number;
}

const props = withDefaults(defineProps<Props>(), {
    isLoading: false,
    error: null,
    columns: 3,
});

defineEmits<{
    select: [product: Product];
    retry: [];
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
}

.product-list__loading-spinner {
    width: 50px;
    height: 50px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid var(--color-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
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
    background: #ffeaea;
    border-radius: 8px;
    border: 1px solid #ffcccc;
}

.product-list__error-message {
    color: #ff4757;
    text-align: center;
    margin: 0;
}

.product-list__error-retry {
    background: #ff4757;
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.product-list__error-retry:hover {
    background: #ff3742;
}

.product-list__empty {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    background: #f8f9fa;
    border-radius: 8px;
    border: 1px dashed #dee2e6;
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
    }
}

@media (max-width: 768px) {
    .product-list__grid {
        grid-template-columns: 1fr;
    }
}
</style>
