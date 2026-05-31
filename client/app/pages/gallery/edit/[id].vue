<template>
    <div class="edit-work-container">
        <div class="header-section">
            <h1 class="page-title">Редактировать работу в галерее</h1>
            <p class="page-subtitle">Измените необходимые поля и сохраните изменения</p>
        </div>

        <!-- Загрузчик -->
        <div v-if="isLoading" class="loading-container">
            <div class="loader" />
            <p>Загрузка данных...</p>
        </div>

        <form v-else class="work-form" @submit.prevent="submitForm">
            <!-- Поле для загрузки изображений -->
            <div class="form-section">
                <h2 class="section-title">Изображения работы</h2>
                <div class="form-group">
                    <label class="form-label"
                        >Выберите новые изображения (опционально)</label
                    >
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop"
                        @click="triggerFileInput"
                    >
                        <UInput
                            ref="fileInput"
                            type="file"
                            multiple
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
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
                                Перетащите изображения сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаются форматы: JPG, JPEG, PNG (макс. 10 файлов)
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
                                :aria-label="`Удалить изображение ${index + 1}`"
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
                <h2 class="section-title">Основная информация</h2>

                <!-- Название картины -->
                <div class="form-group">
                    <label class="form-label"
                        >Название картины по-русски <span class="required">*</span></label
                    >
                    <UInput
                        v-model="work.name_ru"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_ru }"
                        placeholder="Например: 'Звездная ночь'"
                        @blur="validateField('name_ru')"
                    />
                    <div v-if="errors.name_ru" class="error-message">
                        {{ errors.name_ru }}
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label"
                        >Название картины по-английски
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model="work.name_en"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_en }"
                        placeholder="Например: 'Starry Night'"
                        @blur="validateField('name_en')"
                    />
                    <div v-if="errors.name_en" class="error-message">
                        {{ errors.name_en }}
                    </div>
                </div>

                <!-- Размеры картины -->
                <div class="form-group">
                    <label class="form-label">
                        Размеры ({{ units }}) <span class="required">*</span>
                    </label>
                    <div class="size-inputs">
                        <div class="size-input-wrapper">
                            <UInput
                                v-model.number="work.width"
                                type="number"
                                required
                                min="1"
                                class="form-control size-input"
                                :class="{ 'is-invalid': errors.width }"
                                placeholder="Ширина"
                                @blur="validateField('width')"
                            />
                            <div v-if="errors.width" class="error-message">
                                {{ errors.width }}
                            </div>
                        </div>
                        <span class="size-separator">×</span>
                        <div class="size-input-wrapper">
                            <UInput
                                v-model.number="work.height"
                                type="number"
                                required
                                min="1"
                                class="form-control size-input"
                                :class="{ 'is-invalid': errors.height }"
                                placeholder="Высота"
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
                        >Год создания <span class="required">*</span></label
                    >
                    <UInput
                        v-model.number="work.year"
                        type="number"
                        required
                        min="2000"
                        :max="new Date().getFullYear()"
                        class="form-control"
                        :class="{ 'is-invalid': errors.year }"
                        placeholder="Например: 2023"
                        @blur="validateField('year')"
                    />
                    <div v-if="errors.year" class="error-message">
                        {{ errors.year }}
                    </div>
                </div>
            </div>

            <!-- Технические характеристики -->
            <div class="form-section">
                <h2 class="section-title">Технические характеристики</h2>

                <!-- Основа -->
                <div class="form-group">
                    <label class="form-label"
                        >Основа <span class="required">*</span></label
                    >
                    <USelect
                        v-model="work.base_id"
                        :items="baseOptions"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id }"
                        placeholder="Выберите основу"
                        @blur="validateField('base_id')"
                    />
                    <div v-if="errors.base_id" class="error-message">
                        {{ errors.base_id }}
                    </div>
                </div>

                <!-- Материал -->
                <div class="form-group">
                    <label class="form-label"
                        >Материалы
                        <span class="required">{{
                            isMaterialsRequired ? '*' : ''
                        }}</span></label
                    >
                    <USelect
                        v-model="work.materials_ids"
                        :items="materialOptions"
                        multiple
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.materials_ids }"
                        placeholder="Выберите материалы"
                        @blur="validateField('materials_ids')"
                    />
                    <div v-if="errors.materials_ids" class="error-message">
                        {{ errors.materials_ids }}
                    </div>
                </div>
            </div>

            <!-- Описание -->
            <div class="form-section">
                <h2 class="section-title">Дополнительная информация</h2>
                <div class="form-group">
                    <label class="form-label">Описание</label>
                    <UTextarea
                        v-model="work.descr"
                        class="form-control"
                        placeholder="Краткое описание картины"
                        :rows="4"
                        :maxlength="500"
                    />
                    <div class="char-count">{{ work.descr.length }}/500</div>
                </div>
            </div>

            <!-- Тип работы -->
            <div class="form-section">
                <h2 class="section-title">Тип работы</h2>
                <div class="form-group">
                    <label class="form-label"
                        >Тип работы <span class="required">*</span></label
                    >
                    <USelect
                        v-model="work.type"
                        :items="workTypeOptions"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.type }"
                        placeholder="Выберите тип работы"
                        @blur="validateField('type')"
                    />
                    <div v-if="errors.type" class="error-message">
                        {{ errors.type }}
                    </div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <UButton type="button" class="btn btn-secondary" @click="resetForm">
                    Сбросить изменения
                </UButton>
                <UButton
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">Сохранить изменения</span>
                    <span v-else>
                        <span class="spinner" />
                        Сохранение...
                    </span>
                </UButton>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue';
