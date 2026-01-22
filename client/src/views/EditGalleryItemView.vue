<template>
    <div class="edit-work-container">
        <div class="header-section">
            <h1 class="page-title">Редактировать работу в галерее</h1>
            <p class="page-subtitle">Измените необходимые поля и сохраните изменения</p>
        </div>

        <!-- Загрузчик -->
        <div v-if="isLoading" class="loading-container">
            <div class="loader"></div>
            <p>Загрузка данных...</p>
        </div>

        <form v-else @submit.prevent="submitForm" class="work-form">
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
                        <input
                            type="file"
                            @change="handleFileUpload"
                            multiple
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            ref="fileInput"
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
                                ></path>
                                <polyline points="17 8 12 3 7 8"></polyline>
                                <line x1="12" y1="3" x2="12" y2="15"></line>
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
                    <div class="preview-container" v-if="previewImages.length > 0">
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
                            <button
                                type="button"
                                @click="removeImage(index)"
                                class="remove-btn"
                                :aria-label="`Удалить изображение ${index + 1}`"
                            >
                                &times;
                            </button>
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
                    <input
                        type="text"
                        v-model="work.name_ru"
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
                    <input
                        type="text"
                        v-model="work.name_en"
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
                            <input
                                type="number"
                                v-model.number="work.width"
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
                            <input
                                type="number"
                                v-model.number="work.height"
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
                    <input
                        type="number"
                        v-model.number="work.year"
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
                    <select
                        v-model="work.base_id"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id }"
                        @blur="validateField('base_id')"
                    >
                        <option value="" disabled>Выберите основу</option>
                        <option v-for="base in bases" :key="base.id" :value="base.id">
                            {{
                                $i18n.locale === 'RUS'
                                    ? `${base.base_ru}`
                                    : `${base.base_en}`
                            }}
                        </option>
                    </select>
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
                    <div class="multi-select-wrapper">
                        <div
                            class="select-display drop-down-arrow"
                            @click="materialsToggleDropdown"
                            :class="{ 'is-invalid': errors.materials_ids }"
                            tabindex="0"
                            @keydown.enter="materialsToggleDropdown"
                            @blur="validateField('materials_ids')"
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
                                <input
                                    type="checkbox"
                                    :id="'material-' + material.id"
                                    :value="material.id"
                                    v-model="work.materials_ids"
                                />
                                <label :for="'material-' + material.id">
                                    {{
                                        $i18n.locale === 'RUS'
                                            ? material.material_ru
                                            : material.material_en
                                    }}
                                </label>
                            </div>
                        </div>
                        <div v-if="errors.materials_ids" class="error-message">
                            {{ errors.materials_ids }}
                        </div>
                    </div>
                </div>
            </div>

            <!-- Описание -->
            <div class="form-section">
                <h2 class="section-title">Дополнительная информация</h2>
                <div class="form-group">
                    <label class="form-label">Описание</label>
                    <textarea
                        v-model="work.descr"
                        class="form-control textarea"
                        placeholder="Краткое описание картины"
                        rows="4"
                        maxlength="500"
                    ></textarea>
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
                    <select
                        v-model="work.type"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.type }"
                        @blur="validateField('type')"
                    >
                        <option value="1" selected>Картина</option>
                        <option value="2">Иллюстрация</option>
                        <option value="3">3D</option>
                    </select>
                    <div v-if="errors.type" class="error-message">
                        {{ errors.type }}
                    </div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <button type="button" @click="resetForm" class="btn btn-secondary">
                    Сбросить изменения
                </button>
                <button
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">Сохранить изменения</span>
                    <span v-else>
                        <span class="spinner"></span>
                        Сохранение...
                    </span>
                </button>
            </div>
        </form>

        <!-- Success Alert -->
        <Alert
            v-model="showSuccessAlert"
            type="success"
            title="Успешно!"
            message="Работа успешно обновлена в галерее."
            closeButtonText="Закрыть"
        />

        <!-- Danger Alert -->
        <Alert
            v-model="showErrorAlert"
            type="danger"
            title="Ошибка!"
            :message="
                errorMessage ||
                'Не удалось обновить работу в галерее. Пожалуйста, попробуйте снова.'
            "
            closeButtonText="Закрыть"
        />
    </div>
</template>

<script lang="ts">
import axios from 'axios';
import { defineComponent } from 'vue';
import { CreateWorkDto, Base, Material, RequestResult } from '@/types';
import Alert from '@/components/app-ui/Alert.vue';

