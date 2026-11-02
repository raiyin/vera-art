<template>
    <div class="edit-news-container">
        <div class="header-section">
            <h1 class="page-title">
                Редактировать новость
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
            class="news-form"
            @submit.prevent="submitForm"
        >
            <!-- Изображения новости -->
            <div class="form-section">
                <h2 class="section-title">
                    Изображения новости
                </h2>

                <!-- Главное изображение -->
                <div class="form-group">
                    <label class="form-label">
                        Главное изображение
                        <span class="required">*</span>
                    </label>
                    <div
                        class="file-drop-area"
                        :class="{ 'drag-over': isDragOver, }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'main',)"
                        @click="triggerFileInput('mainImgInput',)"
                    >
                        <input
                            ref="mainImgInput"
                            type="file"
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            @change="handleMainImageSelected"
                        >
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
                                />
                            </svg>
                            <p class="upload-text">
                                Перетащите файл сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаемые форматы: JPG, JPEG, PNG (макс. 5MB)
                            </p>
                        </div>
                    </div>
                    <div
                        v-if="errors.main_image"
                        class="error-message"
                    >
                        {{ errors.main_image }}
                    </div>
                    <div
                        v-if="main_image_preview"
                        class="preview-container"
                    >
                        <div class="image-preview">
                            <img
                                :src="main_image_preview.preview"
                                class="preview-image"
                                alt="Главное изображение"
                            >
                            <UButton
                                type="button"
                                class="remove-btn"
                                aria-label="Удалить главное изображение"
                                @click="removeMainImage"
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
                        :class="{ 'drag-over': isDragOver, }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'images',)"
                        @click="triggerFileInput('imagesInput',)"
                    >
                        <input
                            ref="imagesInput"
                            type="file"
                            multiple
                            accept="image/jpg,image/jpeg,image/png"
                            class="file-input"
                            @change="handleImagesSelected"
                        >
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
                                />
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
                    <div
                        v-if="errors.images"
                        class="error-message"
                    >
                        {{ errors.images }}
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
                                :alt="`Изображение ${index + 1}`"
                            >
                            <UButton
                                type="button"
                                class="remove-btn"
                                :aria-label="`Удалить изображение ${index + 1}`"
                                @click="removeImageFromImages(index,)"
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
                        :class="{ 'drag-over': isDragOver, }"
                        @dragover.prevent="handleDragOver"
                        @dragleave.prevent="handleDragLeave"
                        @drop.prevent="handleDrop($event, 'videos',)"
                        @click="triggerFileInput('videosInput',)"
                    >
                        <input
                            ref="videosInput"
                            type="file"
                            multiple
                            accept="video/*"
                            class="file-input"
                            @change="handleVideosSelected"
                        >
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
                                />
                            </svg>
                            <p class="upload-text">
                                Перетащите видеофайлы сюда или нажмите для выбора
                            </p>
                            <p class="upload-hint">
                                Поддерживаемые форматы: MP4, MOV, AVI (макс. 100MB каждый)
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
                        v-if="previewVideos.length > 0"
                        class="video-previews"
                    >
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
                                class="remove-btn"
                                :aria-label="`Удалить видео ${index + 1}`"
                                @click="removeVideoFromVideos(index,)"
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

                <!-- Заголовок -->
                <div class="form-group">
                    <label class="form-label">
                        Заголовок по по-русски
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.title_ru"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.title_ru, },]"
                        placeholder="Например: 'Звездная ночь'"
                        @blur="validateField('title_ru',)"
                    />
                    <div
                        v-if="errors.title_ru"
                        class="error-message"
                    >
                        {{ errors.title_ru }}
                    </div>
                </div>

                <!-- Заголовок по английски -->
                <div class="form-group">
                    <label class="form-label">
                        Заголовок по-английски
                        <span class="required">*</span>
                    </label>
                    <UInput
                        v-model="news.title_en"
                        type="text"
                        :class="['form-control', { 'is-invalid': errors.title_en, },]"
                        placeholder="For example: 'Starry Night'"
                        @blur="validateField('title_en',)"
                    />
                    <div
                        v-if="errors.title_en"
                        class="error-message"
                    >
                        {{ errors.title_en }}
                    </div>
                </div>
            </div>

            <!-- Дополнительная информация -->
            <div class="form-section">
                <h2 class="section-title">
                    Дополнительная информация
                </h2>

                <!-- Дата -->
                <div class="form-group">
                    <label class="form-label">Дата</label>
                    <UInput
                        v-model="newsDate"
                        type="date"
                        required
                        min="2000-01-01"
                        :max="new Date().toISOString().split('T',)[0]"
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
                        :class="[
                            'form-control textarea',
                            { 'is-invalid': errors.text_ru, },
                        ]"
                        placeholder="Введите текст новости на русском языке"
                        rows="6"
                        @blur="validateField('text_ru',)"
                    />
                    <div
                        v-if="errors.text_ru"
                        class="error-message"
                    >
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
                        :class="[
                            'form-control textarea',
                            { 'is-invalid': errors.text_en, },
                        ]"
                        placeholder="Enter the news text in English"
                        rows="6"
                        @blur="validateField('text_en',)"
                    />
                    <div
                        v-if="errors.text_en"
                        class="error-message"
                    >
                        {{ errors.text_en }}
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
import { ref, reactive, computed, onMounted, } from 'vue';
import type { NewsDesc, } from '~/types';
import { getHttpClient, } from '~/api/http-client';

