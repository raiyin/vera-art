<script setup lang="ts">
import type { DashboardStats } from '~/types/dashboard';

const props = defineProps<{
    stats: DashboardStats;
}>();

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
        ...(props.stats.sales_by_month ?? []).map((s) => s.count),
        1
    );
    const percentage = (count / maxCount) * 100;
    return `${Math.max(percentage, 4)}%`;
}

function getCategoryWidth(count: number): string {
    const maxCount = Math.max(
        ...(props.stats.popular_categories ?? []).map((c) => c.count),
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

<template>
    <div class="admin-dashboard__charts-row">
        <!-- Sales Chart -->
        <UCard class="admin-dashboard__chart-card">
            <template #header>
                <div class="admin-dashboard__chart-header">
                    <h3 class="admin-dashboard__chart-title">Продажи по месяцам</h3>
                    <UBadge v-if="stats.revenue_month > 0" color="success" variant="soft">
                        +{{ formatPrice(stats.revenue_month) }} в этом месяце
                    </UBadge>
                </div>
            </template>
            <div class="admin-dashboard__chart-body">
                <div
                    v-if="(stats.sales_by_month ?? []).length === 0"
                    class="admin-dashboard__chart-empty"
                >
                    <UIcon name="i-lucide-bar-chart-3" class="size-8 text-gray-400" />
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
                    <h3 class="admin-dashboard__chart-title">Популярные категории</h3>
                </div>
            </template>
            <div class="admin-dashboard__chart-body">
                <div
                    v-if="(stats.popular_categories ?? []).length === 0"
                    class="admin-dashboard__chart-empty"
                >
                    <UIcon name="i-lucide-pie-chart" class="size-8 text-gray-400" />
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
</template>

<style scoped>
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
</style>
