<template>
    <div class="admin-dashboard">
        <!-- Page Header -->
        <div class="admin-dashboard__header">
            <div>
                <h1 class="admin-dashboard__title">Dashboard</h1>
                <p class="admin-dashboard__subtitle">
                    Добро пожаловать в панель управления
                </p>
            </div>
            <div class="admin-dashboard__header-actions">
                <UButton
                    icon="i-lucide-refresh-cw"
                    color="neutral"
                    variant="outline"
                    :loading="refreshing"
                    @click="refreshData"
                >
                    Обновить
                </UButton>
                <UButton icon="i-lucide-plus" color="primary" to="/admin/gallery/add">
                    Добавить работу
                </UButton>
            </div>
        </div>

        <!-- Loading State -->
        <template v-if="loading">
            <div class="admin-dashboard__stats-grid">
                <div v-for="i in 8" :key="i" class="admin-dashboard__skeleton-card">
                    <div class="admin-dashboard__skeleton-icon" />
                    <div class="admin-dashboard__skeleton-text">
                        <div class="admin-dashboard__skeleton-value" />
                        <div class="admin-dashboard__skeleton-label" />
                    </div>
                </div>
            </div>
            <div class="admin-dashboard__charts-row">
                <div class="admin-dashboard__skeleton-chart" />
                <div class="admin-dashboard__skeleton-chart" />
            </div>
        </template>

        <!-- Error State -->
        <template v-else-if="error">
            <UCard>
                <div class="admin-dashboard__error">
                    <UIcon
                        name="i-lucide-alert-circle"
                        class="admin-dashboard__error-icon"
                    />
                    <p>{{ error }}</p>
                    <UButton color="primary" variant="outline" @click="loadData">
                        Повторить загрузку
                    </UButton>
                </div>
            </UCard>
        </template>

        <!-- Dashboard Content -->
        <template v-else>
            <!-- Stats Grid -->
            <div class="admin-dashboard__stats-grid admin-stagger">
                <UCard
                    v-for="stat in statCards"
                    :key="stat.key"
                    class="admin-dashboard__stat-card"
                    :ui="{ body: 'p-0' }"
                >
                    <div class="admin-dashboard__stat-inner">
                        <div
                            class="admin-dashboard__stat-icon"
                            :class="`admin-dashboard__stat-icon--${stat.color}`"
                        >
                            <UIcon :name="stat.icon" class="size-5" />
                        </div>
                        <div class="admin-dashboard__stat-info">
                            <span class="admin-dashboard__stat-value">{{
                                formatStatValue(stat)
                            }}</span>
                            <span class="admin-dashboard__stat-label">{{
                                stat.label
                            }}</span>
                        </div>
                    </div>
                </UCard>
            </div>

            <!-- Charts Row -->
            <div class="admin-dashboard__charts-row">
                <!-- Sales Chart -->
                <UCard class="admin-dashboard__chart-card">
                    <template #header>
                        <div class="admin-dashboard__chart-header">
                            <h3 class="admin-dashboard__chart-title">
                                Продажи по месяцам
                            </h3>
                            <UBadge
                                v-if="stats.revenue_month > 0"
                                color="success"
                                variant="soft"
                            >
                                +{{ formatPrice(stats.revenue_month) }} в этом месяце
                            </UBadge>
                        </div>
                    </template>
                    <div class="admin-dashboard__chart-body">
                        <div
                            v-if="(stats.sales_by_month ?? []).length === 0"
                            class="admin-dashboard__chart-empty"
                        >
                            <UIcon
                                name="i-lucide-bar-chart-3"
                                class="size-8 text-gray-400"
                            />
                            <p>Данных о продажах пока нет</p>
                        </div>
                        <div v-else class="admin-dashboard__bar-chart">
                            <div
                                v-for="(item, index) in stats.sales_by_month ?? []"
                                :key="item.month"
                                class="admin-dashboard__bar-item"
                            >
                                <div class="admin-dashboard__bar-tooltip">
                                    <span class="admin-dashboard__bar-tooltip-count"
                                        >{{ item.count }} шт.</span
                                    >
                                    <span class="admin-dashboard__bar-tooltip-revenue">{{
                                        formatPrice(item.revenue)
                                    }}</span>
                                </div>
                                <div
                                    class="admin-dashboard__bar"
                                    :style="{ height: getBarHeight(item.count) }"
                                />
                                <span class="admin-dashboard__bar-label">{{
                                    formatMonth(item.month)
                                }}</span>
                            </div>
                        </div>
                    </div>
                </UCard>

                <!-- Popular Categories -->
                <UCard class="admin-dashboard__chart-card">
                    <template #header>
                        <div class="admin-dashboard__chart-header">
                            <h3 class="admin-dashboard__chart-title">
                                Популярные категории
                            </h3>
                        </div>
                    </template>
                    <div class="admin-dashboard__chart-body">
                        <div
                            v-if="(stats.popular_categories ?? []).length === 0"
                            class="admin-dashboard__chart-empty"
                        >
                            <UIcon
                                name="i-lucide-pie-chart"
                                class="size-8 text-gray-400"
                            />
                            <p>Категории не найдены</p>
                        </div>
                        <div v-else class="admin-dashboard__category-list">
                            <div
                                v-for="cat in stats.popular_categories ?? []"
                                :key="cat.name"
                                class="admin-dashboard__category-item"
                            >
                                <div class="admin-dashboard__category-info">
                                    <span class="admin-dashboard__category-name">{{
                                        cat.name
                                    }}</span>
                                    <span class="admin-dashboard__category-count"
                                        >{{ cat.count }}
                                        {{
                                            pluralize(
                                                cat.count,
                                                'продукт',
                                                'продукта',
                                                'продуктов'
                                            )
                                        }}</span
                                    >
                                </div>
                                <div class="admin-dashboard__category-bar-bg">
                                    <div
                                        class="admin-dashboard__category-bar-fill"
                                        :style="{ width: getCategoryWidth(cat.count) }"
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </UCard>
            </div>

            <!-- Navigation Sections -->
            <div class="admin-dashboard__nav-section">
                <h3 class="admin-dashboard__nav-section-title">Навигация по разделам</h3>
                <div
                    v-for="(section, sIdx) in navSections"
                    :key="sIdx"
                    class="admin-dashboard__nav-group"
                >
                    <span class="admin-dashboard__nav-group-label">{{
                        section.label
                    }}</span>
                    <div class="admin-dashboard__nav-grid">
                        <NuxtLink
                            v-for="item in section.items"
                            :key="item.to"
                            :to="item.to"
                            class="admin-dashboard__nav-tile"
                        >
                            <div
                                class="admin-dashboard__nav-tile-icon"
                                :class="`admin-dashboard__nav-tile-icon--${item.color}`"
                            >
                                <UIcon :name="item.icon" class="size-5" />
                            </div>
                            <div class="admin-dashboard__nav-tile-info">
                                <span class="admin-dashboard__nav-tile-label">{{
                                    item.label
                                }}</span>
                                <span class="admin-dashboard__nav-tile-desc">{{
                                    item.desc
                                }}</span>
                            </div>
                            <UIcon
                                name="i-lucide-chevron-right"
                                class="admin-dashboard__nav-tile-arrow size-4"
                            />
                        </NuxtLink>
                    </div>
                </div>
            </div>

            <!-- Bottom Row: Recent Activity -->
            <div class="admin-dashboard__bottom-row">
                <!-- Recent Activity -->
                <UCard class="admin-dashboard__activity-card">
                    <template #header>
                        <div class="admin-dashboard__section-header">
                            <h3 class="admin-dashboard__section-title">
                                Последняя активность
                            </h3>
                            <UButton
                                v-if="activities.length > 0"
                                color="neutral"
                                variant="ghost"
                                icon="i-lucide-refresh-cw"
                                size="xs"
                                @click="loadActivity"
                            />
                        </div>
                    </template>
                    <div
                        v-if="activities.length === 0"
                        class="admin-dashboard__activity-empty"
                    >
                        <UIcon name="i-lucide-clock" class="size-6 text-gray-400" />
                        <p>Активность пока отсутствует</p>
                    </div>
                    <div v-else class="admin-dashboard__activity-list">
                        <div
                            v-for="activity in activities"
                            :key="activity.id"
                            class="admin-dashboard__activity-item"
                        >
                            <div
                                class="admin-dashboard__activity-dot"
                                :class="`admin-dashboard__activity-dot--${activity.type}`"
                            />
                            <div class="admin-dashboard__activity-content">
                                <span class="admin-dashboard__activity-text">{{
                                    activity.text
                                }}</span>
                                <span class="admin-dashboard__activity-time">{{
                                    activity.time
                                }}</span>
                            </div>
                        </div>
                    </div>
                </UCard>
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { fetchDashboardStats, fetchRecentActivity } from '~/api/admin';
import type { DashboardStats, RecentActivityItem } from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const loading = ref(true);
const refreshing = ref(false);
const error = ref<string | null>(null);
const stats = ref<DashboardStats>({
    gallery_works_count: 0,
    shop_items_count: 0,
    news_count: 0,
    users_count: 0,
    courses_count: 0,
    master_classes_count: 0,
    reviews_total: 0,
    reviews_pending: 0,
    purchases_total: 0,
    revenue_total: 0,
    revenue_month: 0,
    active_chats: 0,
    users_registered_month: 0,
    sales_by_month: [],
    popular_categories: [],
});
const activities = ref<RecentActivityItem[]>([]);

