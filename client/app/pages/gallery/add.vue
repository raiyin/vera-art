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
                        {{ $t('admin_gallery_form.labels.dimensions') }} ({{ $t('admin_gallery_form.units.cm') }})
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

            <!-- Основа -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_gallery_form.sections.tech_specs') }}
                </h2>
                <div class="form-group">
                    <label class="form-label"
                        >{{ $t('admin_gallery_form.labels.base') }}
                        <span class="required">*</span></label
                    >
                    <USelect
                        v-model="work.base_id"
                        :items="baseOptions"
                        required
                        class="form-control drop-down-arrow"
                        :class="{ 'is-invalid': errors.base_id }"
                        :placeholder="$t('admin_gallery_form.placeholders.select_base')"
                        @blur="validateField('base_id')"
                    />
                    <div v-if="errors.base_id" class="error-message">
                        {{ errors.base_id }}
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
                        $t('admin_gallery_form.labels.description_ru')
                    }}</label>
                    <UTextarea
                        v-model="work.descr_ru"
                        class="form-control"
                        :placeholder="$t('admin_gallery_form.placeholders.description')"
                        :rows="4"
                        :maxlength="500"
                    />
                    <div class="char-count">{{ work.descr_ru.length }}/500</div>
                </div>
                <div class="form-group">
                    <label class="form-label">{{
                        $t('admin_gallery_form.labels.description_en')
                    }}</label>
                    <UTextarea
                        v-model="work.descr_en"
                        class="form-control"
                        :placeholder="$t('admin_gallery_form.placeholders.description')"
                        :rows="4"
                        :maxlength="500"
                    />
                    <div class="char-count">{{ work.descr_en.length }}/500</div>
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
    </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import { ref, reactive, computed, onMounted } from 'vue';
import type { CreateWorkDto, RequestResult } from '~/types';
import { useMaterialStore } from '~/stores/MaterialStore';

const { t, locale } = useI18n();
const toast = useToast();
const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;

const materialStore = useMaterialStore();

const work = reactive<CreateWorkDto>({
    width: 0,
    height: 0,
    year: new Date().getFullYear(),
    name_ru: '',
    name_en: '',
    base_id: 0,
    descr_ru: '',
    descr_en: '',
});

const isSubmitting = ref(false);
const isLoading = ref(true);

const errors = reactive<Record<string, string>>({
    name_ru: '',
    name_en: '',
    width: '',
    height: '',
    year: '',
    base_id: '',
});

const bases = computed(() => materialStore.bases);

const baseOptions = computed(() => {
    const options = bases.value.map((base) => ({
        label: locale.value === 'ru' ? base.name_ru : base.name_en,
        value: base.id,
    }));
    return [...options];
});

const isFormValid = computed(() => {
    return (
        work.name_ru.trim() !== '' &&
        work.name_en.trim() !== '' &&
        work.width > 0 &&
        work.height > 0 &&
        work.year >= 2000 &&
        work.year <= new Date().getFullYear() &&
        work.base_id > 0
    );
});

function validateField(fieldName: string) {
    switch (fieldName) {
        case 'name_ru':
            errors.name_ru = !work.name_ru.trim() ? t('admin_gallery_form.errors.name_ru_required') : '';
            break;
        case 'name_en':
            errors.name_en = !work.name_en.trim() ? t('admin_gallery_form.errors.name_en_required') : '';
            break;
        case 'width':
            errors.width = work.width <= 0 ? t('admin_gallery_form.errors.width_required') : '';
            break;
        case 'height':
            errors.height = work.height <= 0 ? t('admin_gallery_form.errors.height_required') : '';
            break;
        case 'year':
            if (work.year < 2000 || work.year > new Date().getFullYear()) {
                errors.year = t('admin_gallery_form.errors.year_range', {
                    year: new Date().getFullYear(),
                });
            } else {
                errors.year = '';
            }
            break;
        case 'base_id':
            errors.base_id = work.base_id <= 0 ? t('admin_gallery_form.errors.base_required') : '';
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
    return Object.values(errors).every((error) => error === '');
}

async function submitForm() {
    if (isSubmitting.value) return;

    if (!validateForm()) {
        return;
    }

    try {
        isSubmitting.value = true;

        const response = await axios.post(SERVER_URL + 'works', work, {
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (response.status === 201 || response.status === 200) {
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
                description: t('toast.error.description'),
                icon: 'i-heroicons-exclamation-triangle',
                color: 'error',
                duration: 5000,
            });
        }
    } catch (error: any) {
        console.error('Error submitting form:', error);
        let description = t('toast.error.description');
        if (error.response?.status === 400) {
            description = t('admin_gallery_form.messages.invalid_data');
        } else {
            description = t('admin_gallery_form.messages.general_error');
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
    work.width = 0;
    work.height = 0;
    work.year = new Date().getFullYear();
    work.name_ru = '';
    work.name_en = '';
    work.base_id = 0;
    work.descr_ru = '';
    work.descr_en = '';
    Object.keys(errors).forEach((key) => {
        errors[key] = '';
    });
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

.add-work-container {
    max-width: 800px;
    margin: 3rem auto;
    padding: 2.5rem;
    background-color: var(--color-on-surface);
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

:root.dark .add-work-container {
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

    .page-title {
        font-size: 1.5rem;
    }

    .section-title {
        font-size: 1.2rem;
    }
}
</style>
