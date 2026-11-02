<template>
    <div
        class="admin-uploader"
        :class="{
            'admin-uploader--dark': isDark,
            'admin-uploader--dragover': isDragover,
            'admin-uploader--has-image': !!modelValue,
        }"
        @dragover.prevent="isDragover = true"
        @dragleave.prevent="isDragover = false"
        @drop.prevent="onDrop"
    >
        <!-- Hidden file input -->
        <input
            ref="fileInput"
            type="file"
            :accept="accept"
            class="admin-uploader__input"
            @change="onFileSelected"
        >

        <!-- Image preview -->
        <div
            v-if="modelValue"
            class="admin-uploader__preview"
        >
            <img
                :src="modelValue"
                alt="preview"
                class="admin-uploader__image"
                @error="onImageError"
            >
            <div class="admin-uploader__overlay">
                <button
                    type="button"
                    class="admin-uploader__action-btn"
                    title="Заменить"
                    @click="openFilePicker"
                >
                    <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        width="18"
                        height="18"
                    >
                        <path
                            d="M23 19a2 2 0 01-2 2H3a2 2 0 01-2-2V8a2 2 0 012-2h4l2-3h6l2 3h4a2 2 0 012 2z"
                        />
                        <circle
                            cx="12"
                            cy="13"
                            r="4"
                        />
                    </svg>
                </button>
                <button
                    type="button"
                    class="admin-uploader__action-btn admin-uploader__action-btn--remove"
                    title="Удалить"
                    @click="onRemove"
                >
                    <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        width="18"
                        height="18"
                    >
                        <line
                            x1="18"
                            y1="6"
                            x2="6"
                            y2="18"
                        />
                        <line
                            x1="6"
                            y1="6"
                            x2="18"
                            y2="18"
                        />
                    </svg>
                </button>
            </div>
        </div>

        <!-- Upload placeholder -->
        <div
            v-else
            class="admin-uploader__placeholder"
            @click="openFilePicker"
        >
            <div class="admin-uploader__icon">
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    width="40"
                    height="40"
                >
                    <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
                    <polyline points="17 8 12 3 7 8" />
                    <line
                        x1="12"
                        y1="3"
                        x2="12"
                        y2="15"
                    />
                </svg>
            </div>
            <p class="admin-uploader__text">
                {{ uploadText }}
            </p>
            <p class="admin-uploader__hint">
                {{ hintText }}
            </p>
        </div>

        <!-- Upload progress -->
        <div
            v-if="uploading"
            class="admin-uploader__progress"
        >
            <div
                class="admin-uploader__progress-bar"
                :style="{ width: uploadProgress + '%', }"
            />
        </div>

        <!-- Error message -->
        <p
            v-if="error"
            class="admin-uploader__error"
        >
            {{ error }}
        </p>
    </div>
</template>

<script setup lang="ts">
import { ref, } from 'vue';

const props = withDefaults(
    defineProps<{
        modelValue: string | null
        accept?: string
        maxSize?: number // in bytes
        uploadText?: string
        hintText?: string
        isDark?: boolean
    }>(),
    {
        modelValue: null,
        accept: 'image/jpeg,image/png,image/webp,image/gif',
        maxSize: 10 * 1024 * 1024, // 10MB
        uploadText: 'Перетащите изображение сюда или нажмите для выбора',
        hintText: 'JPEG, PNG, WebP до 10MB',
        isDark: false,
    }
);

const emit = defineEmits<{
    'update:modelValue': [value: string | null,]
    fileSelect: [file: File,]
    uploadStart: []
    uploadComplete: [url: string,]
    uploadError: [error: string,]
}>();

const fileInput = ref<HTMLInputElement | null>(null,);
const isDragover = ref(false,);
const uploading = ref(false,);
const uploadProgress = ref(0,);
const error = ref('',);

function openFilePicker() {
    fileInput.value?.click();
}

function onFileSelected(e: Event,) {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files[0]) {
        validateAndUpload(target.files[0],);
    }
}

