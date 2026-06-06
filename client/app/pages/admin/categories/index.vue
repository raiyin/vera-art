<template>
    <div class="admin-page">
        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Категории</h1>
                <p class="admin-page__subtitle">Управление категориями продуктов</p>
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
                    Создать категорию
                </UButton>
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading && !items.length" class="admin-page__loading">
            <UCard v-for="i in 4" :key="i">
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
                <UIcon name="i-lucide-folder-open" class="admin-page__empty-icon" />
                <h3 class="admin-page__empty-title">Категории не найдены</h3>
                <p class="admin-page__empty-desc">
                    Создайте первую категорию, чтобы начать.
                </p>
            </div>
        </UCard>

        <!-- Categories List -->
        <UCard v-else class="admin-page__table-card">
            <div class="admin-categories__list">
                <div
                    v-for="item in sortedItems"
                    :key="item.id"
                    class="admin-categories__item"
                >
                    <div class="admin-categories__item-main">
                        <div class="admin-categories__item-info">
                            <div class="admin-categories__item-name">
                                <span class="admin-categories__item-name-ru">{{
                                    item.name_ru
                                }}</span>
                                <span class="admin-categories__item-name-en">{{
                                    item.name_en
                                }}</span>
                            </div>
                            <div class="admin-categories__item-meta">
                                <span class="admin-categories__item-slug"
                                    >/{{ item.slug }}</span
                                >
                                <span class="admin-categories__item-sort">
                                    <UIcon
                                        name="i-lucide-arrow-up-down"
                                        class="admin-categories__item-meta-icon"
                                    />
                                    {{ item.sort_order }}
                                </span>
                            </div>
                        </div>
                        <div class="admin-categories__item-badges">
                            <UBadge
                                :color="item.is_active ? 'success' : 'neutral'"
                                variant="soft"
                                size="sm"
                            >
                                {{ item.is_active ? 'Активна' : 'Неактивна' }}
                            </UBadge>
                        </div>
                    </div>
                    <div class="admin-categories__item-actions">
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
                                    ? 'Редактировать категорию'
                                    : 'Создать категорию'
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
                    <UFormGroup label="Название (RU)" required>
                        <UInput
                            v-model="form.name_ru"
                            placeholder="Например: Акварель"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                    <UFormGroup label="Название (EN)" required>
                        <UInput
                            v-model="form.name_en"
                            placeholder="Например: Watercolor"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                    <UFormGroup label="Slug" required>
                        <UInput
                            v-model="form.slug"
                            placeholder="Например: watercolor"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                    <UFormGroup label="Описание (RU)">
                        <UInput
                            v-model="form.description_ru"
                            placeholder="Описание на русском"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                    <UFormGroup label="Описание (EN)">
                        <UInput
                            v-model="form.description_en"
                            placeholder="Description in English"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormGroup>
                    <div class="grid grid-cols-2 gap-4">
                        <UFormGroup label="Порядок сортировки">
                            <UInput
                                v-model="form.sort_order"
                                type="number"
                                placeholder="0"
                                color="neutral"
                                variant="outline"
                                min="0"
                            />
                        </UFormGroup>
                        <UFormGroup label="Активна">
                            <USelect
                                v-model="form.is_active"
                                :items="activeOptions"
                                color="neutral"
                                variant="outline"
                            />
                        </UFormGroup>
                    </div>
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
                            :disabled="!form.name_ru || !form.name_en || !form.slug"
                            @click="saveCategory"
                        >
                            {{ editingItem ? 'Сохранить' : 'Создать' }}
                        </UButton>
                    </div>
                </template>
            </UCard>
        </UModal>

        <!-- Delete Confirmation Modal -->
        <UModal v-model="deleteModalOpen" class="max-w-sm">
            <UCard>
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-lg font-semibold">Подтверждение удаления</h3>
                        <UButton
                            icon="i-lucide-x"
                            color="neutral"
                            variant="ghost"
                            size="sm"
                            @click="deleteModalOpen = false"
                        />
                    </div>
                </template>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                    Вы уверены, что хотите удалить категорию
                    <strong>{{ deleteTarget?.name_ru }}</strong
                    >?
                </p>
                <p v-if="deleteError" class="text-sm text-red-500 mt-2">
                    {{ deleteError }}
                </p>
                <template #footer>
                    <div class="flex justify-end gap-2">
                        <UButton
                            color="neutral"
                            variant="outline"
                            @click="deleteModalOpen = false"
                        >
                            Нет
                        </UButton>
                        <UButton color="error" :loading="deleting" @click="doDelete">
                            Да, удалить
                        </UButton>
                    </div>
                </template>
            </UCard>
        </UModal>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, } from 'vue';
