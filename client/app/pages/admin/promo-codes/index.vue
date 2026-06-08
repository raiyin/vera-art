<template>
    <div class="admin-page">
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Промокоды</span>
        </div>
        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Промокоды</h1>
                <p class="admin-page__subtitle">Управление промокодами и скидками</p>
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
                <UButton icon="i-lucide-plus" color="primary" @click="openCreateModal">
                    Создать промокод
                </UButton>
            </div>
        </div>

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
                <UIcon name="i-lucide-ticket-percent" class="admin-page__empty-icon" />
                <h3 class="admin-page__empty-title">Промокоды не найдены</h3>
                <p class="admin-page__empty-desc">
                    Создайте первый промокод, чтобы начать.
                </p>
            </div>
        </UCard>

        <!-- Promo Codes List -->
        <UCard v-else class="admin-page__table-card">
            <div class="admin-promocodes__list">
                <div v-for="item in items" :key="item.id" class="admin-promocodes__item">
                    <div class="admin-promocodes__item-main">
                        <div class="admin-promocodes__item-code-section">
                            <span class="admin-promocodes__item-code">{{
                                item.code
                            }}</span>
                            <UBadge
                                :color="item.is_active ? 'success' : 'neutral'"
                                variant="soft"
                                size="sm"
                            >
                                {{ item.is_active ? 'Активен' : 'Неактивен' }}
                            </UBadge>
                        </div>
                        <div class="admin-promocodes__item-discount">
                            <span class="admin-promocodes__item-discount-value">
                                {{
                                    item.discount_type === 'percentage'
                                        ? `${item.discount_value}%`
                                        : `${(item.discount_value / 100).toLocaleString(
                                              'ru-RU'
                                          )} ₽`
                                }}
                            </span>
                            <span class="admin-promocodes__item-discount-type">
                                {{
                                    item.discount_type === 'percentage'
                                        ? 'Процент'
                                        : 'Фиксированная'
                                }}
                            </span>
                        </div>
                        <div class="admin-promocodes__item-usage">
                            <span class="admin-promocodes__item-usage-count">
                                {{ item.used_count
                                }}{{ item.max_uses ? ` / ${item.max_uses}` : '' }}
                            </span>
                            <span class="admin-promocodes__item-usage-label"
                                >использований</span
                            >
                        </div>
                        <div class="admin-promocodes__item-dates">
                            <div
                                v-if="item.valid_from"
                                class="admin-promocodes__item-date"
                            >
                                <UIcon
                                    name="i-lucide-calendar-start"
                                    class="admin-promocodes__item-date-icon"
                                />
                                <span>с {{ formatDate(item.valid_from) }}</span>
                            </div>
                            <div
                                v-if="item.valid_until"
                                class="admin-promocodes__item-date"
                            >
                                <UIcon
                                    name="i-lucide-calendar-end"
                                    class="admin-promocodes__item-date-icon"
                                />
                                <span>до {{ formatDate(item.valid_until) }}</span>
                            </div>
                            <div
                                v-if="!item.valid_from && !item.valid_until"
                                class="admin-promocodes__item-date"
                            >
                                <UIcon
                                    name="i-lucide-infinity"
                                    class="admin-promocodes__item-date-icon"
                                />
                                <span>Без ограничений</span>
                            </div>
                        </div>
                    </div>
                    <div class="admin-promocodes__item-actions">
                        <UTooltip text="Редактировать">
                            <UButton
                                icon="i-lucide-pencil"
                                color="neutral"
                                variant="ghost"
                                size="sm"
                                @click="openEditModal(item)"
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
                </div>
            </div>
        </UCard>

        <!-- Create/Edit Modal -->
        <UModal v-model="formModalOpen" class="max-w-lg">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-lg font-semibold">
                            {{
                                editingItem
                                    ? 'Редактировать промокод'
                                    : 'Создать промокод'
                            }}
                        </h3>
                        <UButton
                            icon="i-lucide-x"
                            color="neutral"
                            variant="ghost"
                            size="sm"
                            @click="formModalOpen = false"
                        />
                    </div>
                </template>
                <div class="space-y-4">
                    <UFormGroup label="Код промокода" required>
                        <UInput
                            v-model="form.code"
                            placeholder="Например: SUMMER2024"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                    <div class="grid grid-cols-2 gap-4">
                        <UFormGroup label="Тип скидки" required>
                            <USelect
                                v-model="form.discount_type"
                                :items="discountTypeOptions"
                                color="neutral"
                                variant="outline"
                            />
                        </UFormGroup>
                        <UFormGroup label="Значение" required>
                            <UInput
                                v-model="form.discount_value"
                                type="number"
                                :placeholder="
                                    form.discount_type === 'percentage'
                                        ? 'Процент'
                                        : 'Сумма в рублях'
                                "
                                color="neutral"
                                variant="outline"
                                min="1"
                            />
                        </UFormGroup>
                    </div>
                    <UFormGroup label="Максимум использований">
                        <UInput
                            v-model="form.max_uses"
                            type="number"
                            placeholder="Без ограничений"
                            color="neutral"
                            variant="outline"
                            min="1"
                        />
                    </UFormGroup>
                    <div class="grid grid-cols-2 gap-4">
                        <UFormGroup label="Дата начала">
                            <UInput
                                v-model="form.valid_from"
                                type="date"
                                color="neutral"
                                variant="outline"
                            />
                        </UFormGroup>
                        <UFormGroup label="Дата окончания">
                            <UInput
                                v-model="form.valid_until"
                                type="date"
                                color="neutral"
                                variant="outline"
                            />
                        </UFormGroup>
                    </div>
                    <UFormGroup label="Активен">
                        <USelect
                            v-model="form.is_active"
                            :items="activeOptions"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                </div>
                <template #footer>
                    <div class="flex justify-end gap-2">
                        <UButton
                            color="neutral"
                            variant="outline"
                            @click="formModalOpen = false"
                        >
                            Отмена
                        </UButton>
                        <UButton
                            color="primary"
                            :loading="saving"
                            :disabled="
                                !form.code || !form.discount_type || !form.discount_value
                            "
                            @click="savePromoCode"
                        >
                            {{ editingItem ? 'Сохранить' : 'Создать' }}
                        </UButton>
                    </div>
                </template>
            </UCard>
        </UModal>

        <!-- Delete Confirmation Modal -->
        <AdminConfirmDialog
            :visible="deleteModalOpen"
            title="Подтверждение удаления"
            :message="`Вы уверены, что хотите удалить промокод «${deleteTarget?.code}»?`"
            type="danger"
            confirm-text="Да, удалить"
            cancel-text="Нет"
            loading-text="Удаление..."
            :loading="deleting"
            @confirm="doDelete"
            @cancel="deleteModalOpen = false"
            @update:visible="deleteModalOpen = $event"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, } from 'vue';
