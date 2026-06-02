<template>
    <div class="add-sale-container">
        <div class="header-section">
            <h1 class="page-title">{{ $t('admin_shop_form.page_title') }}</h1>
        </div>

        <!-- Загрузчик -->
        <div v-if="isLoading" class="loading-container">
            <div class="loader" />
            <p>{{ $t('admin_shop_form.loading') }}</p>
        </div>

        <form v-else @submit.prevent="submitForm" class="sale-form">
            <!-- Поле для загрузки изображений -->
            <div class="form-section">
                <h2 class="section-title">{{ $t('admin_shop_form.sections.images') }}</h2>
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.select_images') }}
                        <span class="required">*</span></label
                    >
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop"
                        @click="triggerFileInput"
                    >
                        <input
                            ref="fileInput"
                            type="file"
                            multiple
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            required
                            @change="handleFileUpload"
                        />
                        <div class="file-drop-content">
                            <svg
                                class="upload-icon"
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                                <polyline points="17 8 12 3 7 8" />
                                <line x1="12" y1="3" x2="12" y2="15" />
                            </svg>
                            <p class="upload-text">
                                {{ $t('admin_shop_form.labels.drag_drop_text') }}
                            </p>
                            <p class="upload-hint">
                                {{ $t('admin_shop_form.labels.file_formats') }}
                            </p>
                        </div>
                    </div>
                    <div v-if="fileError" class="error-message">
                        {{ fileError }}
                    </div>
                    <div v-if="previewImages.length > 0" class="preview-container">
                        <div
                            v-for="(image, index) in previewImages"
                            :key="index"
                            class="image-preview"
                        >
                            <img
                                :src="image.preview"
                                class="preview-image"
                                :alt="`Preview ${index + 1}`"
                            />
                            <UButton
                                type="button"
                                class="remove-btn"
                                :aria-label="
                                    $t('admin_shop_form.aria_labels.remove_image', {
                                        index: index + 1,
                                    })
                                "
                                @click="removeImage(index)"
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Основная информация -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_shop_form.sections.main_info') }}
                </h2>

                <!-- Название картины -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.name_ru') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model="sale.name_ru"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_ru }"
                        :placeholder="$t('admin_shop_form.placeholders.name_ru')"
                        @blur="validateField('name_ru')"
                    />
                    <div v-if="errors.name_ru" class="error-message">
                        {{ errors.name_ru }}
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.name_en') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model="sale.name_en"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_en }"
                        :placeholder="$t('admin_shop_form.placeholders.name_en')"
                        @blur="validateField('name_en')"
                    />
                    <div v-if="errors.name_en" class="error-message">
                        {{ errors.name_en }}
                    </div>
                </div>

                <!-- Размеры картины -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.dimensions') }}
                        <span class="required">*</span></label
                    >
                    <div class="size-inputs">
                        <div class="size-input-wrapper">
                            <UInput
                                v-model.number="sale.width"
                                type="number"
                                required
                                min="1"
                                class="form-control size-input"
                                :class="{ 'is-invalid': errors.width }"
                                :placeholder="$t('admin_shop_form.labels.width')"
                                @blur="validateField('width')"
                            />
                            <div v-if="errors.width" class="error-message">
                                {{ errors.width }}
                            </div>
                        </div>
                        <span class="size-separator">×</span>
                        <div class="size-input-wrapper">
                            <UInput
                                v-model.number="sale.height"
                                type="number"
                                required
                                min="1"
                                class="form-control size-input"
                                :class="{ 'is-invalid': errors.height }"
                                :placeholder="$t('admin_shop_form.labels.height')"
                                @blur="validateField('height')"
                            />
                            <div v-if="errors.height" class="error-message">
                                {{ errors.height }}
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Год создания -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.year') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model.number="sale.year"
                        type="number"
                        required
                        min="2000"
                        :max="new Date().getFullYear()"
                        class="form-control"
                        :class="{ 'is-invalid': errors.year }"
                        :placeholder="$t('admin_shop_form.placeholders.year')"
                        @blur="validateField('year')"
                    />
                    <div v-if="errors.year" class="error-message">
                        {{ errors.year }}
                    </div>
                </div>

                <!-- Цена -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.price') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model.number="sale.price"
                        type="number"
                        required
                        min="1"
                        class="form-control"
                        :class="{ 'is-invalid': errors.price }"
                        :placeholder="$t('admin_shop_form.placeholders.price')"
                        @blur="validateField('price')"
                    />
                    <div v-if="errors.price" class="error-message">
                        {{ errors.price }}
                    </div>
                </div>
            </div>

            <!-- Технические характеристики -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_shop_form.sections.tech_specs') }}
                </h2>

                <!-- Основа -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.base') }}
                        <span class="required">*</span></label
                    >
                    <USelect
                        v-model="sale.base_id"
                        :items="baseOptions"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id }"
                        :placeholder="$t('admin_shop_form.placeholders.select_base')"
                        @blur="validateField('base_id')"
                    />
                    <div v-if="errors.base_id" class="error-message">
                        {{ errors.base_id }}
                    </div>
                </div>

                <!-- Материал -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_shop_form.labels.materials') }}
                        <span class="required">*</span></label
                    >
                    <USelect
                        v-model="sale.materials_ids"
                        :items="materialOptions"
                        multiple
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.materials_ids }"
                        :placeholder="$t('admin_shop_form.placeholders.select_materials')"
                        @blur="validateField('materials_ids')"
                    />
                    <div v-if="errors.materials_ids" class="error-message">
                        {{ errors.materials_ids }}
                    </div>
                </div>
            </div>

            <!-- Описание -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_shop_form.sections.additional_info') }}
                </h2>
                <div class="form-group">
                    <label class="form-label">{{
                        $t('admin_shop_form.labels.description')
                    }}</label>
                    <UTextarea
                        v-model="sale.descr"
                        class="form-control"
                        :placeholder="$t('admin_shop_form.placeholders.description')"
                        :rows="4"
                        :maxlength="500"
                    />
                    <div class="char-count">{{ sale.descr.length }}/500</div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <UButton type="button" class="btn btn-secondary" @click="resetForm">
                    {{ $t('admin_shop_form.buttons.clear_form') }}
                </UButton>
                <UButton
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">{{
                        $t('admin_shop_form.buttons.add_work')
                    }}</span>
                    <span v-else>
                        <span class="spinner" />
                        {{ $t('admin_shop_form.buttons.submitting') }}
                    </span>
                </UButton>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import { ref, reactive, computed, onMounted } from 'vue';
