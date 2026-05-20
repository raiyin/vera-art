<script setup lang="ts">
definePageMeta({
    middleware: 'admin-auth',
});
</script>

<script lang="ts">
import { defineComponent } from 'vue';
import type { NewsDesc } from '../../types';
import axios from 'axios';
import { useToast } from '@nuxt/ui/runtime/composables/index.js';

interface PreviewItem {
    file: File;
    preview: string;
}

type NewsForm = Omit<NewsDesc, 'id'>;

export default defineComponent({
    name: 'AddNews',
    components: {},
    data() {
        return {
            news: {
                title_en: '',
                title_ru: '',
                subTitle_en: '',
                subTitle_ru: '',
                img_back: '', // имя файла
                img_backfull: '', // имя файла
                images: [] as string[],
                videos: [] as string[],
                datetime: '',
                text_en: '',
                text_ru: '',
                dir: '',
            } as NewsForm,

            // Для отправки файлов
            images: [] as File[],
            videos: [] as File[],

            img_back_preview: null as PreviewItem | null,
            img_backfull_preview: null as PreviewItem | null,

            previewImages: [] as PreviewItem[],
            previewVideos: [] as string[],

            video_error: '',
            previewWidth: 400,

            isSubmitting: false,
            errors: {
                title_ru: '',
                title_en: '',
                subTitle_ru: '',
                subTitle_en: '',
                img_back: '',
                img_backfull: '',
                images: '',
                text_ru: '',
                text_en: '',
            } as Record<string, string>,
            fileError: null as string | null,
            isDragOver: false,
        };
    },
    computed: {
        server() {
            return import.meta.env.VITE_SERVER_URL;
        },
        isFormValid() {
            return (
                this.news.title_ru.trim() !== '' &&
                this.news.title_en.trim() !== '' &&
                this.news.subTitle_ru.trim() !== '' &&
                this.news.subTitle_en.trim() !== '' &&
                this.img_back_preview !== null &&
                this.img_backfull_preview !== null &&
                this.previewImages.length > 0 &&
                this.news.text_ru.trim() !== '' &&
                this.news.text_en.trim() !== ''
            );
        },
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
            // Handle drop for different file inputs
        },
        triggerFileInput(refName: string) {
            const input = this.$refs[refName] as HTMLInputElement;
            if (input) {
                input.click();
            }
        },
        handleBackImageSelected(event: Event) {
            const target = event.target as HTMLInputElement;
            const selectedImage = target.files?.[0];

            if (!selectedImage) return;

            // Validate file type
            const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
            if (!validTypes.includes(selectedImage.type)) {
                this.fileError = this.$t('admin_news_form.errors.invalid_image_type');
                return;
            }

            // Validate file size (max 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            if (selectedImage.size > maxSize) {
                this.fileError = this.$t('admin_news_form.errors.image_size');
                return;
            }

            // Store file object for form submission, but keep string for type compatibility
            this.news.img_back = selectedImage.name;
            this.fileError = null;

            const reader = new FileReader();
            reader.onload = (e) => {
                this.img_back_preview = {
                    file: selectedImage,
                    preview: e.target?.result as string,
                };
            };
            reader.readAsDataURL(selectedImage);
        },

        removeBackImage() {
            this.img_back_preview = null;
            this.news.img_back = '';
        },

        handleBackFullImageSelected(event: Event) {
            const target = event.target as HTMLInputElement;
            const selectedImage = target.files?.[0];

            if (!selectedImage) return;

            // Validate file type
            const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
            if (!validTypes.includes(selectedImage.type)) {
                this.fileError = this.$t('admin_news_form.errors.invalid_image_type');
                return;
            }

            // Validate file size (max 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            if (selectedImage.size > maxSize) {
                this.fileError = this.$t('admin_news_form.errors.image_size');
                return;
            }

            // Store file name for type compatibility, but keep file object for submission
            this.news.img_backfull = selectedImage.name;
            this.fileError = null;

            const reader = new FileReader();
            reader.onload = (e) => {
                this.img_backfull_preview = {
                    file: selectedImage,
                    preview: e.target?.result as string,
                };
            };
            reader.readAsDataURL(selectedImage);
        },

        removeBackFullImage() {
            this.img_backfull_preview = null;
            this.news.img_backfull = '';
        },

        handleImagesSelected(event: Event) {
            const target = event.target as HTMLInputElement;
            const selectedFiles = Array.from(target.files || []);

            // Проверка на количество файлов
            if (this.images.length + selectedFiles.length > 10) {
                this.fileError = this.$t('admin_news_form.errors.max_images');
                return;
            }

            // Проверка типов файлов
            const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
            const invalidFiles = selectedFiles.filter(
                (file) => !validTypes.includes(file.type)
            );

            if (invalidFiles.length > 0) {
                this.fileError = this.$t('admin_news_form.errors.invalid_image_type');
                return;
            }

            // Проверка размера файлов (макс. 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

            if (largeFiles.length > 0) {
                this.fileError = this.$t('admin_news_form.errors.image_size');
                return;
            }

            this.fileError = null;

            // Добавляем новые файлы
            this.images = [...this.images, ...selectedFiles];

            // Добавляем имена файлов в news.images
            selectedFiles.forEach((file) => {
                this.news.images.push(file.name);
            });

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

        async handleVideosSelected(event: Event) {
            const target = event.target as HTMLInputElement;
            const files = target.files as FileList;
            this.previewVideos.length = 0;
            this.videos = [];

            if (!files || files.length == 0) return;

            // Validate file type
            for (let i = 0; i < files.length; i++) {
                const file = files[i]!;
                if (!file.type.startsWith('video/')) {
                    this.fileError = this.$t('admin_news_form.errors.invalid_video_type');
                    return;
                }
            }

            // Validate file size (max 100MB)
            const maxSize = 100 * 1024 * 1024;
            for (let i = 0; i < files.length; i++) {
                const file = files[i]!;
                if (file.size > maxSize) {
                    this.fileError = this.$t('admin_news_form.errors.video_size');
                    return;
                }
            }

            this.fileError = null;

            // Populate both arrays
            for (let i = 0; i < files.length; i++) {
                const file = files[i]!;
                this.videos.push(file);
                this.news.videos.push(file.name);
                this.previewVideos.push(URL.createObjectURL(file));
            }
        },

        removeImageFromImages(index: number) {
            this.previewImages.splice(index, 1);
            this.images.splice(index, 1);
            this.news.images.splice(index, 1);
        },
        removeVideoFromVideos(index: number) {
            this.previewVideos.splice(index, 1);
            this.videos.splice(index, 1);
            this.news.videos.splice(index, 1);
        },

        onVideoError() {
            // Video error handling
        },

        validateField(fieldName: string) {
            switch (fieldName) {
                case 'title_ru':
                    if (!this.news.title_ru.trim()) {
                        this.errors.title_ru = this.$t(
                            'admin_news_form.errors.title_ru_required'
                        );
                    } else if (
                        this.news.title_ru.trim().length < 3 ||
                        this.news.title_ru.trim().length > 50
                    ) {
                        this.errors.title_ru = this.$t(
                            'admin_news_form.errors.title_ru_length'
                        );
                    } else {
                        this.errors.title_ru = '';
                    }
                    break;
                case 'title_en':
                    if (!this.news.title_en.trim()) {
                        this.errors.title_en = this.$t(
                            'admin_news_form.errors.title_en_required'
                        );
                    } else if (
                        this.news.title_en.trim().length < 3 ||
                        this.news.title_en.trim().length > 50
                    ) {
                        this.errors.title_en = this.$t(
                            'admin_news_form.errors.title_en_length'
                        );
                    } else {
                        this.errors.title_en = '';
                    }
                    break;
                case 'subTitle_ru':
                    if (!this.news.subTitle_ru.trim()) {
                        this.errors.subTitle_ru = this.$t(
                            'admin_news_form.errors.subtitle_ru_required'
                        );
                    } else if (
                        this.news.subTitle_ru.trim().length < 3 ||
                        this.news.subTitle_ru.trim().length > 50
                    ) {
                        this.errors.subTitle_ru = this.$t(
                            'admin_news_form.errors.subtitle_ru_length'
                        );
                    } else {
                        this.errors.subTitle_ru = '';
                    }
                    break;
                case 'subTitle_en':
                    if (!this.news.subTitle_en.trim()) {
                        this.errors.subTitle_en = this.$t(
                            'admin_news_form.errors.subtitle_en_required'
                        );
                    } else if (
                        this.news.subTitle_en.trim().length < 3 ||
                        this.news.subTitle_en.trim().length > 50
                    ) {
                        this.errors.subTitle_en = this.$t(
                            'admin_news_form.errors.subtitle_en_length'
                        );
                    } else {
                        this.errors.subTitle_en = '';
                    }
                    break;
                case 'img_back':
                    if (!this.img_back_preview || !this.news.img_back) {
                        this.errors.img_back = this.$t(
                            'admin_news_form.errors.preview_image_required'
                        );
                    } else {
                        this.errors.img_back = '';
                    }
                    break;
                case 'img_backfull':
                    if (!this.img_backfull_preview || !this.news.img_backfull) {
                        this.errors.img_backfull = this.$t(
                            'admin_news_form.errors.main_image_required'
                        );
                    } else {
                        this.errors.img_backfull = '';
                    }
                    break;
                case 'images':
                    if (this.previewImages.length === 0) {
                        this.errors.images = this.$t(
                            'admin_news_form.errors.images_required'
                        );
                    } else {
                        this.errors.images = '';
                    }
                    break;
                case 'text_ru':
                    if (!this.news.text_ru.trim()) {
                        this.errors.text_ru = this.$t(
                            'admin_news_form.errors.text_ru_required'
                        );
                    } else {
                        this.errors.text_ru = '';
                    }
                    break;
                case 'text_en':
                    if (!this.news.text_en.trim()) {
                        this.errors.text_en = this.$t(
                            'admin_news_form.errors.text_en_required'
                        );
                    } else {
                        this.errors.text_en = '';
                    }
                    break;
            }
        },

        validateForm() {
            this.validateField('title_ru');
            this.validateField('title_en');
            this.validateField('subTitle_ru');
            this.validateField('subTitle_en');
            this.validateField('img_back');
            this.validateField('img_backfull');
            this.validateField('images');
            this.validateField('text_ru');
            this.validateField('text_en');

            // Проверка отсутствия ошибок
            return Object.values(this.errors).every((error) => error === '');
        },

        /**
         * Submit the news form data to the server
         * Sends news information along with images and videos
         */
        async submitForm() {
            if (this.isSubmitting) return;

            if (!this.validateForm()) {
                return;
            }

            try {
                this.isSubmitting = true;

                // Формируем данные для отправки
                const formData = new FormData();

                // Append main images
                if (this.img_back_preview?.file) {
                    formData.append('img_back', this.img_back_preview.file);
                }

                if (this.img_backfull_preview?.file) {
                    formData.append('img_backfull', this.img_backfull_preview.file);
                }
                // Add additional images and videos
                this.images.forEach((image) => {
                    formData.append('images', image);
                });

                this.videos.forEach((video) => {
                    formData.append('videos', video);
                });

                // Add news data as JSON
                const newsData = {
                    ...this.news,
                    datetime:
                        this.news.datetime || new Date().toISOString().split('T')[0],
                };

                formData.append('data', JSON.stringify(newsData));

                // Send data to server
                const response = await axios.post(`${this.server}news`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                });

                if (response.status === 200 || response.status === 201) {
                    const toast = useToast();
                    toast.add({
                        title: this.$t('toast.success.title'),
                        description: this.$t('admin_news_form.messages.submit_success'),
                        icon: 'i-heroicons-check-circle',
                        color: 'success',
                        duration: 5000,
                    });
                    this.resetForm();
                } else {
                    const toast = useToast();
                    toast.add({
                        title: this.$t('toast.error.title'),
                        description: this.$t('admin_news_form.messages.submit_failed'),
                        icon: 'i-heroicons-exclamation-triangle',
                        color: 'error',
                        duration: 5000,
                    });
                }
            } catch (error: any) {
                console.error('Error submitting form:', error);
                const toast = useToast();
                let description = this.$t('admin_news_form.messages.general_error');
                if (error.response?.status === 413) {
                    description = this.$t('admin_news_form.messages.file_too_large');
                } else if (error.response?.status === 400) {
                    description = this.$t('admin_news_form.messages.invalid_data');
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
            this.news = {
                datetime: '',
                title_en: '',
                title_ru: '',
                subTitle_en: '',
                subTitle_ru: '',
                dir: '',
                img_back: '',
                img_backfull: '',
                text_en: '',
                text_ru: '',
                images: [] as string[],
                videos: [] as string[],
            };
            this.images = [];
            this.videos = [];
            this.previewImages = [];
            this.previewVideos = [];
            this.img_back_preview = null;
            this.img_backfull_preview = null;
            // Сброс input файлов
            const fileInputs = this.$el.querySelectorAll('input[type="file"]');
            fileInputs.forEach((input: HTMLInputElement) => {
                input.value = '';
            });

            // Сброс ошибок
            Object.keys(this.errors).forEach((key) => {
                this.errors[key] = '';
            });
            this.fileError = null;
        },
    },
});
</script>

