<script lang="ts">
import axios from 'axios';
import { defineComponent } from 'vue';
import { NewsDescDto } from '@/types';
import Alert from '@/components/app-ui/Alert.vue';

interface PreviewItem {
    file: File;
    preview: string;
}

export default defineComponent({
    name: 'AddNews',
    components: {
        Alert,
    },
    data() {
        return {
            news: {
                id: 0,
                title_en: '',
                title_ru: '',
                subTitle_en: '',
                subTitle_ru: '',
                img_back: '', // имя файла
                img_backfull: '', // имя файла
                images: [],
                videos: [],
                datetime: '',
                text_en: '',
                text_ru: '',
                dir: '',
            } as NewsDescDto,

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
            showSuccessAlert: false,
            showErrorAlert: false,
            errorMessage: '',
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
                this.fileError =
                    'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
                return;
            }

            // Validate file size (max 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            if (selectedImage.size > maxSize) {
                this.fileError = 'Размер файла не должен превышать 5 МБ';
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
                this.fileError =
                    'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
                return;
            }

            // Validate file size (max 5MB)
            const maxSize = 5 * 1024 * 1024; // 5MB
            if (selectedImage.size > maxSize) {
                this.fileError = 'Размер файла не должен превышать 5 МБ';
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

            this.fileError = null;

            // Добавляем новые файлы
            this.images = [...this.images, ...selectedFiles];

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
                const file = files[i];
                if (!file.type.startsWith('video/')) {
                    this.fileError = 'Пожалуйста, выберите только видеофайлы';
                    return;
                }
            }

            // Validate file size (max 100MB)
            const maxSize = 100 * 1024 * 1024;
            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                if (file.size > maxSize) {
                    this.fileError = 'Размер файла не должен превышать 100MB';
                    return;
                }
            }

            this.fileError = null;

            // Populate both arrays
            for (let i = 0; i < files.length; i++) {
                this.videos.push(files[i]);
                this.previewVideos.push(URL.createObjectURL(files[i]));
            }
        },

        removeImageFromImages(index: number) {
            this.previewImages.splice(index, 1);
            this.images.splice(index, 1);
        },
        removeVideoFromVideos(index: number) {
            this.previewVideos.splice(index, 1);
            this.videos.splice(index, 1);
        },

        onVideoError() {
            // Video error handling
        },

        validateField(fieldName: string) {
            switch (fieldName) {
                case 'title_ru':
                    if (!this.news.title_ru.trim()) {
                        this.errors.title_ru = 'Пожалуйста, введите заголовок на русском';
                    } else if (
                        this.news.title_ru.trim().length < 3 ||
                        this.news.title_ru.trim().length > 50
                    ) {
                        this.errors.title_ru =
                            'Заголовок на русском должен быть от 3 до 50 символов';
                    } else {
                        this.errors.title_ru = '';
                    }
                    break;
                case 'title_en':
                    if (!this.news.title_en.trim()) {
                        this.errors.title_en =
                            'Пожалуйста, введите заголовок на английском';
                    } else if (
                        this.news.title_en.trim().length < 3 ||
                        this.news.title_en.trim().length > 50
                    ) {
                        this.errors.title_en =
                            'Заголовок на английском должен быть от 3 до 50 символов';
                    } else {
                        this.errors.title_en = '';
                    }
                    break;
                case 'subTitle_ru':
                    if (!this.news.subTitle_ru.trim()) {
                        this.errors.subTitle_ru =
                            'Пожалуйста, введите подзаголовок на русском';
                    } else if (
                        this.news.subTitle_ru.trim().length < 3 ||
                        this.news.subTitle_ru.trim().length > 50
                    ) {
                        this.errors.subTitle_ru =
                            'Подзаголовок на русском должен быть от 3 до 50 символов';
                    } else {
                        this.errors.subTitle_ru = '';
                    }
                    break;
                case 'subTitle_en':
                    if (!this.news.subTitle_en.trim()) {
                        this.errors.subTitle_en =
                            'Пожалуйста, введите подзаголовок на английском';
                    } else if (
                        this.news.subTitle_en.trim().length < 3 ||
                        this.news.subTitle_en.trim().length > 50
                    ) {
                        this.errors.subTitle_en =
                            'Подзаголовок на английском должен быть от 3 до 50 символов';
                    } else {
                        this.errors.subTitle_en = '';
                    }
                    break;
                case 'img_back':
                    if (!this.img_back_preview || !this.news.img_back) {
                        this.errors.img_back =
                            'Пожалуйста, добавьте предварительное изображение новости';
                    } else {
                        this.errors.img_back = '';
                    }
                    break;
                case 'img_backfull':
                    if (!this.img_backfull_preview || !this.news.img_backfull) {
                        this.errors.img_backfull =
                            'Пожалуйста, добавьте главное изображение новости';
                    } else {
                        this.errors.img_backfull = '';
                    }
                    break;
                case 'images':
                    if (this.previewImages.length === 0) {
                        this.errors.images =
                            'Пожалуйста, добавьте хотя бы одно изображение';
                    } else {
                        this.errors.images = '';
                    }
                    break;
                case 'text_ru':
                    if (!this.news.text_ru.trim()) {
                        this.errors.text_ru =
                            'Пожалуйста, добавьте текст на русском языке';
                    } else {
                        this.errors.text_ru = '';
                    }
                    break;
                case 'text_en':
                    if (!this.news.text_en.trim()) {
                        this.errors.text_en =
                            'Пожалуйста, добавьте текст на английском языке';
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
                this.errorMessage = '';

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
                    this.showSuccessAlert = true;
                    this.resetForm();
                } else {
                    this.showErrorAlert = true;
                    this.errorMessage =
                        'Не удалось добавить новость. Пожалуйста, попробуйте снова.';
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
                        'Произошла ошибка при добавлении новости. Пожалуйста, попробуйте снова.';
                }
            } finally {
                this.isSubmitting = false;
            }
        },

        resetForm() {
            this.news = {
                id: '',
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
                images: [],
                videos: [],
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

        closeAlert() {
            this.showSuccessAlert = false;
            this.showErrorAlert = false;
            this.errorMessage = '';
        },
    },
});
</script>

