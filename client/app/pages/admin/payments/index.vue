<template>
    <div class="admin-page">
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Платежи</span>
        </div>
        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Платежи</h1>
                <p class="admin-page__subtitle">Управление платежами и возвратами</p>
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
                        placeholder="Поиск по пользователю или описанию..."
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
                <UIcon name="i-lucide-credit-card" class="admin-page__empty-icon" />
                <h3 class="admin-page__empty-title">Платежи не найдены</h3>
                <p class="admin-page__empty-desc">
                    <template v-if="searchQuery || statusFilter">
                        По заданным критериям ничего не найдено. Попробуйте изменить
                        параметры поиска.
                    </template>
                    <template v-else> Пока нет ни одного платежа. </template>
                </p>
            </div>
        </UCard>

        <!-- Payments List -->
        <UCard v-else class="admin-page__table-card">
            <div class="admin-payments__list">
                <div v-for="item in items" :key="item.id" class="admin-payments__item">
                    <div class="admin-payments__item-main">
                        <div class="admin-payments__item-user">
                            <UAvatar
                                :text="((item.username || '?')[0] || '?').toUpperCase()"
                                size="sm"
                                color="neutral"
                            />
                            <div class="admin-payments__item-user-info">
                                <span class="admin-payments__item-username">{{
                                    item.username
                                }}</span>
                                <span class="admin-payments__item-email">{{
                                    item.user_email
                                }}</span>
                            </div>
                        </div>
                        <div class="admin-payments__item-amount">
                            <span class="admin-payments__item-amount-value"
                                >{{ (item.amount / 100).toLocaleString('ru-RU') }} ₽</span
                            >
                        </div>
                        <div class="admin-payments__item-status">
                            <UBadge
                                :color="paymentStatusColor(item.status) as 'success' | 'warning' | 'error' | 'neutral' | 'info'"
                                variant="soft"
                                size="sm"
                            >
                                {{ paymentStatusLabel(item.status) }}
                            </UBadge>
                        </div>
                        <div class="admin-payments__item-method">
                            <span
                                v-if="item.payment_method"
                                class="admin-payments__item-method-text"
                            >
                                {{ item.payment_method }}
                            </span>
                            <span
                                v-else
                                class="admin-payments__item-method-text admin-payments__item-method-text--empty"
                            >
                                —
                            </span>
                        </div>
                        <div class="admin-payments__item-date">
                            <UIcon
                                name="i-lucide-calendar"
                                class="admin-payments__item-date-icon"
                            />
                            <span>{{ item.created_at }}</span>
                        </div>
                    </div>
                    <div class="admin-payments__item-actions">
                        <UTooltip text="Детали платежа">
                            <UButton
                                icon="i-lucide-eye"
                                color="neutral"
                                variant="ghost"
                                size="sm"
                                @click="openDetailModal(item)"
                            />
                        </UTooltip>
                        <UTooltip text="Возврат платежа">
                            <UButton
                                icon="i-lucide-undo-2"
                                color="error"
                                variant="ghost"
                                size="sm"
                                :disabled="
                                    item.status !== 'succeeded' &&
                                    item.status !== 'waiting_for_capture'
                                "
                                @click="confirmRefund(item)"
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

        <!-- Detail Modal -->
        <UModal v-model:open="detailModalOpen" class="max-w-lg">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-semibold">
                        Детали платежа #{{ detailTarget?.id }}
                    </h3>
                    <UButton
                        icon="i-lucide-x"
                        color="neutral"
                        variant="ghost"
                        size="sm"
                        @click="detailModalOpen = false"
                    />
                </div>
            </template>
            <template #body>
                <div v-if="detailTarget" class="space-y-3">
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Пользователь</span>
                        <span class="admin-payments__detail-value"
                            >{{ detailTarget.username }} ({{
                                detailTarget.user_email
                            }})</span
                        >
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Сумма</span>
                        <span class="admin-payments__detail-value"
                            >{{ (detailTarget.amount / 100).toLocaleString('ru-RU') }}
                            {{ detailTarget.currency }}</span
                        >
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Статус</span>
                        <UBadge
                            :color="paymentStatusColor(detailTarget.status) as 'success' | 'warning' | 'error' | 'neutral' | 'info'"
                            variant="soft"
                            size="sm"
                        >
                            {{ paymentStatusLabel(detailTarget.status) }}
                        </UBadge>
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Способ оплаты</span>
                        <span class="admin-payments__detail-value">{{
                            detailTarget.payment_method || '—'
                        }}</span>
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Описание</span>
                        <span class="admin-payments__detail-value">{{
                            detailTarget.description || '—'
                        }}</span>
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">External ID</span>
                        <span class="admin-payments__detail-value">{{
                            detailTarget.external_id || '—'
                        }}</span>
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Создан</span>
                        <span class="admin-payments__detail-value">{{
                            detailTarget.created_at
                        }}</span>
                    </div>
                    <div class="admin-payments__detail-row">
                        <span class="admin-payments__detail-label">Обновлён</span>
                        <span class="admin-payments__detail-value">{{
                            detailTarget.updated_at
                        }}</span>
                    </div>
                </div>
            </template>
            <template #footer>
                <div class="flex justify-end">
                    <UButton
                        color="neutral"
                        variant="outline"
                        @click="detailModalOpen = false"
                    >
                        Закрыть
                    </UButton>
                </div>
            </template>
        </UModal>

        <!-- Refund Confirmation Modal -->
        <AdminConfirmDialog
            :visible="refundModalOpen"
            title="Подтверждение возврата"
            :message="
                refundTarget
                    ? `Вы уверены, что хотите инициировать возврат платежа #${
                          refundTarget.id
                      } на сумму ${(refundTarget.amount / 100).toLocaleString(
                          'ru-RU'
                      )} ₽? Доступ к продукту будет отменён.`
                    : ''
            "
            type="danger"
            confirm-text="Да, вернуть"
            cancel-text="Нет"
            loading-text="Выполнение возврата..."
            :loading="refunding"
            @confirm="doRefund"
            @cancel="refundModalOpen = false"
            @update:visible="refundModalOpen = $event"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, watch, } from 'vue';
