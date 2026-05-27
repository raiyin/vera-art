<template>
    <div class="add-work-container">
        <div class="header-section">
            <h1 class="page-title">
                {{ $t('admin_gallery_form.page_title') }}
            </h1>
            <p class="page-subtitle">
                {{ $t('admin_gallery_form.page_subtitle') }}
            </p>
        </div>

        <!-- Загрузчик -->
        <div v-if="isLoading" class="loading-container">
            <div class="loader" />
            <p>{{ $t('admin_gallery_form.loading') }}</p>
        </div>

        <form v-else class="work-form" @submit.prevent="submitForm">
            <!-- Поле для загрузки изображений -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_gallery_form.sections.images') }}
                </h2>
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.select_images') }}
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
                        <UInput
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
                                {{ $t('admin_gallery_form.labels.drag_drop_text') }}
                            </p>
                            <p class="upload-hint">
                                {{ $t('admin_gallery_form.labels.file_formats') }}
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
                                    $t('admin_gallery_form.aria_labels.remove_image', {
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
                    {{ $t('admin_gallery_form.sections.main_info') }}
                </h2>

                <!-- Название картины -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.name_ru') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model="work.name_ru"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_ru }"
                        :placeholder="$t('admin_gallery_form.placeholders.name_ru')"
                        @blur="validateField('name_ru')"
                    />
                    <div v-if="errors.name_ru" class="error-message">
                        {{ errors.name_ru }}
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.name_en') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model="work.name_en"
                        type="text"
                        required
                        class="form-control"
                        :class="{ 'is-invalid': errors.name_en }"
                        :placeholder="$t('admin_gallery_form.placeholders.name_en')"
                        @blur="validateField('name_en')"
                    />
                    <div v-if="errors.name_en" class="error-message">
                        {{ errors.name_en }}
                    </div>
                </div>

                <!-- Размеры картины -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_gallery_form.labels.dimensions') }} ({{ units }})
                        <span class="required">*</span>
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
                                :placeholder="$t('admin_gallery_form.labels.width')"
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
                                :placeholder="$t('admin_gallery_form.labels.height')"
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
                        >{{ $t('admin_gallery_form.labels.year') }}
                        <span class="required">*</span></label
                    >
                    <UInput
                        v-model.number="work.year"
                        type="number"
                        required
                        min="2000"
                        :max="new Date().getFullYear()"
                        class="form-control"
                        :class="{ 'is-invalid': errors.year }"
                        :placeholder="$t('admin_gallery_form.placeholders.year')"
                        @blur="validateField('year')"
                    />
                    <div v-if="errors.year" class="error-message">
                        {{ errors.year }}
                    </div>
                </div>
            </div>

            <!-- Технические характеристики -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_gallery_form.sections.tech_specs') }}
                </h2>

                <!-- Основа -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.base') }}
                        <span class="required">*</span></label
                    >
                    <USelect
                        v-model="work.base_id"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id }"
                        @blur="validateField('base_id')"
                    >
                        <option value="" disabled>
                            {{ $t('admin_gallery_form.placeholders.select_base') }}
                        </option>
                        <option v-for="base in bases" :key="base.id" :value="base.id">
                            {{
                                $i18n.locale === 'ru'
                                    ? `${base.base_ru}`
                                    : `${base.base_en}`
                            }}
                        </option>
                    </USelect>
                    <div v-if="errors.base_id" class="error-message">
                        {{ errors.base_id }}
                    </div>
                </div>

                <!-- Материал -->
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.materials') }}
                        <span class="required">{{
                            isMaterialsRequired ? '*' : ''
                        }}</span></label
                    >
                    <div class="multi-select-wrapper">
                        <div
                            class="select-display drop-down-arrow"
                            :class="{ 'is-invalid': errors.materials_ids }"
                            tabindex="0"
                            @click="materialsToggleDropdown"
                            @keydown.enter="materialsToggleDropdown"
                            @blur="validateField('materials_ids')"
                        >
                            {{
                                selectedMaterialsDisplay ||
                                $t('admin_gallery_form.placeholders.select_materials')
                            }}
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
                        <div v-if="errors.materials_ids" class="error-message">
                            {{ errors.materials_ids }}
                        </div>
                    </div>
                </div>
            </div>

            <!-- Описание -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_gallery_form.sections.additional_info') }}
                </h2>
                <div class="form-group">
                    <label class="form-label">{{
                        $t('admin_gallery_form.labels.description')
                    }}</label>
                    <textarea
                        v-model="work.descr"
                        class="form-control textarea"
                        :placeholder="$t('admin_gallery_form.placeholders.description')"
                        rows="4"
                        maxlength="500"
                    />
                    <div class="char-count">{{ work.descr.length }}/500</div>
                </div>
            </div>

            <!-- Тип работы -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_gallery_form.sections.work_type') }}
                </h2>
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.work_type') }}
                        <span class="required">*</span></label
                    >
                    <USelect
                        v-model.number="work.type"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.work_type }"
                        @blur="validateField('work_type')"
                    >
                        <option value="1" selected>
                            {{ $t('admin_gallery_form.work_types.painting') }}
                        </option>
                        <option value="2">
                            {{ $t('admin_gallery_form.work_types.illustration') }}
                        </option>
                        <option value="3">
                            {{ $t('admin_gallery_form.work_types.3d') }}
                        </option>
                    </USelect>
                    <div v-if="errors.work_type" class="error-message">
                        {{ errors.work_type }}
                    </div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <UButton type="button" class="btn btn-secondary" @click="resetForm">
                    {{ $t('admin_gallery_form.buttons.clear_form') }}
                </UButton>
                <UButton
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">{{
                        $t('admin_gallery_form.buttons.add_work')
                    }}</span>
                    <span v-else>
                        <span class="spinner" />
                        {{ $t('admin_gallery_form.buttons.submitting') }}
                    </span>
                </UButton>
            </div>
        </form>

        <!-- Toast notifications are handled via useToast() -->
    </div>
