<script setup lang="ts">
import type { PropType } from 'vue';
import type { CommonGetWorkDto } from '~/types';
import { useMaterialStore } from '../stores/MaterialStore';
import { useAuthStore } from '../stores/AuthStore';
import type { CommonTypedGetWorkDto } from '~/types/common_work';
import type { SelectItem } from '@nuxt/ui';
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from '#imports';
import { useRouter } from 'vue-router';

const props = defineProps({
    images: {
        type: [Array] as PropType<CommonTypedGetWorkDto[]>,
        required: true,
    },
});

const emit = defineEmits<{
    'work-deleted': [id: string];
}>();

// Reactive state
const selectedWorkType = ref('all');
const items = ref<SelectItem[]>([
    { value: 'all', label: 'Все работы' },
    { value: '1', label: 'Картины' },
    { value: '2', label: 'Иллюстрации' },
    { value: '3', label: '3D работы' },
]);
const value = ref('all');
// Image modal state
const showModal = ref(false);
const selectedWork = ref<CommonGetWorkDto | null>(null);
const currentImageIndex = ref(0);
// Delete confirmation modal state
const showDeleteModal = ref(false);
const workToDelete = ref<CommonGetWorkDto | null>(null);
const deletingId = ref<string | null>(null);

// Computed properties
const filteredImages = computed(() => {
    if (selectedWorkType.value === 'all') {
        return props.images;
    }
    return props.images.filter((image) => {
        return image.type === parseInt(selectedWorkType.value);
    });
});

const { locale } = useI18n();

const isAuthenticated = computed(() => {
    // Simple auth check - can be enhanced with proper auth store
    return !!localStorage.getItem('token');
});

const materialStore = useMaterialStore();
const authStore = useAuthStore();
const router = useRouter();

const isAdmin = computed(() => {
    return authStore.isAuthenticated && authStore.isAdmin;
});

// Methods
const handleWorkDeleted = (id: string) => {
    // Emit event to parent component to update the list
    emit('work-deleted', id);
};

const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;

const navigateToEdit = (work: CommonGetWorkDto) => {
    const id = work.id.toString();
    if (work.__type === 'GetSaleDto') {
        router.push(`/shop/edit/${id}`);
    } else {
        router.push(`/gallery/edit/${id}`);
    }
};

const confirmDelete = (work: CommonGetWorkDto) => {
    workToDelete.value = work;
    showDeleteModal.value = true;
};