import {
    fetchAdminPayments,
    refundAdminPayment,
} from '~/api/admin';
import type { AdminPaymentItem } from '~/types';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const items = ref<AdminPaymentItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const page = ref(1);
const perPage = ref(20);
const total = ref(0);
const searchQuery = ref('');
const statusFilter = ref<string | undefined>(undefined);

const statusOptions = [
    { label: 'Ожидает', value: 'pending' },
    { label: 'Успешен', value: 'succeeded' },
    { label: 'Отменён', value: 'canceled' },
    { label: 'Возвращён', value: 'refunded' },
    { label: 'Ожидает подтверждения', value: 'waiting_for_capture' },
];

// Detail modal
const detailModalOpen = ref(false);
const detailTarget = ref<AdminPaymentItem | null>(null);

// Refund modal
const refundModalOpen = ref(false);
const refundTarget = ref<AdminPaymentItem | null>(null);
const refunding = ref(false);

function paymentStatusColor(status: string): string {
    const colors: Record<string, string> = {
        pending: 'warning',
        waiting_for_capture: 'info',
        succeeded: 'success',
        canceled: 'error',
        refunded: 'neutral',
    };
    return colors[status] || 'neutral';
}

function paymentStatusLabel(status: string): string {
    const labels: Record<string, string> = {
        pending: 'Ожидает',
        waiting_for_capture: 'Ожидает подтверждения',
        succeeded: 'Успешен',
        canceled: 'Отменён',
        refunded: 'Возвращён',
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
        const response = await fetchAdminPayments({
            page: page.value,
            per_page: perPage.value,
            search: searchQuery.value || undefined,
            status: statusFilter.value || undefined,
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

function openDetailModal(item: AdminPaymentItem) {
    detailTarget.value = item;
    detailModalOpen.value = true;
}

function confirmRefund(item: AdminPaymentItem) {
    refundTarget.value = item;
    refundModalOpen.value = true;
}

async function doRefund() {
    if (!refundTarget.value) return;
    refunding.value = true;
    try {
        await refundAdminPayment(refundTarget.value.id);
        refundModalOpen.value = false;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка возврата платежа';
    } finally {
        refunding.value = false;
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

/* Payments list */
.admin-payments__list {
    display: flex;
    flex-direction: column;
}

.admin-payments__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
}

.admin-payments__item:last-child {
    border-bottom: none;
}

.admin-payments__item-main {
    display: flex;
    align-items: center;
    gap: 20px;
    flex: 1;
    flex-wrap: wrap;
}

.admin-payments__item-user {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 160px;
}

.admin-payments__item-user-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-payments__item-username {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-payments__item-email {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-payments__item-amount {
    min-width: 100px;
}

.admin-payments__item-amount-value {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    white-space: nowrap;
}

.admin-payments__item-status {
    min-width: 100px;
}

.admin-payments__item-method {
    min-width: 120px;
}

.admin-payments__item-method-text {
    font-size: 13px;
    color: var(--admin-text-primary, #2d3436);
}

.admin-payments__item-method-text--empty {
    color: var(--admin-text-secondary, #636e72);
}

.admin-payments__item-date {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    min-width: 120px;
}

.admin-payments__item-date-icon {
    font-size: 14px;
    flex-shrink: 0;
}

.admin-payments__item-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
}

/* Detail modal */
.admin-payments__detail-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
}

.admin-payments__detail-row:last-child {
    border-bottom: none;
}

.admin-payments__detail-label {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
    flex-shrink: 0;
}

.admin-payments__detail-value {
    font-size: 14px;
    color: var(--admin-text-primary, #2d3436);
    text-align: right;
    word-break: break-all;
}

@media (max-width: 768px) {
    .admin-payments__item {
        flex-direction: column;
        align-items: flex-start;
    }

    .admin-payments__item-main {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
    }

    .admin-payments__item-actions {
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
