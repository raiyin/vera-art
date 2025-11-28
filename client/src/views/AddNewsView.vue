<script lang="ts">
import axios from 'axios';
import { defineComponent, Ref, ref } from 'vue';
import { NewsDescDto, RequestResult } from '@/types';

interface PreviewItem {
    file: File;
    preview: string;
}

export default defineComponent({
    name: 'AddNews',
    data() {
        return {
            news: {
                title_en: '',
                title_ru: '',
                subTitle_en: '',
                subTitle_ru: '',
                img_back: '', // имя файла
                img_backfull: '', // имя файла
                imagescount: 0,
                videoscount: 0,
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
            validated: true,
            validationErrorMessage: '' as string,
            requestResult: 'unknown' as RequestResult,
        };
    },
    computed: {
        server() {
            return import.meta.env.VITE_SERVER_URL;
        },
    },
    methods: {
        handleBackImageSelected(event) {
            const target = event.target as HTMLInputElement;
            const selectedImage = target.files?.[0];

            if (!selectedImage) return;

            // Store file object for form submission, but keep string for type compatibility
            this.news.img_back = selectedImage.name;

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

        handleBackFullImageSelected(event) {
            const target = event.target as HTMLInputElement;
            const selectedImage = target.files?.[0];

            if (!selectedImage) return;

            // Store file name for type compatibility, but keep file object for submission
            this.news.img_backfull = selectedImage.name;

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

        handleImagesSelected(event) {
            const target = event.target as HTMLInputElement;
            const selectedFiles = Array.from(target.files || []);

            if (this.images.length + selectedFiles.length > 10) {
                alert('Можно загрузить не более 10 изображений');
                return;
            }

            this.images = [...this.images, ...selectedFiles];

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

        async handleVideosSelected(event) {
            const files = event.target.files as FileList;
            this.previewVideos.length = 0;
            this.videos = [];

            if (!files || files.length == 0) return;

            // Validate file type
            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                if (!file.type.startsWith('video/')) {
                    this.video_error = 'Пожалуйста, выберите только видеофайлы';
                    return;
                }
            }

            // Validate file size (max 100MB)
            const maxSize = 100 * 1024 * 1024;
            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                if (file.size > maxSize) {
                    this.video_error = 'Размер файла не должен превышать 100MB';
                    return;
                }
            }

            // Populate both arrays
            for (let i = 0; i < files.length; i++) {
                this.videos.push(files[i]);
                this.previewVideos.push(URL.createObjectURL(files[i]));
            }

            this.video_error = '';
        },

        removeImageFromImages(index) {
            this.previewImages.splice(index, 1);
            this.images.splice(index, 1);
        },
        removeVideoFromVideos(index) {
            this.previewVideos.splice(index, 1);
            this.videos.splice(index, 1);
        },

        onVideoError() {
            // Video error handling
        },

        formatFileSize(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        },

        formatDuration(seconds) {
            const minutes = Math.floor(seconds / 60);
            const remainingSeconds = Math.floor(seconds % 60);
            return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
        },

        /**
         * Submit the news form data to the server
         * Sends news information along with images and videos
         */
        async submitForm() {
            if (!this.validateForm()) {
                return;
            }

            if (this.isSubmitting) return;

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
                    imagescount: this.images.length,
                    videoscount: this.videos.length,
                    datetime:
                        this.news.datetime || new Date().toISOString().split('T')[0],
                };

                formData.append('data', JSON.stringify(newsData));

                // Send data to server
                const response = await axios.post(`${this.server}news`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data', // Important for file uploads
                        // Add authorization header if needed
                        // Authorization: `Bearer ${yourAuthToken}`,
                    },
                    // Optional: track upload progress
                    onUploadProgress: (progressEvent) => {
                        const percentCompleted = Math.round(
                            (progressEvent.loaded * 100) / progressEvent.total
                        );
                        console.log(percentCompleted);
                        // You can update a progress bar here
                    },
                });

                if (response.status === 200 || response.status === 201) {
                    console.log('News added successfully');
                    this.resetForm();
                    this.requestResult = 'success';
                } else {
                    console.error('Error adding news');
                    this.requestResult = 'error';
                }
            } catch (error) {
                console.error('Error submitting form:', error);
                this.requestResult = 'error';
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
                imagescount: 0,
                videoscount: 0,
                text_en: '',
                text_ru: '',
            };
            this.images = [];
            this.videos = [];
            this.previewImages = [];
            this.img_back_preview = null;
            this.img_backfull_preview = null;
            // Сброс input файлов
            const fileInputs = this.$el.querySelectorAll('input[type="file"]');
            fileInputs.forEach((input: HTMLInputElement) => {
                input.value = '';
            });
        },

        validateForm(): boolean {
            if (
                !this.news.title_ru.trim() ||
                this.news.title_ru.trim().length < 3 ||
                this.news.title_ru.trim().length > 50
            ) {
                this.validated = false;
                this.validationErrorMessage =
                    'Заголовок на русском должен быть от 3 до 50 символов';
                return false;
            }

            if (
                !this.news.title_en.trim() ||
                this.news.title_en.trim().length < 3 ||
                this.news.title_en.trim().length > 50
            ) {
                this.validated = false;
                this.validationErrorMessage =
                    'Заголовок на английском должен быть от 3 до 50 символов';
                return false;
            }

            if (
                !this.news.subTitle_ru.trim() ||
                this.news.subTitle_ru.trim().length < 3 ||
                this.news.subTitle_ru.trim().length > 50
            ) {
                this.validated = false;
                this.validationErrorMessage =
                    'Подзаголовок на русском должен быть от 3 до 50 символов';
                return false;
            }
            if (
                !this.news.subTitle_en.trim() ||
                this.news.subTitle_en.trim().length < 3 ||
                this.news.subTitle_en.trim().length > 50
            ) {
                this.validated = false;
                this.validationErrorMessage =
                    'Подзаголовок на английском должен быть от 3 до 50 символов';
                return false;
            }

            if (!this.img_back_preview || !this.news.img_back) {
                this.validated = false;
                this.validationErrorMessage =
                    'Необходимо добавить предварительное изображение новости';
                return false;
            }

            if (!this.img_backfull_preview || !this.news.img_backfull) {
                this.validated = false;
                this.validationErrorMessage =
                    'Необходимо добавить главное изображение новости';
                return false;
            }

            if (this.previewImages.length === 0) {
                this.validated = false;
                this.validationErrorMessage =
                    'Необходимо добавить хотя бы одно изображение';
                return false;
            }

            if (!this.news.text_ru.trim()) {
                this.validated = false;
                this.validationErrorMessage =
                    'Необходимо добавить текст на русском языке';
                return false;
            }

            if (!this.news.text_en.trim()) {
                this.validated = false;
                this.validationErrorMessage =
                    'Необходимо добавить текст на английском языке';
                return false;
            }

            this.validated = true;
            this.validationErrorMessage = '';
            return true;
        },

        closeErrorMessage() {
            this.validated = true;
            this.validationErrorMessage = '';
        },
    },
});
</script>