import type { CreateSaleDto, RequestResult } from '../../types';
import { useMaterialStore } from '../../stores/MaterialStore';

definePageMeta({
    middleware: 'admin-auth',
});

const { t, locale } = useI18n();
const toast = useToast();
const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;

const materialStore = useMaterialStore();

// Reactive state
const sale = reactive<CreateSaleDto>({
    width: 0,
    height: 0,
    year: new Date().getFullYear(),
    price: 0,
    name_ru: '',
    name_en: '',
    base_id: (null as unknown) as number,
    materials_ids: [],
    descr: '',
    images: [],
});

const files = ref<File[]>([]);
const previewImages = ref<{ file: File; preview: string }[]>([]);
const isSubmitting = ref(false);
const isLoading = ref(true);
const isDragOver = ref(false);
const fileError = ref<string | null>(null);
const errorMessage = ref('');

const errors = reactive<Record<string, string>>({
    name_ru: '',
    name_en: '',
    width: '',
    height: '',
    year: '',
    price: '',
    base_id: '',
    materials_ids: '',
});

const fileInput = ref<HTMLInputElement | null>(null);

// Computed
const bases = computed(() => materialStore.bases);
const materials = computed(() => materialStore.materials);

const baseOptions = computed(() => {
    return bases.value.map((base) => ({
        label: locale.value === 'ru' ? base.base_ru : base.base_en,
        value: base.id,
    }));
});

const materialOptions = computed(() => {
    return materials.value.map((material) => ({
        label: locale.value === 'ru' ? material.material_ru : material.material_en,
        value: material.id,
    }));
});

const isFormValid = computed(() => {
    return (
        sale.name_ru.trim() !== '' &&
        sale.name_en.trim() !== '' &&
        sale.width > 0 &&
        sale.height > 0 &&
        sale.year >= 2000 &&
        sale.year <= new Date().getFullYear() &&
        sale.price > 0 &&
        sale.base_id > 0 &&
        sale.materials_ids.length > 0 &&
        files.value.length > 0
    );
});

