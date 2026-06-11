<template>
    <div class="admin-page">
        <!-- Breadcrumbs -->
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Пользователи</span>
        </div>

        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Пользователи</h1>
                <p class="admin-page__subtitle">Управление пользователями платформы</p>
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
        <UCard class="admin-page__filters-card" :ui="{ body: 'p-4' }">
            <div class="admin-page__filters">
                <div class="admin-page__search">
                    <UInput
                        v-model="searchQuery"
                        placeholder="Поиск по имени, email или username..."
                        icon="i-lucide-search"
                        color="neutral"
                        variant="outline"
                        class="w-full"
                    />
                </div>
                <USelect
                    v-model="roleFilter"
                    :items="roleOptions"
                    color="neutral"
                    variant="outline"
                    class="admin-page__filter-select"
                    @change="loadData"
                />
                <USelect
                    v-model="blockedFilter"
                    :items="blockedOptions"
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
                <UIcon name="i-lucide-users" class="admin-page__empty-icon" />
                <h3 class="admin-page__empty-title">Пользователи не найдены</h3>
                <p class="admin-page__empty-desc">
                    <template v-if="searchQuery || roleFilter || blockedFilter">
                        По заданным критериям ничего не найдено
                    </template>
                    <template v-else>
                        На платформе пока нет зарегистрированных пользователей
                    </template>
                </p>
            </div>
        </UCard>

        <!-- Data Table -->
        <UCard v-else class="admin-page__table-card">
            <div class="admin-page__table-wrapper">
                <table class="admin-page__table">
                    <thead>
                        <tr>
                            <th class="admin-page__cell admin-page__cell--head">ID</th>
                            <th
                                class="admin-page__cell admin-page__cell--head admin-page__cell--sortable"
                                @click="toggleSort('u.username')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Username</span>
                                    <UIcon
                                        v-if="sortBy === 'u.username'"
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
                                @click="toggleSort('u.full_name')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Имя</span>
                                    <UIcon
                                        v-if="sortBy === 'u.full_name'"
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
                                @click="toggleSort('u.email')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Email</span>
                                    <UIcon
                                        v-if="sortBy === 'u.email'"
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
                                @click="toggleSort('u.role')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Роль</span>
                                    <UIcon
                                        v-if="sortBy === 'u.role'"
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
                                Покупки
                            </th>
                            <th class="admin-page__cell admin-page__cell--head">
                                Отзывы
                            </th>
                            <th
                                class="admin-page__cell admin-page__cell--head admin-page__cell--sortable"
                                @click="toggleSort('u.created_at')"
                            >
                                <div class="admin-page__head-content">
                                    <span>Дата рег.</span>
                                    <UIcon
                                        v-if="sortBy === 'u.created_at'"
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
                            :class="{ 'admin-page__row--blocked': item.blocked }"
                        >
                            <td class="admin-page__cell admin-page__cell--mono">
                                #{{ item.id }}
                            </td>
                            <td class="admin-page__cell">
                                <div class="admin-page__name-cell">
                                    <span class="admin-page__name-ru">{{
                                        item.username
                                    }}</span>
                                </div>
                            </td>
                            <td class="admin-page__cell">
                                {{ item.full_name || '—' }}
                            </td>
                            <td class="admin-page__cell">
                                {{ item.email || '—' }}
                            </td>
                            <td class="admin-page__cell">
                                <UBadge
                                    :color="item.role === 'admin' ? 'warning' : 'neutral'"
                                    variant="subtle"
                                    size="sm"
                                >
                                    {{ item.role === 'admin' ? 'Админ' : 'Пользователь' }}
                                </UBadge>
                            </td>
                            <td class="admin-page__cell">
                                <UBadge
                                    v-if="item.blocked"
                                    color="error"
                                    variant="soft"
                                    size="sm"
                                >
                                    Заблокирован
                                </UBadge>
                                <UBadge v-else color="success" variant="soft" size="sm">
                                    Активен
                                </UBadge>
                            </td>
                            <td class="admin-page__cell admin-page__cell--mono">
                                {{ item.purchases_count }}
                            </td>
                            <td class="admin-page__cell admin-page__cell--mono">
                                {{ item.reviews_count }}
                            </td>
                            <td class="admin-page__cell admin-page__cell--mono">
                                {{ formatDate(item.created_at) }}
                            </td>
                            <td class="admin-page__cell admin-page__cell--actions">
                                <div class="admin-page__actions">
                                    <UTooltip text="Просмотр">
                                        <UButton
                                            icon="i-lucide-eye"
                                            color="neutral"
                                            variant="ghost"
                                            size="sm"
                                            @click="viewUser(item)"
                                        />
                                    </UTooltip>
                                    <UTooltip
                                        :text="
                                            item.role === 'admin'
                                                ? 'Снять админа'
                                                : 'Назначить админом'
                                        "
                                    >
                                        <UButton
                                            :icon="
                                                item.role === 'admin'
                                                    ? 'i-lucide-shield-off'
                                                    : 'i-lucide-shield'
                                            "
                                            color="neutral"
                                            variant="ghost"
                                            size="sm"
                                            @click="confirmRoleChange(item)"
                                        />
                                    </UTooltip>
                                    <UTooltip
                                        :text="
                                            item.blocked
                                                ? 'Разблокировать'
                                                : 'Заблокировать'
                                        "
                                    >
                                        <UButton
                                            :icon="
                                                item.blocked
                                                    ? 'i-lucide-unlock'
                                                    : 'i-lucide-lock'
                                            "
                                            :color="item.blocked ? 'success' : 'error'"
                                            variant="ghost"
                                            size="sm"
                                            @click="confirmBlockToggle(item)"
                                        />
                                    </UTooltip>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Pagination -->
            <template #footer>
                <div class="admin-page__table-footer">
                    <div class="admin-page__bulk-actions" />
                    <div class="admin-page__pagination">
                        <span class="admin-page__pagination-info">
                            {{ paginationInfo }}
                        </span>
                        <UPagination
                            v-if="totalPages > 1"
                            v-model:page="currentPage"
                            :total="total"
                            :items-per-page="perPage"
                            :max="5"
                            size="sm"
                            @update:page="onPageChange"
                        />
                    </div>
                </div>
            </template>
        </UCard>

        <!-- User Detail Modal -->
        <UModal v-model="showDetailModal" class="max-w-3xl">
            <UCard v-if="selectedUser" :ui="{ body: 'p-0' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-lg font-semibold">
                            Пользователь: {{ selectedUser.username }}
                        </h3>
                        <UBadge
                            :color="selectedUser.blocked ? 'error' : 'success'"
                            variant="soft"
                            size="sm"
                        >
                            {{ selectedUser.blocked ? 'Заблокирован' : 'Активен' }}
                        </UBadge>
                    </div>
                </template>

                <div class="p-6 space-y-6">
                    <!-- User Info -->
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <span class="text-sm text-gray-500 dark:text-gray-400"
                                >ID</span
                            >
                            <p class="font-medium">#{{ selectedUser.id }}</p>
                        </div>
                        <div>
                            <span class="text-sm text-gray-500 dark:text-gray-400"
                                >Username</span
                            >
                            <p class="font-medium">{{ selectedUser.username }}</p>
                        </div>
                        <div>
                            <span class="text-sm text-gray-500 dark:text-gray-400"
                                >Полное имя</span
                            >
                            <p class="font-medium">{{ selectedUser.full_name || '—' }}</p>
                        </div>
                        <div>
                            <span class="text-sm text-gray-500 dark:text-gray-400"
                                >Email</span
                            >
                            <p class="font-medium">{{ selectedUser.email || '—' }}</p>
                        </div>
                        <div>
                            <span class="text-sm text-gray-500 dark:text-gray-400"
                                >Роль</span
                            >
                            <p class="font-medium">
                                <UBadge
                                    :color="
                                        selectedUser.role === 'admin'
                                            ? 'warning'
                                            : 'neutral'
                                    "
                                    variant="subtle"
                                    size="sm"
                                >
                                    {{
                                        selectedUser.role === 'admin'
                                            ? 'Админ'
                                            : 'Пользователь'
                                    }}
                                </UBadge>
                            </p>
                        </div>
                        <div>
                            <span class="text-sm text-gray-500 dark:text-gray-400">
                                Дата регистрации
                            </span>
                            <p class="font-medium">
                                {{ formatDate(selectedUser.created_at) }}
                            </p>
                        </div>
                    </div>

                    <!-- Purchases Section -->
                    <div>
                        <h4 class="text-md font-semibold mb-3 flex items-center gap-2">
                            <UIcon name="i-lucide-shopping-cart" class="size-4" />
                            Покупки ({{ selectedUser.purchases.length }})
                        </h4>
                        <div
                            v-if="selectedUser.purchases.length === 0"
                            class="text-sm text-gray-500 dark:text-gray-400 py-2"
                        >
                            Нет покупок
                        </div>
                        <div v-else class="space-y-2">
                            <div
                                v-for="purchase in selectedUser.purchases"
                                :key="purchase.id"
                                class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
                            >
                                <div class="flex-1 min-w-0">
                                    <p class="text-sm font-medium truncate">
                                        {{
                                            purchase.product_title_ru ||
                                            purchase.product_title_en ||
                                            '—'
                                        }}
                                    </p>
                                    <p class="text-xs text-gray-500 dark:text-gray-400">
                                        {{ formatDate(purchase.purchase_date) }}
                                        <span v-if="purchase.product_type" class="ml-2">
                                            ·
                                            {{
                                                purchase.product_type === 'course'
                                                    ? 'Курс'
                                                    : 'Мастер-класс'
                                            }}
                                        </span>
                                    </p>
                                </div>
                                <div class="flex items-center gap-3 ml-4">
                                    <span class="text-sm font-medium">
                                        {{ formatPrice(purchase.price_paid) }}
                                    </span>
                                    <UBadge
                                        :color="purchaseStatusColor(purchase.status)"
                                        variant="soft"
                                        size="sm"
                                    >
                                        {{ purchaseStatusLabel(purchase.status) }}
                                    </UBadge>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Reviews Section -->
                    <div>
                        <h4 class="text-md font-semibold mb-3 flex items-center gap-2">
                            <UIcon name="i-lucide-star" class="size-4" />
                            Отзывы ({{ selectedUser.reviews.length }})
                        </h4>
                        <div
                            v-if="selectedUser.reviews.length === 0"
                            class="text-sm text-gray-500 dark:text-gray-400 py-2"
                        >
                            Нет отзывов
                        </div>
                        <div v-else class="space-y-2">
                            <div
                                v-for="review in selectedUser.reviews"
                                :key="review.id"
                                class="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
                            >
                                <div class="flex items-center justify-between mb-1">
                                    <p class="text-sm font-medium truncate">
                                        {{
                                            review.product_title_ru ||
                                            review.product_title_en ||
                                            '—'
                                        }}
                                    </p>
                                    <div class="flex items-center gap-2">
                                        <div class="flex items-center">
                                            <UIcon
                                                v-for="s in 5"
                                                :key="s"
                                                :name="
                                                    s <= review.rating
                                                        ? 'i-lucide-star'
                                                        : 'i-lucide-star'
                                                "
                                                class="size-3"
                                                :class="
                                                    s <= review.rating
                                                        ? 'text-yellow-400'
                                                        : 'text-gray-300 dark:text-gray-600'
                                                "
                                            />
                                        </div>
                                        <UBadge
                                            v-if="review.is_approved"
                                            color="success"
                                            variant="soft"
                                            size="sm"
                                        >
                                            Одобрен
                                        </UBadge>
                                        <UBadge
                                            v-else
                                            color="warning"
                                            variant="soft"
                                            size="sm"
                                        >
                                            На модерации
                                        </UBadge>
                                    </div>
                                </div>
                                <p class="text-xs text-gray-500 dark:text-gray-400">
                                    {{ formatDate(review.created_at) }}
                                </p>
                                <p
                                    v-if="review.comment_ru || review.comment_en"
                                    class="text-sm mt-1"
                                >
                                    {{ review.comment_ru || review.comment_en }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <template #footer>
                    <div class="flex justify-end gap-3">
                        <UButton
                            color="neutral"
                            variant="outline"
                            @click="showDetailModal = false"
                        >
                            Закрыть
                        </UButton>
                    </div>
                </template>
            </UCard>
        </UModal>

        <!-- Role Change Confirmation Modal -->
        <AdminConfirmDialog
            :visible="showRoleModal"
            :title="
                roleTarget?.role === 'admin'
                    ? 'Снять права администратора'
                    : 'Назначить администратором'
            "
            :message="
                roleTarget?.role === 'admin'
                    ? `Вы уверены, что хотите снять права администратора с пользователя «${roleTarget?.username}»?`
                    : `Вы уверены, что хотите назначить пользователя «${roleTarget?.username}» администратором?`
            "
            :type="roleTarget?.role === 'admin' ? 'warning' : 'info'"
            :confirm-text="roleTarget?.role === 'admin' ? 'Снять права' : 'Назначить'"
            cancel-text="Отмена"
            loading-text="Выполнение..."
            :loading="updatingRole"
            @confirm="executeRoleChange"
            @cancel="showRoleModal = false"
            @update:visible="showRoleModal = $event"
        />

        <!-- Block/Unblock Confirmation Modal -->
        <AdminConfirmDialog
            :visible="showBlockModal"
            :title="
                blockTarget?.blocked
                    ? 'Разблокировать пользователя'
                    : 'Заблокировать пользователя'
            "
            :message="
                blockTarget?.blocked
                    ? `Вы уверены, что хотите разблокировать пользователя «${blockTarget?.username}»? Он снова сможет войти в систему.`
                    : `Вы уверены, что хотите заблокировать пользователя «${blockTarget?.username}»? Он не сможет войти в систему.`
            "
            :type="blockTarget?.blocked ? 'info' : 'danger'"
            :confirm-text="blockTarget?.blocked ? 'Разблокировать' : 'Заблокировать'"
            cancel-text="Отмена"
            loading-text="Выполнение..."
            :loading="updatingBlock"
            @confirm="executeBlockToggle"
            @cancel="showBlockModal = false"
            @update:visible="showBlockModal = $event"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, } from 'vue';
