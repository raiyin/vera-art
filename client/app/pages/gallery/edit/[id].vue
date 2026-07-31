<template>
    <div class="edit-work-container">
        <div class="header-section">
            <h1 class="page-title">
                Редактировать работу в галерее
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
            class="work-form"
            @submit.prevent="submitForm"
        >
            <!-- Основная информация -->
            <div class="form-section">
                <h2 class="section-title">
                    Основная информация
                </h2>

                <!-- Название картины -->
                <div class="form-group">
                    <label class="form-label">Название картины по-русски <span class="required">*</span></label>
                    <UInput
                        v-model="work.name_ru"
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
                        v-model="work.name_en"
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
                                v-model.number="work.width"
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
                                v-model.number="work.height"
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
                        v-model.number="work.year"
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
            </div>

            <!-- Основа -->
            <div class="form-section">
                <h2 class="section-title">
                    Технические характеристики
                </h2>
                <div class="form-group">
                    <label class="form-label">Основа <span class="required">*</span></label>
                    <USelect
                        v-model="work.base_id"
                        :items="baseOptions"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id, }"
                        placeholder="Выберите основу"
                        @blur="validateField('base_id',)"
                    />
                    <div
                        v-if="errors.base_id"
                        class="error-message"
                    >
                        {{ errors.base_id }}
                    </div>
                </div>
            </div>

            <!-- Описание -->
            <div class="form-section">
                <h2 class="section-title">
                    Дополнительная информация
                </h2>
                <div class="form-group">
                    <label class="form-label">Описание (русский)</label>
                    <UTextarea
                        v-model="work.descr_ru"
                        class="form-control"
                        placeholder="Краткое описание картины на русском"
                        :rows="4"
                        :maxlength="500"
                    />
                    <div class="char-count">
                        {{ work.descr_ru.length }}/500
                    </div>
                </div>
                <div class="form-group">
                    <label class="form-label">Описание (английский)</label>
                    <UTextarea
                        v-model="work.descr_en"
                        class="form-control"
                        placeholder="Краткое описание картины на английском"
                        :rows="4"
                        :maxlength="500"
                    />
                    <div class="char-count">
                        {{ work.descr_en.length }}/500
                    </div>
                </div>
            </div>

            <!-- Изображения -->
            <div class="form-section">
                <h2 class="section-title">
                    Изображения
                </h2>

                <!-- Текущие изображения -->
                <div
                    v-if="work.images.length > 0"
                    class="form-group"
                >
                    <label class="form-label">Текущие изображения</label>
                    <div class="preview-container">
                        <div
                            v-for="(img, idx) in work.images"
                            :key="idx"
                            class="image-preview"
                        >
                            <img
                                :src="work.dir + img"
                                :alt="'Image ' + (idx + 1)"
                                class="preview-image"
                            >
                        </div>
                    </div>
                </div>

                <!-- Загрузка нового изображения -->
                <div class="form-group">
                    <label class="form-label">Новое изображение</label>
                    <div
                        class="file-drop-area"
                        @click="triggerFileInput"
                        @dragover.prevent="onDragOver"
                        @dragleave.prevent="onDragLeave"
                        @drop.prevent="onDrop"
                    >
                        <input
                            ref="fileInputRef"
                            type="file"
                            accept="image/*"
                            class="file-input"
                            @change="onFileSelected"
                        >
                        <div class="file-drop-content">
                            <svg
                                class="upload-icon"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                                />
                            </svg>
                            <p class="upload-text">
                                Нажмите или перетащите файл для загрузки
                            </p>
                            <p class="upload-hint">
                                PNG, JPG до 10MB
                            </p>
                        </div>
                    </div>
                    <div
                        v-if="selectedFile"
                        class="preview-container"
                    >
                        <div class="image-preview">
                            <img
                                :src="selectedFilePreview"
                                alt="New image preview"
                                class="preview-image"
                            >
                            <button
                                type="button"
                                class="remove-btn"
                                @click="selectedFile = null; selectedFilePreview = '';"
                            >
                                ×
                            </button>
                        </div>
                        <p class="upload-hint">
                            {{ selectedFile.name }}
                        </p>
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
    </div>
</template>