interface StatCard {
    key: keyof DashboardStats;
    label: string;
    icon: string;
    color: string;
    prefix?: string;
    suffix?: string;
    format?: 'number' | 'price';
}

const statCards = computed<StatCard[]>(() => [
    {
        key: 'gallery_works_count',
        label: 'Работ в галерее',
        icon: 'i-lucide-image',
        color: 'purple',
    },
    {
        key: 'shop_items_count',
        label: 'Товаров в магазине',
        icon: 'i-lucide-shopping-bag',
        color: 'green',
    },
    { key: 'news_count', label: 'Новостей', icon: 'i-lucide-newspaper', color: 'blue' },
    {
        key: 'users_count',
        label: 'Пользователей',
        icon: 'i-lucide-users',
        color: 'orange',
    },
    {
        key: 'courses_count',
        label: `Курсов / МК (${stats.value.master_classes_count} МК)`,
        icon: 'i-lucide-graduation-cap',
        color: 'pink',
        format: 'number',
    },
    {
        key: 'reviews_pending',
        label: `Отзывов ожидает (всего ${stats.value.reviews_total})`,
        icon: 'i-lucide-star',
        color: 'red',
    },
    {
        key: 'purchases_total',
        label: 'Покупок',
        icon: 'i-lucide-shopping-cart',
        color: 'teal',
    },
    {
        key: 'revenue_total',
        label: 'Выручка',
        icon: 'i-lucide-circle-dollar-sign',
        color: 'yellow',
        format: 'price',
    },
]);

