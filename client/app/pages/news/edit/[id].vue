<template>
    <div class="edit-news-container">
        <div class="header-section">
            <h1 class="page-title">Редактировать новость</h1>
            <p class="page-subtitle">Измените необходимые поля и сохраните изменения</p>
        </div>

        <!-- Загрузчик -->
        <div v-if="isLoading" class="loading-container">
            <div class="loader"></div>
            <p>Загрузка данных...</p>
        </div>

        <form v-else @submit.prevent="submitForm" class="news-form">
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
            </div>

            <!-- Дополнительная информация -->
            <div class="form-section">
                <h2 class="section-title">Дополнительная информация</h2>

                <!-- Дата -->
                <div class="form-group">
                    <label class="form-label">Дата</label>
                    <UInput
                        v-model="newsDate"
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
                    Сбросить изменения
                </UButton>
                <UButton
                    type="submit"
                    class="btn btn-primary"
                    :disabled="isSubmitting || !isFormValid"
                >
                    <span v-if="!isSubmitting">Сохранить изменения</span>
                    <span v-else>
                        <span class="spinner"></span>
                        Сохранение...
                    </span>
                </UButton>
            </div>
        </form>

        <!-- Success Alert -->
        <UAlert
            v-if="showSuccessAlert"
            title="Успешно!"
            description="Новость успешно обновлена."
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
                errorMessage ||
                'Не удалось обновить новость. Пожалуйста, попробуйте снова.'
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
import type { NewsDescDto } from '../../../types';

const route = useRoute();
const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;

interface PreviewItem {
    file: File;
    preview: string;
}