// Methods
function handleDragOver() {
    isDragOver.value = true;
}

function handleDragLeave() {
    isDragOver.value = false;
}

function handleDrop(event: DragEvent) {
    isDragOver.value = false;
    if (event.dataTransfer && event.dataTransfer.files.length) {
        const droppedFiles = Array.from(event.dataTransfer.files);
        addImages(droppedFiles);
    }
}

function triggerFileInput() {
    fileInput.value?.click();
}

function handleFileUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length) {
        const selectedFiles = Array.from(target.files);
        addImages(selectedFiles);
    }
}

function addImages(selectedFiles: File[]) {
    fileError.value = null;

    // Check max file count
    if (files.value.length + selectedFiles.length > 10) {
        fileError.value = t('admin_shop_form.errors.max_files');
        return;
    }

    // Check file types
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    const invalidFiles = selectedFiles.filter((file) => !validTypes.includes(file.type));

    if (invalidFiles.length > 0) {
        fileError.value = t('admin_shop_form.errors.invalid_file_type');
        return;
    }

    // Check file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

    if (largeFiles.length > 0) {
        fileError.value = t('admin_shop_form.errors.file_size');
        return;
    }

    // Add new files
    files.value = [...files.value, ...selectedFiles];
    sale.images = files.value.map((file) => file.name);

    // Create previews for new images
    selectedFiles.forEach((file) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            previewImages.value.push({
                file,
                preview: e.target?.result as string,
            });
        };
        reader.readAsDataURL(file);
    });
}

function removeImage(index: number) {
    previewImages.value.splice(index, 1);
    files.value.splice(index, 1);
    sale.images.splice(index, 1);
}

function validateField(fieldName: string) {
    switch (fieldName) {
        case 'name_ru':
            if (!sale.name_ru.trim()) {
                errors.name_ru = t('admin_shop_form.errors.name_ru_required');
            } else {
                errors.name_ru = '';
            }
            break;
        case 'name_en':
            if (!sale.name_en.trim()) {
                errors.name_en = t('admin_shop_form.errors.name_en_required');
            } else {
                errors.name_en = '';
            }
            break;
        case 'width':
            if (sale.width <= 0) {
                errors.width = t('admin_shop_form.errors.width_required');
            } else {
                errors.width = '';
            }
            break;
        case 'height':
            if (sale.height <= 0) {
                errors.height = t('admin_shop_form.errors.height_required');
            } else {
                errors.height = '';
            }
            break;
        case 'year':
            if (sale.year < 2000 || sale.year > new Date().getFullYear()) {
                errors.year = t('admin_shop_form.errors.year_range', {
                    year: new Date().getFullYear(),
                });
            } else {
                errors.year = '';
            }
            break;
        case 'price':
            if (sale.price <= 0) {
                errors.price = t('admin_shop_form.errors.price_required');
            } else {
                errors.price = '';
            }
            break;
        case 'base_id':
            if (sale.base_id <= 0) {
                errors.base_id = t('admin_shop_form.errors.base_required');
            } else {
                errors.base_id = '';
            }
            break;
        case 'materials_ids':
            if (sale.materials_ids.length === 0) {
                errors.materials_ids = t('admin_shop_form.errors.materials_required');
            } else {
                errors.materials_ids = '';
            }
            break;
    }
}

function validateForm() {
    validateField('name_ru');
    validateField('name_en');
    validateField('width');
    validateField('height');
    validateField('year');
    validateField('price');
    validateField('base_id');
    validateField('materials_ids');

    // Check images
    if (files.value.length === 0) {
        fileError.value = t('admin_shop_form.errors.images_required');
        return false;
    }

    // Check no errors
    return Object.values(errors).every((error) => error === '');
}

