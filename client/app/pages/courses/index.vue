<template>
    <div class="courses-page">
        <!-- Hero Section -->
        <section class="courses-hero">
            <div class="courses-hero__container">
                <h1 class="courses-hero__title">Онлайн-курсы по живописи и рисунку</h1>
                <p class="courses-hero__subtitle">
                    Профессиональное обучение от ведущих художников. Изучайте в удобном
                    темпе с поддержкой преподавателя.
                </p>
                <div class="courses-hero__stats">
                    <div class="courses-hero__stat">
                        <span class="courses-hero__stat-number">{{ totalCourses }}</span>
                        <span class="courses-hero__stat-label">курсов</span>
                    </div>
                    <div class="courses-hero__stat">
                        <span class="courses-hero__stat-number">{{ totalLessons }}</span>
                        <span class="courses-hero__stat-label">уроков</span>
                    </div>
                    <div class="courses-hero__stat">
                        <span class="courses-hero__stat-number">100%</span>
                        <span class="courses-hero__stat-label">практики</span>
                    </div>
                </div>
            </div>
        </section>

        <!-- Categories Filter -->
        <section class="courses-categories">
            <div class="courses-categories__container">
                <h2 class="courses-categories__title">Категории курсов</h2>
                <div class="courses-categories__list">
                    <button
                        v-for="category in categories"
                        :key="category.id"
                        class="courses-categories__item"
                        :class="{
                            'courses-categories__item--active':
                                activeCategoryId === category.id,
                        }"
                        @click="setActiveCategory(category.id)"
                    >
                        <span class="courses-categories__item-name">{{
                            category.name_ru
                        }}</span>
                        <span class="courses-categories__item-count">{{
                            getProductCountByCategory(category.id)
                        }}</span>
                    </button>
                    <button
                        class="courses-categories__item"
                        :class="{
                            'courses-categories__item--active': activeCategoryId === null,
                        }"
                        @click="setActiveCategory(null)"
                    >
                        <span class="courses-categories__item-name">Все курсы</span>
                        <span class="courses-categories__item-count">{{
                            products.length
                        }}</span>
                    </button>
                </div>
            </div>
        </section>

        <!-- Courses Grid -->
        <section class="courses-grid">
            <div class="courses-grid__container">
                <div class="courses-grid__header">
                    <h2 class="courses-grid__title">{{ activeCategoryName }} курсы</h2>
                    <div class="courses-grid__controls">
                        <select v-model="sortBy" class="courses-grid__sort">
                            <option value="created_at">По новизне</option>
                            <option value="price">По цене</option>
                            <option value="name">По названию</option>
                        </select>
                        <select v-model="sortOrder" class="courses-grid__sort">
                            <option value="desc">По убыванию</option>
                            <option value="asc">По возрастанию</option>
                        </select>
                    </div>
                </div>

                <ProductList
                    :products="filteredProducts"
                    :isLoading="isLoading"
                    :error="error"
                    :columns="3"
                    @select="handleProductSelect"
                    @retry="fetchData"
                    class="courses-grid__list"
                />

                <div
                    v-if="!isLoading && filteredProducts.length === 0"
                    class="courses-empty"
                >
                    <p class="courses-empty__message">Курсы не найдены</p>
                    <button @click="resetFilters" class="courses-empty__button">
                        Сбросить фильтры
                    </button>
                </div>
            </div>
        </section>

        <!-- FAQ Section -->
        <section class="courses-faq">
            <div class="courses-faq__container">
                <h2 class="courses-faq__title">Частые вопросы</h2>
                <div class="courses-faq__list">
                    <div class="courses-faq__item">
                        <h3 class="courses-faq__question">
                            Как долго длится доступ к курсу?
                        </h3>
                        <p class="courses-faq__answer">
                            Доступ к курсу предоставляется навсегда. Вы можете
                            пересматривать материалы в любое время.
                        </p>
                    </div>
                    <div class="courses-faq__item">
                        <h3 class="courses-faq__question">
                            Можно ли смотреть курсы на телефоне?
                        </h3>
                        <p class="courses-faq__answer">
                            Да, все наши курсы адаптированы для просмотра на любых
                            устройствах: компьютерах, планшетах и смартфонах.
                        </p>
                    </div>
                    <div class="courses-faq__item">
                        <h3 class="courses-faq__question">
                            Есть ли обратная связь от преподавателя?
                        </h3>
                        <p class="courses-faq__answer">
                            Да, у каждого курса есть закрытый чат с преподавателем, где вы
                            можете задавать вопросы и получать обратную связь.
                        </p>
                    </div>
                    <div class="courses-faq__item">
                        <h3 class="courses-faq__question">Как происходит оплата?</h3>
                        <p class="courses-faq__answer">
                            Оплатить курс можно банковской картой, через PayPal или
                            Яндекс.Кассу. Доступ открывается мгновенно после оплаты.
                        </p>
                    </div>
                </div>
            </div>
        </section>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useProductStore } from '~/stores/ProductStore';