import type { RequestResult, UpdateWorkRequest, UpdateWorkResponse } from '../../../types';
import { useMaterialStore } from '../../../stores/MaterialStore';

definePageMeta({
    middleware: 'admin-auth',
});

const toast = useToast();
const config = useRuntimeConfig();
const route = useRoute();
const SERVER_URL = config.public.serverUrl;

const materialStore = useMaterialStore();

// Reactive state
const work = reactive<UpdateWorkResponse>({
    id: 0,
    str_id: '',
    dir: '',
    width: 0,
    height: 0,
    year: new Date().getFullYear(),
    name_ru: '',
    name_en: '',
    base_id: 0,
    descr: '',
    type: 0,
    materials_ids: [],
    images: [],
});

// For form reset
const originalWork = reactive<UpdateWorkResponse>({
    id: 0,
    str_id: '',
    dir: '',
    width: 0,
    height: 0,
    year: new Date().getFullYear(),
    name_ru: '',
    name_en: '',
    base_id: 0,
    descr: '',
    type: 0,
    materials_ids: [],
    images: [],
});

const files = ref<File[]>([]);
const previewImages = ref<
    { file?: File; preview: string; isExisting?: boolean; filename?: string }[]
>([]);
const imagesToDelete = ref<string[]>([]);
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
    base_id: '',
    materials_ids: '',
    type: '',
});

const fileInput = ref<HTMLInputElement | null>(null);

// Computed
const bases = computed(() => materialStore.bases);
const materials = computed(() => materialStore.materials);

const baseOptions = computed(() => {
    return bases.value.map((base) => ({
        label: base.base_ru,
        value: base.id,
    }));
});

const workTypeOptions = computed(() => {
    return [
        { label: 'Картина', value: 1 },
        { label: 'Иллюстрация', value: 2 },
        { label: '3D', value: 3 },
    ];
});

const materialOptions = computed(() => {
    return materials.value.map((material) => ({
        label: material.material_ru,
        value: material.id,
    }));
});

const isFormValid = computed(() => {
    return (
        work.name_ru.trim() !== '' &&
        work.name_en.trim() !== '' &&
        work.width > 0 &&
        work.height > 0 &&
        work.year >= 2000 &&
        work.year <= new Date().getFullYear() &&
        work.base_id > 0 &&
        (work.type < 3
            ? work.materials_ids.length > 0
            : work.materials_ids.length === 0) &&
        work.type > 0
    );
});

const units = computed(() => {
    return work.type <= 1 ? 'см' : 'px';
});

const isMaterialsRequired = computed(() => {
    return work.type === 1 || work.type === 2;
});

// Methods
async function loadWork() {
    try {
        const id = route.params.id;
        const response = await axios.get(`${SERVER_URL}works/${id}/edit`);
        Object.assign(work, response.data);
        Object.assign(originalWork, { ...response.data });

        // Load existing images as previews
        loadPreviewImages();

        isLoading.value = false;
    } catch (error) {
        console.error('Ошибка при загрузке работы:', error);
        isLoading.value = false;
        toast.add({
            title: 'Ошибка!',
            description: 'Не удалось загрузить данные. Пожалуйста, попробуйте позже.',
            icon: 'i-heroicons-exclamation-triangle',
            color: 'error',
            duration: 5000,
        });
    }
}