function onDrop(e: DragEvent,) {
    isDragover.value = false;
    if (e.dataTransfer?.files && e.dataTransfer.files[0]) {
        validateAndUpload(e.dataTransfer.files[0],);
    }
}

function validateAndUpload(file: File,) {
    error.value = '';

    // Validate type
    const allowedTypes = props.accept.split(',',);
    const isAllowed = allowedTypes.some((t,) => {
        const type = t.trim();
        if (type.endsWith('/*',)) {
            return file.type.startsWith(type.replace('/*', '/',),);
        }
        return file.type === type;
    });

    if (!isAllowed) {
        error.value = 'Неподдерживаемый формат файла';
        return;
    }

    // Validate size
    if (file.size > props.maxSize) {
        const maxMB = Math.round(props.maxSize / (1024 * 1024),);
        error.value = `Файл слишком большой. Максимальный размер: ${maxMB}MB`;
        return;
    }

    emit('fileSelect', file,);

    // Create local preview
    const reader = new FileReader();
    reader.onload = (e,) => {
        emit('update:modelValue', e.target?.result as string,);
    };
    reader.readAsDataURL(file,);

    // Simulate upload (override in parent)
    emit('uploadStart',);
}

function onRemove() {
    emit('update:modelValue', null,);
    if (fileInput.value) {
        fileInput.value.value = '';
    }
}

function onImageError() {
    error.value = 'Не удалось загрузить изображение';
}
</script>

<style scoped>
.admin-uploader {
    position: relative;
    border: 2px dashed var(--admin-border, #e0e0e0);
    border-radius: 12px;
    overflow: hidden;
    transition: all 0.2s;
    cursor: pointer;
    min-height: 200px;
}

.admin-uploader--dark {
    border-color: var(--admin-border, #2d2d3d);
}

.admin-uploader:hover {
    border-color: var(--admin-primary, #6c5ce7);
}

.admin-uploader--dragover {
    border-color: var(--admin-primary, #6c5ce7);
    background: rgba(108, 92, 231, 0.05);
}

.admin-uploader--dark.admin-uploader--dragover {
    background: rgba(108, 92, 231, 0.1);
}

.admin-uploader--has-image {
    border-style: solid;
    cursor: default;
}

.admin-uploader__input {
    display: none;
}

/* Placeholder */
.admin-uploader__placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
    cursor: pointer;
    min-height: 200px;
}

.admin-uploader__icon {
    color: var(--admin-text-secondary, #636e72);
    opacity: 0.4;
    margin-bottom: 12px;
}

.admin-uploader__text {
    font-size: 14px;
    color: var(--admin-text-primary, #2d3436);
    margin: 0 0 6px;
    text-align: center;
}

.admin-uploader--dark .admin-uploader__text {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-uploader__hint {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    margin: 0;
}

/* Preview */
.admin-uploader__preview {
    position: relative;
    width: 100%;
    min-height: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.admin-uploader__image {
    max-width: 100%;
    max-height: 300px;
    object-fit: contain;
    display: block;
}

.admin-uploader__overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    opacity: 0;
    transition: opacity 0.2s;
}

.admin-uploader__preview:hover .admin-uploader__overlay {
    opacity: 1;
}

.admin-uploader__action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.2);
    border: none;
    color: white;
    cursor: pointer;
    transition: background 0.2s;
}

.admin-uploader__action-btn:hover {
    background: rgba(255, 255, 255, 0.35);
}

.admin-uploader__action-btn--remove:hover {
    background: rgba(225, 112, 85, 0.8);
}

/* Progress */
.admin-uploader__progress {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: var(--admin-bg, #f0f2f5);
}

.admin-uploader--dark .admin-uploader__progress {
    background: rgba(255, 255, 255, 0.1);
}

.admin-uploader__progress-bar {
    height: 100%;
    background: var(--admin-primary, #6c5ce7);
    transition: width 0.3s ease;
    border-radius: 0 2px 2px 0;
}

/* Error */
.admin-uploader__error {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    margin: 0;
    padding: 8px 12px;
    background: rgba(225, 112, 85, 0.9);
    color: white;
    font-size: 12px;
    text-align: center;
}
</style>
