<template>
    <div class="add-sale-container">
        <h1 class="page-title">Добавить новую работу в магазин</h1>

        <!-- Загрузчик -->
        <div v-if="isLoading" class="loading-container">
            <div class="loader"></div>
            <p>Загрузка данных...</p>
        </div>

        <form v-else @submit.prevent="submitForm" class="sale-form">
            <!-- Поле для загрузки изображений -->
            <div class="form-group">
                <label class="form-label">Изображения картины</label>
                <input
                    type="file"
                    @change="handleFileUpload"
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
                            @click="removeImage(index)"
                            class="remove-btn"
                        >
                            &times;
                        </button>
                    </div>
                </div>
            </div>

            <!-- Название картины -->
            <div class="form-group">
                <label class="form-label">Название картины по-русски</label>
                <input
                    type="text"
                    v-model="sale.name_ru"
                    required
                    class="form-control"
                    placeholder="Например: 'Звездная ночь'"
                />
            </div>

            <div class="form-group">
                <label class="form-label">Название картины по-английски</label>
                <input
                    type="text"
                    v-model="sale.name_en"
                    required
                    class="form-control"
                    placeholder="Например: 'Звездная ночь'"
                />
            </div>

            <!-- Размеры картины -->
            <div class="form-group">
                <label class="form-label">Размеры (см)</label>
                <div class="size-inputs">
                    <input
                        type="number"
                        v-model.number="sale.width"
                        required
                        min="1"
                        class="form-control size-input"
                        placeholder="Ширина"
                    />
                    <span class="size-separator">×</span>
                    <input
                        type="number"
                        v-model.number="sale.height"
                        required
                        min="1"
                        class="form-control size-input"
                        placeholder="Высота"
                    />
                </div>
            </div>

            <!-- Год создания -->
            <div class="form-group">
                <label class="form-label">Год создания</label>
                <input
                    type="number"
                    v-model.number="sale.year"
                    required
                    min="2000"
                    :max="new Date().getFullYear()"
                    class="form-control"
                    placeholder="Например: 1889"
                />
            </div>

            <!-- Цена -->
            <div class="form-group">
                <label class="form-label">Цена</label>
                <input
                    type="number"
                    v-model.number="sale.price"
                    required
                    min="1"
                    class="form-control"
                    placeholder="Например: 1000"
                />
            </div>

            <!-- Основа -->
            <div class="form-group">
                <label class="form-label">Основа</label>
                <select
                    v-model="sale.base_id"
                    required
                    class="form-control drop-down-arrow"
                >
                    <option value="" disabled selected>Выберите основу</option>
                    <option v-for="base in bases" :key="base.id" :value="base.id">
                        {{
                            $i18n.locale === 'RUS' ? `${base.base_ru}` : `${base.base_en}`
                        }}
                    </option>
                </select>
            </div>

            <!-- Материал -->
            <div class="form-group">
                <label class="form-label">Материал</label>
                <div class="multi-select-wrapper">
                    <div
                        class="select-display drop-down-arrow"
                        @click="materialsToggleDropdown"
                    >
                        {{ selectedMaterialsDisplay || 'Выберите материал' }}
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
                                v-model="sale.materials_ids"
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
                </div>
            </div>

            <!-- Описание -->
            <div class="form-group">
                <label class="form-label">Описание (необязательно)</label>
                <textarea
                    v-model="sale.descr"
                    class="form-control textarea"
                    placeholder="Краткое описание картины"
                    rows="4"
                ></textarea>
            </div>

            <!-- Кнопки -->
            <div class="form-actions">
                <button type="button" @click="resetForm" class="btn btn-secondary">
                    Очистить форму
                </button>
                <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
                    <span v-if="!isSubmitting">Добавить картину</span>
                    <span v-else>Отправка...</span>
                </button>
            </div>
        </form>

        <!-- Success Alert -->
        <div
            v-if="requestResult == 'success'"
            class="alert alert-success alert-dismissible fade show"
            role="alert"
        >
            <strong>Успешно!</strong> Работы успешно загружены в базу данных.
            <button
                type="button"
                class="btn-close"
                data-bs-dismiss="alert"
                aria-label="Close"
            ></button>
        </div>

        <!-- Danger Alert -->
        <div
            v-if="requestResult == 'error'"
            class="alert alert-danger alert-dismissible fade show"
            role="alert"
        >
            <strong>Ошибка!</strong> Не удалось загрузить работы в базу данных.
            <button
                type="button"
                class="btn-close"
                data-bs-dismiss="alert"
                aria-label="Close"
            ></button>
        </div>
    </div>
