<script setup lang="ts">
definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

import { ref, reactive, computed } from 'vue';
import type { NewsDesc } from '~/types';
import { getHttpClient } from '~/api/http-client';

const api = getHttpClient();

const { t } = useI18n();
const toast = useToast();
const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;

interface PreviewItem {
    file: File;
    preview: string;
}

type NewsForm = Omit<NewsDesc, 'id'>;

// Reactive state
const news = reactive<NewsForm>({
    title_en: '',
    title_ru: '',
    img_back: '',
    img_backfull: '',
    images: [] as string[],
    videos: [] as string[],
    datetime: '',
    text_en: '',
    text_ru: '',
    dir: '',
});

const images = ref<File[]>([]);
const videos = ref<File[]>([]);

const img_back_preview = ref<PreviewItem | null>(null);
const img_backfull_preview = ref<PreviewItem | null>(null);

const previewImages = ref<PreviewItem[]>([]);
const previewVideos = ref<string[]>([]);

const previewWidth = ref(400);
const isSubmitting = ref(false);
const fileError = ref<string | null>(null);
const isDragOver = ref(false);

const errors = reactive<Record<string, string>>({
    title_ru: '',
    title_en: '',
    img_back: '',
    img_backfull: '',
    images: '',
    text_ru: '',
    text_en: '',
});

// Template refs
const backFullInput = ref<HTMLInputElement | null>(null);
const backInput = ref<HTMLInputElement | null>(null);
const imagesInput = ref<HTMLInputElement | null>(null);
const videosInput = ref<HTMLInputElement | null>(null);

// Computed
const isFormValid = computed(() => {
    return (
        news.title_ru.trim() !== '' &&
        news.title_en.trim() !== '' &&
        img_back_preview.value !== null &&
        img_backfull_preview.value !== null &&
        previewImages.value.length > 0 &&
        news.text_ru.trim() !== '' &&
        news.text_en.trim() !== ''
    );
});

// Methods
function handleDragOver() {
    isDragOver.value = true;
}

function handleDragLeave() {
    isDragOver.value = false;
}

function handleDrop(event: DragEvent, target: 'backFull' | 'back' | 'images' | 'videos') {
    isDragOver.value = false;
    const files = Array.from(event.dataTransfer?.files || []);
    if (files.length === 0) return;

    switch (target) {
        case 'backFull':
            handleBackFullImageDrop(files[0]!);
            break;
        case 'back':
            handleBackImageDrop(files[0]!);
            break;
        case 'images':
            handleImagesDrop(files);
            break;
        case 'videos':
            handleVideosDrop(files);
            break;
    }
}

