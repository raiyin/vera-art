<template>
    <div class="admin-page">
        <!-- Breadcrumbs -->
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Мастер-классы</span>
        </div>

        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Мастер-классы</h1>
                <p class="admin-page__subtitle">Управление мастер-классами</p>
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
                <UButton
                    icon="i-lucide-plus"
                    color="primary"
                    to="/admin/master-classes/add"
                >
                    Создать МК
                </UButton>
            </div>
        </div>

        <!-- Search & Filters -->
        <UCard class="admin-page__filters-card" :ui="{ body: 'p-4' }">
            <div class="admin-page__filters">
                <div class="admin-page__search">
                    <UInput
                        v-model="searchQuery"
                        placeholder="Поиск по названию..."
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
                    color="neutral"
                    variant="outline"
                    class="admin-page__filter-select"
                    @change="loadData"
                />
                <USelect
                    v-model="difficultyFilter"
                    :items="difficultyOptions"
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
                    <div class="admin-page__skeleton-image" />
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
                <UIcon name="i-lucide-video" class="admin-page__empty-icon" />
                <h3 class="admin-page__empty-title">Мастер-классы не найдены</h3>
                <p class="admin-page__empty-desc">
                    По заданным критериям ничего не найдено
                </p>
                <UButton color="primary" to="/admin/master-classes/add">
                    Создать первый МК
                </UButton>
            </div>
        </UCard>

        <!-- Data Table -->
        <UCard v-else class="admin-page__table-card">
            <div class="admin-page__table-wrapper">
                <table class="admin-page__table">
                    <thead>
                        <tr>
                            <th class="admin-page__cell admin-page__cell--checkbox">
                                <UCheckbox
                                    :model-value="allSelected"
                                    :indeterminate="someSelected"
                                    @change="toggleSelectAll"
                                />
                            </th>
                            <th class="admin-page__cell admin-page__cell--head">ID</th>
                            <th
                                class="admin-page__cell admin-page__cell--head admin-page__cell--sortable"
                                @click="toggleSort('title_ru')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Название</span>
                                    <UIcon
                                        v-if="sortBy === 'title_ru'"
                                        :name="
                                            sortDir === 'asc'
                                                ? 'i-lucide-arrow-up'
                                                : 'i-lucide-arrow-down'
                                        "
                                        class="size-3"
                                    />
                                </div>
                            </th>
                            <th
                                class="admin-page__cell admin-page__cell--head admin-page__cell--sortable"
                                @click="toggleSort('price')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Цена</span>
                                    <UIcon
                                        v-if="sortBy === 'price'"
                                        :name="
                                            sortDir === 'asc'
                                                ? 'i-lucide-arrow-up'
                                                : 'i-lucide-arrow-down'
                                        "
                                        class="size-3"
                                    />
                                </div>
                            </th>
                            <th class="admin-page__cell admin-page__cell--head">
                                Статус
                            </th>
                            <th class="admin-page__cell admin-page__cell--head">
                                Сложность
                            </th>
                            <th class="admin-page__cell admin-page__cell--head">
                                Категория
                            </th>
                            <th class="admin-page__cell admin-page__cell--head">
                                Просмотры
                            </th>
                            <th
                                class="admin-page__cell admin-page__cell--head admin-page__cell--actions"
                            >
                                Действия
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                            v-for="item in items"
                            :key="item.id"
                            class="admin-page__row"
                            :class="{ 'admin-page__row--selected': isSelected(item.id) }"
                        >
                            <td class="admin-page__cell admin-page__cell--checkbox">
                                <UCheckbox
                                    :model-value="isSelected(item.id)"
                                    @change="toggleSelect(item.id)"
                                />
                            </td>
                            <td class="admin-page__cell admin-page__cell--mono">
                                #{{ item.id }}
                            </td>
                            <td class="admin-page__cell">
                                <div class="admin-page__name-cell">
                                    <span class="admin-page__name-ru">{{
                                        item.title_ru || '—'
                                    }}</span>
                                    <span class="admin-page__name-en">{{
                                        item.title_en || '—'
                                    }}</span>
                                </div>
                            </td>
                            <td class="admin-page__cell">
                                <span class="admin-page__price">{{
                                    formatPrice(item.price)
                                }}</span>
                            </td>
                            <td class="admin-page__cell">
                                <UBadge
                                    :color="statusColor(item.status)"
                                    variant="soft"
                                    size="sm"
                                >
                                    {{ statusLabel(item.status) }}
                                </UBadge>
                            </td>
                            <td class="admin-page__cell">
                                <UBadge
                                    :color="difficultyColor(item.difficulty)"
                                    variant="subtle"
                                    size="sm"
                                >
                                    {{ difficultyLabel(item.difficulty) }}
                                </UBadge>
                            </td>
                            <td class="admin-page__cell">
                                {{ item.category_name_ru || '—' }}
                            </td>
                            <td class="admin-page__cell admin-page__cell--mono">
                                {{ item.view_count }}
                            </td>
                            <td class="admin-page__cell admin-page__cell--actions">
                                <div class="admin-page__actions">
                                    <UTooltip text="Редактировать">
                                        <UButton
                                            icon="i-lucide-pencil"
                                            color="neutral"
                                            variant="ghost"
                                            size="sm"
                                            :to="`/admin/master-classes/edit/${item.id}`"
                                        />
                                    </UTooltip>
                                    <UTooltip text="Удалить">
                                        <UButton
                                            icon="i-lucide-trash-2"
                                            color="error"
                                            variant="ghost"
                                            size="sm"
                                            @click="confirmDelete(item)"
                                        />
                                    </UTooltip>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Bulk Actions & Pagination -->
            <template #footer>
                <div class="admin-page__table-footer">
                    <div class="admin-page__bulk-actions">
                        <UButton
                            v-if="selectedIds.length > 0"
                            color="error"
                            variant="outline"
                            size="sm"
                            :loading="deleting"
                            @click="confirmBulkDelete"
                        >
                            Удалить выбранные ({{ selectedIds.length }})
                        </UButton>
                    </div>
                    <div class="admin-page__pagination">
                        <span class="admin-page__pagination-info">
                            {{ paginationInfo }}
                        </span>
                        <UPagination
                            v-if="totalPages > 1"
                            v-model="currentPage"
                            :total="total"
                            :page-size="perPage"
                            :max="5"
                            size="sm"
                            @update:model-value="onPageChange"
                        />
                    </div>
                </div>
            </template>
        </UCard>

        <!-- Delete Confirmation Modal -->
        <UModal v-model="showDeleteModal">
            <UCard>
                <template #header>
                    <h3 class="text-lg font-semibold">
                        {{ deletingSingle ? 'Удаление МК' : 'Удаление мастер-классов' }}
                    </h3>
                </template>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                    {{ deleteConfirmMessage }}
                </p>
                <template #footer>
                    <div class="flex justify-end gap-3">
                        <UButton
                            color="neutral"
                            variant="outline"
                            @click="showDeleteModal = false"
                        >
                            Отмена
                        </UButton>
                        <UButton color="error" :loading="deleting" @click="executeDelete">
                            Удалить
                        </UButton>
                    </div>
                </template>
            </UCard>
        </UModal>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, } from 'vue';