import {
    fetchAdminUsers,
    fetchAdminUserDetail,
    updateAdminUserRole,
    toggleAdminUserBlock,
    type AdminUserItem,
    type AdminUserDetail,
} from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

// State
const items = ref<AdminUserItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);
const searchQuery = ref('');
const roleFilter = ref<string | null>(null);
const blockedFilter = ref<string | null>(null);
const currentPage = ref(1);
const perPage = ref(20);
const total = ref(0);
const totalPages = ref(0);
const sortBy = ref('u.id');
const sortDir = ref<'asc' | 'desc'>('desc');

// Detail modal
const showDetailModal = ref(false);
const selectedUser = ref<AdminUserDetail | null>(null);
const loadingDetail = ref(false);

// Role change
const showRoleModal = ref(false);
const roleTarget = ref<AdminUserItem | null>(null);
const updatingRole = ref(false);

// Block toggle
const showBlockModal = ref(false);
const blockTarget = ref<AdminUserItem | null>(null);
const updatingBlock = ref(false);

// Filter options
const roleOptions = [
    { label: 'Все роли', value: null },
    { label: 'Администраторы', value: 'admin' },
    { label: 'Пользователи', value: 'user' },
];

const blockedOptions = [
    { label: 'Все статусы', value: null },
    { label: 'Активные', value: 'false' },
    { label: 'Заблокированные', value: 'true' },
];