const api = getHttpClient();

const route = useRoute();
const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;
const toast = useToast();
const { t, } = useI18n();

interface PreviewItem {
    file: File
    preview: string
}

const news = reactive<NewsDesc>({
    id: '',
    title_en: '',
    title_ru: '',
    main_image: '',
    datetime: '',
    text_en: '',
    text_ru: '',
    dir: '',
    images: [],
    videos: [],
});

const images = ref<File[]>([],);
const videos = ref<File[]>([],);

const main_image_preview = ref<PreviewItem | null>(null,);

const previewImages = ref<PreviewItem[]>([],);
const previewVideos = ref<string[]>([],);

const previewWidth = ref(400,);
const isSubmitting = ref(false,);
const fileError = ref<string | null>(null,);
const isDragOver = ref(false,);
const isLoading = ref(true,);

const errors = reactive<Record<string, string>>({
    title_ru: '',
    title_en: '',
    main_image: '',
    images: '',
    text_ru: '',
    text_en: '',
});

const originalNews = ref<NewsDesc>({} as NewsDesc,);

const mainImgInput = ref<HTMLInputElement | null>(null,);
const imagesInput = ref<HTMLInputElement | null>(null,);
const videosInput = ref<HTMLInputElement | null>(null,);

const newsDate = computed({
    get: () => (news.datetime || '').split('T',)[0],
    set: (val: string,) => {
        news.datetime = val;
    },
});

const isFormValid = computed(() => {
    return (
        news.title_ru.trim() !== ''
            && news.title_en.trim() !== ''
        && (main_image_preview.value !== null || news.main_image !== '')
        && news.text_ru.trim() !== ''
            && news.text_en.trim() !== ''
    );
});

onMounted(async () => {
    await loadNews();
});

async function loadNews() {
    try {
        const id = route.params.id;
        const response = await api.get(SERVER_URL + 'news/' + id,);
        Object.assign(news, response.data,);
        originalNews.value = { ...response.data, };

        // Load existing main image with full server URL
        if (news.main_image) {
            main_image_preview.value = {
                file: new File([], news.main_image,),
                preview: `${news.dir}${news.main_image}`,
            };
        }

        // Load existing gallery images
        if (news.images && news.images.length > 0) {
            news.images.forEach((imgName,) => {
                previewImages.value.push({
                    file: new File([], imgName,),
                    preview: `${news.dir}${imgName}`,
                });
            });
        }

        // Load existing videos
        if (news.videos && news.videos.length > 0) {
            news.videos.forEach((videoName,) => {
                previewVideos.value.push(`${news.dir}videos/${videoName}/${videoName}.mp4`,);
            });
        }

        isLoading.value = false;
    } catch (error) {
        console.error('Ошибка при загрузке новости:', error,);
        isLoading.value = false;
        toast.add({
            title: t('toast.error.title',),
            description: 'Не удалось загрузить данные. Пожалуйста, попробуйте позже.',
            icon: 'i-heroicons-exclamation-triangle',
            color: 'error',
            duration: 5000,
        });
    }
}

function handleDragOver() {
    isDragOver.value = true;
}

function handleDragLeave() {
    isDragOver.value = false;
}

function handleDrop(event: DragEvent, target: 'main' | 'images' | 'videos',) {
    isDragOver.value = false;
    const files = Array.from(event.dataTransfer?.files || [],);
    if (files.length === 0) return;

    switch (target) {
        case 'main':
            handleMainImageDrop(files[0]!,);
            break;
        case 'images':
            handleImagesDrop(files,);
            break;
        case 'videos':
            handleVideosDrop(files,);
            break;
    }
}