const news = reactive<NewsDescDto>({
    id: '',
    title_en: '',
    title_ru: '',
    img_back: '',
    img_backfull: '',
    datetime: '',
    text_en: '',
    text_ru: '',
    dir: '',
    images: [],
    videos: [],
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
const isLoading = ref(true);
const showSuccessAlert = ref(false);
const showErrorAlert = ref(false);
const errorMessage = ref('');

const errors = reactive<Record<string, string>>({
    title_ru: '',
    title_en: '',
    img_back: '',
    img_backfull: '',
    images: '',
    text_ru: '',
    text_en: '',
});

const originalNews = ref<NewsDescDto>({} as NewsDescDto);

const backFullInput = ref<HTMLInputElement | null>(null);
const backInput = ref<HTMLInputElement | null>(null);
const imagesInput = ref<HTMLInputElement | null>(null);
const videosInput = ref<HTMLInputElement | null>(null);

const newsDate = computed({
    get: () => news.datetime.split('T')[0],
    set: (val) => {
        news.datetime = val;
    },
});

const isFormValid = computed(() => {
    return (
        news.title_ru.trim() !== '' &&
        news.title_en.trim() !== '' &&
        (img_back_preview.value !== null || news.img_back !== '') &&
        (img_backfull_preview.value !== null || news.img_backfull !== '') &&
        news.text_ru.trim() !== '' &&
        news.text_en.trim() !== ''
    );
});

onMounted(async () => {
    await loadNews();
});

async function loadNews() {
    try {
        const id = route.params.id;
        const response = await axios.get(SERVER_URL + 'news/' + id);
        Object.assign(news, response.data);
        originalNews.value = { ...response.data };

        if (news.img_back) {
            img_back_preview.value = {
                file: new File([], news.img_back),
                preview: `${news.dir}${news.img_back}`,
            };
        }

        if (news.img_backfull) {
            img_backfull_preview.value = {
                file: new File([], news.img_backfull),
                preview: `${news.dir}${news.img_backfull}`,
            };
        }

        isLoading.value = false;
    } catch (error) {
        console.error('Ошибка при загрузке новости:', error);
        isLoading.value = false;
        showErrorAlert.value = true;
        errorMessage.value = 'Не удалось загрузить данные. Пожалуйста, попробуйте позже.';
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

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    if (!validTypes.includes(selectedImage.type)) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    if (selectedImage.size > maxSize) {
        fileError.value = 'Размер файла не должен превышать 5 МБ';
        return;
    }

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

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    if (!validTypes.includes(selectedImage.type)) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    if (selectedImage.size > maxSize) {
        fileError.value = 'Размер файла не должен превышать 5 МБ';
        return;
    }

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

    if (images.value.length + selectedFiles.length > 10) {
        fileError.value = 'Можно загрузить не более 10 изображений';
        return;
    }

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png'];
    const invalidFiles = selectedFiles.filter((file) => !validTypes.includes(file.type));

    if (invalidFiles.length > 0) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    const largeFiles = selectedFiles.filter((file) => file.size > maxSize);

    if (largeFiles.length > 0) {
        fileError.value = 'Размер каждого файла не должен превышать 5 МБ';
        return;
    }

    fileError.value = null;

    images.value = [...images.value, ...selectedFiles];

    selectedFiles.forEach((file) => {
        news.images.push(file.name);
    });

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
    previewVideos.value.length = 0;
    videos.value = [];

    if (!files || files.length == 0) return;

    for (let i = 0; i < files.length; i++) {
        const file = files[i]!;
        if (!file.type.startsWith('video/')) {
            fileError.value = 'Пожалуйста, выберите только видеофайлы';
            return;
        }
    }

    const maxSize = 100 * 1024 * 1024;
    for (let i = 0; i < files.length; i++) {
        const file = files[i]!;
        if (file.size > maxSize) {
            fileError.value = 'Размер файла не должен превышать 100MB';
            return;
        }
    }

    fileError.value = null;

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
                errors.title_ru = 'Пожалуйста, введите заголовок на русском';
            } else if (
                news.title_ru.trim().length < 3 ||
                news.title_ru.trim().length > 50
            ) {
                errors.title_ru = 'Заголовок на русском должен быть от 3 до 50 символов';
            } else {
                errors.title_ru = '';
            }
            break;
        case 'title_en':
            if (!news.title_en.trim()) {
                errors.title_en = 'Пожалуйста, введите заголовок на английском';
            } else if (
                news.title_en.trim().length < 3 ||
                news.title_en.trim().length > 50
            ) {
                errors.title_en =
                    'Заголовок на английском должен быть от 3 до 50 символов';
            } else {
                errors.title_en = '';
            }
            break;
        case 'img_back':
            if (!img_back_preview.value && !news.img_back) {
                errors.img_back =
                    'Пожалуйста, добавьте предварительное изображение новости';
            } else {
                errors.img_back = '';
            }
            break;
        case 'img_backfull':
            if (!img_backfull_preview.value && !news.img_backfull) {
                errors.img_backfull = 'Пожалуйста, добавьте главное изображение новости';
            } else {
                errors.img_backfull = '';
            }
            break;
        case 'images':
            break;
        case 'text_ru':
            if (!news.text_ru.trim()) {
                errors.text_ru = 'Пожалуйста, добавьте текст на русском языке';
            } else {
                errors.text_ru = '';
            }
            break;
        case 'text_en':
            if (!news.text_en.trim()) {
                errors.text_en = 'Пожалуйста, добавьте текст на английском языке';
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

        const formData = new FormData();

        if (img_back_preview.value?.file) {
            formData.append('img_back', img_back_preview.value.file);
        }

        if (img_backfull_preview.value?.file) {
            formData.append('img_backfull', img_backfull_preview.value.file);
        }

        images.value.forEach((image) => {
            formData.append('images', image);
        });

        videos.value.forEach((video) => {
            formData.append('videos', video);
        });

        const newsData = {
            ...news,
            datetime: news.datetime || new Date().toISOString().split('T')[0],
        };

        formData.append('data', JSON.stringify(newsData));

        const id = route.params.id;
        const response = await axios.put(`${SERVER_URL}news/${id}`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (response.status === 200) {
            showSuccessAlert.value = true;
            originalNews.value = { ...news };
        } else {
            showErrorAlert.value = true;
            errorMessage.value =
                'Не удалось обновить новость. Пожалуйста, попробуйте снова.';
        }
    } catch (error: any) {
        console.error('Error submitting form:', error);
        showErrorAlert.value = true;
        if (error.response?.status === 413) {
            errorMessage.value =
                'Файлы слишком большие. Пожалуйста, загрузите меньшие изображения.';
        } else if (error.response?.status === 400) {
            errorMessage.value =
                'Некорректные данные. Пожалуйста, проверьте введенные значения.';
        } else {
            errorMessage.value =
                'Произошла ошибка при обновлении новости. Пожалуйста, попробуйте снова.';
        }
    } finally {
        isSubmitting.value = false;
    }
}

function resetForm() {
    Object.assign(news, originalNews.value);
    images.value = [];
    videos.value = [];
    previewImages.value = [];
    previewVideos.value = [];
    img_back_preview.value = null;
    img_backfull_preview.value = null;

    if (backFullInput.value) backFullInput.value.value = '';
    if (backInput.value) backInput.value.value = '';
    if (imagesInput.value) imagesInput.value.value = '';
    if (videosInput.value) videosInput.value.value = '';

    Object.keys(errors).forEach((key) => {
        errors[key] = '';
    });
    fileError.value = null;
}
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: #999;
}

.edit-news-container {
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
    .edit-news-container {
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

/* ===== Dark Mode Overrides ===== */
:root.dark .edit-news-container {
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

:root.dark .video-preview {
    border-color: #334155;
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

:root.dark .select-display {
    background-color: #1e293b;
    border-color: #334155;
    color: #e2e8f0;
}

:root.dark .select-display:hover {
    border-color: #3b82f6;
}

:root.dark .dropdown-options {
    background: #1e293b;
    border-color: #334155;
}

:root.dark .option-item {
    color: #e2e8f0;
}

:root.dark .option-item:hover {
    background-color: #334155;
}

:root.dark .arrow {
    color: #94a3b8;
}
</style>