async function submitForm() {
    if (isSubmitting.value) return;

    if (!validateForm()) {
        return;
    }

    try {
        isSubmitting.value = true;
        errorMessage.value = '';

        // Build form data
        const formData = new FormData();

        // Add files
        files.value.forEach((file) => {
            formData.append('images', file);
        });

        // Add other data
        formData.append('data', JSON.stringify(sale));

        const response = await axios.post(SERVER_URL + 'sales', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });

        if (response.status === 200) {
            toast.add({
                title: t('toast.success.title'),
                description: t('toast.success.description'),
                icon: 'i-heroicons-check-circle',
                color: 'success',
                duration: 5000,
            });
            resetForm();
        } else {
            toast.add({
                title: t('toast.error.title'),
                description: t('admin_shop_form.messages.submit_failed'),
                icon: 'i-heroicons-exclamation-triangle',
                color: 'error',
                duration: 5000,
            });
        }
    } catch (error: any) {
        console.error('Error submitting form:', error);
        let description = t('admin_shop_form.messages.general_error');
        if (error.response?.status === 413) {
            description = t('admin_shop_form.messages.file_too_large');
        } else if (error.response?.status === 400) {
            description = t('admin_shop_form.messages.invalid_data');
        }
        toast.add({
            title: t('toast.error.title'),
            description,
            icon: 'i-heroicons-exclamation-triangle',
            color: 'error',
            duration: 5000,
        });
    } finally {
        isSubmitting.value = false;
    }
}

function resetForm() {
    sale.width = 0;
    sale.height = 0;
    sale.year = new Date().getFullYear();
    sale.price = 0;
    sale.name_ru = '';
    sale.name_en = '';
    sale.base_id = 0;
    sale.materials_ids = [];
    sale.descr = '';
    sale.images = [];

    files.value = [];
    previewImages.value = [];

    if (fileInput.value) {
        fileInput.value.value = '';
    }

    // Reset errors
    Object.keys(errors).forEach((key) => {
        errors[key] = '';
    });
    fileError.value = null;
}

// Lifecycle
onMounted(async () => {
    if (materialStore.materials.length === 0 || materialStore.bases.length === 0) {
        await materialStore.fetchAll();
    }
    isLoading.value = false;
});
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: #999;
}

.add-sale-container {
    max-width: 800px;
    margin: 3rem auto;
    padding: 2.5rem;
    background-color: var(--color-on-surface);
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

:root.dark .add-sale-container {
    background-color: #1e293b;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.header-section {
    text-align: center;
    margin-bottom: 2rem;
}

.page-title {
    color: #333;
    margin-bottom: 0.5rem;
    font-size: 1.8rem;
}

:root.dark .page-title {
    color: #f1f5f9;
}

.page-subtitle {
    color: #333;
    font-size: 1rem;
    margin: 0;
}

:root.dark .page-subtitle {
    color: #cbd5e1;
}

.sale-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.form-section {
    background-color: #f8f9fa;
    border-radius: 6px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

:root.dark .form-section {
    background-color: #0f172a;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.section-title {
    font-size: 1.3rem;
    color: #333;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #4a90e2;
}

:root.dark .section-title {
    color: #e2e8f0;
    border-bottom-color: #3b82f6;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
}

.form-label {
    font-weight: 600;
    color: #444;
    display: flex;
    align-items: center;
}

:root.dark .form-label {
    color: #cbd5e1;
}

.required {
    color: #e74c3c;
    margin-left: 0.25rem;
}

.form-control {
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    min-height: 46px;
    box-sizing: border-box;
    transition: border-color 0.3s, box-shadow 0.3s;
}

:root.dark .form-control {
    background-color: #1e293b;
    border-color: #334155;
    color: #e2e8f0;
}

.form-control:hover {
    border-color: #4a90e2;
}

:root.dark .form-control:hover {
    border-color: #3b82f6;
}

.form-control:focus {
    border-color: #4a90e2;
    outline: none;
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
}

:root.dark .form-control:focus {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.25);
}

/* Consistent placeholder color across all input fields */
.form-control::placeholder,
.form-control input::placeholder,
.form-control [data-placeholder],
.form-control [data-slot='placeholder'] {
    color: #999 !important;
    opacity: 1;
}

:root.dark .form-control::placeholder,
:root.dark .form-control input::placeholder,
:root.dark .form-control [data-placeholder],
:root.dark .form-control [data-slot='placeholder'] {
    color: #64748b !important;
}

.is-invalid {
    border-color: #e74c3c !important;
    box-shadow: 0 0 0 3px rgba(231, 76, 60, 0.1) !important;
}

:root.dark .is-invalid {
    border-color: #ef4444 !important;
    box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2) !important;
}