function handleMainImageDrop(file: File,) {
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png',];
    if (!validTypes.includes(file.type,)) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    if (file.size > maxSize) {
        fileError.value = 'Размер файла не должен превышать 5 МБ';
        return;
    }

    news.main_image = file.name;
    fileError.value = null;

    const reader = new FileReader();
    reader.onload = (e,) => {
        main_image_preview.value = {
            file,
            preview: e.target?.result as string,
        };
    };
    reader.readAsDataURL(file,);
}

function handleImagesDrop(files: File[],) {
    if (images.value.length + files.length > 10) {
        fileError.value = 'Можно загрузить не более 10 изображений';
        return;
    }

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png',];
    const invalidFiles = files.filter(file => !validTypes.includes(file.type,),);
    if (invalidFiles.length > 0) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    const largeFiles = files.filter(file => file.size > maxSize,);
    if (largeFiles.length > 0) {
        fileError.value = 'Размер каждого файла не должен превышать 5 МБ';
        return;
    }

    fileError.value = null;

    images.value = [...images.value, ...files,];

    files.forEach((file,) => {
        news.images.push(file.name,);
    });

    files.forEach((file,) => {
        const reader = new FileReader();
        reader.onload = (e,) => {
            previewImages.value.push({
                file,
                preview: e.target?.result as string,
            });
        };
        reader.readAsDataURL(file,);
    });
}

function handleVideosDrop(files: File[],) {
    for (const file of files) {
        if (!file.type.startsWith('video/',)) {
            fileError.value = 'Пожалуйста, выберите только видеофайлы';
            return;
        }
    }

    const maxSize = 100 * 1024 * 1024;
    for (const file of files) {
        if (file.size > maxSize) {
            fileError.value = 'Размер файла не должен превышать 100MB';
            return;
        }
    }

    fileError.value = null;

    for (const file of files) {
        videos.value.push(file,);
        news.videos.push(storedVideoName(file.name,),);
        previewVideos.value.push(URL.createObjectURL(file,),);
    }
}