import ProductList from '~/components/ProductList.vue';
import type { Product } from '~/stores/ProductStore';

const productStore = useProductStore();

// State
const activeCategoryId = ref<number | null>(null);
const sortBy = ref<'created_at' | 'price' | 'name'>('created_at');
const sortOrder = ref<'asc' | 'desc'>('desc');

// Computed
const categories = computed(() => productStore.activeCategories);
const products = computed(() => productStore.publishedProducts);
const isLoading = computed(() => productStore.isLoading);
const error = computed(() => productStore.error);

const totalCourses = computed(() => {
    return products.value.filter((p) => p.type === 'course').length;
});

const totalLessons = computed(() => {
    return products.value.reduce((sum, product) => sum + product.total_lessons, 0);
});

const activeCategoryName = computed(() => {
    if (activeCategoryId.value === null) return 'Все';
    const category = categories.value.find((c) => c.id === activeCategoryId.value);
    return category?.name_ru || 'Все';
});

const filteredProducts = computed(() => {
    let filtered = products.value.filter((p) => p.type === 'course');

    // Filter by category
    if (activeCategoryId.value !== null) {
        filtered = filtered.filter((p) => p.category_id === activeCategoryId.value);
    }

    // Sort
    filtered = [...filtered].sort((a, b) => {
        let aValue: any, bValue: any;

        if (sortBy.value === 'price') {
            aValue = a.discount_price || a.price;
            bValue = b.discount_price || b.price;
        } else if (sortBy.value === 'name') {
            aValue = a.name_ru || a.name_en;
            bValue = b.name_ru || b.name_en;
        } else {
            aValue = new Date(a.created_at).getTime();
            bValue = new Date(b.created_at).getTime();
        }

        if (sortOrder.value === 'asc') {
            return aValue > bValue ? 1 : -1;
        } else {
            return aValue < bValue ? 1 : -1;
        }
    });

    return filtered;
});

// Methods
const getProductCountByCategory = (categoryId: number) => {
    return products.value.filter(
        (p) => p.category_id === categoryId && p.type === 'course'
    ).length;
};

const setActiveCategory = (categoryId: number | null) => {
    activeCategoryId.value = categoryId;
};

const handleProductSelect = (product: Product) => {
    // Navigate to product detail page
    console.log('Selected product:', product);
    // In a real app: router.push(`/courses/${product.slug}`)
};

const resetFilters = () => {
    activeCategoryId.value = null;
    sortBy.value = 'created_at';
    sortOrder.value = 'desc';
};

const fetchData = async () => {
    try {
        await Promise.all([
            productStore.fetchCategories(),
            productStore.fetchProducts({
                type: 'course',
                status: 'published',
            }),
        ]);
    } catch (err) {
        console.error('Failed to fetch data:', err);
    }
};

// Lifecycle
onMounted(() => {
    fetchData();
});
</script>

<style scoped>
.courses-page {
    min-height: 100vh;
    background: var(--color-background);
}

/* Hero Section */
.courses-hero {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 80px 20px;
    text-align: center;
}

.courses-hero__container {
    max-width: 1200px;
    margin: 0 auto;
}

.courses-hero__title {
    font-size: 48px;
    font-weight: 700;
    margin-bottom: 20px;
    line-height: 1.2;
}