<template>
    <div class="add-painting-container">
        <h1 class="page-title">Добавить новость</h1>

        <div class="error-message" v-if="!validated">
            <div class="error-message__icon">⚠️</div>
            <div class="error-message__content">
                <h3>Некоторые поля заполнены некорректно</h3>
                <p>
                    {{ validationErrorMessage }}
                </p>
            </div>
            <button class="error-message__close" @click="closeErrorMessage">
                &times;
            </button>
        </div>

        <form @submit.prevent="submitForm" class="painting-form">
            <!-- Заголовок -->
            <div class="form-group">
                <label class="form-label">Заголовок по по-русски</label>
                <input
                    v-model="news.title_ru"
                    type="text"
                    required
                    class="form-control"
                    placeholder="Например: 'Звездная ночь'"
                />
            </div>

            <!-- Заголовок по английски-->
            <div class="form-group">
                <label class="form-label">Заголовок по-английски</label>
                <input
                    v-model="news.title_en"
                    type="text"
                    required
                    class="form-control"
                    placeholder="For example: 'Starry Night'"
                />
            </div>

            <!-- Подзаголовок -->
            <div class="form-group">
                <label class="form-label">Подзаголовок по-русски</label>
                <input
                    v-model="news.subTitle_ru"
                    type="text"
                    required
                    class="form-control"
                    placeholder="Например: 'Звездная ночь'"
                />
            </div>

            <!-- Подзаголовок по английски-->
            <div class="form-group">
                <label class="form-label">Подзаголовок по-английски</label>
                <input
                    v-model="news.subTitle_en"
                    type="text"
                    required
                    class="form-control"
                    placeholder="For example: 'Starry Night'"
                />
            </div>

            <!-- Главное изображение -->
            <div class="form-group">
                <label class="form-label">Главное изображение</label>
                <input
                    type="file"
                    @change="handleBackFullImageSelected"
                    accept="image/jpg,image/jpeg"
                    class="file-input"
                    ref="fileInput"
                    required
                />
                <div class="preview-container" v-if="img_backfull_preview">
                    <div class="image-preview">
                        <img :src="img_backfull_preview.preview" class="preview-image" />
                        <button
                            type="button"
                            @click="removeBackFullImage"
                            class="remove-btn"
                        >
                            &times;
                        </button>
                    </div>
                </div>
            </div>

            <!-- Превью изображение новости -->
            <div class="form-group">
                <label class="form-label">Превью изображение новости</label>
                <input
                    type="file"
                    @change="handleBackImageSelected"
                    accept="image/jpg,image/jpeg"
                    class="file-input"
                    ref="fileInput"
                    required
                />
                <div class="preview-container" v-if="img_back_preview">
                    <div class="image-preview">
                        <img :src="img_back_preview.preview" class="preview-image" />
                        <button type="button" @click="removeBackImage" class="remove-btn">
                            &times;
                        </button>
                    </div>
                </div>
            </div>

            <div class="form-group">
                <label class="form-label">Фотогаллерея</label>
                <input
                    type="file"
                    @change="handleImagesSelected"
                    multiple
                    accept="image/jpg,image/jpeg"
                    class="file-input"
                    ref="fileInput"
                    required
                />
                <div class="preview-container" v-if="previewImages.length > 0">
                    <div
                        v-for="(image, index) in previewImages"
                        :key="index"
                        class="image-preview"
                    >
                        <img :src="image.preview" class="preview-image" />
                        <button
                            type="button"
                            @click="removeImageFromImages(index)"
                            class="remove-btn"
                        >
                            &times;
                        </button>
                    </div>
                </div>
            </div>

            <!-- Поле для выбора видео -->
            <div class="form-group">
                <label class="form-label">Видеогаллерея</label>
                <input
                    type="file"
                    @change="handleVideosSelected"
                    multiple
                    accept="video"
                    class="file-input"
                    ref="fileInput"
                />
                <div v-if="previewVideos" class="video-previews">
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
                        <button
                            type="button"
                            @click="removeVideoFromVideos(index)"
                            class="remove-btn"
                        >
                            &times;
                        </button>
                    </div>
                </div>
            </div>

            <!-- Дата -->
            <div class="form-group form-date">
                <label class="form-label">Дата</label>
                <input
                    v-model="news.datetime"
                    type="date"
                    required
                    min="2000"
                    :max="new Date().getFullYear()"
                    class="form-control"
                    placeholder="Например: 1889"
                />
            </div>

            <!-- Текст -->
            <div class="form-group">
                <label class="form-label">Текст по-русски</label>
                <textarea
                    v-model="news.text_ru"
                    class="form-control textarea"
                    placeholder="Краткое описание картины"
                    rows="4"
                    required
                ></textarea>
            </div>

            <div class="form-group">
                <label class="form-label">Текст по-английски</label>
                <textarea
                    v-model="news.text_en"
                    class="form-control textarea"
                    placeholder="Краткое описание картины"
                    rows="4"
                    required
                ></textarea>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <button type="button" @click="resetForm" class="btn btn-secondary">
                    Очистить форму
                </button>
                <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
                    {{ isSubmitting ? 'Отправка...' : 'Добавить новость' }}
                </button>
            </div>
        </form>
    </div>
