<template>
    <div class="admin-page">
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Покупки</span>
        </div>
        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Покупки</h1>
                <p class="admin-page__subtitle">
                    Управление покупками курсов и мастер-классов
                </p>
            </div>
            <div class="admin-page__header-actions">
                <UButton
                    icon="i-lucide-refresh-cw"
                    color="neutral"
                    variant="outline"
                    :loading="loading"
                    @click="loadData"
                >
                    Обновить
                </UButton>
            </div>
        </div>

        <!-- Search & Filters -->
        <UCard class="admin-page__filters-card">
            <div class="admin-page__filters">
                <div class="admin-page__search">
                    <UInput
                        v-model="searchQuery"
                        placeholder="Поиск по пользователю или продукту..."
                        icon="i-lucide-search"
                        color="neutral"
                        variant="outline"
                        class="w-full"
                        @keyup.enter="loadData"
                    />
                </div>
                <USelect
                    v-model="statusFilter"
                    :items="statusOptions"
                    placeholder="Все статусы"
                    color="neutral"
                    variant="outline"
                    class="admin-page__filter-select"
                    @change="loadData"
                />
                <USelect
                    v-model="productTypeFilter"
                    :items="productTypeOptions"
                    placeholder="Все типы"
                    color="neutral"
                    variant="outline"
                    class="admin-page__filter-select"
                    @change="loadData"
                />
            </div>
        </UCard>

        <!-- Loading State -->
        <div v-if="loading && !items.length" class="admin-page__loading">
            <UCard v-for="i in 5" :key="i">
                <div class="admin-page__skeleton-row">
                    <div class="admin-page__skeleton-lines">
                        <div class="admin-page__skeleton-line w-1/2" />
                        <div class="admin-page__skeleton-line w-1/3" />
                    </div>
                </div>
            </UCard>
        </div>

        <!-- Error State -->
        <UCard v-else-if="error" class="admin-page__error-card">
            <div class="admin-page__error">
                <UIcon name="i-lucide-alert-circle" class="admin-page__error-icon" />
                <p>{{ error }}</p>
                <UButton color="primary" variant="outline" @click="loadData">
                    Повторить загрузку
                </UButton>
            </div>
        </UCard>

        <!-- Empty State -->
        <UCard v-else-if="!items.length && !loading">
            <div class="admin-page__empty">
                <UIcon name="i-lucide-shopping-cart" class="admin-page__empty-icon" />
                <h3 class="admin-page__empty-title">Покупки не найдены</h3>
                <p class="admin-page__empty-desc">
                    <template v-if="searchQuery || statusFilter || productTypeFilter">
                        По заданным критериям ничего не найдено. Попробуйте изменить
                        параметры поиска.
                    </template>
                    <template v-else> Пока нет ни одной покупки. </template>
                </p>
            </div>
        </UCard>

        <!-- Purchases List -->
        <UCard v-else class="admin-page__table-card">
            <div class="admin-purchases__list">
                <div v-for="item in items" :key="item.id" class="admin-purchases__item">
                    <div class="admin-purchases__item-main">
                        <div class="admin-purchases__item-user">
                            <UAvatar
                                :text="((item.username || '?')[0] || '?').toUpperCase()"
                                size="sm"
                                color="neutral"
                            />
                            <div class="admin-purchases__item-user-info">
                                <span class="admin-purchases__item-username">{{
                                    item.username
                                }}</span>
                                <span class="admin-purchases__item-email">{{
                                    item.user_email
                                }}</span>
                            </div>
                        </div>
                        <div class="admin-purchases__item-product">
                            <span class="admin-purchases__item-product-title">{{
                                item.product_title_ru || item.product_title_en
                            }}</span>
                            <UBadge
                                :color="(item.product_type === 'course' ? 'info' : 'warning') as 'info' | 'warning'"
                                variant="soft"
                                size="sm"
                            >
                                {{
                                    item.product_type === 'course'
                                        ? 'Курс'
                                        : 'Мастер-класс'
                                }}
                            </UBadge>
                        </div>
                        <div class="admin-purchases__item-price">
                            <span class="admin-purchases__item-price-value"
                                >{{
                                    (item.price_paid / 100).toLocaleString('ru-RU')
                                }}
                                ₽</span
                            >
                        </div>
                        <div class="admin-purchases__item-status">
                            <UBadge
                                :color="purchaseStatusColor(item.status) as 'success' | 'error' | 'neutral'"
                                variant="soft"
                                size="sm"
                            >
                                {{ purchaseStatusLabel(item.status) }}
                            </UBadge>
                        </div>
                        <div class="admin-purchases__item-dates">
                            <div class="admin-purchases__item-date">
                                <UIcon
                                    name="i-lucide-calendar"
                                    class="admin-purchases__item-date-icon"
                                />
                                <span>{{ item.purchase_date }}</span>
                            </div>
                            <div
                                v-if="item.days_remaining >= 0"
                                class="admin-purchases__item-days"
                            >
                                <UIcon
                                    name="i-lucide-clock"
                                    class="admin-purchases__item-date-icon"
                                />
                                <span
                                    :class="{ 'text-red-500': item.days_remaining <= 7 }"
                                >
                                    {{ item.days_remaining }} дн.
                                </span>
                            </div>
                            <div v-else class="admin-purchases__item-days">
                                <UIcon
                                    name="i-lucide-infinity"
                                    class="admin-purchases__item-date-icon"
                                />
                                <span>Бессрочно</span>
                            </div>
                        </div>
                    </div>
                    <div class="admin-purchases__item-actions">
                        <UTooltip text="Продлить доступ">
                            <UButton
                                icon="i-lucide-clock-arrow-up"
                                color="neutral"
                                variant="ghost"
                                size="sm"
                                @click="openExtendModal(item)"
                            />
                        </UTooltip>
                        <UTooltip text="Отменить доступ">
                            <UButton
                                icon="i-lucide-ban"
                                color="error"
                                variant="ghost"
                                size="sm"
                                :disabled="item.status !== 'active'"
                                @click="confirmCancel(item)"
                            />
                        </UTooltip>
                    </div>
                </div>
            </div>

            <!-- Pagination -->
            <template #footer>
                <div class="admin-page__pagination">
                    <span class="admin-page__pagination-info">
                        {{ (page - 1) * perPage + 1 }}–{{
                            Math.min(page * perPage, total)
                        }}
                        из {{ total }}
                    </span>
                    <UPagination
                        v-model:page="page"
                        :total="total"
                        :items-per-page="perPage"
                        :max="5"
                        @update:page="loadData"
                    />
                </div>
            </template>
        </UCard>

        <!-- Extend Access Modal -->
        <UModal v-model:open="extendModalOpen" class="max-w-md">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-semibold">Продлить доступ</h3>
                    <UButton
                        icon="i-lucide-x"
                        color="neutral"
                        variant="ghost"
                        size="sm"
                        @click="extendModalOpen = false"
                    />
                </div>
            </template>
            <template #body>
                <div class="space-y-4">
                    <p class="text-sm text-gray-600 dark:text-gray-400">
                        Пользователь: <strong>{{ extendTarget?.username }}</strong
                        ><br />
                        Продукт: <strong>{{ extendTarget?.product_title_ru }}</strong>
                    </p>
                    <UInput
                        v-model="extendDays"
                        type="number"
                        placeholder="Количество дней"
                        color="neutral"
                        variant="outline"
                        min="1"
                    />
                </div>
            </template>
            <template #footer>
                <div class="flex justify-end gap-2">
                    <UButton
                        color="neutral"
                        variant="outline"
                        @click="extendModalOpen = false"
                    >
                        Отмена
                    </UButton>
                    <UButton
                        color="primary"
                        :loading="extending"
                        :disabled="!extendDays || extendDays < 1"
                        @click="doExtend"
                    >
                        Продлить
                    </UButton>
                </div>
            </template>
        </UModal>

        <!-- Cancel Confirmation Modal -->
        <AdminConfirmDialog
            :visible="cancelModalOpen"
            title="Подтверждение"
            :message="
                cancelTarget
                    ? `Вы уверены, что хотите отменить доступ для покупки #${cancelTarget.id}?`
                    : ''
            "
            type="danger"
            confirm-text="Да, отменить"
            cancel-text="Нет"
            loading-text="Выполнение..."
            :loading="cancelling"
            @confirm="doCancel"
            @cancel="cancelModalOpen = false"
            @update:visible="cancelModalOpen = $event"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, } from 'vue';