function triggerFileInput(refName: string,) {
    let input: HTMLInputElement | null = null;
    switch (refName) {
        case 'mainImgInput':
            input = mainImgInput.value;
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

function handleMainImageSelected(event: Event,) {
    const target = event.target as HTMLInputElement;
    const selectedImage = target.files?.[0];

    if (!selectedImage) return;

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png',];
    if (!validTypes.includes(selectedImage.type,)) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    if (selectedImage.size > maxSize) {
        fileError.value = 'Размер файла не должен превышать 5 МБ';
        return;
    }

    news.main_image = selectedImage.name;
    fileError.value = null;

    const reader = new FileReader();
    reader.onload = (e,) => {
        main_image_preview.value = {
            file: selectedImage,
            preview: e.target?.result as string,
        };
    };
    reader.readAsDataURL(selectedImage,);
}

function removeMainImage() {
    main_image_preview.value = null;
    news.main_image = '';
}

function handleImagesSelected(event: Event,) {
    const target = event.target as HTMLInputElement;
    const selectedFiles = Array.from(target.files || [],);

    if (images.value.length + selectedFiles.length > 10) {
        fileError.value = 'Можно загрузить не более 10 изображений';
        return;
    }

    const validTypes = ['image/jpeg', 'image/jpg', 'image/png',];
    const invalidFiles = selectedFiles.filter(file => !validTypes.includes(file.type,),);

    if (invalidFiles.length > 0) {
        fileError.value = 'Пожалуйста, загружайте только изображения (JPG, JPEG, PNG)';
        return;
    }

    const maxSize = 5 * 1024 * 1024;
    const largeFiles = selectedFiles.filter(file => file.size > maxSize,);

    if (largeFiles.length > 0) {
        fileError.value = 'Размер каждого файла не должен превышать 5 МБ';
        return;
    }

    fileError.value = null;

    images.value = [...images.value, ...selectedFiles,];

    selectedFiles.forEach((file,) => {
        news.images.push(file.name,);
    });

    selectedFiles.forEach((file,) => {
        const reader = new FileReader();
        reader.onload = (e,) => {
            previewImages.value.push({
                file,
                preview: e.target?.result as string,
            });
        };
        reader.readAsDataURL(file,);
    });
}

async function handleVideosSelected(event: Event,) {
    const target = event.target as HTMLInputElement;
    const files = target.files as FileList;

    if (!files || files.length == 0) return;

    for (let i = 0; i < files.length; i++) {
        const file = files[i]!;
        if (!file.type.startsWith('video/',)) {
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
        videos.value.push(file,);
        news.videos.push(storedVideoName(file.name,),);
        previewVideos.value.push(URL.createObjectURL(file,),);
    }
}

function removeImageFromImages(index: number,) {
    // images.value only contains newly uploaded files, while previewImages
    // contains both existing (server) and new files. Remove the matching file.
    const existingCount = news.images.length - images.value.length;
    if (index >= existingCount) {
        images.value.splice(index - existingCount, 1,);
    }
    previewImages.value.splice(index, 1,);
    news.images.splice(index, 1,);
}

function removeVideoFromVideos(index: number,) {
    // videos.value only contains newly uploaded files, while previewVideos
    // contains both existing (server) and new files. Remove the matching file.
    const existingCount = news.videos.length - videos.value.length;
    if (index >= existingCount) {
        videos.value.splice(index - existingCount, 1,);
    }
    previewVideos.value.splice(index, 1,);
    news.videos.splice(index, 1,);
}

function onVideoError() {
    // Video error handling
}

function storedVideoName(filename: string,): string {
    // The server stores videos in a folder named after the file without extension.
    return filename.replace(/\.[^/.]+$/, '',);
}

function validateField(fieldName: string,) {
    switch (fieldName) {
        case 'title_ru':
            if (!news.title_ru.trim()) {
                errors.title_ru = 'Пожалуйста, введите заголовок на русском';
            } else if (
                news.title_ru.trim().length < 3
            || news.title_ru.trim().length > 50
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
                news.title_en.trim().length < 3
            || news.title_en.trim().length > 50
            ) {
                errors.title_en
                        = 'Заголовок на английском должен быть от 3 до 50 символов';
            } else {
                errors.title_en = '';
            }
            break;
        case 'main_image':
            if (!main_image_preview.value && !news.main_image) {
                errors.main_image = 'Пожалуйста, добавьте главное изображение новости';
            } else {
                errors.main_image = '';
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
    validateField('title_ru',);
    validateField('title_en',);
    validateField('main_image',);
    validateField('images',);
    validateField('text_ru',);
    validateField('text_en',);

    return Object.values(errors,).every(error => error === '',);
}

async function submitForm() {
    if (isSubmitting.value) return;

    if (!validateForm()) {
        return;
    }

    try {
        isSubmitting.value = true;

        const formData = new FormData();

        // Only append the file that was actually selected by the user (has a real File object with size > 0)
        if (main_image_preview.value?.file && main_image_preview.value.file.size > 0) {
            formData.append('main_image', main_image_preview.value.file,);
        }

        images.value.forEach((image,) => {
            formData.append('images', image,);
        });

        videos.value.forEach((video,) => {
            formData.append('videos', video,);
        });

        const newsData = {
            ...news,
            datetime: news.datetime || new Date().toISOString().split('T',)[0],
        };

        formData.append('data', JSON.stringify(newsData,),);

        const id = route.params.id;
        const response = await api.put(`${SERVER_URL}news/${id}`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });

        if (response.status === 200) {
            toast.add({
                title: t('toast.success.title',),
                description: 'Новость успешно обновлена.',
                icon: 'i-heroicons-check-circle',
                color: 'success',
                duration: 5000,
            });
            originalNews.value = { ...news, };
        } else {
            toast.add({
                title: t('toast.error.title',),
                description: 'Не удалось обновить новость. Пожалуйста, попробуйте снова.',
                icon: 'i-heroicons-exclamation-triangle',
                color: 'error',
                duration: 5000,
            });
        }
    } catch (error: any) {
        console.error('Error submitting form:', error,);
        let description
                = 'Произошла ошибка при обновлении новости. Пожалуйста, попробуйте снова.';
        if (error.response?.status === 413) {
            description
                    = 'Файлы слишком большие. Пожалуйста, загрузите меньшие изображения.';
        } else if (error.response?.status === 400) {
            description
                    = 'Некорректные данные. Пожалуйста, проверьте введенные значения.';
        }
        toast.add({
            title: t('toast.error.title',),
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
    Object.assign(news, originalNews.value,);
    images.value = [];
    videos.value = [];
    previewImages.value = [];
    previewVideos.value = [];
    main_image_preview.value = null;

    if (mainImgInput.value) mainImgInput.value.value = '';
    if (imagesInput.value) imagesInput.value.value = '';
    if (videosInput.value) videosInput.value.value = '';

    Object.keys(errors,).forEach((key,) => {
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

:deep(.remove-btn) {
    position: absolute;
    top: 0;
    right: 0;
    background-color: rgba(255, 0, 0, 0.7);
    color: white;
    border: none;
    width: 24px;
    height: 24px;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-size: 1rem;
    line-height: 1;
    border-radius: 0 0 0 4px;
}

:deep(.remove-btn:hover) {
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