const quickActions = [
    {
        label: 'Добавить работу',
        icon: 'i-lucide-plus',
        color: 'primary' as const,
        to: '/admin/gallery/add',
    },
    {
        label: 'Добавить товар',
        icon: 'i-lucide-plus',
        color: 'success' as const,
        to: '/admin/shop/add',
    },
    {
        label: 'Новая новость',
        icon: 'i-lucide-plus',
        color: 'info' as const,
        to: '/admin/news/add',
    },
    {
        label: 'Создать курс',
        icon: 'i-lucide-plus',
        color: 'warning' as const,
        to: '/admin/courses/add',
    },
    {
        label: 'Создать МК',
        icon: 'i-lucide-plus',
        color: 'error' as const,
        to: '/admin/master-classes/add',
    },
    {
        label: 'Модерация отзывов',
        icon: 'i-lucide-message-square',
        color: 'neutral' as const,
        to: '/admin/reviews',
    },
];

interface NavItem {
    to: string;
    label: string;
    desc: string;
    icon: string;
    color: string;
}

interface NavSection {
    label: string;
    items: NavItem[];
}

const navSections = computed<NavSection[]>(() => [
    {
        label: 'Контент',
        items: [
            {
                to: '/admin/gallery',
                label: 'Галерея',
                desc: 'Управление работами в галерее',
                icon: 'i-lucide-image',
                color: 'purple',
            },
            {
                to: '/admin/shop',
                label: 'Магазин',
                desc: 'Управление товарами в магазине',
                icon: 'i-lucide-shopping-bag',
                color: 'green',
            },
            {
                to: '/admin/news',
                label: 'Новости',
                desc: 'Управление новостями',
                icon: 'i-lucide-newspaper',
                color: 'blue',
            },
            {
                to: '/admin/courses',
                label: 'Курсы',
                desc: 'Управление курсами',
                icon: 'i-lucide-graduation-cap',
                color: 'pink',
            },
            {
                to: '/admin/master-classes',
                label: 'Мастер-классы',
                desc: 'Управление мастер-классами',
                icon: 'i-lucide-video',
                color: 'orange',
            },
            {
                to: '/admin/lessons',
                label: 'Уроки',
                desc: 'Управление уроками курсов и МК',
                icon: 'i-lucide-book-open',
                color: 'teal',
            },
        ],
    },
    {
        label: 'Пользователи',
        items: [
            {
                to: '/admin/users',
                label: 'Пользователи',
                desc: 'Управление пользователями',
                icon: 'i-lucide-users',
                color: 'orange',
            },
            {
                to: '/admin/reviews',
                label: 'Отзывы',
                desc: 'Модерация отзывов пользователей',
                icon: 'i-lucide-star',
                color: 'red',
            },
        ],
    },
    {
        label: 'Финансы',
        items: [
            {
                to: '/admin/purchases',
                label: 'Покупки',
                desc: 'Управление покупками',
                icon: 'i-lucide-shopping-cart',
                color: 'teal',
            },
            {
                to: '/admin/payments',
                label: 'Платежи',
                desc: 'Управление платежами и возвратами',
                icon: 'i-lucide-credit-card',
                color: 'yellow',
            },
            {
                to: '/admin/promo-codes',
                label: 'Промокоды',
                desc: 'Управление промокодами и скидками',
                icon: 'i-lucide-ticket-percent',
                color: 'green',
            },
        ],
    },
    {
        label: 'Коммуникация',
        items: [
            {
                to: '/admin/chats',
                label: 'Чаты',
                desc: 'Управление обращениями пользователей',
                icon: 'i-lucide-message-square',
                color: 'blue',
            },
        ],
    },
    {
        label: 'Справочники',
        items: [
            {
                to: '/admin/categories',
                label: 'Категории',
                desc: 'Управление категориями продуктов',
                icon: 'i-lucide-folder-tree',
                color: 'purple',
            },
            {
                to: '/admin/tags',
                label: 'Теги',
                desc: 'Управление тегами продуктов',
                icon: 'i-lucide-tags',
                color: 'pink',
            },
            {
                to: '/admin/settings',
                label: 'Настройки',
                desc: 'Управление справочниками',
                icon: 'i-lucide-settings',
                color: 'neutral',
            },
        ],
    },
]);