interface Work extends CreateWorkDto {
    id: number;
    str_id: string;
    dir: string;
    base_ru: string;
    base_en: string;
    type: number;
    removed_indices?: number[];
}

export default defineComponent({
    name: 'EditGalleryItemView',
    components: {
        Alert,
    },
    data() {
        return {
            work: {
                id: 0,
                str_id: '',
                dir: '',
                width: 0,
                height: 0,
                year: new Date().getFullYear(),
                name_ru: '',
                name_en: '',
                base_id: 0,
                base_ru: '',
                base_en: '',
                materials_ids: [] as number[],
                img_count: 0,
                descr: '',
                type: 0,
            } as Work,
            files: [] as File[],
            existingImageIndices: [] as number[], // Track indices of existing images
            previewImages: [] as {
                file?: File;
                preview: string;
                isExisting: boolean;
                index?: number;
            }[],
            isSubmitting: false,
            bases: [] as Base[],
            materials: [] as Material[],
            isLoading: true,
            loadError: null as string | null,
            server: import.meta.env.VITE_SERVER_URL,
            requestResult: 'unknown' as RequestResult,
            materialsDropdownOpen: false,
            basesDropdownOpen: false,
            isDragOver: false,
            fileError: null as string | null,
            errors: {
                name_ru: '',
                name_en: '',
                width: '',
                height: '',
                year: '',
                base_id: '',
                materials_ids: '',
                type: '',
            } as Record<string, string>,
            errorMessage: '',
            showSuccessAlert: false,
            showErrorAlert: false,
            originalWork: {} as Work,
        };
    },
    async created() {
        await this.loadBases();
        await this.loadMaterials();
        await this.loadWork();
    },
    methods: {
        async loadBases() {
            try {
                const response = await axios.get(this.server + 'bases');
                this.bases = response.data;
                this.isLoading = false;
            } catch (error) {
                console.error('Ошибка при загрузке основ:', error);
                this.loadError = 'Не удалось загрузить список основ';
                this.isLoading = false;
                this.showErrorAlert = true;
                this.errorMessage =
                    'Не удалось загрузить данные. Пожалуйста, попробуйте позже.';
            }
        },
        async loadMaterials() {
            try {
                const response = await axios.get(this.server + 'materials');
                this.materials = response.data;
                this.isLoading = false;
            } catch (error) {
                console.error('Ошибка при загрузке материалов:', error);
                this.loadError = 'Не удалось загрузить список материалов';
                this.isLoading = false;
                this.showErrorAlert = true;
                this.errorMessage =
                    'Не удалось загрузить данные. Пожалуйста, попробуйте позже.';
            }
        },
        async loadWork() {
            try {
                const id = this.$route.params.id;
                const response = await axios.get(this.server + 'works/' + id);
                this.work = response.data;
                this.originalWork = { ...response.data };

                // Load existing images as previews
                this.loadPreviewImages();

                this.isLoading = false;
            } catch (error) {
                console.error('Ошибка при загрузке работы:', error);
                this.loadError = 'Не удалось загрузить работу';
                this.isLoading = false;
                this.showErrorAlert = true;
                this.errorMessage =
                    'Не удалось загрузить данные. Пожалуйста, попробуйте позже.';
            }
        },
        loadPreviewImages() {
            this.previewImages = [];
            this.existingImageIndices = [];
            for (let i = 1; i <= this.work.img_count; i++) {
                const imageUrl = `${this.work.dir}${i}.jpg`;
                this.previewImages.push({
                    preview: imageUrl,
                    isExisting: true,
                    index: i,
                });
                this.existingImageIndices.push(i);
            }
        },
        handleDragOver() {
            this.isDragOver = true;
        },
        handleDragLeave() {
            this.isDragOver = false;
        },
        handleDrop(event: DragEvent) {
            this.isDragOver = false;
            if (event.dataTransfer && event.dataTransfer.files.length) {
                const files = Array.from(event.dataTransfer.files);
                this.addImages(files);
            }
        },
        triggerFileInput() {
            (this.$refs.fileInput as HTMLInputElement)?.click();
        },
        handleFileUpload(event: Event) {
            const target = event.target as HTMLInputElement;
            if (target.files && target.files.length) {
                const files = Array.from(target.files);
                this.addImages(files);
            }
        },
        addImages(selectedFiles: File[]) {
            this.fileError = null;

            // Проверка на количество файлов
            if (this.files.length + selectedFiles.length > 10) {
                this.fileError = 'Можно загрузить не более 10 изображений';
                return;
            }

            // Проверка типов файлов
            const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
            const invalidFiles = selectedFiles.filter(
                (file) => !validTypes.includes(file.type)
            );

            if (invalidFiles.length > 0) {
                this.fileError =
                    'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
                return;
            }

            // Проверка размера файлов (макс. 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

            if (largeFiles.length > 0) {
                this.fileError = 'Размер каждого файла не должен превышать 5 МБ';
                return;
            }

            // Добавляем новые файлы
            this.files = [...this.files, ...selectedFiles];

            // Создаем превью для новых изображений
            selectedFiles.forEach((file) => {
                const reader = new FileReader();
                reader.onload = (e) => {
                    this.previewImages.push({
                        file,
                        preview: e.target?.result as string,
                        isExisting: false,
                    });
                };
                reader.readAsDataURL(file);
            });
        },
        removeImage(index: number) {
            const image = this.previewImages[index];

            // If it's an existing image, remove it from existingImageIndices
            if (image.isExisting && image.index !== undefined) {
                const existingIndex = this.existingImageIndices.indexOf(image.index);
                if (existingIndex !== -1) {
                    this.existingImageIndices.splice(existingIndex, 1);
                }
            } else {
                // If it's a new image, remove it from files array
                const fileIndex = this.previewImages
                    .slice(0, index)
                    .filter((img) => !img.isExisting).length;
                this.files.splice(fileIndex, 1);
            }

            // Remove from previewImages
            this.previewImages.splice(index, 1);
        },
        validateField(fieldName: string) {
            switch (fieldName) {
                case 'name_ru':
                    if (!this.work.name_ru.trim()) {
                        this.errors.name_ru = 'Пожалуйста, введите название на русском';
                    } else {
                        this.errors.name_ru = '';
                    }
                    break;
                case 'name_en':
                    if (!this.work.name_en.trim()) {
                        this.errors.name_en =
                            'Пожалуйста, введите название на английском';
                    } else {
                        this.errors.name_en = '';
                    }
                    break;
                case 'width':
                    if (this.work.width <= 0) {
                        this.errors.width = 'Ширина должна быть больше 0';
                    } else {
                        this.errors.width = '';
                    }
                    break;
                case 'height':
                    if (this.work.height <= 0) {
                        this.errors.height = 'Высота должна быть больше 0';
                    } else {
                        this.errors.height = '';
                    }
                    break;
                case 'year':
                    if (
                        this.work.year < 2000 ||
                        this.work.year > new Date().getFullYear()
                    ) {
                        this.errors.year = `Год должен быть между 2000 и ${new Date().getFullYear()}`;
                    } else {
                        this.errors.year = '';
                    }
                    break;
                case 'base_id':
                    if (this.work.base_id <= 0) {
                        this.errors.base_id = 'Пожалуйста, выберите основу';
                    } else {
                        this.errors.base_id = '';
                    }
                    break;
                case 'materials_ids':
                    if (this.work.materials_ids.length === 0) {
                        this.errors.materials_ids =
                            'Пожалуйста, выберите хотя бы один материал';
                    } else {
                        this.errors.materials_ids = '';
                    }
                    break;
                case 'type':
                    if (this.work.type < 0) {
                        this.errors.type = 'Пожалуйста, выберите тип работы';
                    } else {
                        this.errors.type = '';
                    }
                    break;
            }
        },
        validateForm() {
            this.validateField('name_ru');
            this.validateField('name_en');
            this.validateField('width');
            this.validateField('height');
            this.validateField('year');
            this.validateField('base_id');
            this.validateField('materials_ids');
            this.validateField('type');

            // Проверка отсутствия ошибок
            return Object.values(this.errors).every((error) => error === '');
        },
        async submitForm() {
            if (this.isSubmitting) return;

            if (!this.validateForm()) {
                return;
            }

            try {
                this.isSubmitting = true;
                this.errorMessage = '';

                // Формируем данные для отправки
                const formData = new FormData();

                // Добавляем файлы, если есть
                if (this.files.length > 0) {
                    this.files.forEach((file) => {
                        formData.append('images', file);
                    });
                }

                // Добавляем информацию об удаленных существующих изображениях
                const removedIndices = [];
                for (let i = 1; i <= this.originalWork.img_count; i++) {
                    if (!this.existingImageIndices.includes(i)) {
                        removedIndices.push(i);
                    }
                }

                // Добавляем остальные данные
                const workData = {
                    ...this.work,
                    img_count: this.existingImageIndices.length + this.files.length,
                    removed_indices: removedIndices, // Добавляем информацию об удаленных изображениях
                };
                workData.type = parseInt(workData.type as any);

                formData.append('data', JSON.stringify(workData));

                const strId = this.$route.params.id;
                const response = await axios.put(
                    this.server + 'works/' + strId,
                    formData,
                    {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                            Authorization: `Bearer ${localStorage.getItem('token')}`,
                        },
                    }
                );

                if (response.status === 200) {
                    this.showSuccessAlert = true;
                    // Обновляем оригинальную работу
                    this.originalWork = { ...this.work };
                } else {
                    this.showErrorAlert = true;
                    this.errorMessage =
                        'Не удалось обновить работу. Пожалуйста, попробуйте снова.';
                }
            } catch (error: any) {
                console.error('Error submitting form:', error);
                this.showErrorAlert = true;
                if (error.response?.status === 413) {
                    this.errorMessage =
                        'Файлы слишком большие. Пожалуйста, загрузите меньшие изображения.';
                } else if (error.response?.status === 400) {
                    this.errorMessage =
                        'Некорректные данные. Пожалуйста, проверьте введенные значения.';
                } else {
                    this.errorMessage =
                        'Произошла ошибка при обновлении работы. Пожалуйста, попробуйте снова.';
                }
            } finally {
                this.isSubmitting = false;
            }
        },
        resetForm() {
            this.work = { ...this.originalWork };
            this.files = [];
            this.previewImages = [];
            this.existingImageIndices = [];
            this.loadPreviewImages();
            if (this.$refs.fileInput) {
                (this.$refs.fileInput as HTMLInputElement).value = '';
            }

            // Сброс ошибок
            Object.keys(this.errors).forEach((key) => {
                this.errors[key] = '';
            });
            this.fileError = null;
        },
        materialsToggleDropdown() {
            this.materialsDropdownOpen = !this.materialsDropdownOpen;
        },
        basesToggleDropdown() {
            this.basesDropdownOpen = !this.basesDropdownOpen;
        },
        closeAlert() {
            this.showSuccessAlert = false;
            this.showErrorAlert = false;
            this.errorMessage = '';
        },
    },
    computed: {
        selectedMaterialsDisplay() {
            if (this.work.materials_ids.length === 0) return '';
            const selectedNames = this.materials
                .filter((material) => this.work.materials_ids.includes(material.id))
                .map((material) =>
                    this.$i18n.locale === 'RUS'
                        ? material.material_ru
                        : material.material_en
                );
            return selectedNames.join(', ');
        },
        isFormValid() {
            return (
                this.work.name_ru.trim() !== '' &&
                this.work.name_en.trim() !== '' &&
                this.work.width > 0 &&
                this.work.height > 0 &&
                this.work.year >= 2000 &&
                this.work.year <= new Date().getFullYear() &&
                this.work.base_id > 0 &&
                (this.work.type < 3
                    ? this.work.materials_ids.length > 0
                    : this.work.materials_ids.length == 0) &&
                this.work.type > 0
            );
        },
        units(): 'см' | 'px' {
            return this.work.type <= 1 ? 'см' : 'px';
        },
        isMaterialsRequired() {
            return this.work.type === '1' || this.work.type === '2';
        },
    },
});
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: red;
}

.edit-work-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 1rem;
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
    transition: border-color 0.3s, box-shadow 0.3s;
}

.form-control:focus {
    border-color: #4a90e2;
    outline: none;
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
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

.alert {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    padding: 1rem;
    border-radius: 6px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    max-width: 350px;
    z-index: 1000;
}

.alert-content {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    flex: 1;
}

.alert-icon {
    width: 1.5rem;
    height: 1.5rem;
    flex-shrink: 0;
}

.alert-success {
    background-color: #d4edda;
    border: 1px solid #c3e6cb;
    color: #155724;
}

.alert-success .alert-icon {
    color: #28a745;
}

.alert-danger {
    background-color: #f8d7da;
    border: 1px solid #f5c6cb;
    color: #721c24;
}

.alert-danger .alert-icon {
    color: #dc3545;
}

.btn-close {
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0;
    width: 1.5rem;
    height: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    color: inherit;
    opacity: 0.7;
}

.btn-close:hover {
    opacity: 1;
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

    .alert {
        right: 0.5rem;
        left: 0.5rem;
        max-width: none;
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