function loadPreviewImages() {
    previewImages.value = [];
    for (let i = 0; i < work.images.length; i++) {
        const imageUrl = `${work.dir}${work.images[i]}`;
        previewImages.value.push({
            preview: imageUrl,
            isExisting: true,
            filename: work.images[i],
        });
    }
}

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
        fileError.value = 'Можно загрузить не более 10 изображений';
        return;
    }

    // Check file types
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    const invalidFiles = selectedFiles.filter((file) => !validTypes.includes(file.type));

    if (invalidFiles.length > 0) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    // Check file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

    if (largeFiles.length > 0) {
        fileError.value = 'Размер каждого файла не должен превышать 5 МБ';
        return;
    }

    // Add new files
    files.value = [...files.value, ...selectedFiles];

    // Create previews for new images
    selectedFiles.forEach((file) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            previewImages.value.push({
                file,
                preview: e.target?.result as string,
                isExisting: false,
                filename: file.name,
            });
        };
        reader.readAsDataURL(file);
    });
}

function removeImage(index: number) {
    const imageToRemove = previewImages.value[index];

    if (imageToRemove.isExisting && imageToRemove.filename) {
        // Mark existing image for deletion
        if (!imagesToDelete.value.includes(imageToRemove.filename)) {
            imagesToDelete.value.push(imageToRemove.filename);
        }
    } else if (imageToRemove.file) {
        // Remove from files array if it's a newly uploaded file
        const fileIndex = files.value.indexOf(imageToRemove.file);
        if (fileIndex > -1) {
            files.value.splice(fileIndex, 1);
        }
    }

    // Remove from preview images
    previewImages.value.splice(index, 1);
}