onMounted(async () => {
    await loadData();
});

async function loadData() {
    loading.value = true;
    error.value = null;
    try {
        const [statsData, activityData] = await Promise.all([
            fetchDashboardStats(),
            fetchRecentActivity(),
        ]);
        stats.value = statsData;
        activities.value = activityData;
    } catch (e) {
        console.error('Error loading dashboard data:', e);
        error.value =
            'Не удалось загрузить данные дашборда. Проверьте подключение к серверу.';
    } finally {
        loading.value = false;
    }
}

async function refreshData() {
    refreshing.value = true;
    try {
        const [statsData, activityData] = await Promise.all([
            fetchDashboardStats(),
            fetchRecentActivity(),
        ]);
        stats.value = statsData;
        activities.value = activityData;
    } catch (e) {
        console.error('Error refreshing dashboard data:', e);
    } finally {
        refreshing.value = false;
    }
}

async function loadActivity() {
    try {
        activities.value = await fetchRecentActivity();
    } catch (e) {
        console.error('Error loading activity:', e);
    }
}

function formatStatValue(stat: StatCard): string {
    const value = stats.value[stat.key] as number;
    if (stat.format === 'price') {
        return formatPrice(value);
    }
    return value.toLocaleString('ru-RU');
}

function formatPrice(value: number): string {
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(value);
}