import {
    fetchAdminPurchases,
    extendAdminPurchaseAccess,
    cancelAdminPurchase,
} from '~/api/admin';
import type { AdminPurchaseItem } from '~/types';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const items = ref<AdminPurchaseItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const page = ref(1);
const perPage = ref(20);
const total = ref(0);
const searchQuery = ref('');
const statusFilter = ref<string | undefined>(undefined);
const productTypeFilter = ref<string | undefined>(undefined);

const statusOptions = [
    { label: 'Активен', value: 'active' },
    { label: 'Истёк', value: 'expired' },
    { label: 'Отменён', value: 'cancelled' },
];

const productTypeOptions = [
    { label: 'Курсы', value: 'course' },
    { label: 'Мастер-классы', value: 'masterclass' },
];

// Extend modal
const extendModalOpen = ref(false);
const extendTarget = ref<AdminPurchaseItem | null>(null);
const extendDays = ref<number>(30);
const extending = ref(false);

// Cancel modal
const cancelModalOpen = ref(false);
const cancelTarget = ref<AdminPurchaseItem | null>(null);
const cancelling = ref(false);

function purchaseStatusColor(status: string): string {
    const colors: Record<string, string> = {
        active: 'success',
        expired: 'neutral',
        cancelled: 'error',
    };
    return colors[status] || 'neutral';
}