import {
    fetchAdminCategories,
    createAdminCategory,
    updateAdminCategory,
    deleteAdminCategory,
    type AdminCategoryItem,
} from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const items = ref<AdminCategoryItem[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

// Form modal
const formModalOpen = ref(false);
const editingItem = ref<AdminCategoryItem | null>(null);
const saving = ref(false);

const form = reactive({
    name_ru: '',
    name_en: '',
    slug: '',
    description_ru: '',
    description_en: '',
    sort_order: 0,
    is_active: 'true',
});

const activeOptions = [
    { label: 'Да', value: 'true' },
    { label: 'Нет', value: 'false' },
];

// Delete modal
const deleteModalOpen = ref(false);
const deleteTarget = ref<AdminCategoryItem | null>(null);
const deleting = ref(false);
const deleteError = ref<string | null>(null);

const sortedItems = computed(() => {
    return [...items.value].sort((a, b) => a.sort_order - b.sort_order);
});

function resetForm() {
    form.name_ru = '';
    form.name_en = '';
    form.slug = '';
    form.description_ru = '';
    form.description_en = '';
    form.sort_order = 0;
    form.is_active = 'true';
}

function openCreateModal() {
    editingItem.value = null;
    resetForm();
    formModalOpen.value = true;
}

function openEditModal(item: AdminCategoryItem) {
    editingItem.value = item;
    form.name_ru = item.name_ru;
    form.name_en = item.name_en;
    form.slug = item.slug;
    form.description_ru = item.description_ru || '';
    form.description_en = item.description_en || '';
    form.sort_order = item.sort_order;
    form.is_active = item.is_active ? 'true' : 'false';
    formModalOpen.value = true;
}

async function saveCategory() {
    saving.value = true;
    try {
        const data = {
            name_ru: form.name_ru,
            name_en: form.name_en,
            slug: form.slug,
            description_ru: form.description_ru || undefined,
            description_en: form.description_en || undefined,
            sort_order: Number(form.sort_order),
            is_active: form.is_active === 'true',
        };

        if (editingItem.value) {
            await updateAdminCategory(editingItem.value.id, data);
        } else {
            await createAdminCategory(data);
        }

        formModalOpen.value = false;
        await loadData();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка сохранения';
    } finally {
        saving.value = false;
    }
}

function confirmDelete(item: AdminCategoryItem) {
    deleteTarget.value = item;
    deleteError.value = null;
    deleteModalOpen.value = true;
}

async function doDelete() {
    if (!deleteTarget.value) return;
    deleting.value = true;
    deleteError.value = null;
    try {
        await deleteAdminCategory(deleteTarget.value.id);
        deleteModalOpen.value = false;
        await loadData();
    } catch (e) {
        deleteError.value = e instanceof Error ? e.message : 'Ошибка удаления';
    } finally {
        deleting.value = false;
    }
}

async function loadData() {
    loading.value = true;
    error.value = null;
    try {
        items.value = await fetchAdminCategories();
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

/* Categories list */
.admin-categories__list {
    display: flex;
    flex-direction: column;
}

.admin-categories__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
}

.admin-categories__item:last-child {
    border-bottom: none;
}

.admin-categories__item-main {
    display: flex;
    align-items: center;
    gap: 24px;
    flex: 1;
    flex-wrap: wrap;
}

.admin-categories__item-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 200px;
}

.admin-categories__item-name {
    display: flex;
    align-items: center;
    gap: 8px;
}

.admin-categories__item-name-ru {
    font-size: 16px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-categories__item-name-en {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-categories__item-meta {
    display: flex;
    align-items: center;
    gap: 12px;
}

.admin-categories__item-slug {
    font-size: 12px;
    font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
    color: var(--admin-primary, #6c5ce7);
}

.admin-categories__item-sort {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    display: flex;
    align-items: center;
    gap: 2px;
}

.admin-categories__item-meta-icon {
    font-size: 12px;
}

.admin-categories__item-badges {
    display: flex;
    gap: 4px;
}

.admin-categories__item-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
}

@media (max-width: 768px) {
    .admin-categories__item {
        flex-direction: column;
        align-items: flex-start;
    }

    .admin-categories__item-main {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
    }

    .admin-categories__item-actions {
        align-self: flex-end;
    }
}
</style>
