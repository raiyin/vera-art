<template>
    <div class="admin-page">
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Теги</span>
        </div>
        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">
                    Теги
                </h1>
                <p class="admin-page__subtitle">
                    Управление тегами продуктов
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
                <UButton
                    icon="i-lucide-plus"
                    color="primary"
                    @click="openCreateModal"
                >
                    Создать тег
                </UButton>
            </div>
        </div>

        <!-- Loading State -->
        <div
            v-if="loading && !items.length"
            class="admin-page__loading"
        >
            <UCard
                v-for="i in 4"
                :key="i"
            >
                <div class="admin-page__skeleton-row">
                    <div class="admin-page__skeleton-lines">
                        <div class="admin-page__skeleton-line w-1/2" />
                        <div class="admin-page__skeleton-line w-1/3" />
                    </div>
                </div>
            </UCard>
        </div>

        <!-- Error State -->
        <UCard
            v-else-if="error"
            class="admin-page__error-card"
        >
            <div class="admin-page__error">
                <UIcon
                    name="i-lucide-alert-circle"
                    class="admin-page__error-icon"
                />
                <p>{{ error }}</p>
                <UButton
                    color="primary"
                    variant="outline"
                    @click="loadData"
                >
                    Повторить загрузку
                </UButton>
            </div>
        </UCard>

        <!-- Empty State -->
        <UCard v-else-if="!items.length && !loading">
            <div class="admin-page__empty">
                <UIcon
                    name="i-lucide-tags"
                    class="admin-page__empty-icon"
                />
                <h3 class="admin-page__empty-title">
                    Теги не найдены
                </h3>
                <p class="admin-page__empty-desc">
                    Создайте первый тег, чтобы начать.
                </p>
            </div>
        </UCard>

        <!-- Tags List -->
        <UCard
            v-else
            class="admin-page__table-card"
        >
            <div class="admin-tags__list">
                <div
                    v-for="item in items"
                    :key="item.id"
                    class="admin-tags__item"
                >
                    <div class="admin-tags__item-main">
                        <div class="admin-tags__item-info">
                            <div class="admin-tags__item-name">
                                <span class="admin-tags__item-name-ru">{{
                                    item.name_ru
                                }}</span>
                                <span class="admin-tags__item-name-en">{{
                                    item.name_en
                                }}</span>
                            </div>
                            <div class="admin-tags__item-slug">
                                /{{ item.slug }}
                            </div>
                        </div>
                    </div>
                    <div class="admin-tags__item-actions">
                        <UTooltip text="Редактировать">
                            <UButton
                                icon="i-lucide-pencil"
                                color="neutral"
                                variant="ghost"
                                size="sm"
                                @click="openEditModal(item,)"
                            />
                        </UTooltip>
                        <UTooltip text="Удалить">
                            <UButton
                                icon="i-lucide-trash-2"
                                color="error"
                                variant="ghost"
                                size="sm"
                                @click="confirmDelete(item,)"
                            />
                        </UTooltip>
                    </div>
                </div>
            </div>
        </UCard>

        <!-- Create/Edit Modal -->
        <UModal
            v-model:open="formModalOpen"
            class="max-w-lg"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-lg font-semibold">
                        {{ editingItem ? 'Редактировать тег' : 'Создать тег' }}
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
            <template #body>
                <div class="space-y-4">
                    <UFormField
                        label="Название (RU)"
                        required
                    >
                        <UInput
                            v-model="form.name_ru"
                            placeholder="Например: Акварель"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormField>
                    <UFormField
                        label="Название (EN)"
                        required
                    >
                        <UInput
                            v-model="form.name_en"
                            placeholder="Например: Watercolor"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormField>
                    <UFormField
                        label="Slug"
                        required
                    >
                        <UInput
                            v-model="form.slug"
                            placeholder="Например: watercolor"
                            color="neutral"
                            variant="outline"
                        />
                    </UFormField>
                </div>
            </template>
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
                        @click="saveTag"
                    >
                        {{ editingItem ? 'Сохранить' : 'Создать' }}
                    </UButton>
                </div>
            </template>
        </UModal>

        <!-- Delete Confirmation Modal -->
        <AdminConfirmDialog
            :visible="deleteModalOpen"
            title="Подтверждение удаления"
            :message="`Вы уверены, что хотите удалить тег «${deleteTarget?.name_ru}»?`"
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
        fetchAdminTags,
        createAdminTag,
        updateAdminTag,
        deleteAdminTag,
    } from '~/api/admin';
    import type { AdminTagItem, } from '~/types';

    definePageMeta({
        layout: 'admin',
        middleware: 'admin-auth',
    });

    const items = ref<AdminTagItem[]>([],);
    const loading = ref(false,);
    const error = ref<string | null>(null,);

    // Form modal
    const formModalOpen = ref(false,);
    const editingItem = ref<AdminTagItem | null>(null,);
    const saving = ref(false,);

    const form = reactive({
        name_ru: '',
        name_en: '',
        slug: '',
    });

    // Delete modal
    const deleteModalOpen = ref(false,);
    const deleteTarget = ref<AdminTagItem | null>(null,);
    const deleting = ref(false,);

    function resetForm() {
        form.name_ru = '';
        form.name_en = '';
        form.slug = '';
    }

    function openCreateModal() {
        editingItem.value = null;
        resetForm();
        formModalOpen.value = true;
    }

    function openEditModal(item: AdminTagItem,) {
        editingItem.value = item;
        form.name_ru = item.name_ru;
        form.name_en = item.name_en;
        form.slug = item.slug;
        formModalOpen.value = true;
    }

    async function saveTag() {
        saving.value = true;
        try {
            const data = {
                name_ru: form.name_ru,
                name_en: form.name_en,
                slug: form.slug,
            };

            if (editingItem.value) {
                await updateAdminTag(editingItem.value.id, data,);
            } else {
                await createAdminTag(data,);
            }

            formModalOpen.value = false;
            await loadData();
        } catch (e) {
            error.value = e instanceof Error ? e.message : 'Ошибка сохранения';
        } finally {
            saving.value = false;
        }
    }

    function confirmDelete(item: AdminTagItem,) {
        deleteTarget.value = item;
        deleteModalOpen.value = true;
    }

    async function doDelete() {
        if (!deleteTarget.value) return;
        deleting.value = true;
        try {
            await deleteAdminTag(deleteTarget.value.id,);
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
            items.value = await fetchAdminTags();
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

/* Tags list */
.admin-tags__list {
    display: flex;
    flex-direction: column;
}

.admin-tags__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
}

.admin-tags__item:last-child {
    border-bottom: none;
}

.admin-tags__item-main {
    display: flex;
    align-items: center;
    gap: 24px;
    flex: 1;
    flex-wrap: wrap;
}

.admin-tags__item-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 200px;
}

.admin-tags__item-name {
    display: flex;
    align-items: center;
    gap: 8px;
}

.admin-tags__item-name-ru {
    font-size: 16px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-tags__item-name-en {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-tags__item-slug {
    font-size: 12px;
    font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
    color: var(--admin-primary, #6c5ce7);
}

.admin-tags__item-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
}

@media (max-width: 768px) {
    .admin-tags__item {
        flex-direction: column;
        align-items: flex-start;
    }

    .admin-tags__item-main {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
    }

    .admin-tags__item-actions {
        align-self: flex-end;
    }
}
</style>