function formatMonth(month: string): string {
    const [year, m] = month.split('-');
    const months = [
        'Янв',
        'Фев',
        'Мар',
        'Апр',
        'Май',
        'Июн',
        'Июл',
        'Авг',
        'Сен',
        'Окт',
        'Ноя',
        'Дек',
    ];
    return `${months[parseInt(m || '0') - 1] || ''} ${year}`;
}

function getBarHeight(count: number): string {
    const maxCount = Math.max(
        ...(stats.value.sales_by_month ?? []).map((s) => s.count),
        1
    );
    const percentage = (count / maxCount) * 100;
    return `${Math.max(percentage, 4)}%`;
}

function getCategoryWidth(count: number): string {
    const maxCount = Math.max(
        ...(stats.value.popular_categories ?? []).map((c) => c.count),
        1
    );
    return `${(count / maxCount) * 100}%`;
}

function pluralize(count: number, one: string, few: string, many: string): string {
    const mod10 = count % 10;
    const mod100 = count % 100;
    if (mod100 >= 11 && mod100 <= 19) return many;
    if (mod10 === 1) return one;
    if (mod10 >= 2 && mod10 <= 4) return few;
    return many;
}
</script>

<style scoped>
.admin-dashboard {
    max-width: 1400px;
    margin: 0 auto;
}

/* Header */
.admin-dashboard__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
    gap: 16px;
}

@media (max-width: 640px) {
    .admin-dashboard__header {
        flex-direction: column;
    }
}

.admin-dashboard__title {
    font-size: 28px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    margin: 0 0 4px;
}

.admin-dashboard__subtitle {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    margin: 0;
}

.admin-dashboard__header-actions {
    display: flex;
    gap: 8px;
    flex-shrink: 0;
}

/* Stats Grid */
.admin-dashboard__stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
    margin-bottom: 24px;
}

.admin-dashboard__stat-card {
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.admin-dashboard__stat-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.admin-dashboard__stat-inner {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
}

.admin-dashboard__stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}

.admin-dashboard__stat-icon--purple {
    background: rgba(108, 92, 231, 0.1);
    color: #6c5ce7;
}

.admin-dashboard__stat-icon--green {
    background: rgba(0, 184, 148, 0.1);
    color: #00b894;
}

.admin-dashboard__stat-icon--blue {
    background: rgba(116, 185, 255, 0.1);
    color: #74b9ff;
}

.admin-dashboard__stat-icon--orange {
    background: rgba(253, 203, 110, 0.15);
    color: #e17055;
}

.admin-dashboard__stat-icon--pink {
    background: rgba(232, 67, 147, 0.1);
    color: #e84393;
}

.admin-dashboard__stat-icon--red {
    background: rgba(225, 112, 85, 0.1);
    color: #e17055;
}

.admin-dashboard__stat-icon--teal {
    background: rgba(0, 206, 201, 0.1);
    color: #00cec9;
}

.admin-dashboard__stat-icon--yellow {
    background: rgba(253, 203, 110, 0.15);
    color: #fdcb6e;
}

.admin-dashboard__stat-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-dashboard__stat-value {
    font-size: 24px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.2;
}

.admin-dashboard__stat-label {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
    line-height: 1.3;
}

/* Charts Row */
.admin-dashboard__charts-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    margin-bottom: 24px;
}

@media (max-width: 900px) {
    .admin-dashboard__charts-row {
        grid-template-columns: 1fr;
    }
}

.admin-dashboard__chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
}

.admin-dashboard__chart-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin: 0;
}

.admin-dashboard__chart-body {
    min-height: 200px;
}

