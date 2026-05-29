<template>
    <div class="edit-sale-container">
        <div class="header-section">
            <h1 class="page-title">
                Редактировать работу в магазине
            </h1>
            <p class="page-subtitle">
                Измените необходимые поля и сохраните изменения
            </p>
        </div>

        <!-- Загрузчик -->
        <div
            v-if="isLoading"
            class="loading-container"
        >
            <div class="loader" />
            <p>Загрузка данных...</p>
        </div>

        <form
            v-else
            class="sale-form"
            @submit.prevent="submitForm"
        >
            <!-- Поле для загрузки изображений -->
            <div class="form-section">
                <h2 class="section-title">
                    Изображения работы
                </h2>
                <div class="form-group">
                    <label class="form-label">Выберите новые изображения (опционально)</label>
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver, }"
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
                                <path
                                    d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"
                                />
                                <polyline points="17 8 12 3 7 8" />
                                <line
                                    x1="12"
                                    y1="3"
                                    x2="12"
                                    y2="15"
                                />
                            </svg>
                            <p class="upload-text">
                                Перетащите изображения сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаются форматы: JPG, JPEG, PNG (макс. 10 файлов)
                            </p>
                        </div>
                    </div>
                    <div
                        v-if="fileError"
                        class="error-message"
                    >
                        {{ fileError }}
                    </div>
                    <div
                        v-if="previewImages.length > 0"
                        class="preview-container"
                    >
                        <div
                            v-for="(image, index) in previewImages"
                            :key="index"
                            class="image-preview"
                        >
                            <img
                                :src="image.preview"
                                class="preview-image"
                                :alt="`Preview ${index + 1}`"
                            >
                            <UButton
                                type="button"
                                class="remove-btn"
                                :aria-label="`Удалить изображение ${index + 1}`"
                                @click="removeImage(index,)"
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
                    Основная информация
                </h2>

                <!-- Название картины -->
                <div class="form-group">
                    <label class="form-label">Название картины по-русски <span class="required">*</span></label>
                    <UInput
                        v-model="sale.name_ru"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_ru, }"
                        placeholder="Например: 'Звездная ночь'"
                        @blur="validateField('name_ru',)"
                    />
                    <div
                        v-if="errors.name_ru"
                        class="error-message"
                    >
                        {{ errors.name_ru }}
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">Название картины по-английски
                        <span class="required">*</span></label>
                    <UInput
                        v-model="sale.name_en"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_en, }"
                        placeholder="Например: 'Starry Night'"
                        @blur="validateField('name_en',)"
                    />
                    <div
                        v-if="errors.name_en"
                        class="error-message"
                    >
                        {{ errors.name_en }}
                    </div>
                </div>

                <!-- Размеры картины -->
                <div class="form-group">
                    <label class="form-label">
                        Размеры (см) <span class="required">*</span>
                    </label>
                    <div class="size-inputs">
                        <div class="size-input-wrapper">
                            <UInput
                                v-model.number="sale.width"
                                type="number"
                                required
                                min="1"
                                class="form-control size-input"
                                :class="{ 'is-invalid': errors.width, }"
                                placeholder="Ширина"
                                @blur="validateField('width',)"
                            />
                            <div
                                v-if="errors.width"
                                class="error-message"
                            >
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
                                :class="{ 'is-invalid': errors.height, }"
                                placeholder="Высота"
                                @blur="validateField('height',)"
                            />
                            <div
                                v-if="errors.height"
                                class="error-message"
                            >
                                {{ errors.height }}
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Год создания -->
                <div class="form-group">
                    <label class="form-label">Год создания <span class="required">*</span></label>
                    <UInput
                        v-model.number="sale.year"
                        type="number"
                        required
                        min="2000"
                        :max="new Date().getFullYear()"
                        class="form-control"
                        :class="{ 'is-invalid': errors.year, }"
                        placeholder="Например: 2023"
                        @blur="validateField('year',)"
                    />
                    <div
                        v-if="errors.year"
                        class="error-message"
                    >
                        {{ errors.year }}
                    </div>
                </div>

                <!-- Цена -->
                <div class="form-group">
                    <label class="form-label">Цена <span class="required">*</span></label>
                    <UInput
                        v-model.number="sale.price"
                        type="number"
                        required
                        min="1"
                        class="form-control"
                        :class="{ 'is-invalid': errors.price, }"
                        placeholder="Например: 1000"
                        @blur="validateField('price',)"
                    />
                    <div
                        v-if="errors.price"
                        class="error-message"
                    >
                        {{ errors.price }}
                    </div>
                </div>
            </div>

            <!-- Технические характеристики -->
            <div class="form-section">
                <h2 class="section-title">
                    Технические характеристики
                </h2>

                <!-- Основа -->
                <div class="form-group">
                    <label class="form-label">Основа <span class="required">*</span></label>
                    <USelect
                        v-model="sale.base_id"
                        :items="baseOptions"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id, }"
                        placeholder="Выберите основу"
                        @blur="validateField('base_id')"
                    />
                    <div
                        v-if="errors.base_id"
                        class="error-message"
                    >
                        {{ errors.base_id }}
                    </div>
                </div>

                <!-- Материал -->
                <div class="form-group">
                    <label class="form-label">Материалы <span class="required">*</span></label>
                    <div class="multi-select-wrapper">
                        <div
                            class="select-display drop-down-arrow"
                            :class="{ 'is-invalid': errors.materials_ids, }"
                            tabindex="0"
                            @click="materialsToggleDropdown"
                            @keydown.enter="materialsToggleDropdown"
                            @blur="validateField('materials_ids',)"
                        >
                            {{ selectedMaterialsDisplay || 'Выберите материалы' }}
                        </div>
                        <div
                            v-if="materialsDropdownOpen"
                            class="dropdown-options form-control"
                        >
                            <div
                                v-for="material in materials"
                                :key="material.id"
                                class="option-item"
                            >
                                <UInput
                                    :id="'material-' + material.id"
                                    type="checkbox"
                                    :value="material.id"
                                    :model-value="
                                        sale.materials_ids.includes(material.id,)
                                    "
                                    @update:model-value="
                                        (checked,) => toggleMaterial(material.id, checked,)
                                    "
                                />
                                <label :for="'material-' + material.id">
                                    {{
                                        $i18n.locale === 'ru'
                                            ? material.material_ru
                                            : material.material_en
                                    }}
                                </label>
                            </div>
                        </div>
                        <div
                            v-if="errors.materials_ids"
                            class="error-message"
                        >
                            {{ errors.materials_ids }}
                        </div>
                    </div>
                </div>
            </div>

            <!-- Описание -->
            <div class="form-section">
                <h2 class="section-title">
                    Дополнительная информация
                </h2>
                <div class="form-group">
                    <label class="form-label">Описание</label>
                    <textarea
                        v-model="sale.descr"
                        class="form-control textarea"
                        placeholder="Краткое описание картины"
                        rows="4"
                        maxlength="500"
                    />
                    <div class="char-count">
                        {{ sale.descr.length }}/500
                    </div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <UButton
                    type="button"
                    class="btn btn-secondary"
                    @click="resetForm"
                >
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

        <!-- Success Alert -->
        <UAlert
            v-if="showSuccessAlert"
            title="Успешно!"
            description="Работа успешно обновлена в магазине."
            color="success"
            icon="i-heroicons-check-circle"
            closable
            @close="showSuccessAlert = false"
        />

        <!-- Danger Alert -->
        <UAlert
            v-if="showErrorAlert"
            title="Ошибка!"
            :description="
                errorMessage
                    || 'Не удалось обновить работу в магазине. Пожалуйста, попробуйте снова.'
            "
            color="error"
            icon="i-heroicons-exclamation-triangle"
            closable
            @close="showErrorAlert = false"
        />
    </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import type {
    UpdateSaleRequest,
    UpdateSaleResponse,
    RequestResult,
} from '../../types';
import { useMaterialStore } from '../../stores/MaterialStore';