<script setup lang="ts">
    import axios from 'axios';
    import { ref, reactive, computed, onMounted, } from 'vue';
    import { useMaterialStore, } from '~/stores/MaterialStore';

    interface EditWorkData {
        str_id: string
        name_ru: string
        name_en: string
        width: number
        height: number
        year: number
        base_id: number
        descr_ru: string
        descr_en: string
        work_path: string
        dir: string
        images: string[]
    }

    definePageMeta({
        layout: 'admin',
        middleware: 'admin-auth',
    });

    const toast = useToast();
    const config = useRuntimeConfig();
    const route = useRoute();
    const SERVER_URL = config.public.serverUrl;

    const materialStore = useMaterialStore();

    const work = reactive<EditWorkData>({
        str_id: '',
        name_ru: '',
        name_en: '',
        width: 0,
        height: 0,
        year: new Date().getFullYear(),
        base_id: 0,
        descr_ru: '',
        descr_en: '',
        work_path: '',
        dir: '',
        images: [],
    });

    const originalWork = reactive<EditWorkData>({
        str_id: '',
        name_ru: '',
        name_en: '',
        width: 0,
        height: 0,
        year: new Date().getFullYear(),
        base_id: 0,
        descr_ru: '',
        descr_en: '',
        work_path: '',
        dir: '',
        images: [],
    });

    const isSubmitting = ref(false,);
    const isLoading = ref(true,);
    const selectedFile = ref<File | null>(null);
    const selectedFilePreview = ref<string>('');

    const errors = reactive<Record<string, string>>({
        name_ru: '',
        name_en: '',
        width: '',
        height: '',
        year: '',
        base_id: '',
    });

    const bases = computed(() => materialStore.bases,);

    const baseOptions = computed(() => {
        return bases.value.map(base => ({
            label: base.name_ru,
            value: base.id,
        }),);
    });

    const isFormValid = computed(() => {
        return (
            work.name_ru?.trim() !== ''
        && work.name_en?.trim() !== ''
            && (work.width ?? 0) > 0
        && (work.height ?? 0) > 0
            && (work.year ?? 0) >= 2000
        && (work.year ?? 0) <= new Date().getFullYear()
            && (work.base_id ?? 0) > 0
        );
    });

    async function loadWork() {
        try {
            const id = route.params.id;
            const response = await axios.get(`${SERVER_URL}works/${id}`,);
            const data = response.data;
            work.str_id = data.str_id ?? '';
            work.name_ru = data.name_ru;
            work.name_en = data.name_en;
            work.width = data.width;
            work.height = data.height;
            work.year = data.year;
            work.base_id = data.base_id;
            work.descr_ru = data.descr_ru ?? '';
            work.descr_en = data.descr_en ?? '';
            work.work_path = data.work_path ?? '';
            work.dir = data.dir ?? '';
            work.images = data.images ?? [];
            Object.assign(originalWork, { ...work, },);
            isLoading.value = false;
        } catch (error) {
            console.error('Ошибка при загрузке работы:', error,);
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

    function validateField(fieldName: string,) {
        switch (fieldName) {
    case 'name_ru':
        errors.name_ru = !work.name_ru?.trim() ? 'Пожалуйста, введите название на русском' : '';
        break;
    case 'name_en':
        errors.name_en = !work.name_en?.trim() ? 'Пожалуйста, введите название на английском' : '';
        break;
    case 'width':
        errors.width = (work.width ?? 0) <= 0 ? 'Ширина должна быть больше 0' : '';
        break;
    case 'height':
        errors.height = (work.height ?? 0) <= 0 ? 'Высота должна быть больше 0' : '';
        break;
    case 'year':
        if ((work.year ?? 0) < 2000 || (work.year ?? 0) > new Date().getFullYear()) {
            errors.year = `Год должен быть между 2000 и ${new Date().getFullYear()}`;
        } else {
            errors.year = '';
        }
        break;
    case 'base_id':
        errors.base_id = (work.base_id ?? 0) <= 0 ? 'Пожалуйста, выберите основу' : '';
        break;
        }
    }

    function validateForm() {
        validateField('name_ru',);
        validateField('name_en',);
        validateField('width',);
        validateField('height',);
        validateField('year',);
        validateField('base_id',);
        return Object.values(errors,).every(error => error === '',);
    }

    async function submitForm() {
        if (isSubmitting.value) return;

        if (!validateForm()) {
            return;
        }

        try {
            isSubmitting.value = true;

            const id = route.params.id;
            const formData = new FormData();
            formData.append('str_id', work.str_id);
            formData.append('name_ru', work.name_ru);
            formData.append('name_en', work.name_en);
            formData.append('width', String(work.width));
            formData.append('height', String(work.height));
            formData.append('year', String(work.year));
            formData.append('base_id', String(work.base_id));
            formData.append('descr_ru', work.descr_ru);
            formData.append('descr_en', work.descr_en);
            if (selectedFile.value) {
                formData.append('image', selectedFile.value);
            }

            const response = await axios.put(SERVER_URL + 'works/' + id, formData, {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token',)}`,
                },
            });

            if (response.status === 200) {
                const resData = response.data;
                work.work_path = resData.work_path ?? work.work_path;
                work.dir = resData.dir ?? work.dir;
                work.images = resData.images ?? work.images;
                work.str_id = resData.str_id ?? work.str_id;
                selectedFile.value = null;
                selectedFilePreview.value = '';
                Object.assign(originalWork, { ...work, },);
                toast.add({
                    title: 'Успешно!',
                    description: 'Работа успешно обновлена в галерее.',
                    icon: 'i-heroicons-check-circle',
                    color: 'success',
                    duration: 5000,
                });
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
            console.error('Error submitting form:', error,);
            let description = 'Произошла ошибка при обновлении работы. Пожалуйста, попробуйте снова.';
            if (error.response?.status === 400) {
                description = 'Некорректные данные. Пожалуйста, проверьте введенные значения.';
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

    const fileInputRef = ref<HTMLInputElement | null>(null);

    function triggerFileInput() {
        fileInputRef.value?.click();
    }

    function onDragOver(e: DragEvent) {
        const target = e.currentTarget as HTMLElement;
        target.classList.add('drag-over');
    }

    function onDragLeave(e: DragEvent) {
        const target = e.currentTarget as HTMLElement;
        target.classList.remove('drag-over');
    }

    function onDrop(e: DragEvent) {
        const target = e.currentTarget as HTMLElement;
        target.classList.remove('drag-over');
        const file = e.dataTransfer?.files?.[0];
        if (file && fileInputRef.value) {
            const dt = new DataTransfer();
            dt.items.add(file);
            fileInputRef.value.files = dt.files;
            onFileSelected({ target: fileInputRef.value } as unknown as Event);
        }
    }

    function onFileSelected(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0];
        if (file) {
            selectedFile.value = file;
            const reader = new FileReader();
            reader.onload = (e) => {
                selectedFilePreview.value = e.target?.result as string;
            };
            reader.readAsDataURL(file);
        }
    }

    function resetForm() {
        Object.assign(work, { ...originalWork, },);
        Object.keys(errors,).forEach((key,) => {
            errors[key] = '';
        });
        selectedFile.value = null;
        selectedFilePreview.value = '';
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