.drop-down-arrow {
    background-image: url("data:image/svg+xml;charset=UTF-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 1rem;
    padding-right: 2.5rem;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
}

:root.dark .drop-down-arrow {
    background-image: url("data:image/svg+xml;charset=UTF-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
}

.file-drop-area {
    border: 2px dashed #ccc;
    border-radius: 6px;
    padding: 2rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
    background-color: #fafafa;
}

:root.dark .file-drop-area {
    border-color: #475569;
    background-color: #0f172a;
}

.file-drop-area:hover,
.file-drop-area.drag-over {
    border-color: #4a90e2;
    background-color: #f0f8ff;
}

:root.dark .file-drop-area:hover,
:root.dark .file-drop-area.drag-over {
    border-color: #3b82f6;
    background-color: #1e293b;
}

.file-drop-area.drag-over {
    transform: scale(1.02);
}

.file-input {
    display: none;
}

.file-drop-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
}

.upload-icon {
    color: #4a90e2;
    width: 2rem;
    height: 2rem;
}

:root.dark .upload-icon {
    color: #60a5fa;
}

.upload-text {
    font-weight: 500;
    color: #333;
    margin: 0;
}

:root.dark .upload-text {
    color: #e2e8f0;
}

.upload-hint {
    color: #666;
    font-size: 0.875rem;
    margin: 0;
}

:root.dark .upload-hint {
    color: #94a3b8;
}

.preview-container {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 1rem;
}

.image-preview {
    position: relative;
    width: 100px;
    height: 100px;
    border: 1px solid #eee;
    border-radius: 4px;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

:root.dark .image-preview {
    border-color: #334155;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.preview-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.remove-btn {
    position: absolute;
    top: 0;
    right: 0;
    background-color: rgba(255, 0, 0, 0.7);
    color: white;
    border: none;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-size: 1rem;
    line-height: 1;
    border-radius: 0 0 0 4px;
}

.remove-btn:hover {
    background-color: rgba(255, 0, 0, 0.9);
}

.size-inputs {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.size-input-wrapper {
    flex: 1;
}

.size-input {
    flex: 1;
}

.size-separator {
    font-size: 1.2rem;
    color: #666;
}

:root.dark .size-separator {
    color: #94a3b8;
}

.char-count {
    font-size: 0.875rem;
    color: #666;
    text-align: right;
    margin-top: 0.25rem;
}

:root.dark .char-count {
    color: #94a3b8;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1rem;
}

.btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.3s;
    font-weight: 500;
}

.btn-primary {
    background-color: #4a90e2;
    color: white;
}

.btn-primary:hover:not(:disabled) {
    background-color: #3a7bc8;
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-primary:disabled {
    background-color: #a0c4ff;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
}

:root.dark .btn-primary:disabled {
    background-color: #1e3a5f;
    color: #64748b;
}

.btn-secondary {
    background-color: #f0f0f0;
    color: #333;
    border: 1px solid #ddd;
}

:root.dark .btn-secondary {
    background-color: #334155;
    color: #e2e8f0;
    border-color: #475569;
}

.btn-secondary:hover {
    background-color: #e0e0e0;
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

:root.dark .btn-secondary:hover {
    background-color: #475569;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

.error-message {
    color: #e74c3c;
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

:root.dark .error-message {
    color: #f87171;
}

.loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    gap: 1rem;
}

:root.dark .loading-container {
    color: #cbd5e1;
}

.loader {
    width: 2rem;
    height: 2rem;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #4a90e2;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

:root.dark .loader {
    border-color: #334155;
    border-top-color: #3b82f6;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}

.spinner {
    display: inline-block;
    width: 1rem;
    height: 1rem;
    border: 2px solid #ffffff;
    border-radius: 50%;
    border-top-color: transparent;
    animation: spin 1s linear infinite;
    margin-right: 0.5rem;
    vertical-align: middle;
}

@media (max-width: 768px) {
    .add-sale-container {
        padding: 1rem;
    }

    .form-section {
        padding: 1rem;
    }

    .form-actions {
        flex-direction: column;
    }

    .btn {
        width: 100%;
    }

    .page-title {
        font-size: 1.5rem;
    }

    .section-title {
        font-size: 1.2rem;
    }
}
</style>