<template>
    <div class="add-news-container">
        <div class="header-section">
            <div class="header-content">
                <h1 class="page-title">{{ $t('admin_news_form.page_title') }}</h1>
                <p class="page-subtitle">
                    {{ $t('admin_news_form.page_subtitle') }}
                </p>
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
                        @drop.prevent="handleDrop"
                        @click="triggerFileInput('backFullInput')"
                    >
                        <UInput
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
                        @drop.prevent="handleDrop"
                        @click="triggerFileInput('backInput')"
                    >
                        <UInput
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
                        @drop.prevent="handleDrop"
                        @click="triggerFileInput('imagesInput')"
                    >
                        <UInput
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
                        @drop.prevent="handleDrop"
                        @click="triggerFileInput('videosInput')"
                    >
                        <UInput
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

                <!-- Подзаголовок -->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.subtitle_ru') }}
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.subTitle_ru"
                        @blur="validateField('subTitle_ru')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.subTitle_ru }]"
                        :placeholder="$t('admin_news_form.placeholders.subtitle_ru')"
                    />
                    <div class="error-message" v-if="errors.subTitle_ru">
                        {{ errors.subTitle_ru }}
                    </div>
                </div>

                <!-- Подзаголовок по английски-->
                <div class="form-group">
                    <label class="form-label">
                        {{ $t('admin_news_form.labels.subtitle_en') }}
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.subTitle_en"
                        @blur="validateField('subTitle_en')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.subTitle_en }]"
                        :placeholder="$t('admin_news_form.placeholders.subtitle_en')"
                    />
                    <div class="error-message" v-if="errors.subTitle_en">
                        {{ errors.subTitle_en }}
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
    color: red;
}

.add-news-container {
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
    color: #666;
    font-size: 1rem;
    margin: 0;
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