.admin-dashboard__chart-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 40px 20px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-dashboard__chart-empty p {
    margin: 0;
    font-size: 14px;
}

/* Bar Chart */
.admin-dashboard__bar-chart {
    display: flex;
    align-items: flex-end;
    gap: 8px;
    height: 180px;
    padding-top: 24px;
}

.admin-dashboard__bar-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    position: relative;
}

.admin-dashboard__bar-tooltip {
    position: absolute;
    bottom: 100%;
    left: 50%;
    transform: translateX(-50%);
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 6px;
    padding: 4px 8px;
    font-size: 11px;
    white-space: nowrap;
    opacity: 0;
    transition: opacity 0.2s;
    pointer-events: none;
    z-index: 10;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.admin-dashboard__bar-item:hover .admin-dashboard__bar-tooltip {
    opacity: 1;
}

.admin-dashboard__bar-tooltip-count {
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-dashboard__bar-tooltip-revenue {
    color: var(--admin-text-secondary, #636e72);
}

.admin-dashboard__bar {
    width: 100%;
    max-width: 48px;
    background: linear-gradient(180deg, #6c5ce7 0%, #a29bfe 100%);
    border-radius: 4px 4px 0 0;
    min-height: 4px;
    transition: height 0.3s ease;
    cursor: pointer;
}

.admin-dashboard__bar-label {
    font-size: 10px;
    color: var(--admin-text-secondary, #636e72);
    text-align: center;
    white-space: nowrap;
}

/* Category List */
.admin-dashboard__category-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.admin-dashboard__category-item {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.admin-dashboard__category-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.admin-dashboard__category-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--admin-text-primary, #2d3436);
}

.admin-dashboard__category-count {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-dashboard__category-bar-bg {
    height: 8px;
    background: var(--admin-border, #e0e0e0);
    border-radius: 4px;
    overflow: hidden;
}

.admin-dashboard__category-bar-fill {
    height: 100%;
    background: linear-gradient(90deg, #6c5ce7, #a29bfe);
    border-radius: 4px;
    transition: width 0.3s ease;
}

/* Navigation Section */
.admin-dashboard__nav-section {
    margin-bottom: 24px;
}

.admin-dashboard__nav-section-title {
    font-size: 18px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    margin: 0 0 16px;
}

.admin-dashboard__nav-group {
    margin-bottom: 20px;
}

.admin-dashboard__nav-group:last-child {
    margin-bottom: 0;
}

.admin-dashboard__nav-group-label {
    display: block;
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--admin-text-secondary, #636e72);
    margin-bottom: 8px;
    padding-left: 4px;
}

.admin-dashboard__nav-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 8px;
}

.admin-dashboard__nav-tile {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 10px;
    text-decoration: none;
    transition: all 0.2s ease;
    position: relative;
}

.admin-dashboard__nav-tile:hover {
    border-color: var(--admin-primary, #6c5ce7);
    box-shadow: 0 2px 8px rgba(108, 92, 231, 0.1);
    transform: translateY(-1px);
}

.admin-dashboard__nav-tile:active {
    transform: translateY(0);
}

.admin-dashboard__nav-tile-icon {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}

.admin-dashboard__nav-tile-icon--purple {
    background: rgba(108, 92, 231, 0.1);
    color: #6c5ce7;
}

.admin-dashboard__nav-tile-icon--green {
    background: rgba(0, 184, 148, 0.1);
    color: #00b894;
}

.admin-dashboard__nav-tile-icon--blue {
    background: rgba(116, 185, 255, 0.1);
    color: #74b9ff;
}

.admin-dashboard__nav-tile-icon--orange {
    background: rgba(253, 203, 110, 0.15);
    color: #e17055;
}

.admin-dashboard__nav-tile-icon--pink {
    background: rgba(232, 67, 147, 0.1);
    color: #e84393;
}

.admin-dashboard__nav-tile-icon--red {
    background: rgba(225, 112, 85, 0.1);
    color: #e17055;
}

.admin-dashboard__nav-tile-icon--teal {
    background: rgba(0, 206, 201, 0.1);
    color: #00cec9;
}

.admin-dashboard__nav-tile-icon--yellow {
    background: rgba(253, 203, 110, 0.15);
    color: #fdcb6e;
}

.admin-dashboard__nav-tile-icon--neutral {
    background: rgba(99, 110, 114, 0.1);
    color: #636e72;
}

.admin-dashboard__nav-tile-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-dashboard__nav-tile-label {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.3;
}

.admin-dashboard__nav-tile-desc {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    line-height: 1.3;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.admin-dashboard__nav-tile-arrow {
    color: var(--admin-text-secondary, #636e72);
    flex-shrink: 0;
    opacity: 0;
    transition: opacity 0.2s ease, transform 0.2s ease;
}

.admin-dashboard__nav-tile:hover .admin-dashboard__nav-tile-arrow {
    opacity: 1;
    transform: translateX(2px);
}

@media (max-width: 640px) {
    .admin-dashboard__nav-grid {
        grid-template-columns: 1fr;
    }
}

/* Bottom Row */
.admin-dashboard__bottom-row {
    display: grid;
    grid-template-columns: 1fr;
    gap: 16px;
}

.admin-dashboard__section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.admin-dashboard__section-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin: 0;
}

/* Activity */
.admin-dashboard__activity-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 32px 20px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-dashboard__activity-empty p {
    margin: 0;
    font-size: 14px;
}

.admin-dashboard__activity-list {
    display: flex;
    flex-direction: column;
}

.admin-dashboard__activity-item {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    transition: background 0.15s;
}

.admin-dashboard__activity-item:last-child {
    border-bottom: none;
}

.admin-dashboard__activity-item:hover {
    background: rgba(108, 92, 231, 0.03);
}

.admin-dashboard__activity-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-top: 6px;
    flex-shrink: 0;
}

.admin-dashboard__activity-dot--gallery {
    background: #6c5ce7;
}
.admin-dashboard__activity-dot--shop {
    background: #00b894;
}
.admin-dashboard__activity-dot--news {
    background: #74b9ff;
}
.admin-dashboard__activity-dot--review {
    background: #fdcb6e;
}
.admin-dashboard__activity-dot--user {
    background: #e17055;
}
.admin-dashboard__activity-dot--purchase {
    background: #00cec9;
}

.admin-dashboard__activity-content {
    display: flex;
    flex-direction: column;
    gap: 2px;
    flex: 1;
    min-width: 0;
}

.admin-dashboard__activity-text {
    font-size: 13px;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
}

.admin-dashboard__activity-time {
    font-size: 11px;
    color: var(--admin-text-secondary, #636e72);
}

/* Quick Actions */
.admin-dashboard__quick-actions {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.admin-dashboard__quick-btn {
    justify-content: flex-start !important;
}

/* Error */
.admin-dashboard__error {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 40px 20px;
    text-align: center;
    color: var(--admin-text-secondary, #636e72);
}

.admin-dashboard__error-icon {
    width: 48px;
    height: 48px;
    color: #e17055;
}

.admin-dashboard__error p {
    margin: 0;
    font-size: 14px;
}

/* Skeleton Loading */
.admin-dashboard__skeleton-card {
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 12px;
    padding: 20px;
    display: flex;
    align-items: center;
    gap: 16px;
}

.admin-dashboard__skeleton-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    background: var(--admin-border, #e0e0e0);
    animation: pulse 1.5s ease-in-out infinite;
}

.admin-dashboard__skeleton-text {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.admin-dashboard__skeleton-value {
    width: 60%;
    height: 24px;
    border-radius: 4px;
    background: var(--admin-border, #e0e0e0);
    animation: pulse 1.5s ease-in-out infinite;
}

.admin-dashboard__skeleton-label {
    width: 80%;
    height: 14px;
    border-radius: 4px;
    background: var(--admin-border, #e0e0e0);
    animation: pulse 1.5s ease-in-out infinite;
}

.admin-dashboard__skeleton-chart {
    height: 260px;
    border-radius: 12px;
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
    0%,
    100% {
        opacity: 0.5;
    }
    50% {
        opacity: 0.2;
    }
}
</style>