const route = useRoute();
const { locale } = useI18n();
const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;
const materialStore = useMaterialStore();

// Reactive state
const sale = reactive<UpdateSaleResponse>({
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
    price: 0,
    materials_ids: [],
    images: [],
});

const originalSale = reactive<UpdateSaleResponse>({
    id: 0,
    str_id: '',
    dir: '',
    name_ru: '',
    name_en: '',
    base_id: 0,
    year: new Date().getFullYear(),
    descr: '',
    width: 0,
    height: 0,
    price: 0,
    images: [],
    materials_ids: [],
});

const addedFiles = ref<File[]>([]);
const previewImages = ref<{ file?: File; preview: string; isExisting?: boolean; filename?: string }[]>([]);
const imagesToDelete = ref<string[]>([]);
const isSubmitting = ref(false);
const isLoading = ref(true);
const loadError = ref<string | null>(null);
const requestResult = ref<RequestResult>('unknown');
const materialsDropdownOpen = ref(false);
const basesDropdownOpen = ref(false);
const isDragOver = ref(false);
const fileError = ref<string | null>(null);
const errorMessage = ref('');
const showSuccessAlert = ref(false);
const showErrorAlert = ref(false);

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

const selectedMaterialsDisplay = computed(() => {
    if (sale.materials_ids.length === 0) return '';
    const selectedNames = materialStore.materials
        .filter(material => sale.materials_ids.includes(material.id))
        .map(material =>
            locale.value === 'ru' ? material.material_ru : material.material_en
        );
    return selectedNames.join(', ');
});

