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
                <AdminStatCard
                    v-for="stat in statCards"
                    :key="stat.key"
                    :stat="stat"
                    :stats="stats"
                />
            </div>

            <!-- Charts Row -->
            <AdminChartsSection :stats="stats" />

            <!-- Navigation Sections -->
            <AdminQuickNav />

            <!-- Bottom Row: Recent Activity -->
            <div class="admin-dashboard__bottom-row">
                <AdminActivityFeed
                    :activities="activities"
                    :loading="refreshing"
                    @refresh="loadActivity"
                />
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { fetchDashboardStats, fetchRecentActivity } from '~/api/admin';
import type { DashboardStats, RecentActivityItem } from '~/types/dashboard';
import AdminStatCard from '~/components/admin/AdminStatCard.vue';
import AdminChartsSection from '~/components/admin/AdminChartsSection.vue';
import AdminQuickNav from '~/components/admin/AdminQuickNav.vue';
import AdminActivityFeed from '~/components/admin/AdminActivityFeed.vue';

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

/* Charts Row (skeleton) */
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

/* Bottom Row */
.admin-dashboard__bottom-row {
    display: grid;
    grid-template-columns: 1fr;
    gap: 16px;
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