function handleBackImageDrop(file: File) {
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    if (!validTypes.includes(file.type)) {
        fileError.value = t('admin_news_form.errors.invalid_image_type');
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    if (file.size > maxSize) {
        fileError.value = t('admin_news_form.errors.image_size');
        return;
    }

    news.img_back = file.name;
    fileError.value = null;

    const reader = new FileReader();
    reader.onload = (e) => {
        img_back_preview.value = {
            file,
            preview: e.target?.result as string,
        };
    };
    reader.readAsDataURL(file);
}

function handleBackFullImageDrop(file: File) {
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    if (!validTypes.includes(file.type)) {
        fileError.value = t('admin_news_form.errors.invalid_image_type');
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    if (file.size > maxSize) {
        fileError.value = t('admin_news_form.errors.image_size');
        return;
    }

    news.img_backfull = file.name;
    fileError.value = null;

    const reader = new FileReader();
    reader.onload = (e) => {
        img_backfull_preview.value = {
            file,
            preview: e.target?.result as string,
        };
    };
    reader.readAsDataURL(file);
}

function handleImagesDrop(files: File[]) {
    if (images.value.length + files.length > 10) {
        fileError.value = t('admin_news_form.errors.max_images');
        return;
    }

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    const invalidFiles = files.filter((file) => !validTypes.includes(file.type));
    if (invalidFiles.length > 0) {
        fileError.value = t('admin_news_form.errors.invalid_image_type');
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    const largeFiles = files.filter((file) => file.size > maxSize);
    if (largeFiles.length > 0) {
        fileError.value = t('admin_news_form.errors.image_size');
        return;
    }

    fileError.value = null;

    images.value = [...images.value, ...files];

    files.forEach((file) => {
        news.images.push(file.name);
    });

    files.forEach((file) => {
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

function handleVideosDrop(files: File[]) {
    for (const file of files) {
        if (!file.type.startsWith('video/')) {
            fileError.value = t('admin_news_form.errors.invalid_video_type');
            return;
        }
    }

    const maxSize = 100 * 1024 * 1024;
    for (const file of files) {
        if (file.size > maxSize) {
            fileError.value = t('admin_news_form.errors.video_size');
            return;
        }
    }

    fileError.value = null;

    for (const file of files) {
        videos.value.push(file);
        news.videos.push(file.name);
        previewVideos.value.push(URL.createObjectURL(file));
    }
}

function triggerFileInput(refName: string) {
    let input: HTMLInputElement | null = null;
    switch (refName) {
        case 'backFullInput':
            input = backFullInput.value;
            break;
        case 'backInput':
            input = backInput.value;
            break;
        case 'imagesInput':
            input = imagesInput.value;
            break;
        case 'videosInput':
            input = videosInput.value;
            break;
    }
    if (input) {
        input.click();
    }
}

function handleBackImageSelected(event: Event) {
    const target = event.target as HTMLInputElement;
    const selectedImage = target.files?.[0];

    if (!selectedImage) return;

    // Validate file type
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    if (!validTypes.includes(selectedImage.type)) {
        fileError.value = t('admin_news_form.errors.invalid_image_type');
        return;
    }

    // Validate file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (selectedImage.size > maxSize) {
        fileError.value = t('admin_news_form.errors.image_size');
        return;
    }

    // Store file object for form submission, but keep string for type compatibility
    news.img_back = selectedImage.name;
    fileError.value = null;

    const reader = new FileReader();
    reader.onload = (e) => {
        img_back_preview.value = {
            file: selectedImage,
            preview: e.target?.result as string,
        };
    };
    reader.readAsDataURL(selectedImage);
}

function removeBackImage() {
    img_back_preview.value = null;
    news.img_back = '';
}

function handleBackFullImageSelected(event: Event) {
    const target = event.target as HTMLInputElement;
    const selectedImage = target.files?.[0];

    if (!selectedImage) return;

    // Validate file type
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    if (!validTypes.includes(selectedImage.type)) {
        fileError.value = t('admin_news_form.errors.invalid_image_type');
        return;
    }

    // Validate file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (selectedImage.size > maxSize) {
        fileError.value = t('admin_news_form.errors.image_size');
        return;
    }

    // Store file name for type compatibility, but keep file object for submission
    news.img_backfull = selectedImage.name;
    fileError.value = null;

    const reader = new FileReader();
    reader.onload = (e) => {
        img_backfull_preview.value = {
            file: selectedImage,
            preview: e.target?.result as string,
        };
    };
    reader.readAsDataURL(selectedImage);
}

function removeBackFullImage() {
    img_backfull_preview.value = null;
    news.img_backfull = '';
}

function handleImagesSelected(event: Event) {
    const target = event.target as HTMLInputElement;
    const selectedFiles = Array.from(target.files || []);

    // Проверка на количество файлов
    if (images.value.length + selectedFiles.length > 10) {
        fileError.value = t('admin_news_form.errors.max_images');
        return;
    }

    // Проверка типов файлов
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    const invalidFiles = selectedFiles.filter((file) => !validTypes.includes(file.type));

    if (invalidFiles.length > 0) {
        fileError.value = t('admin_news_form.errors.invalid_image_type');
        return;
    }

    // Проверка размера файлов (макс. 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

    if (largeFiles.length > 0) {
        fileError.value = t('admin_news_form.errors.image_size');
        return;
    }

    fileError.value = null;

    // Добавляем новые файлы
    images.value = [...images.value, ...selectedFiles];

    // Добавляем имена файлов в news.images
    selectedFiles.forEach((file) => {
        news.images.push(file.name);
    });

    // Создаем превью для новых изображений
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

async function handleVideosSelected(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files as FileList;

    if (!files || files.length == 0) return;

    // Validate file type
    for (let i = 0; i < files.length; i++) {
        const file = files[i]!;
        if (!file.type.startsWith('video/')) {
            fileError.value = t('admin_news_form.errors.invalid_video_type');
            return;
        }
    }

    // Validate file size (max 100MB)
    const maxSize = 100 * 1024 * 1024;
    for (let i = 0; i < files.length; i++) {
        const file = files[i]!;
        if (file.size > maxSize) {
            fileError.value = t('admin_news_form.errors.video_size');
            return;
        }
    }

    fileError.value = null;

    // Populate both arrays
    for (let i = 0; i < files.length; i++) {
        const file = files[i]!;
        videos.value.push(file);
        news.videos.push(file.name);
        previewVideos.value.push(URL.createObjectURL(file));
    }
}

function removeImageFromImages(index: number) {
    previewImages.value.splice(index, 1);
    images.value.splice(index, 1);
    news.images.splice(index, 1);
}

function removeVideoFromVideos(index: number) {
    previewVideos.value.splice(index, 1);
    videos.value.splice(index, 1);
    news.videos.splice(index, 1);
}

function onVideoError() {
    // Video error handling
}

function validateField(fieldName: string) {
    switch (fieldName) {
        case 'title_ru':
            if (!news.title_ru.trim()) {
                errors.title_ru = t('admin_news_form.errors.title_ru_required');
            } else if (
                news.title_ru.trim().length < 3 ||
                news.title_ru.trim().length > 50
            ) {
                errors.title_ru = t('admin_news_form.errors.title_ru_length');
            } else {
                errors.title_ru = '';
            }
            break;
        case 'title_en':
            if (!news.title_en.trim()) {
                errors.title_en = t('admin_news_form.errors.title_en_required');
            } else if (
                news.title_en.trim().length < 3 ||
                news.title_en.trim().length > 50
            ) {
                errors.title_en = t('admin_news_form.errors.title_en_length');
            } else {
                errors.title_en = '';
            }
            break;
        case 'img_back':
            if (!img_back_preview.value || !news.img_back) {
                errors.img_back = t('admin_news_form.errors.preview_image_required');
            } else {
                errors.img_back = '';
            }
            break;
        case 'img_backfull':
            if (!img_backfull_preview.value || !news.img_backfull) {
                errors.img_backfull = t('admin_news_form.errors.main_image_required');
            } else {
                errors.img_backfull = '';
            }
            break;
        case 'images':
            if (previewImages.value.length === 0) {
                errors.images = t('admin_news_form.errors.images_required');
            } else {
                errors.images = '';
            }
            break;
        case 'text_ru':
            if (!news.text_ru.trim()) {
                errors.text_ru = t('admin_news_form.errors.text_ru_required');
            } else {
                errors.text_ru = '';
            }
            break;
        case 'text_en':
            if (!news.text_en.trim()) {
                errors.text_en = t('admin_news_form.errors.text_en_required');
            } else {
                errors.text_en = '';
            }
            break;
    }
}

function validateForm() {
    validateField('title_ru');
    validateField('title_en');
    validateField('img_back');
    validateField('img_backfull');
    validateField('images');
    validateField('text_ru');
    validateField('text_en');

    // Проверка отсутствия ошибок
    return Object.values(errors).every((error) => error === '');
}

/**
 * Submit the news form data to the server
 * Sends news information along with images and videos
 */
async function submitForm() {
    if (isSubmitting.value) return;

    if (!validateForm()) {
        return;
    }

    try {
        isSubmitting.value = true;

        // Формируем данные для отправки
        const formData = new FormData();

        // Append main images
        if (img_back_preview.value?.file) {
            formData.append('img_back', img_back_preview.value.file);
        }

        if (img_backfull_preview.value?.file) {
            formData.append('img_backfull', img_backfull_preview.value.file);
        }
        // Add additional images and videos
        images.value.forEach((image) => {
            formData.append('images', image);
        });

        videos.value.forEach((video) => {
            formData.append('videos', video);
        });

        // Add news data as JSON
        const newsData = {
            ...news,
            datetime: news.datetime || new Date().toISOString().split('T')[0],
        };

        formData.append('data', JSON.stringify(newsData));

        // Send data to server
        const response = await api.post(`${SERVER_URL}news`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });

        if (response.status === 200 || response.status === 201) {
            toast.add({
                title: t('toast.success.title'),
                description: t('admin_news_form.messages.submit_success'),
                icon: 'i-heroicons-check-circle',
                color: 'success',
                duration: 5000,
            });
            resetForm();
        } else {
            toast.add({
                title: t('toast.error.title'),
                description: t('admin_news_form.messages.submit_failed'),
                icon: 'i-heroicons-exclamation-triangle',
                color: 'error',
                duration: 5000,
            });
        }
    } catch (error: any) {
        console.error('Error submitting form:', error);
        let description = t('admin_news_form.messages.general_error');
        if (error.response?.status === 413) {
            description = t('admin_news_form.messages.file_too_large');
        } else if (error.response?.status === 400) {
            description = t('admin_news_form.messages.invalid_data');
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
    news.datetime = '';
    news.title_en = '';
    news.title_ru = '';
    news.dir = '';
    news.img_back = '';
    news.img_backfull = '';
    news.text_en = '';
    news.text_ru = '';
    news.images = [];
    news.videos = [];

    images.value = [];
    videos.value = [];
    previewImages.value = [];
    previewVideos.value = [];
    img_back_preview.value = null;
    img_backfull_preview.value = null;

    // Сброс input файлов
    if (backFullInput.value) backFullInput.value.value = '';
    if (backInput.value) backInput.value.value = '';
    if (imagesInput.value) imagesInput.value.value = '';
    if (videosInput.value) videosInput.value.value = '';

    // Сброс ошибок
    Object.keys(errors).forEach((key) => {
        errors[key] = '';
    });
    fileError.value = null;
}
</script>

<template>
    <div class="add-news-container">
        <div class="header-section">
            <div class="header-content">
                <h1 class="page-title">{{ $t('admin_news_form.page_title') }}</h1>
            </div>
        </div>

        <form @submit.prevent="submitForm" class="news-form">
            <!-- Изображения новости -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_news_form.sections.images') }}
                </h2>

                <!-- Главное изображение -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.main_image') }}
                        <span class="required">*</span>
                    </label>
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'backFull')"
                        @click="triggerFileInput('backFullInput')"
                    >
                        <input
                            type="file"
                            @change="handleBackFullImageSelected"
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            ref="backFullInput"
                        />
                        <div class="file-drop-content">
                            <svg
                                class="upload-icon"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                                ></path>
                            </svg>
                            <p class="upload-text">
                                {{ $t('admin_news_form.labels.drag_drop_single') }}
                            </p>
                            <p class="upload-hint">
                                {{ $t('admin_news_form.hints.image_formats') }}
                            </p>
                        </div>
                    </div>
                    <div class="error-message" v-if="errors.img_backfull">
                        {{ errors.img_backfull }}
                    </div>
                    <div class="preview-container" v-if="img_backfull_preview">
                        <div class="image-preview">
                            <img
                                :src="img_backfull_preview.preview"
                                class="preview-image"
                                :alt="$t('admin_news_form.labels.main_image')"
                            />
                            <UButton
                                type="button"
                                @click="removeBackFullImage"
                                class="remove-btn"
                                :aria-label="
                                    $t('admin_news_form.aria_labels.remove_main_image')
                                "
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>

                <!-- Превью изображение новости -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.preview_image') }}
                        <span class="required">*</span>
                    </label>
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'back')"
                        @click="triggerFileInput('backInput')"
                    >
                        <input
                            type="file"
                            @change="handleBackImageSelected"
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            ref="backInput"
                        />
                        <div class="file-drop-content">
                            <svg
                                class="upload-icon"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                                ></path>
                            </svg>
                            <p class="upload-text">
                                {{ $t('admin_news_form.labels.drag_drop_single') }}
                            </p>
                            <p class="upload-hint">
                                {{ $t('admin_news_form.hints.image_formats') }}
                            </p>
                        </div>
                    </div>
                    <div class="error-message" v-if="errors.img_back">
                        {{ errors.img_back }}
                    </div>
                    <div class="preview-container" v-if="img_back_preview">
                        <div class="image-preview">
                            <img
                                :src="img_back_preview.preview"
                                class="preview-image"
                                :alt="$t('admin_news_form.labels.preview_image')"
                            />
                            <UButton
                                type="button"
                                @click="removeBackImage"
                                class="remove-btn"
                                :aria-label="
                                    $t('admin_news_form.aria_labels.remove_preview_image')
                                "
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>

                <!-- Фотогаллерея -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.photo_gallery') }}
                        <span class="required">*</span>
                    </label>
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'images')"
                        @click="triggerFileInput('imagesInput')"
                    >
                        <input
                            type="file"
                            @change="handleImagesSelected"
                            multiple
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            ref="imagesInput"
                        />
                        <div class="file-drop-content">
                            <svg
                                class="upload-icon"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                                ></path>
                            </svg>
                            <p class="upload-text">
                                {{ $t('admin_news_form.labels.drag_drop_multiple') }}
                            </p>
                            <p class="upload-hint">
                                {{ $t('admin_news_form.hints.image_formats_multiple') }}
                            </p>
                        </div>
                    </div>
                    <div class="error-message" v-if="errors.images">
                        {{ errors.images }}
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
                                :alt="`${$t('admin_news_form.labels.photo_gallery')} ${
                                    index + 1
                                }`"
                            />
                            <UButton
                                type="button"
                                @click="removeImageFromImages(index)"
                                class="remove-btn"
                                :aria-label="
                                    $t('admin_news_form.aria_labels.remove_image', {
                                        index: index + 1,
                                    })
                                "
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>

                <!-- Видеогаллерея -->
                <div class="form-group">
                    <label class="form-label">{{
                        $t('admin_news_form.labels.video_gallery')
                    }}</label>
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'videos')"
                        @click="triggerFileInput('videosInput')"
                    >
                        <input
                            type="file"
                            @change="handleVideosSelected"
                            multiple
                            accept="video/*"
                            class="file-input"
                            ref="videosInput"
                        />
                        <div class="file-drop-content">
                            <svg
                                class="upload-icon"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
                                ></path>
                            </svg>
                            <p class="upload-text">
                                {{ $t('admin_news_form.labels.drag_drop_videos') }}
                            </p>
                            <p class="upload-hint">
                                {{ $t('admin_news_form.hints.video_formats') }}
                            </p>
                        </div>
                    </div>
                    <div class="error-message" v-if="fileError">
                        {{ fileError }}
                    </div>
                    <div v-if="previewVideos.length > 0" class="video-previews">
                        <div
                            v-for="(videoUrl, index) in previewVideos"
                            :key="index"
                            class="video-preview"
                        >
                            <video
                                class="video-preview"
                                :src="videoUrl"
                                controls
                                :width="previewWidth"
                                @error="onVideoError"
                            />
                            <UButton
                                type="button"
                                @click="removeVideoFromVideos(index)"
                                class="remove-btn"
                                :aria-label="
                                    $t('admin_news_form.aria_labels.remove_video', {
                                        index: index + 1,
                                    })
                                "
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
                    {{ $t('admin_news_form.sections.main_info') }}
                </h2>

                <!-- Заголовок -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.title_ru') }}
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.title_ru"
                        @blur="validateField('title_ru')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.title_ru }]"
                        :placeholder="$t('admin_news_form.placeholders.title_ru')"
                    />
                    <div class="error-message" v-if="errors.title_ru">
                        {{ errors.title_ru }}
                    </div>
                </div>

                <!-- Заголовок по английски-->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.title_en') }}
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.title_en"
                        @blur="validateField('title_en')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.title_en }]"
                        :placeholder="$t('admin_news_form.placeholders.title_en')"
                    />
                    <div class="error-message" v-if="errors.title_en">
                        {{ errors.title_en }}
                    </div>
                </div>
            </div>

            <!-- Дополнительная информация -->
            <div class="form-section">
                <h2 class="section-title">
                    {{ $t('admin_news_form.sections.additional_info') }}
                </h2>

                <!-- Дата -->
                <div class="form-group">
                    <label class="form-label">{{
                        $t('admin_news_form.labels.date')
                    }}</label>
                    <UInput
                        v-model="news.datetime"
                        type="date"
                        required
                        min="2000-01-01"
                        :max="new Date().toISOString().split('T')[0]"
                        class="form-control"
                    />
                </div>

                <!-- Текст -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.text_ru') }}
                        <span class="required">*</span>
                    </label>
                    <UTextarea
                        v-model="news.text_ru"
                        @blur="validateField('text_ru')"
                        :class="['form-control', { 'is-invalid': errors.text_ru }]"
                        :placeholder="$t('admin_news_form.placeholders.text_ru')"
                        :rows="6"
                        autoresize
                    />
                    <div class="error-message" v-if="errors.text_ru">
                        {{ errors.text_ru }}
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.text_en') }}
                        <span class="required">*</span>
                    </label>
                    <UTextarea
                        v-model="news.text_en"
                        @blur="validateField('text_en')"
                        :class="['form-control', { 'is-invalid': errors.text_en }]"
                        :placeholder="$t('admin_news_form.placeholders.text_en')"
                        :rows="6"
                        autoresize
                    />
                    <div class="error-message" v-if="errors.text_en">
                        {{ errors.text_en }}
                    </div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <UButton type="button" @click="resetForm" class="btn btn-secondary">
                    {{ $t('admin_news_form.buttons.clear_form') }}
                </UButton>
                <UButton
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">{{
                        $t('admin_news_form.buttons.add_news')
                    }}</span>
                    <span v-else>
                        <span class="spinner"></span>
                        {{ $t('admin_news_form.buttons.submitting') }}
                    </span>
                </UButton>
            </div>
        </form>
    </div>
</template>

<style scoped>
select:has(option.placeholder:checked) {
    color: #999;
}

.add-news-container {
    max-width: 800px;
    margin: 3rem auto;
    padding: 2.5rem;
    background-color: var(--color-on-surface);
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

:root.dark .add-news-container {
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

.news-form {
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
    padding-right: 2.5rem; /* Make space for the arrow */
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

.video-previews {
    width: 100%;
    display: flex;
    gap: 1rem;
}

.video-preview {
    position: relative;
    width: 150px;
    height: 150px;
    border: 1px solid #eee;
    border-radius: 4px;
    overflow: hidden;
    box-sizing: border-box;
}

:root.dark .video-preview {
    border-color: #334155;
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

.textarea {
    resize: vertical;
    min-height: 100px;
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

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}

@media (max-width: 768px) {
    .add-news-container {
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

:root.dark .select-display {
    background-color: #1e293b;
    border-color: #334155;
    color: #e2e8f0;
}

.select-display:hover {
    border-color: #4a90e2;
}

:root.dark .select-display:hover {
    border-color: #3b82f6;
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

:root.dark .dropdown-options {
    background: #1e293b;
    border-color: #334155;
}

.option-item {
    padding: 8px 12px;
    display: flex;
    align-items: center;
    cursor: pointer;
}

:root.dark .option-item {
    color: #e2e8f0;
}

.option-item:hover {
    background-color: #f8f9fa;
}

:root.dark .option-item:hover {
    background-color: #334155;
}

.option-item input {
    margin-right: 8px;
}

.arrow {
    float: right;
    transition: transform 0.3s;
}

:root.dark .arrow {
    color: #94a3b8;
}

.arrow-up {
    transform: rotate(180deg);
}
</style>