</template>

<script lang="ts">
import axios from 'axios';
import { defineComponent } from 'vue';
import { AddSaleDto, Sale, Base, Material, RequestResult } from '@/types';
export default defineComponent({
    name: 'AddSale',
    data() {
        return {
            sale: {
                width: 0,
                height: 0,
                year: 2000,
                price: 0,
                name_ru: '',
                name_en: '',
                base_id: 0,
                materials_ids: [],
                img_count: 0,
                descr: '',
            } as Sale,
            files: [],
            previewImages: [],
            isSubmitting: false,
            bases: [] as Base[], // Здесь будут храниться основы с сервера
            materials: [] as Material[], // Здесь будут храниться материалы с сервера
            isLoading: true,
            loadError: null,
            server: import.meta.env.VITE_SERVER_URL,
            requestResult: 'unknown' as RequestResult,
            materialsDropdownOpen: false,
            basesDropdownOpen: false,
        };
    },
    async created() {
        await this.loadBases();
        await this.loadMaterials();
    },
    methods: {
        async loadBases() {
            try {
                // Запрос к API для получения списка основ
                const response = await axios.get(this.server + 'bases');
                this.bases = response.data;
                this.isLoading = false;
            } catch (error) {
                console.error('Ошибка при загрузке основ:', error);
                this.loadError = 'Не удалось загрузить список основ';
                this.isLoading = false;
            }
        },
        async loadMaterials() {
            try {
                // Запрос к API для получения списка материалов
                const response = await axios.get(this.server + 'materials');
                this.materials = response.data;
                this.isLoadingM = false;
            } catch (error) {
                console.error('Ошибка при загрузке материалов:', error);
                this.loadError = 'Не удалось загрузить список материалов';
                this.isLoading = false;
            }
        },
        handleFileUpload(event) {
            const selectedFiles = Array.from(event.target.files);

            // Проверка на количество файлов
            if (this.files.length + selectedFiles.length > 10) {
                alert('Можно загрузить не более 10 изображений');
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
                        preview: e.target.result,
                    });
                };
                reader.readAsDataURL(file);
            });
        },
        removeImage(index) {
            this.previewImages.splice(index, 1);
            this.files.splice(index, 1);
        },
        async submitForm() {
            if (this.isSubmitting) return;

            try {
                this.isSubmitting = true;

                // Формируем данные для отправки
                const formData = new FormData();

                // Добавляем файлы
                this.files.forEach((file) => {
                    formData.append('images', file);
                });

                // Добавляем остальные данные
                const saleData = {
                    ...this.sale,
                    img_count: this.previewImages.length,
                };

                formData.append('data', JSON.stringify(saleData));
                const response = await axios.post(this.server + 'sales', formData, {
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

                if (response.status === 200) {
                    console.log('Sales added successfully');
                    this.resetForm();
                    this.requestResult = 'success';
                } else {
                    console.error('Error adding sales');
                    this.requestResult = 'error';
                }

                this.resetForm();
            } catch (error) {
                console.error('Error submitting form:', error);
                this.requestResult = 'error';
            } finally {
                this.isSubmitting = false;
            }
        },
        resetForm() {
            this.sale = {
                width: 0,
                height: 0,
                year: 2000,
                price: 0,
                name_ru: '',
                name_en: '',
                base_id: 0,
                materials_ids: [],
                img_count: 0,
                descr: '',
            };
            this.files = [];
            this.previewImages = [];
            this.$refs.fileInput.value = '';
        },
        materialsToggleDropdown() {
            this.materialsDropdownOpen = !this.materialsDropdownOpen;
        },
        basesToggleDropdown() {
            this.basesDropdownOpen = !this.basesDropdownOpen;
        },
    },
    computed: {
        selectedMaterialsDisplay() {
            if (this.sale.materials_ids.length === 0) return '';
            const selectedNames = this.materials
                .filter((material) => this.sale.materials_ids.includes(material.id))
                .map((material) =>
                    this.$i18n.locale === 'RUS'
                        ? material.material_ru
                        : material.material_en
                );
            return selectedNames.join(', ');
        },
    },
});
</script>

<style scoped>
select:has(option.placeholder:checked) {
    color: red;
}

.add-sale-container {
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

.sale-form {
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
    .add-sale-container {
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