function purchaseStatusLabel(status: string): string {
    const labels: Record<string, string> = {
        active: 'Активен',
        expired: 'Истёк',
        cancelled: 'Отменён',
    };
    return labels[status] || status;
}

let searchTimeout: ReturnType<typeof setTimeout> | null = null;
watch(searchQuery, () => {
    if (searchTimeout) clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
        page.value = 1;
        loadData();
    }, 400);
});

async function loadData() {
    loading.value = true;
    error.value = null;
    try {
        const response = await fetchAdminPurchases({
            page: page.value,
            per_page: perPage.value,
            search: searchQuery.value || undefined,
            status: statusFilter.value || undefined,
            product_type: productTypeFilter.value || undefined,
        });
        items.value = response.items;
        total.value = response.total;
        page.value = response.page;
        perPage.value = response.per_page;
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка загрузки';
    } finally {
        loading.value = false;
    }
}

function openExtendModal(item: AdminPurchaseItem) {
    extendTarget.value = item;
    extendDays.value = 30;
    extendModalOpen.value = true;
}

async function doExtend() {
    if (!extendTarget.value || !extendDays.value || extendDays.value < 1) return;
    extending.value = true;
    try {
        await extendAdminPurchaseAccess(extendTarget.value.id, extendDays.value);
        extendModalOpen.value = false;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка продления доступа';
    } finally {
        extending.value = false;
    }
}

function confirmCancel(item: AdminPurchaseItem) {
    cancelTarget.value = item;
    cancelModalOpen.value = true;
}

async function doCancel() {
    if (!cancelTarget.value) return;
    cancelling.value = true;
    try {
        await cancelAdminPurchase(cancelTarget.value.id);
        cancelModalOpen.value = false;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка отмены доступа';
    } finally {
        cancelling.value = false;
    }
}

// Initial load
loadData();
</script>

<style scoped>
@import '../_shared.css';

.admin-page__subtitle {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    margin: 4px 0 0;
}

.admin-page__filters-card {
    margin-bottom: 16px;
}

.admin-page__filters {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}

.admin-page__search {
    flex: 1;
    min-width: 200px;
}

.admin-page__filter-select {
    min-width: 160px;
}

.admin-page__loading {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.admin-page__skeleton-row {
    padding: 12px 0;
}

.admin-page__skeleton-lines {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.admin-page__skeleton-line {
    height: 14px;
    background: linear-gradient(
        90deg,
        var(--admin-border, #e0e0e0) 25%,
        transparent 50%,
        var(--admin-border, #e0e0e0) 75%
    );
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-radius: 4px;
}

@keyframes shimmer {
    0% {
        background-position: 200% 0;
    }
    100% {
        background-position: -200% 0;
    }
}

.admin-page__error-card {
    margin-bottom: 16px;
}

.admin-page__error {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 40px 20px;
    text-align: center;
    color: var(--admin-text-secondary, #636e72);
}

.admin-page__error-icon {
    font-size: 40px;
    color: #e17055;
}

.admin-page__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 60px 20px;
    text-align: center;
}

.admin-page__empty-icon {
    font-size: 48px;
    color: var(--admin-text-secondary, #636e72);
    opacity: 0.5;
}

.admin-page__empty-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin: 0;
}

.admin-page__empty-desc {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    max-width: 400px;
    margin: 0;
}

.admin-page__table-card {
    margin-bottom: 16px;
}

.admin-page__pagination {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
}

.admin-page__pagination-info {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

/* Purchases list */
.admin-purchases__list {
    display: flex;
    flex-direction: column;
}

.admin-purchases__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
}

.admin-purchases__item:last-child {
    border-bottom: none;
}

.admin-purchases__item-main {
    display: flex;
    align-items: center;
    gap: 20px;
    flex: 1;
    flex-wrap: wrap;
}

.admin-purchases__item-user {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 160px;
}

.admin-purchases__item-user-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-purchases__item-username {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-purchases__item-email {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-purchases__item-product {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 180px;
    flex: 1;
}

.admin-purchases__item-product-title {
    font-size: 14px;
    font-weight: 500;
    color: var(--admin-text-primary, #2d3436);
}

.admin-purchases__item-price {
    min-width: 80px;
}

.admin-purchases__item-price-value {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    white-space: nowrap;
}

.admin-purchases__item-status {
    min-width: 80px;
}

.admin-purchases__item-dates {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 120px;
}

.admin-purchases__item-date,
.admin-purchases__item-days {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-purchases__item-date-icon {
    font-size: 14px;
    flex-shrink: 0;
}

.admin-purchases__item-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
}

@media (max-width: 768px) {
    .admin-purchases__item {
        flex-direction: column;
        align-items: flex-start;
    }

    .admin-purchases__item-main {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
    }

    .admin-purchases__item-actions {
        align-self: flex-end;
    }

    .admin-page__filters {
        flex-direction: column;
    }

    .admin-page__filter-select {
        width: 100%;
    }
}
</style>