import {
    fetchAdminProducts,
    deleteAdminProducts,
    type AdminProductItem,
} from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

// State
const items = ref<AdminProductItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const searchQuery = ref('');
const statusFilter = ref<string | null>(null);
const difficultyFilter = ref<string | null>(null);
const currentPage = ref(1);
const perPage = ref(20);
const total = ref(0);
const totalPages = ref(0);
const sortBy = ref('created_at');
const sortDir = ref<'asc' | 'desc'>('desc');
const selectedIds = ref<number[]>([]);
const showDeleteModal = ref(false);
const deleting = ref(false);
const deletingSingle = ref<AdminProductItem | null>(null);

// Filter options
const statusOptions = [
    { label: 'Все статусы', value: null },
    { label: 'Опубликован', value: 'published' },
    { label: 'Черновик', value: 'draft' },
    { label: 'Архивный', value: 'archived' },
];

const difficultyOptions = [
    { label: 'Любая сложность', value: null },
    { label: 'Начальный', value: 'beginner' },
    { label: 'Средний', value: 'intermediate' },
    { label: 'Продвинутый', value: 'advanced' },
];

// Computed
const allSelected = computed(() => {
    if (!items.value.length) return false;
    return items.value.every((item) => selectedIds.value.includes(item.id));
});