import {
    fetchAdminPromoCodes,
    createAdminPromoCode,
    updateAdminPromoCode,
    deleteAdminPromoCode,
    type AdminPromoCodeItem,
} from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const items = ref<AdminPromoCodeItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

// Form modal
const formModalOpen = ref(false);
const editingItem = ref<AdminPromoCodeItem | null>(null);
const saving = ref(false);

const form = reactive({
    code: '',
    discount_type: 'percentage',
    discount_value: 10,
    max_uses: null as number | null,
    valid_from: '',
    valid_until: '',
    is_active: 'true',
});

const discountTypeOptions = [
    { label: 'Процент', value: 'percentage' },
    { label: 'Фиксированная сумма', value: 'fixed' },
];

const activeOptions = [
    { label: 'Да', value: 'true' },
    { label: 'Нет', value: 'false' },
];

// Delete modal
const deleteModalOpen = ref(false);
const deleteTarget = ref<AdminPromoCodeItem | null>(null);
const deleting = ref(false);

function formatDate(dateStr: string): string {
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric', });
}

function resetForm() {
    form.code = '';
    form.discount_type = 'percentage';
    form.discount_value = 10;
    form.max_uses = null;
    form.valid_from = '';
    form.valid_until = '';
    form.is_active = 'true';
}

function openCreateModal() {
    editingItem.value = null;
    resetForm();
    formModalOpen.value = true;
}

function openEditModal(item: AdminPromoCodeItem) {
    editingItem.value = item;
    form.code = item.code;
    form.discount_type = item.discount_type;
    form.discount_value = item.discount_value;
    form.max_uses = item.max_uses;
    form.valid_from = item.valid_from ? item.valid_from.split('T')[0] : '';
    form.valid_until = item.valid_until ? item.valid_until.split('T')[0] : '';
    form.is_active = item.is_active ? 'true' : 'false';
    formModalOpen.value = true;
}

async function savePromoCode() {
    saving.value = true;
    try {
        const data = {
            code: form.code,
            discount_type: form.discount_type,
            discount_value: Number(form.discount_value),
            max_uses: form.max_uses || null,
            valid_from: form.valid_from ? new Date(form.valid_from).toISOString() : null,
            valid_until: form.valid_until ? new Date(form.valid_until).toISOString() : null,
            is_active: form.is_active === 'true',
        };

        if (editingItem.value) {
            await updateAdminPromoCode(editingItem.value.id, data);
        } else {
            await createAdminPromoCode(data);
        }

        formModalOpen.value = false;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка сохранения';
    } finally {
        saving.value = false;
    }
}

function confirmDelete(item: AdminPromoCodeItem) {
    deleteTarget.value = item;
    deleteModalOpen.value = true;
}

async function doDelete() {
    if (!deleteTarget.value) return;
    deleting.value = true;
    try {
        await deleteAdminPromoCode(deleteTarget.value.id);
        deleteModalOpen.value = false;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка удаления';
    } finally {
        deleting.value = false;
    }
}

async function loadData() {
    loading.value = true;
    error.value = null;
    try {
        items.value = await fetchAdminPromoCodes();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка загрузки';
    } finally {
        loading.value = false;
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

/* Promo codes list */
.admin-promocodes__list {
    display: flex;
    flex-direction: column;
}

.admin-promocodes__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
}

.admin-promocodes__item:last-child {
    border-bottom: none;
}

.admin-promocodes__item-main {
    display: flex;
    align-items: center;
    gap: 24px;
    flex: 1;
    flex-wrap: wrap;
}

.admin-promocodes__item-code-section {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 160px;
}

.admin-promocodes__item-code {
    font-size: 16px;
    font-weight: 700;
    font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
    color: var(--admin-primary, #6c5ce7);
    letter-spacing: 0.5px;
}

.admin-promocodes__item-discount {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 100px;
}

.admin-promocodes__item-discount-value {
    font-size: 16px;
    font-weight: 700;
    color: #00b894;
}

.admin-promocodes__item-discount-type {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-promocodes__item-usage {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 80px;
}

.admin-promocodes__item-usage-count {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-promocodes__item-usage-label {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-promocodes__item-dates {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 140px;
}

.admin-promocodes__item-date {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-promocodes__item-date-icon {
    font-size: 14px;
    flex-shrink: 0;
}

.admin-promocodes__item-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
}

@media (max-width: 768px) {
    .admin-promocodes__item {
        flex-direction: column;
        align-items: flex-start;
    }

    .admin-promocodes__item-main {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
    }

    .admin-promocodes__item-actions {
        align-self: flex-end;
    }
}
</style>