</template>

<style scoped>
select:has(option.placeholder:checked) {
    color: red;
}

.add-painting-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    background-color: var(--color-on-surface);
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.page-title {
    text-align: center;
    color: #333;
    margin-bottom: 2rem;
}

/* Validation error */
.error-message {
    background: #fff;
    border-left: 4px solid #ff4757;
    border-radius: 4px;
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    display: flex;
    align-items: flex-start;
    animation: slideIn 0.3s ease-out;
    transform-origin: center top;
}

.error-message__icon {
    margin-right: 15px;
    flex-shrink: 0;
    color: #ff4757;
    font-size: 24px;
}

.error-message__content h3 {
    margin: 0 0 8px 0;
    color: #2f3542;
    font-weight: 600;
    font-size: 18px;
}

.error-message__content p {
    margin: 0;
    color: #747d8c;
    line-height: 1.5;
}

/* Close button */
.error-message__close {
    margin-left: auto;
    background: none;
    border: none;
    color: #a4b0be;
    cursor: pointer;
    font-size: 18px;
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s ease;
}

.error-message__close:hover {
    background: #f1f2f6;
    color: #747d8c;
}
/* End of validation error */

.painting-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.form-label {
    font-weight: 600;
    color: #444;
}

.form-date {
    width: 50%;
}

.form-control {
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    transition: border-color 0.3s;
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

.form-control:focus {
    border-color: #4a90e2;
    outline: none;
}

.file-input {
    padding: 0.5rem;
    border: 1px dashed #ccc;
    border-radius: 4px;
    background-color: #f9f9f9;
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
}

.video-preview {
    position: relative;
    width: 100px;
    height: 100px;
    /* border: 1px solid #eee; */
    border-radius: 4px;
    overflow: hidden;
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
}

.size-inputs {
    display: flex;
    align-items: center;
    gap: 0.5rem;
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
    transition: background-color 0.3s;
}

.btn-primary {
    background-color: #4a90e2;
    color: white;
}

.btn-primary:hover {
    background-color: #3a7bc8;
}

.btn-primary:disabled {
    background-color: #a0c4ff;
    cursor: not-allowed;
}

.btn-secondary {
    background-color: #f0f0f0;
    color: #333;
}

.btn-secondary:hover {
    background-color: #e0e0e0;
}

.alert {
    position: fixed;
    bottom: 1rem;
    left: 1rem;
    padding: 1rem;
    width: 30rem;
}

@media (max-width: 768px) {
    .add-painting-container {
        padding: 1rem;
    }

    .form-actions {
        flex-direction: column;
    }

    .btn {
        width: 100%;
    }
}

.multi-select-wrapper {
    position: relative;
    width: 100%;
}

.select-display {
    padding: 0.5rem;
    border: 1px solid #ced4da;
    border-radius: 0.25rem;
    cursor: pointer;
    background-color: white;
    min-height: 38px;
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
}

.option-item {
    padding: 8px 12px;
    display: flex;
    align-items: center;
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