const someSelected = computed(() => {
    if (!items.value.length) return false;
    return items.value.some((item) => selectedIds.value.includes(item.id)) && !allSelected.value;
});

const paginationInfo = computed(() => {
    const start = (currentPage.value - 1) * perPage.value + 1;
    const end = Math.min(currentPage.value * perPage.value, total.value);
    return `${start}–${end} из ${total.value}`;
});

const deleteConfirmMessage = computed(() => {
    if (deletingSingle.value) {
        return `Вы уверены, что хотите удалить мастер-класс «${deletingSingle.value.title_ru || deletingSingle.value.title_en}»? Это действие нельзя отменить.`;
    }
    return `Вы уверены, что хотите удалить ${selectedIds.value.length} мастер-класс(ов)? Это действие нельзя отменить.`;
});

// Methods
function formatPrice(priceKopecks: number): string {
    const rubles = priceKopecks / 100;
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(rubles);
}

function statusColor(status: string): 'success' | 'neutral' | 'error' {
    switch (status) {
        case 'published': return 'success';
        case 'draft': return 'neutral';
        case 'archived': return 'error';
        default: return 'neutral';
    }
}

function statusLabel(status: string): string {
    switch (status) {
        case 'published': return 'Опубликован';
        case 'draft': return 'Черновик';
        case 'archived': return 'Архивный';
        default: return status || '—';
    }
}

function difficultyColor(difficulty: string): 'success' | 'warning' | 'error' | 'neutral' {
    switch (difficulty) {
        case 'beginner': return 'success';
        case 'intermediate': return 'warning';
        case 'advanced': return 'error';
        default: return 'neutral';
    }
}

function difficultyLabel(difficulty: string): string {
    switch (difficulty) {
        case 'beginner': return 'Начальный';
        case 'intermediate': return 'Средний';
        case 'advanced': return 'Продвинутый';
        default: return difficulty || '—';
    }
}

function isSelected(id: number): boolean {
    return selectedIds.value.includes(id);
}

function toggleSelect(id: number) {
    const idx = selectedIds.value.indexOf(id);
    if (idx === -1) {
        selectedIds.value.push(id);
    } else {
        selectedIds.value.splice(idx, 1);
    }
}

function toggleSelectAll() {
    if (allSelected.value) {
        selectedIds.value = [];
    } else {
        selectedIds.value = items.value.map((item) => item.id);
    }
}

function toggleSort(field: string) {
    if (sortBy.value === field) {
        sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc';
    } else {
        sortBy.value = field;
        sortDir.value = 'asc';
    }
    loadData();
}

function onPageChange(page: number) {
    currentPage.value = page;
    loadData();
}

async function loadData() {
    loading.value = true;
    error.value = null;
    try {
        const result = await fetchAdminProducts({
            page: currentPage.value,
            per_page: perPage.value,
            search: searchQuery.value || undefined,
            type: 'masterclass',
            status: statusFilter.value || undefined,
            difficulty: difficultyFilter.value || undefined,
            sort_by: sortBy.value,
            sort_dir: sortDir.value,
        });
        items.value = result.items;
        total.value = result.total;
        totalPages.value = result.total_pages;
        selectedIds.value = [];
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка загрузки данных';
    } finally {
        loading.value = false;
    }
}