// Computed
const paginationInfo = computed(() => {
    const start = (currentPage.value - 1) * perPage.value + 1;
    const end = Math.min(currentPage.value * perPage.value, total.value);
    return `${start}–${end} из ${total.value}`;
});

// Methods
function formatDate(dateStr: string): string {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    });
}

function formatPrice(kopecks: number): string {
    if (!kopecks && kopecks !== 0) return '—';
    return `${(kopecks / 100).toLocaleString('ru-RU')} ₽`;
}

function purchaseStatusColor(status: string): 'success' | 'warning' | 'error' | 'neutral' {
    switch (status) {
        case 'active': return 'success';
        case 'expired': return 'warning';
        case 'cancelled': return 'error';
        default: return 'neutral';
    }
}

function purchaseStatusLabel(status: string): string {
    switch (status) {
        case 'active': return 'Активен';
        case 'expired': return 'Истёк';
        case 'cancelled': return 'Отменён';
        default: return status || '—';
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
        const result = await fetchAdminUsers({
            page: currentPage.value,
            per_page: perPage.value,
            search: searchQuery.value || undefined,
            role: roleFilter.value || undefined,
            blocked: blockedFilter.value || undefined,
            sort_by: sortBy.value,
            sort_dir: sortDir.value,
        });
        items.value = result.items;
        total.value = result.total;
        totalPages.value = result.total_pages;
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка загрузки данных';
    } finally {
        loading.value = false;
    }
}