const deleteWork = async () => {
    if (!workToDelete.value) return;

    const id = workToDelete.value.id.toString();
    deletingId.value = id;
    showDeleteModal.value = false;

    try {
        const response = await fetch(`${SERVER_URL}works/${id}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        handleWorkDeleted(id);
    } catch (e) {
        console.error('Error deleting work:', e);
    } finally {
        deletingId.value = null;
        workToDelete.value = null;
    }
};

const cancelDelete = () => {
    showDeleteModal.value = false;
    workToDelete.value = null;
};

const getWorkDescription = (work: CommonGetWorkDto): string => {
    const descr = locale.value === 'ru' ? (work as any).descr_ru : (work as any).descr_en;
    return descr || '';
};

const hasDescription = (work: CommonGetWorkDto): boolean => {
    const descrRu = (work as any).descr_ru;
    const descrEn = (work as any).descr_en;
    return !!(descrRu && descrRu.trim() !== '' && descrEn && descrEn.trim() !== '');
};

const getWorkName = (work: CommonGetWorkDto): string => {
    return locale.value === 'ru' ? work.name_ru : work.name_en;
};

const getWorkBase = (work: CommonGetWorkDto): string => {
    return locale.value === 'ru' ? work.base_ru || '' : work.base_en || '';
};

const getMainImageUrl = (work: CommonGetWorkDto): string => {
    return work.dir + (work.images?.[0] || '');
};

const getWorkDimensions = (work: CommonGetWorkDto): string => {
    if (work.width && work.width !== 0 && work.height && work.height !== 0) {
        return `${work.width}×${work.height}`;
    }
    return '';
};

const getWorkMaterials = (work: CommonGetWorkDto): string[] => {
    // If materials arrays are already present in the work object, use them
    if (locale.value === 'ru' && work.materials_ru && work.materials_ru.length > 0) {
        return work.materials_ru;
    }
    if (locale.value === 'en' && work.materials_en && work.materials_en.length > 0) {
        return work.materials_en;
    }
    // Otherwise fall back to MaterialStore
    return materialStore.getMaterialNames(work.materials_ids || [], locale.value);
};

const hasMaterials = (work: CommonGetWorkDto): boolean => {
    const materials = getWorkMaterials(work);
    return materials.length > 0;
};

// Modal methods
const openModal = (work: CommonGetWorkDto) => {
    console.log('openModal called with work:', work);
    console.log('Work images:', work.images);
    console.log('Work dir:', work.dir);
    selectedWork.value = work;
    currentImageIndex.value = 0;
    showModal.value = true;
    console.log('showModal set to:', showModal.value);
};

const closeModal = () => {
    showModal.value = false;
    selectedWork.value = null;
    currentImageIndex.value = 0;
};

// Carousel navigation methods
const prevImage = () => {
    if (selectedWork.value && currentImageIndex.value > 0) {
        currentImageIndex.value--;
    }
};

const nextImage = () => {
    if (
        selectedWork.value &&
        currentImageIndex.value < selectedWork.value.images.length - 1
    ) {
        currentImageIndex.value++;
    }
};

const goToImage = (index: number) => {
    if (selectedWork.value && index >= 0 && index < selectedWork.value.images.length) {
        currentImageIndex.value = index;
    }
};

// Keyboard navigation for modal
const handleKeydown = (event: KeyboardEvent) => {
    if (!showModal.value || !selectedWork.value) return;

    switch (event.key) {
        case 'Escape':
            closeModal();
            break;
        case 'ArrowLeft':
            if (currentImageIndex.value > 0) {
                currentImageIndex.value--;
            }
            break;
        case 'ArrowRight':
            if (
                selectedWork.value.images &&
                currentImageIndex.value < selectedWork.value.images.length - 1
            ) {
                currentImageIndex.value++;
            }
            break;
    }
};

// Watcher
watch(showModal, (newVal) => {
    if (newVal) {
        window.addEventListener('keydown', handleKeydown);
    } else {
        window.removeEventListener('keydown', handleKeydown);
    }
});

onBeforeUnmount(() => {
    window.removeEventListener('keydown', handleKeydown);
});
</script>

<template>
    <UContainer class="text-left main-content px-0">
        <div class="my-4">
            <USelect :items="items" v-model="value" />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 gap-6">
            <UCard
                v-for="work in filteredImages"
                :key="work.id"
                :ui="{ body: 'flex-1 flex flex-col' }"
                class="work-card flex flex-col overflow-hidden hover:shadow-xl transition-shadow duration-300 h-full relative"
            >
                <template #header>
                    <div class="flex flex-col justify-between items-start">
                        <h3 class="text-lg font-semibold truncate">
                            {{ getWorkName(work) }}
                        </h3>
                    </div>
                </template>

                <div class="flex flex-col grow">
                    <!-- Image -->
                    <div
                        class="relative aspect-square overflow-hidden rounded-lg bg-gray-100 dark:bg-gray-800 cursor-pointer"
                        @click="() => openModal(work)"
                    >
                        <img
                            :src="getMainImageUrl(work)"
                            :alt="getWorkName(work)"
                            class="w-full h-full object-cover hover:scale-105 transition-transform duration-500"
                            loading="lazy"
                        />
                        <div
                            class="absolute inset-0 bg-linear-to-t from-black/20 to-transparent opacity-0 hover:opacity-100 transition-opacity duration-300"
                        />

                        <!-- Description info icon -->
                        <div
                            v-if="hasDescription(work)"
                            class="absolute top-2 right-2 z-10"
                        >
                            <UPopover
                                mode="hover"
                                :content="{ side: 'bottom', sideOffset: 8 }"
                            >
                                <UIcon
                                    name="i-heroicons-information-circle"
                                    class="w-6 h-6 text-white drop-shadow-lg cursor-pointer hover:text-primary-300 transition-colors duration-200"
                                />

                                <template #content>
                                    <div
                                        class="p-4 max-w-xs text-sm leading-relaxed text-gray-800 dark:text-gray-200 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700"
                                    >
                                        <p>{{ getWorkDescription(work) }}</p>
                                    </div>
                                </template>
                            </UPopover>
                        </div>
                    </div>

                    <!-- Details -->
                    <div class="space-y-2 text-sm my-3">
                        <div
                            v-if="getWorkDimensions(work)"
                            class="flex items-center text-gray-600 dark:text-gray-400"
                        >
                            <UIcon
                                name="i-heroicons-arrows-pointing-out"
                                class="w-4 h-4 mr-2 shrink-0"
                            />
                            <span>{{ getWorkDimensions(work) }}</span>
                        </div>

                        <div
                            v-if="getWorkBase(work)"
                            class="flex items-center text-gray-600 dark:text-gray-400"
                        >
                            <UIcon
                                name="i-heroicons-document-text"
                                class="w-4 h-4 mr-2 shrink-0"
                            />
                            <span class="truncate">{{ getWorkBase(work) }}</span>
                        </div>

                        <!-- Materials -->
                        <div
                            v-if="hasMaterials(work)"
                            class="flex items-center flex-wrap [&:not(:first-child)]:ml-4] mt-2 text-gray-600 dark:text-gray-400"
                        >
                            <UIcon
                                name="i-heroicons-paint-brush"
                                class="w-4 h-4 mr-2 shrink-0"
                            />
                            <span>{{ getWorkMaterials(work).join(', ') }}</span>
                        </div>

                        <div
                            v-if="work.year"
                            class="flex items-center text-gray-600 dark:text-gray-400"
                        >
                            <UIcon
                                name="i-heroicons-calendar"
                                class="w-4 h-4 mr-2 shrink-0"
                            />
                            <span>{{ work.year }}</span>
                        </div>
                    </div>

                    <!-- Price badge positioned at bottom right -->
                    <UBadge
                        v-if="work.__type === 'GetSaleDto' && work.price"
                        color="primary"
                        size="xl"
                        class="absolute bottom-6 right-4 font-semibold"
                    >
                        {{ work.price }} ₽
                    </UBadge>

                    <!-- Admin Controls (at the bottom of the card) -->
                    <ClientOnly>
                        <div
                            v-if="isAdmin"
                            class="mt-auto pt-4 border-t border-gray-200 dark:border-gray-700"
                        >
                            <div class="flex gap-2">
                                <UButton
                                    size="sm"
                                    color="primary"
                                    variant="outline"
                                    @click.stop="navigateToEdit(work)"
                                >
                                    {{ $t('admin.edit_work') }}
                                </UButton>
                                <UButton
                                    size="sm"
                                    color="error"
                                    variant="outline"
                                    :loading="deletingId === work.id.toString()"
                                    :disabled="deletingId === work.id.toString()"
                                    @click.stop="confirmDelete(work)"
                                >
                                    {{ $t('admin.delete_work') }}
                                </UButton>
                            </div>
                        </div>
                    </ClientOnly>
                </div>
            </UCard>
        </div>

        <!-- Delete Confirmation Modal -->
        <UModal
            v-model:open="showDeleteModal"
            :dismissible="false"
            :close="false"
            :transition="true"
            class="delete-modal"
        >
            <template #header="{ close }">
                <div class="delete-modal-header">
                    <div class="delete-modal-icon-wrapper">
                        <UIcon
                            name="i-heroicons-exclamation-triangle"
                            class="delete-modal-icon"
                        />
                    </div>
                    <h3 class="delete-modal-title">
                        {{ $t('admin.delete_work_confirm_title') }}
                    </h3>
                </div>
            </template>

            <template #body>
                <div class="delete-modal-body">
                    <p class="delete-modal-text">
                        {{
                            $t('admin.delete_work_confirm_text', {
                                name: workToDelete ? getWorkName(workToDelete) : '',
                            })
                        }}
                    </p>
                    <p class="delete-modal-warning">
                        {{ $t('admin.delete_work_confirm_warning') }}
                    </p>
                </div>
            </template>

            <template #footer>
                <div class="delete-modal-footer">
                    <UButton
                        size="md"
                        color="neutral"
                        variant="outline"
                        @click="cancelDelete"
                    >
                        {{ $t('admin.delete_work_cancel') }}
                    </UButton>
                    <UButton
                        size="md"
                        color="error"
                        :loading="deletingId !== null"
                        :disabled="deletingId !== null"
                        @click="deleteWork"
                    >
                        {{ $t('admin.delete_work_confirm') }}
                    </UButton>
                </div>
            </template>
        </UModal>

        <!-- Custom Image Modal (matching news carousel) -->
        <transition name="modal-fade">
            <div
                v-if="showModal && selectedWork"
                class="custom-modal-overlay"
                @click.self="closeModal"
            >
                <div class="custom-modal-container">
                    <!-- Close Button -->
                    <button
                        class="custom-modal-close"
                        @click="closeModal"
                        aria-label="Close"
                    >
                        <svg
                            class="w-6 h-6"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>

                    <!-- Image Container -->
                    <div class="custom-modal-image-container">
                        <transition name="image-slide" mode="out-in">
                            <img
                                :key="currentImageIndex"
                                :src="
                                    selectedWork.dir +
                                    selectedWork.images[currentImageIndex]
                                "
                                :alt="
                                    getWorkName(selectedWork) +
                                    ' - Image ' +
                                    (currentImageIndex + 1)
                                "
                                class="custom-modal-image"
                                @click="closeModal"
                            />
                        </transition>

                        <!-- Navigation Buttons -->
                        <button
                            v-if="currentImageIndex > 0"
                            @click="prevImage"
                            class="custom-modal-nav-button left"
                            aria-label="Previous image"
                        >
                            <svg
                                class="w-8 h-8"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 19l-7-7 7-7"
                                />
                            </svg>
                        </button>
                        <button
                            v-if="currentImageIndex < selectedWork.images.length - 1"
                            @click="nextImage"
                            class="custom-modal-nav-button right"
                            aria-label="Next image"
                        >
                            <svg
                                class="w-8 h-8"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 5l7 7-7 7"
                                />
                            </svg>
                        </button>

                        <!-- Image Counter -->
                        <div class="custom-modal-counter">
                            <span class="counter-current">{{
                                currentImageIndex + 1
                            }}</span>
                            <span class="counter-separator">/</span>
                            <span class="counter-total">{{
                                selectedWork.images.length
                            }}</span>
                        </div>
                    </div>

                    <!-- Thumbnail Strip -->
                    <div
                        v-if="selectedWork.images.length > 1"
                        class="custom-modal-thumbnails"
                    >
                        <div
                            v-for="(image, index) in selectedWork.images"
                            :key="index"
                            class="thumbnail-item"
                            :class="{ active: index === currentImageIndex }"
                            @click="goToImage(index)"
                        >
                            <img
                                :src="selectedWork.dir + image"
                                :alt="'Thumbnail ' + (index + 1)"
                                class="thumbnail-image"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </transition>
    </UContainer>
</template>

<style scoped>
.work-card {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.work-card :deep(.u-card-body) {
    display: flex;
    flex-direction: column;
    flex: 1;
}

.work-types-selector {
    display: flex;
    padding: 1rem 1rem;
    align-items: center;
    justify-content: start;
    background-color: var(--color-surface);
    border-radius: 0.5rem;
    box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
    margin-bottom: 2rem;
    width: fit-content;
}

.selector-title {
    color: var(--color-on-surface);
    font-family: 'Montserrat', sans-serif;
    font-weight: 500;
    font-size: 1.3rem;
}

.selector-container {
    position: relative;
    margin-left: 1rem;
}

.work-type-dropdown {
    background-color: var(--color-surface-secondary-solid);
    color: var(--color-on-surface);
    border: 1px solid #ced4da;
    border-radius: 0.25rem;
    padding: 0.75rem;
    font-family: 'Montserrat', sans-serif;
    font-size: 1rem;
    width: 100%;
    transition: border-color 0.3s, box-shadow 0.3s;
}

.work-type-dropdown:focus {
    border-color: #4a90e2;
    outline: none;
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
}

.work-type-dropdown option {
    background-color: var(--color-surface);
    color: var(--color-on-surface);
}

.drop-down-arrow {
    background-image: url("data:image/svg+xml;charset=UTF-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 1rem;
    padding-right: 2.5rem; /* Make space for the arrow */
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
}

/* Custom Modal Styles (matching news carousel) */
.modal-fade-enter-active,
.modal-fade-leave-active {
    transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
    opacity: 0;
}

.image-slide-enter-active,
.image-slide-leave-active {
    transition: transform 0.3s ease, opacity 0.3s ease;
}

.image-slide-enter-from {
    transform: translateX(30px);
    opacity: 0;
}

.image-slide-leave-to {
    transform: translateX(-30px);
    opacity: 0;
}

.custom-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    padding: 1rem;
}

.custom-modal-container {
    position: relative;
    background: var(--color-surface);
    border-radius: 1.5rem;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    max-width: 90vw;
    max-height: 90vh;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.dark .custom-modal-container {
    background: #1e293b;
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.custom-modal-close {
    position: absolute;
    top: 1rem;
    right: 1rem;
    z-index: 10;
    background: rgba(0, 0, 0, 0.5);
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: white;
    transition: all 0.2s ease;
}

.custom-modal-close:hover {
    background: rgba(0, 0, 0, 0.8);
    transform: scale(1.1);
}

.custom-modal-image-container {
    position: relative;
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    min-height: 60vh;
}

.custom-modal-image {
    max-width: 100%;
    max-height: 70vh;
    object-fit: contain;
    border-radius: 0.75rem;
    cursor: zoom-out;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3);
}

.custom-modal-nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(0, 0, 0, 0.6);
    border: none;
    border-radius: 50%;
    width: 56px;
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: white;
    transition: all 0.2s ease;
    z-index: 5;
}

.custom-modal-nav-button:hover {
    background: rgba(0, 0, 0, 0.9);
    transform: translateY(-50%) scale(1.1);
}

.custom-modal-nav-button.left {
    left: 1rem;
}

.custom-modal-nav-button.right {
    right: 1rem;
}

.custom-modal-counter {
    position: absolute;
    top: 1rem;
    left: 1rem;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 2rem;
    font-size: 0.875rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 0.25rem;
}

.counter-current {
    color: #10b981;
}

.counter-separator {
    opacity: 0.7;
}

.counter-total {
    opacity: 0.9;
}

.custom-modal-thumbnails {
    display: flex;
    gap: 0.5rem;
    padding: 1rem;
    background: rgba(0, 0, 0, 0.3);
    overflow-x: auto;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.thumbnail-item {
    flex: 0 0 auto;
    width: 80px;
    height: 80px;
    border-radius: 0.5rem;
    overflow: hidden;
    cursor: pointer;
    border: 3px solid transparent;
    transition: all 0.2s ease;
    opacity: 0.7;
}

.thumbnail-item:hover {
    opacity: 1;
    transform: translateY(-2px);
}

.thumbnail-item.active {
    border-color: #10b981;
    opacity: 1;
    transform: scale(1.05);
}

.thumbnail-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

/* Responsive Modal Styles */
@media (max-width: 768px) {
    .custom-modal-container {
        max-width: 95vw;
        max-height: 95vh;
        border-radius: 1rem;
    }

    .custom-modal-image-container {
        padding: 1.5rem;
        min-height: 50vh;
    }

    .custom-modal-image {
        max-height: 60vh;
    }

    .custom-modal-nav-button {
        width: 44px;
        height: 44px;
    }

    .custom-modal-nav-button.left {
        left: 0.5rem;
    }

    .custom-modal-nav-button.right {
        right: 0.5rem;
    }

    .custom-modal-counter {
        top: 0.5rem;
        left: 0.5rem;
        padding: 0.375rem 0.75rem;
        font-size: 0.75rem;
    }

    .custom-modal-thumbnails {
        padding: 0.75rem;
        gap: 0.375rem;
    }

    .thumbnail-item {
        width: 60px;
        height: 60px;
    }

    .custom-modal-close {
        top: 0.5rem;
        right: 0.5rem;
        width: 36px;
        height: 36px;
    }
}

@media (max-width: 480px) {
    .custom-modal-image-container {
        padding: 1rem;
        min-height: 40vh;
    }

    .custom-modal-image {
        max-height: 50vh;
    }

    .thumbnail-item {
        width: 50px;
        height: 50px;
    }

    .custom-modal-nav-button {
        width: 36px;
        height: 36px;
    }

    .custom-modal-nav-button svg {
        width: 20px;
        height: 20px;
    }
}

@media (max-width: 576px) {
    .work-types-selector {
        padding: 1rem 0;
    }

    .selector-container {
        max-width: 100%;
        padding: 0 1rem;
    }
}

/* Delete Confirmation Modal Styles */
.delete-modal {
    --modal-max-width: 420px;
}

.delete-modal :deep(.ui-modal) {
    border-radius: 16px;
    overflow: hidden;
}

.delete-modal-header {
    text-align: center;
    padding: 1.5rem 1.5rem 0;
}

.delete-modal-icon-wrapper {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background: linear-gradient(135deg, #fef2f2 0%, #fee2e2 100%);
    margin-bottom: 1rem;
}

.dark .delete-modal-icon-wrapper {
    background: linear-gradient(135deg, #450a0a 0%, #7f1d1d 100%);
}

.delete-modal-icon {
    width: 32px;
    height: 32px;
    color: #dc2626;
}

.dark .delete-modal-icon {
    color: #fca5a5;
}

.delete-modal-title {
    font-size: 1.25rem;
    font-weight: 700;
    color: #1e293b;
    margin: 0;
}

.dark .delete-modal-title {
    color: #f1f5f9;
}

.delete-modal-body {
    padding: 1rem 1.5rem;
    text-align: center;
}

.delete-modal-text {
    font-size: 0.95rem;
    color: #475569;
    line-height: 1.6;
    margin: 0 0 0.5rem;
}

.dark .delete-modal-text {
    color: #94a3b8;
}

.delete-modal-warning {
    font-size: 0.85rem;
    color: #ef4444;
    font-weight: 500;
    margin: 0;
}

.dark .delete-modal-warning {
    color: #fca5a5;
}

.delete-modal-footer {
    display: flex;
    justify-content: center;
    gap: 0.75rem;
    padding: 0 1.5rem 1.5rem;
}

@media (max-width: 480px) {
    .delete-modal-footer {
        flex-direction: column-reverse;
    }

    .delete-modal-footer .UButton {
        width: 100%;
    }
}
</style>