<template>
    <div class="add-news-container">
        <div class="header-section">
            <h1 class="page-title">Добавить новость</h1>
            <p class="page-subtitle">
                Заполните все обязательные поля для добавления новости
            </p>
        </div>

        <!-- Success Alert -->
        <UAlert
            v-model="showSuccessAlert"
            type="success"
            title="Успешно!"
            message="Новость успешно добавлена."
            closeButtonText="Закрыть"
        />

        <!-- Danger Alert -->
        <UAlert
            v-model="showErrorAlert"
            type="danger"
            title="Ошибка!"
            :message="
                errorMessage ||
                'Не удалось добавить новость. Пожалуйста, попробуйте снова.'
            "
            closeButtonText="Закрыть"
        />

        <form @submit.prevent="submitForm" class="news-form">
            <!-- Изображения новости -->
            <div class="form-section">
                <h2 class="section-title">Изображения новости</h2>

                <!-- Главное изображение -->
                <div class="form-group">
                    <label class="form-label">
                        Главное изображение
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
                                Перетащите файл сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаемые форматы: JPG, JPEG, PNG (макс. 5MB)
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
                                alt="Главное изображение"
                            />
                            <UButton
                                type="button"
                                @click="removeBackFullImage"
                                class="remove-btn"
                                aria-label="Удалить главное изображение"
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>

                <!-- Превью изображение новости -->
                <div class="form-group">
                    <label class="form-label">
                        Превью изображение новости
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
                                Перетащите файл сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаемые форматы: JPG, JPEG, PNG (макс. 5MB)
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
                                alt="Превью изображение"
                            />
                            <UButton
                                type="button"
                                @click="removeBackImage"
                                class="remove-btn"
                                aria-label="Удалить превью изображение"
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>

                <!-- Фотогаллерея -->
                <div class="form-group">
                    <label class="form-label">
                        Фотогаллерея
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
                                Перетащите файлы сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаемые форматы: JPG, JPEG, PNG (макс. 5MB каждый,
                                макс. 10 файлов)
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
                                :alt="`Изображение ${index + 1}`"
                            />
                            <UButton
                                type="button"
                                @click="removeImageFromImages(index)"
                                class="remove-btn"
                                :aria-label="`Удалить изображение ${index + 1}`"
                            >
                                &times;
                            </UButton>
                        </div>
                    </div>
                </div>

                <!-- Видеогаллерея -->
                <div class="form-group">
                    <label class="form-label">Видеогаллерея</label>
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
                                Перетащите видеофайлы сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаемые форматы: MP4, MOV, AVI (макс. 100MB каждый)
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
                                :aria-label="`Удалить видео ${index + 1}`"
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

                <!-- Заголовок -->
                <div class="form-group">
                    <label class="form-label">
                        Заголовок по по-русски
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.title_ru"
                        @blur="validateField('title_ru')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.title_ru }]"
                        placeholder="Например: 'Звездная ночь'"
                    />
                    <div class="error-message" v-if="errors.title_ru">
                        {{ errors.title_ru }}
                    </div>
                </div>

                <!-- Заголовок по английски-->
                <div class="form-group">
                    <label class="form-label">
                        Заголовок по-английски
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.title_en"
                        @blur="validateField('title_en')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.title_en }]"
                        placeholder="For example: 'Starry Night'"
                    />
                    <div class="error-message" v-if="errors.title_en">
                        {{ errors.title_en }}
                    </div>
                </div>

                <!-- Подзаголовок -->
                <div class="form-group">
                    <label class="form-label">
                        Подзаголовок по-русски
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.subTitle_ru"
                        @blur="validateField('subTitle_ru')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.subTitle_ru }]"
                        placeholder="Например: 'Звездная ночь'"
                    />
                    <div class="error-message" v-if="errors.subTitle_ru">
                        {{ errors.subTitle_ru }}
                    </div>
                </div>

                <!-- Подзаголовок по английски-->
                <div class="form-group">
                    <label class="form-label">
                        Подзаголовок по-английски
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.subTitle_en"
                        @blur="validateField('subTitle_en')"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.subTitle_en }]"
                        placeholder="For example: 'Starry Night'"
                    />
                    <div class="error-message" v-if="errors.subTitle_en">
                        {{ errors.subTitle_en }}
                    </div>
                </div>
            </div>

            <!-- Дополнительная информация -->
            <div class="form-section">
                <h2 class="section-title">Дополнительная информация</h2>

                <!-- Дата -->
                <div class="form-group">
                    <label class="form-label">Дата</label>
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
                        Текст по-русски
                        <span class="required">*</span>
                    </label>
                    <textarea
                        v-model="news.text_ru"
                        @blur="validateField('text_ru')"
                        :class="[
                            'form-control textarea',
                            { 'is-invalid': errors.text_ru },
                        ]"
                        placeholder="Введите текст новости на русском языке"
                        rows="6"
                    ></textarea>
                    <div class="error-message" v-if="errors.text_ru">
                        {{ errors.text_ru }}
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">
                        Текст по-английски
                        <span class="required">*</span>
                    </label>
                    <textarea
                        v-model="news.text_en"
                        @blur="validateField('text_en')"
                        :class="[
                            'form-control textarea',
                            { 'is-invalid': errors.text_en },
                        ]"
                        placeholder="Enter the news text in English"
                        rows="6"
                    ></textarea>
                    <div class="error-message" v-if="errors.text_en">
                        {{ errors.text_en }}
                    </div>
                </div>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <UButton type="button" @click="resetForm" class="btn btn-secondary">
                    Очистить форму
                </UButton>
                <UButton
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">Добавить новость</span>
                    <span v-else>
                        <span class="spinner"></span>
                        Отправка...
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