function confirmDelete(item: AdminProductItem) {
    deletingSingle.value = item;
    showDeleteModal.value = true;
}

function confirmBulkDelete() {
    deletingSingle.value = null;
    showDeleteModal.value = true;
}

async function executeDelete() {
    deleting.value = true;
    try {
        if (deletingSingle.value) {
            await deleteAdminProducts([deletingSingle.value.id]);
        } else {
            await deleteAdminProducts(selectedIds.value);
        }
        showDeleteModal.value = false;
        selectedIds.value = [];
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка при удалении';
    } finally {
        deleting.value = false;
        deletingSingle.value = null;
    }
}

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout>;
watch(searchQuery, () => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
        currentPage.value = 1;
        loadData();
    }, 400);
});

onMounted(async () => {
    await loadData();
});
</script>

<style scoped>
@import '../_shared.css';

.admin-page__subtitle {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    margin: 4px 0 0;
}

.admin-page__header-actions {
    display: flex;
    gap: 8px;
}

.admin-page__filters-card {
    margin-bottom: 16px;
}

.admin-page__filters {
    display: flex;
    gap: 12px;
    align-items: center;
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
    display: flex;
    gap: 16px;
    align-items: center;
}

.admin-page__skeleton-image {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    background: var(--admin-border, #e0e0e0);
    animation: pulse 1.5s ease-in-out infinite;
}

.admin-page__skeleton-lines {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.admin-page__skeleton-line {
    height: 12px;
    border-radius: 4px;
    background: var(--admin-border, #e0e0e0);
    animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
    0%,
    100% {
        opacity: 1;
    }
    50% {
        opacity: 0.5;
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
    width: 48px;
    height: 48px;
    color: #e17055;
}

.admin-page__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 60px 20px;
    text-align: center;
}

.admin-page__empty-icon {
    width: 48px;
    height: 48px;
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
    margin: 0;
}

.admin-page__table-card {
    overflow: hidden;
}

.admin-page__table-wrapper {
    overflow-x: auto;
}

.admin-page__table {
    width: 100%;
    border-collapse: collapse;
    font-size: 14px;
}

.admin-page__cell {
    padding: 12px 16px;
    color: var(--admin-text-primary, #2d3436);
    vertical-align: middle;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
}

.admin-page__cell--head {
    padding: 10px 16px;
    font-weight: 600;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--admin-text-secondary, #636e72);
    white-space: nowrap;
    background: var(--admin-bg, #f0f2f5);
    user-select: none;
}

.admin-page__cell--sortable {
    cursor: pointer;
    transition: color 0.2s;
}

.admin-page__cell--sortable:hover {
    color: var(--admin-primary, #6c5ce7);
}

.admin-page__cell--checkbox {
    width: 48px;
    text-align: center;
}

.admin-page__cell--mono {
    font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
}

.admin-page__cell--actions {
    text-align: right;
    white-space: nowrap;
}

.admin-page__head-content {
    display: flex;
    align-items: center;
    gap: 4px;
}

.admin-page__row {
    transition: background 0.15s ease;
}

.admin-page__row:hover {
    background: rgba(108, 92, 231, 0.03);
}

.admin-page__row--selected {
    background: rgba(108, 92, 231, 0.06) !important;
}

.admin-page__name-cell {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-page__name-ru {
    font-weight: 500;
}

.admin-page__name-en {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-page__price {
    font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', monospace;
    font-weight: 600;
    font-size: 14px;
    color: var(--admin-primary, #6c5ce7);
    white-space: nowrap;
}

.admin-page__actions {
    display: flex;
    gap: 4px;
    justify-content: flex-end;
}

.admin-page__table-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
}

.admin-page__bulk-actions {
    flex: 1;
}

.admin-page__pagination {
    display: flex;
    align-items: center;
    gap: 12px;
}

.admin-page__pagination-info {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
    white-space: nowrap;
}
</style>