function validateField(fieldName: string) {
    switch (fieldName) {
        case 'name_ru':
            if (!work.name_ru.trim()) {
                errors.name_ru = 'Пожалуйста, введите название на русском';
            } else {
                errors.name_ru = '';
            }
            break;
        case 'name_en':
            if (!work.name_en.trim()) {
                errors.name_en = 'Пожалуйста, введите название на английском';
            } else {
                errors.name_en = '';
            }
            break;
        case 'width':
            if (work.width <= 0) {
                errors.width = 'Ширина должна быть больше 0';
            } else {
                errors.width = '';
            }
            break;
        case 'height':
            if (work.height <= 0) {
                errors.height = 'Высота должна быть больше 0';
            } else {
                errors.height = '';
            }
            break;
        case 'year':
            if (work.year < 2000 || work.year > new Date().getFullYear()) {
                errors.year = `Год должен быть между 2000 и ${new Date().getFullYear()}`;
            } else {
                errors.year = '';
            }
            break;
        case 'base_id':
            if (work.base_id <= 0) {
                errors.base_id = 'Пожалуйста, выберите основу';
            } else {
                errors.base_id = '';
            }
            break;
        case 'materials_ids':
            if (work.materials_ids.length === 0) {
                errors.materials_ids = 'Пожалуйста, выберите хотя бы один материал';
            } else {
                errors.materials_ids = '';
            }
            break;
        case 'type':
            if (work.type <= 0) {
                errors.type = 'Пожалуйста, выберите тип работы';
            } else {
                errors.type = '';
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
    validateField('base_id');
    validateField('materials_ids');
    validateField('type');

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

        const finalImages: string[] = [];

        // Add existing images that are not marked for deletion
        for (const imageName of work.images) {
            if (!imagesToDelete.value.includes(imageName)) {
                finalImages.push(imageName);
            }
        }

        // Add new image filenames
        for (const file of files.value) {
            finalImages.push(file.name);
        }

        // Add files if any
        if (files.value.length > 0) {
            files.value.forEach((file) => {
                formData.append('images', file);
            });
        }

        const workDataToUpdate: UpdateWorkRequest = {
            ...work,
            images: finalImages,
        };
        workDataToUpdate.type = parseInt(workDataToUpdate.type as any);

        formData.append('data', JSON.stringify(workDataToUpdate));

        const strId = route.params.id;
        const response = await axios.put(SERVER_URL + 'works/' + strId, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (response.status === 200) {
            toast.add({
                title: 'Успешно!',
                description: 'Работа успешно обновлена в галерее.',
                icon: 'i-heroicons-check-circle',
                color: 'success',
                duration: 5000,
            });
            // Update work images with final list
            work.images = finalImages;
            // Clear deletion list and files
            imagesToDelete.value = [];
            files.value = [];
            // Update original work
            Object.assign(originalWork, { ...work });
        } else {
            toast.add({
                title: 'Ошибка!',
                description: 'Не удалось обновить работу. Пожалуйста, попробуйте снова.',
                icon: 'i-heroicons-exclamation-triangle',
                color: 'error',
                duration: 5000,
            });
        }
    } catch (error: any) {
        console.error('Error submitting form:', error);
        let description =
            'Произошла ошибка при обновлении работы. Пожалуйста, попробуйте снова.';
        if (error.response?.status === 413) {
            description =
                'Файлы слишком большие. Пожалуйста, загрузите меньшие изображения.';
        } else if (error.response?.status === 400) {
            description =
                'Некорректные данные. Пожалуйста, проверьте введенные значения.';
        }
        toast.add({
            title: 'Ошибка!',
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
    Object.assign(work, { ...originalWork });
    files.value = [];
    previewImages.value = [];
    imagesToDelete.value = [];
    loadPreviewImages();
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
    await loadWork();
});
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: #999;
}

.edit-work-container {
    max-width: 800px;
    margin: 3rem auto;
    padding: 2.5rem;
    background-color: var(--color-on-surface);
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
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

.page-subtitle {
    color: #333;
    font-size: 1rem;
    margin: 0;
}

.work-form {
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

.section-title {
    font-size: 1.3rem;
    color: #333;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #4a90e2;
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

.form-control:hover {
    border-color: #4a90e2;
}

.form-control:focus {
    border-color: #4a90e2;
    outline: none;
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
}

/* Consistent placeholder color across all input fields */
.form-control::placeholder,
.form-control input::placeholder,
.form-control [data-placeholder],
.form-control [data-slot='placeholder'] {
    color: #999 !important;
    opacity: 1;
}

.is-invalid {
    border-color: #e74c3c !important;
    box-shadow: 0 0 0 3px rgba(231, 76, 60, 0.1) !important;
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

.file-drop-area {
    border: 2px dashed #ccc;
    border-radius: 6px;
    padding: 2rem;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s;
    background-color: #fafafa;
}

.file-drop-area:hover,
.file-drop-area.drag-over {
    border-color: #4a90e2;
    background-color: #f0f8ff;
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

.upload-text {
    font-weight: 500;
    color: #333;
    margin: 0;
}

.upload-hint {
    color: #666;
    font-size: 0.875rem;
    margin: 0;
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

.char-count {
    font-size: 0.875rem;
    color: #666;
    text-align: right;
    margin-top: 0.25rem;
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

.btn-secondary {
    background-color: #f0f0f0;
    color: #333;
    border: 1px solid #ddd;
}

.btn-secondary:hover {
    background-color: #e0e0e0;
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.error-message {
    color: #e74c3c;
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

.loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    gap: 1rem;
}

.loader {
    width: 2rem;
    height: 2rem;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #4a90e2;
    border-radius: 50%;
    animation: spin 1s linear infinite;
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
    .edit-work-container {
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

/* ===== Dark Mode Overrides ===== */
:root.dark .edit-work-container {
    background-color: #1e293b;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

:root.dark .page-title {
    color: #f1f5f9;
}

:root.dark .page-subtitle {
    color: #cbd5e1;
}

:root.dark .form-section {
    background-color: #0f172a;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

:root.dark .section-title {
    color: #e2e8f0;
    border-bottom-color: #3b82f6;
}

:root.dark .form-label {
    color: #cbd5e1;
}

:root.dark .form-control {
    background-color: #1e293b;
    border-color: #334155;
    color: #e2e8f0;
}

:root.dark .form-control:hover {
    border-color: #3b82f6;
}

:root.dark .form-control:focus {
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.25);
}

:root.dark .form-control::placeholder,
:root.dark .form-control input::placeholder,
:root.dark .form-control [data-placeholder],
:root.dark .form-control [data-slot='placeholder'] {
    color: #64748b !important;
}

:root.dark .is-invalid {
    border-color: #ef4444 !important;
    box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2) !important;
}

:root.dark .drop-down-arrow {
    background-image: url("data:image/svg+xml;charset=UTF-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
}

:root.dark .file-drop-area {
    border-color: #475569;
    background-color: #0f172a;
}

:root.dark .file-drop-area:hover,
:root.dark .file-drop-area.drag-over {
    border-color: #3b82f6;
    background-color: #1e293b;
}

:root.dark .upload-icon {
    color: #60a5fa;
}

:root.dark .upload-text {
    color: #e2e8f0;
}

:root.dark .upload-hint {
    color: #94a3b8;
}

:root.dark .image-preview {
    border-color: #334155;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

:root.dark .size-separator {
    color: #94a3b8;
}

:root.dark .char-count {
    color: #94a3b8;
}

:root.dark .btn-primary:disabled {
    background-color: #1e3a5f;
    color: #64748b;
}

:root.dark .btn-secondary {
    background-color: #334155;
    color: #e2e8f0;
    border-color: #475569;
}

:root.dark .btn-secondary:hover {
    background-color: #475569;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
}

:root.dark .error-message {
    color: #f87171;
}

:root.dark .loading-container {
    color: #cbd5e1;
}

:root.dark .loader {
    border-color: #334155;
    border-top-color: #3b82f6;
}
</style>