</template>

<script lang="ts">
import axios from 'axios';
import { defineComponent } from 'vue';
import type { CreateWorkDto, RequestResult } from '../../types';
import { useToast } from '@nuxt/ui/runtime/composables/index.js';
import { useMaterialStore } from '../../stores/MaterialStore';

const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;
const materialStore = useMaterialStore();

export default defineComponent({
    name: 'AddWork',
    data() {
        return {
            work: {
                width: 0,
                height: 0,
                year: new Date().getFullYear(),
                name_ru: '',
                name_en: '',
                base_id: 0,
                materials_ids: [] as number[],
                descr: '',
                images: [] as string[],
                type: 1,
            } as CreateWorkDto,
            files: [] as File[],
            previewImages: [] as { file: File; preview: string }[],
            isSubmitting: false,
            isLoading: true,
            loadError: null as string | null,
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
                work_type: '',
            } as Record<string, string>,
            errorMessage: '',
        };
    },
    computed: {
        bases() {
            return materialStore.bases;
        },
        materials() {
            return materialStore.materials;
        },
        selectedMaterialsDisplay() {
            if (this.work.materials_ids.length === 0) return '';
            const selectedNames = materialStore.materials
                .filter((material) => this.work.materials_ids.includes(material.id))
                .map((material) =>
                    this.$i18n.locale === 'ru'
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
                this.files.length > 0 &&
                this.work.type > 0
            );
        },
        units(): string {
            return this.work.type <= 1
                ? this.$t('admin_gallery_form.units.cm')
                : this.$t('admin_gallery_form.units.px');
        },
        isMaterialsRequired() {
            return this.work.type === 1 || this.work.type === 2;
        },
    },
    async created() {
        if (materialStore.materials.length === 0 || materialStore.bases.length === 0) {
            await materialStore.fetchAll();
        }
    },
    methods: {
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
                this.fileError = this.$t('admin_gallery_form.errors.max_files');
                return;
            }

            // Проверка типов файлов
            const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
            const invalidFiles = selectedFiles.filter(
                (file) => !validTypes.includes(file.type)
            );

            if (invalidFiles.length > 0) {
                this.fileError = this.$t('admin_gallery_form.errors.invalid_file_type');
                return;
            }

            // Проверка размера файлов (макс. 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

            if (largeFiles.length > 0) {
                this.fileError = this.$t('admin_gallery_form.errors.file_size');
                return;
            }

            // Добавляем новые файлы
            this.files = [...this.files, ...selectedFiles];
            this.work.images = this.files.map((file) => file.name);

            // Создаем превью для новых изображений
            selectedFiles.forEach((file) => {
                const reader = new FileReader();
                reader.onload = (e) => {
                    this.previewImages.push({
                        file,
                        preview: e.target?.result as string,
                    });
                };
                reader.readAsDataURL(file);
            });
        },
        removeImage(index: number) {
            this.previewImages.splice(index, 1);
            this.files.splice(index, 1);
            this.work.images.splice(index, 1);
        },
        validateField(fieldName: string) {
            switch (fieldName) {
                case 'name_ru':
                    if (!this.work.name_ru.trim()) {
                        this.errors.name_ru = this.$t(
                            'admin_gallery_form.errors.name_ru_required'
                        );
                    } else {
                        this.errors.name_ru = '';
                    }
                    break;
                case 'name_en':
                    if (!this.work.name_en.trim()) {
                        this.errors.name_en = this.$t(
                            'admin_gallery_form.errors.name_en_required'
                        );
                    } else {
                        this.errors.name_en = '';
                    }
                    break;
                case 'width':
                    if (this.work.width <= 0) {
                        this.errors.width = this.$t(
                            'admin_gallery_form.errors.width_required'
                        );
                    } else {
                        this.errors.width = '';
                    }
                    break;
                case 'height':
                    if (this.work.height <= 0) {
                        this.errors.height = this.$t(
                            'admin_gallery_form.errors.height_required'
                        );
                    } else {
                        this.errors.height = '';
                    }
                    break;
                case 'year':
                    if (
                        this.work.year < 2000 ||
                        this.work.year > new Date().getFullYear()
                    ) {
                        this.errors.year = this.$t(
                            'admin_gallery_form.errors.year_range',
                            { year: new Date().getFullYear() }
                        );
                    } else {
                        this.errors.year = '';
                    }
                    break;
                case 'base_id':
                    if (this.work.base_id <= 0) {
                        this.errors.base_id = this.$t(
                            'admin_gallery_form.errors.base_required'
                        );
                    } else {
                        this.errors.base_id = '';
                    }
                    break;
                case 'materials_ids':
                    if (this.work.materials_ids.length === 0) {
                        this.errors.materials_ids = this.$t(
                            'admin_gallery_form.errors.materials_required'
                        );
                    } else {
                        this.errors.materials_ids = '';
                    }
                    break;
                case 'work_type':
                    if (this.work.type < 0) {
                        this.errors.work_type = this.$t(
                            'admin_gallery_form.errors.work_type_required'
                        );
                    } else {
                        this.errors.work_type = '';
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
            this.validateField('work_type');

            // Проверка наличия изображений
            if (this.files.length === 0) {
                this.fileError = this.$t('admin_gallery_form.errors.images_required');
                return false;
            }

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

                // Добавляем файлы
                this.files.forEach((file) => {
                    formData.append('images', file);
                });

                // Добавляем остальные данные
                const workData = {
                    ...this.work,
                };

                formData.append('data', JSON.stringify(workData));

                const response = await axios.post(SERVER_URL + 'works', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                });

                if (response.status === 200) {
                    const toast = useToast();
                    toast.add({
                        title: this.$t('toast.success.title'),
                        description: this.$t('toast.success.description'),
                        icon: 'i-heroicons-check-circle',
                        color: 'success',
                        duration: 5000,
                    });
                    this.resetForm();
                } else {
                    const toast = useToast();
                    toast.add({
                        title: this.$t('toast.error.title'),
                        description: this.$t('toast.error.description'),
                        icon: 'i-heroicons-exclamation-triangle',
                        color: 'error',
                        duration: 5000,
                    });
                    this.errorMessage = this.$t(
                        'admin_gallery_form.messages.submit_failed'
                    );
                }
            } catch (error: any) {
                console.error('Error submitting form:', error);
                const toast = useToast();
                let description = this.$t('toast.error.description');
                if (error.response?.status === 413) {
                    description = this.$t('admin_gallery_form.messages.file_too_large');
                } else if (error.response?.status === 400) {
                    description = this.$t('admin_gallery_form.messages.invalid_data');
                } else {
                    description = this.$t('admin_gallery_form.messages.general_error');
                }
                toast.add({
                    title: this.$t('toast.error.title'),
                    description,
                    icon: 'i-heroicons-exclamation-triangle',
                    color: 'error',
                    duration: 5000,
                });
            } finally {
                this.isSubmitting = false;
            }
        },
        resetForm() {
            this.work = {
                width: 0,
                height: 0,
                year: new Date().getFullYear(),
                name_ru: '',
                name_en: '',
                base_id: 0,
                materials_ids: [],
                descr: '',
                type: 0,
                images: [],
            };
            this.files = [];
            this.previewImages = [];
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
    },
});
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: red;
}

.add-work-container {
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
    .add-work-container {
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