.courses-hero__subtitle {
    font-size: 20px;
    opacity: 0.9;
    max-width: 600px;
    margin: 0 auto 40px;
    line-height: 1.6;
}

.courses-hero__stats {
    display: flex;
    justify-content: center;
    gap: 40px;
    flex-wrap: wrap;
}

.courses-hero__stat {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.courses-hero__stat-number {
    font-size: 36px;
    font-weight: 700;
    margin-bottom: 8px;
}

.courses-hero__stat-label {
    font-size: 16px;
    opacity: 0.8;
}

/* Categories Section */
.courses-categories {
    padding: 60px 20px;
    background: #f8f9fa;
}

.courses-categories__container {
    max-width: 1200px;
    margin: 0 auto;
}

.courses-categories__title {
    font-size: 32px;
    font-weight: 600;
    margin-bottom: 30px;
    text-align: center;
    color: var(--color-text);
}

.courses-categories__list {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: center;
}

.courses-categories__item {
    background: white;
    border: 2px solid #e9ecef;
    border-radius: 50px;
    padding: 12px 24px;
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 16px;
    font-weight: 500;
}

.courses-categories__item:hover {
    border-color: var(--color-primary);
    transform: translateY(-2px);
}

.courses-categories__item--active {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: white;
}

.courses-categories__item-count {
    background: #e9ecef;
    border-radius: 12px;
    padding: 2px 8px;
    font-size: 14px;
    font-weight: 500;
}

.courses-categories__item--active .courses-categories__item-count {
    background: rgba(255, 255, 255, 0.2);
}

/* Courses Grid Section */
.courses-grid {
    padding: 60px 20px;
}

.courses-grid__container {
    max-width: 1200px;
    margin: 0 auto;
}

.courses-grid__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 40px;
    flex-wrap: wrap;
    gap: 20px;
}

.courses-grid__title {
    font-size: 32px;
    font-weight: 600;
    color: var(--color-text);
}

.courses-grid__controls {
    display: flex;
    gap: 12px;
}

.courses-grid__sort {
    padding: 10px 16px;
    border: 1px solid #e9ecef;
    border-radius: 8px;
    background: white;
    font-size: 14px;
    color: var(--color-text);
    cursor: pointer;
    transition: border-color 0.2s ease;
}

.courses-grid__sort:hover {
    border-color: var(--color-primary);
}

.courses-grid__list {
    margin-bottom: 40px;
}

.courses-empty {
    text-align: center;
    padding: 60px 20px;
    background: #f8f9fa;
    border-radius: 12px;
    border: 1px dashed #dee2e6;
}

.courses-empty__message {
    font-size: 18px;
    color: #6c757d;
    margin-bottom: 20px;
}

.courses-empty__button {
    background: var(--color-primary);
    color: white;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.courses-empty__button:hover {
    background: var(--color-primary-dark);
}

/* FAQ Section */
.courses-faq {
    padding: 60px 20px;
    background: #f8f9fa;
}

.courses-faq__container {
    max-width: 800px;
    margin: 0 auto;
}

.courses-faq__title {
    font-size: 32px;
    font-weight: 600;
    margin-bottom: 40px;
    text-align: center;
    color: var(--color-text);
}

.courses-faq__list {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.courses-faq__item {
    background: white;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.courses-faq__question {
    font-size: 18px;
    font-weight: 600;
    margin-bottom: 12px;
    color: var(--color-text);
}

.courses-faq__answer {
    font-size: 16px;
    line-height: 1.6;
    color: #6c757d;
}

/* Responsive */
@media (max-width: 768px) {
    .courses-hero__title {
        font-size: 36px;
    }

    .courses-hero__subtitle {
        font-size: 18px;
    }

    .courses-hero__stats {
        gap: 30px;
    }

    .courses-hero__stat-number {
        font-size: 28px;
    }

    .courses-grid__header {
        flex-direction: column;
        align-items: stretch;
    }

    .courses-grid__controls {
        justify-content: center;
    }

    .courses-categories__list {
        justify-content: flex-start;
        overflow-x: auto;
        padding-bottom: 10px;
    }

    .courses-categories__item {
        white-space: nowrap;
    }
}
</style>