async function viewUser(item: AdminUserItem) {
    loadingDetail.value = true;
    showDetailModal.value = true;
    selectedUser.value = null;
    try {
        const detail = await fetchAdminUserDetail(item.id);
        selectedUser.value = detail;
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка загрузки данных пользователя';
        showDetailModal.value = false;
    } finally {
        loadingDetail.value = false;
    }
}

function confirmRoleChange(item: AdminUserItem) {
    roleTarget.value = item;
    showRoleModal.value = true;
}

async function executeRoleChange() {
    if (!roleTarget.value) return;
    updatingRole.value = true;
    try {
        const newRole = roleTarget.value.role === 'admin' ? 'user' : 'admin';
        await updateAdminUserRole(roleTarget.value.id, newRole);
        showRoleModal.value = false;
        roleTarget.value = null;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка при изменении роли';
    } finally {
        updatingRole.value = false;
    }
}

function confirmBlockToggle(item: AdminUserItem) {
    blockTarget.value = item;
    showBlockModal.value = true;
}

async function executeBlockToggle() {
    if (!blockTarget.value) return;
    updatingBlock.value = true;
    try {
        await toggleAdminUserBlock(blockTarget.value.id, !blockTarget.value.blocked);
        showBlockModal.value = false;
        blockTarget.value = null;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка при изменении статуса';
    } finally {
        updatingBlock.value = false;
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
    max-width: 320px;
}

.admin-page__filter-select {
    min-width: 180px;
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
}

.admin-page__cell {
    padding: 12px 16px;
    text-align: left;
    font-size: 14px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    white-space: nowrap;
}

.admin-page__cell--head {
    font-weight: 600;
    color: var(--admin-text-secondary, #636e72);
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    background: var(--admin-surface-secondary, #f8f9fa);
    position: sticky;
    top: 0;
    z-index: 1;
}

.admin-page__cell--sortable {
    cursor: pointer;
    user-select: none;
}

.admin-page__cell--sortable:hover {
    color: var(--admin-primary, #6c5ce7);
}

.admin-page__cell--mono {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
}

.admin-page__cell--actions {
    text-align: right;
}

.admin-page__head-content {
    display: flex;
    align-items: center;
    gap: 4px;
}

.admin-page__row {
    transition: background 0.15s;
}

.admin-page__row:hover {
    background: var(--admin-hover, #f0f0f0);
}

.admin-page__row--blocked {
    opacity: 0.6;
}

.admin-page__name-cell {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-page__name-ru {
    font-weight: 500;
    color: var(--admin-text-primary, #2d3436);
}

.admin-page__name-en {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-page__actions {
    display: flex;
    gap: 4px;
    justify-content: flex-end;
}

.admin-page__table-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 4px 0;
}

.admin-page__pagination {
    display: flex;
    align-items: center;
    gap: 12px;
}

.admin-page__pagination-info {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-page__loading {
    display: flex;
    flex-direction: column;
    gap: 8px;
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
    background: linear-gradient(90deg, #e0e0e0 25%, #f0f0f0 50%, #e0e0e0 75%);
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
    margin-top: 16px;
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
    font-size: 48px;
    color: #e17055;
}

.admin-page__empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 60px 20px;
    text-align: center;
}

.admin-page__empty-icon {
    font-size: 48px;
    color: var(--admin-text-secondary, #636e72);
    margin-bottom: 16px;
}

.admin-page__empty-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin: 0 0 8px;
}

.admin-page__empty-desc {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    margin: 0 0 20px;
}

.admin-page__text-muted {
    color: var(--admin-text-secondary, #636e72);
    font-size: 13px;
}

@media (max-width: 768px) {
    .admin-page__filters {
        flex-direction: column;
        align-items: stretch;
    }

    .admin-page__search {
        max-width: none;
    }

    .admin-page__filter-select {
        min-width: auto;
    }
}
</style>