const isFormValid = computed(() => {
    return (
        sale.name_ru.trim() !== ''
        && sale.name_en.trim() !== ''
        && sale.width > 0
        && sale.height > 0
        && sale.year >= 2000
        && sale.year <= new Date().getFullYear()
        && sale.price > 0
        && sale.base_id > 0
        && sale.materials_ids.length > 0
    );
});

// Methods
async function loadSale() {
    try {
        const id = route.params.id;
        const response = await axios.get(`${SERVER_URL}sales/${id}/edit`);
        Object.assign(sale, response.data);
        Object.assign(originalSale, { ...response.data });
        loadPreviewImages();
        isLoading.value = false;
    } catch (error) {
        console.error('Ошибка при загрузке работы:', error);
        loadError.value = 'Не удалось загрузить работу';
        isLoading.value = false;
        showErrorAlert.value = true;
        errorMessage.value = 'Не удалось загрузить данные. Пожалуйста, попробуйте позже.';
    }
}

function loadPreviewImages() {
    previewImages.value = [];
    for (let i = 0; i < sale.images.length; i++) {
        const imageUrl = `${sale.dir}${sale.images[i]}`;
        previewImages.value.push({
            preview: imageUrl,
            filename: sale.images[i],
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
        const files = Array.from(event.dataTransfer.files);
        addImages(files);
    }
}

function triggerFileInput() {
    fileInput.value?.click();
}

function handleFileUpload(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length) {
        const files = Array.from(target.files);
        addImages(files);
    }
}

function addImages(selectedFiles: File[]) {
    fileError.value = null;

    if (addedFiles.value.length + selectedFiles.length > 10) {
        fileError.value = 'Можно загрузить не более 10 изображений';
        return;
    }

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    const invalidFiles = selectedFiles.filter(
        (file) => !validTypes.includes(file.type)
    );

    if (invalidFiles.length > 0) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    const largeFiles = selectedFiles.filter(file => file.size > maxSize);

    if (largeFiles.length > 0) {
        fileError.value = 'Размер каждого файла не должен превышать 5 МБ';
        return;
    }

    addedFiles.value = [...addedFiles.value, ...selectedFiles];

    selectedFiles.forEach((file) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            previewImages.value.push({
                file,
                preview: e.target?.result as string,
                filename: file.name,
            });
        };
        reader.readAsDataURL(file);
    });
}

function removeImage(index: number) {
    const imageToRemove = previewImages.value[index];
    if (!imageToRemove) return;

    if (imageToRemove.file) {
        const fileIndex = addedFiles.value.indexOf(imageToRemove.file);
        if (fileIndex > -1) {
            addedFiles.value.splice(fileIndex, 1);
        }
    }

    previewImages.value.splice(index, 1);
}

