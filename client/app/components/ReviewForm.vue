<template>
    <form
        class="review-form"
        @submit.prevent="handleSubmit"
    >
        <div class="review-form__section">
            <label class="review-form__label">Оценка</label>
            <div class="review-form__rating">
                <button
                    v-for="star in 5"
                    :key="star"
                    type="button"
                    class="review-form__star"
                    :class="{ 'review-form__star--active': star <= model.rating, }"
                    @click="model.rating = star"
                >
                    ★
                </button>
            </div>
            <div class="review-form__rating-value">
                {{ model.rating }}/5
            </div>
            <div
                v-if="errors.rating"
                class="review-form__error"
            >
                {{ errors.rating }}
            </div>
        </div>

        <div class="review-form__section">
            <label
                class="review-form__label"
                for="title"
            >Заголовок отзыва (необязательно)</label>
            <input
                id="title"
                v-model="model.title"
                type="text"
                class="review-form__input"
                :placeholder="titlePlaceholder"
                maxlength="100"
            >
            <div
                v-if="errors.title"
                class="review-form__error"
            >
                {{ errors.title }}
            </div>
        </div>

        <div class="review-form__section">
            <label
                class="review-form__label"
                for="comment"
            >Комментарий (необязательно)</label>
            <textarea
                id="comment"
                v-model="model.comment"
                class="review-form__textarea"
                :placeholder="commentPlaceholder"
                rows="4"
                maxlength="2000"
            />
            <div
                v-if="errors.comment"
                class="review-form__error"
            >
                {{ errors.comment }}
            </div>
        </div>

        <div class="review-form__actions">
            <button
                type="button"
                class="review-form__button review-form__button--cancel"
                @click="$emit('cancel',)"
            >
                Отмена
            </button>
            <button
                type="submit"
                class="review-form__button review-form__button--submit"
                :disabled="isSubmitting"
            >
                <span v-if="isSubmitting">Отправка...</span>
                <span v-else>{{ submitText }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
    import { reactive, computed, watch, } from 'vue';
    import { useI18n, } from 'vue-i18n';

    interface Props {
        initialData?: {
            rating?: number
            title_ru?: string | null
            title_en?: string | null
            comment_ru?: string | null
            comment_en?: string | null
        }
        isEditing?: boolean
        isSubmitting?: boolean
    }

    const props = withDefaults(defineProps<Props>(), {
        initialData: () => ({}),
        isEditing: false,
        isSubmitting: false,
    });

    const emit = defineEmits<{
        submit: [data: { rating: number, title?: string, comment?: string },]
        cancel: []
    }>();

    const { locale, } = useI18n();

    const model = reactive({
        rating: props.initialData.rating || 5,
        title:
            locale.value === 'ru'
                ? props.initialData.title_ru || ''
                : props.initialData.title_en || '',
        comment:
            locale.value === 'ru'
                ? props.initialData.comment_ru || ''
                : props.initialData.comment_en || '',
    });

    const errors = reactive({
        rating: '',
        title: '',
        comment: '',
    });

    const titlePlaceholder = computed(() => {
        return locale.value === 'ru' ? 'Напишите краткий заголовок' : 'Write a short title';
    });

    const commentPlaceholder = computed(() => {
        return locale.value === 'ru'
            ? 'Поделитесь вашим опытом прохождения курса'
            : 'Share your experience with the course';
    });

    const submitText = computed(() => {
        if (props.isEditing) {
            return locale.value === 'ru' ? 'Сохранить изменения' : 'Save changes';
        }
        return locale.value === 'ru' ? 'Отправить отзыв' : 'Submit review';
    });

    watch(
        () => locale.value,
        (newLocale,) => {
            // При смене языка обновляем поля из initialData
            model.title
            = newLocale === 'ru'
                    ? props.initialData.title_ru || ''
                    : props.initialData.title_en || '';
            model.comment
            = newLocale === 'ru'
                    ? props.initialData.comment_ru || ''
                    : props.initialData.comment_en || '';
        }
    );

    function validate(): boolean {
        let valid = true;
        if (model.rating < 1 || model.rating > 5) {
            errors.rating = 'Оценка должна быть от 1 до 5';
            valid = false;
        } else {
            errors.rating = '';
        }
        if (model.title.length > 100) {
            errors.title = 'Заголовок не должен превышать 100 символов';
            valid = false;
        } else {
            errors.title = '';
        }
        if (model.comment.length > 2000) {
            errors.comment = 'Комментарий не должен превышать 2000 символов';
            valid = false;
        } else {
            errors.comment = '';
        }
        return valid;
    }

    function handleSubmit() {
        if (!validate()) return;

        const submitData = {
            rating: model.rating,
            title: model.title.trim() || undefined,
            comment: model.comment.trim() || undefined,
        };
        emit('submit', submitData,);
    }
</script>

<style scoped>
.review-form {
    background: white;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.review-form__section {
    margin-bottom: 24px;
}

.review-form__label {
    display: block;
    font-weight: 600;
    color: #333;
    margin-bottom: 8px;
    font-size: 15px;
}

.review-form__rating {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
}

.review-form__star {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    border: 2px solid #ddd;
    background: white;
    color: #ddd;
    font-size: 28px;
    line-height: 1;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
}

.review-form__star:hover {
    border-color: #ffc107;
    color: #ffc107;
}

.review-form__star--active {
    border-color: #ffc107;
    background: #ffc107;
    color: white;
}

.review-form__rating-value {
    font-size: 16px;
    color: #666;
    font-weight: 500;
}

.review-form__input,
.review-form__textarea {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 15px;
    font-family: inherit;
    transition: border-color 0.2s;
}

.review-form__input:focus,
.review-form__textarea:focus {
    outline: none;
    border-color: #6a11cb;
    box-shadow: 0 0 0 3px rgba(106, 17, 203, 0.1);
}

.review-form__textarea {
    resize: vertical;
    min-height: 100px;
}

.review-form__error {
    color: #d32f2f;
    font-size: 14px;
    margin-top: 6px;
}

.review-form__actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 32px;
    padding-top: 20px;
    border-top: 1px solid #eee;
}

.review-form__button {
    padding: 10px 24px;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
}

.review-form__button--cancel {
    background: #f5f5f5;
    color: #666;
}

.review-form__button--cancel:hover {
    background: #e0e0e0;
}

.review-form__button--submit {
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
}

.review-form__button--submit:hover:not(:disabled) {
    opacity: 0.9;
}

.review-form__button--submit:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}
</style>