function validateField(fieldName: string) {
    switch (fieldName) {
        case 'name_ru':
            if (!sale.name_ru.trim()) {
                errors.name_ru = 'Пожалуйста, введите название на русском';
            } else {
                errors.name_ru = '';
            }
            break;
        case 'name_en':
            if (!sale.name_en.trim()) {
                errors.name_en = 'Пожалуйста, введите название на английском';
            } else {
                errors.name_en = '';
            }
            break;
        case 'width':
            if (sale.width <= 0) {
                errors.width = 'Ширина должна быть больше 0';
            } else {
                errors.width = '';
            }
            break;
        case 'height':
            if (sale.height <= 0) {
                errors.height = 'Высота должна быть больше 0';
            } else {
                errors.height = '';
            }
            break;
        case 'year':
            if (
                sale.year < 2000
                || sale.year > new Date().getFullYear()
            ) {
                errors.year = `Год должен быть между 2000 и ${new Date().getFullYear()}`;
            } else {
                errors.year = '';
            }
            break;
        case 'price':
            if (sale.price <= 0) {
                errors.price = 'Цена должна быть больше 0';
            } else {
                errors.price = '';
            }
            break;
        case 'base_id':
            if (sale.base_id <= 0) {
                errors.base_id = 'Пожалуйста, выберите основу';
            } else {
                errors.base_id = '';
            }
            break;
        case 'materials_ids':
            if (sale.materials_ids.length === 0) {
                errors.materials_ids = 'Пожалуйста, выберите хотя бы один материал';
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
    validateField('base_id');
    validateField('materials_ids');
    validateField('price');

    return Object.values(errors).every(error => error === '');
}

async function submitForm() {
    if (isSubmitting.value) return;

    if (!validateForm()) {
        return;
    }

    try {
        isSubmitting.value = true;
        errorMessage.value = '';

        const formData = new FormData();

        const finalImages: string[] = [];

        for (const imageName of sale.images) {
            finalImages.push(imageName);
        }

        for (const file of addedFiles.value) {
            finalImages.push(file.name);
        }

        if (addedFiles.value.length > 0) {
            addedFiles.value.forEach((file) => {
                formData.append('images', file);
            });
        }

        let updatedImages: string[];
        if (previewImages.value && previewImages.value.length > 0) {
            updatedImages = previewImages.value
                .map(image => image.filename)
                .filter((filename): filename is string => !!filename);
        } else {
            updatedImages = [];
        }

        const saleDataToUpdate: UpdateSaleRequest = {
            ...sale,
            images: updatedImages,
        };

        formData.append('data', JSON.stringify(saleDataToUpdate));

        const strId = route.params.id;
        const response = await axios.put(
            SERVER_URL + 'sales/' + strId,
            formData,
            {
                headers: {
                    'Content-Type': 'multipart/form-data',
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
            }
        );

        if (response.status === 200) {
            showSuccessAlert.value = true;
            sale.images = finalImages;
            imagesToDelete.value = [];
            addedFiles.value = [];
            Object.assign(originalSale, { ...sale });
        } else {
            showErrorAlert.value = true;
            errorMessage.value = 'Не удалось обновить работу. Пожалуйста, попробуйте снова.';
        }
    } catch (error: any) {
        console.error('Error submitting form:', error);
        let errorMsg = 'Произошла ошибка при обновлении работы. Пожалуйста, попробуйте снова.';

        if (error.response?.status === 413) {
            errorMessage.value = 'Файлы слишком большие. Пожалуйста, загрузите меньшие изображения.';
        } else if (error.response?.status === 400) {
            errorMessage.value = 'Некорректные данные. Пожалуйста, проверьте введенные значения.';
        } else {
            errorMessage.value = 'Произошла ошибка при обновлении работы. Пожалуйста, попробуйте снова.';
        }
        showErrorAlert.value = true;
    } finally {
        isSubmitting.value = false;
    }
}

function resetForm() {
    Object.assign(sale, { ...originalSale });
    addedFiles.value = [];
    previewImages.value = [];
    imagesToDelete.value = [];
    loadPreviewImages();
    if (fileInput.value) {
        fileInput.value.value = '';
    }

    Object.keys(errors).forEach((key) => {
        errors[key] = '';
    });
    fileError.value = null;
}

function materialsToggleDropdown() {
    materialsDropdownOpen.value = !materialsDropdownOpen.value;
}

function toggleMaterial(materialId: number, checked: boolean) {
    if (checked) {
        if (!sale.materials_ids.includes(materialId)) {
            sale.materials_ids.push(materialId);
        }
    } else {
        const index = sale.materials_ids.indexOf(materialId);
        if (index > -1) {
            sale.materials_ids.splice(index, 1);
        }
    }
}

function basesToggleDropdown() {
    basesDropdownOpen.value = !basesDropdownOpen.value;
}

// Lifecycle
onMounted(async () => {
    if (materialStore.materials.length === 0 || materialStore.bases.length === 0) {
        await materialStore.fetchAll();
    }
    await loadSale();
});
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: #999;
}

.edit-sale-container {
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
    padding-right: 2.5rem; /* Make space for the arrow */
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

.textarea {
    resize: vertical;
    min-height: 100px;
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
    .edit-sale-container {
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

.multi-select-wrapper {
    position: relative;
    width: 100%;
}

.select-display {
    padding: 0.75rem;
    border: 1px solid #ced4da;
    border-radius: 0.25rem;
    cursor: pointer;
    background-color: white;
    min-height: 46px;
    display: flex;
    align-items: center;
    transition: all 0.3s;
}

.select-display:hover {
    border-color: #4a90e2;
}

.dropdown-options {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 1000;
    background: white;
    border: 1px solid #ced4da;
    border-radius: 0.25rem;
    max-height: 200px;
    overflow-y: auto;
    margin-top: 0.25rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.option-item {
    padding: 8px 12px;
    display: flex;
    align-items: center;
    cursor: pointer;
}

.option-item:hover {
    background-color: #f8f9fa;
}

.option-item input {
    margin-right: 8px;
}

.arrow {
    float: right;
    transition: transform 0.3s;
}

.arrow-up {
    transform: rotate(180deg);
}
</style>
