<template>
    <div class="review-write-page">
        <div class="review-write-page__container">
            <div class="review-write-page__header">
                <h1 class="review-write-page__title">
                    {{ $t('review.write.title',) }}
                </h1>
                <p class="review-write-page__subtitle">
                    {{ $t('review.write.subtitle',) }}
                </p>
            </div>

            <div
                v-if="product"
                class="review-write-page__product"
            >
                <div class="review-write-page__product-image">
                    <img
                        v-if="product.thumbnail_url"
                        :src="product.thumbnail_url"
                        :alt="productTitle"
                        class="review-write-page__product-thumbnail"
                    >
                    <div
                        v-else
                        class="review-write-page__product-placeholder"
                    >
                        {{ productTypeIcon }}
                    </div>
                </div>
                <div class="review-write-page__product-info">
                    <h2 class="review-write-page__product-title">
                        {{ productTitle }}
                    </h2>
                    <div class="review-write-page__product-category">
                        {{ productCategory }}
                    </div>
                    <div class="review-write-page__product-difficulty">
                        {{ $t('review.write.difficulty',) }}: {{ difficultyLabel }}
                    </div>
                </div>
            </div>

            <div
                v-if="error"
                class="review-write-page__error"
            >
                {{ error }}
            </div>

            <div
                v-if="!hasPurchase"
                class="review-write-page__warning"
            >
                <p>{{ $t('review.write.noPurchase',) }}</p>
                <NuxtLink
                    :to="`/products/${productId}`"
                    class="review-write-page__warning-link"
                >
                    {{ $t('review.write.goToProduct',) }}
                </NuxtLink>
            </div>

            <div v-else>
                <ReviewForm
                    ref="formRef"
                    :is-submitting="isSubmitting"
                    @submit="handleSubmit"
                    @cancel="handleCancel"
                />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, computed, onMounted, } from 'vue';
    import { useRoute, useRouter, } from 'vue-router';
    import { useI18n, } from 'vue-i18n';
    import { useAuthStore, } from '~/stores/AuthStore';
    import { useProductStore, } from '~/stores/ProductStore';
    import ReviewForm from '~/components/ReviewForm.vue';
    import type { CreateReviewDto, } from '~/types';

    const route = useRoute();
    const router = useRouter();
    const { t, } = useI18n();
    const authStore = useAuthStore();
    const productStore = useProductStore();

    const productId = computed(() => {
        const id = route.query.product_id || route.params.product_id;
        return Number(id,);
    });

    const product = ref<any>(null,);
    const hasPurchase = ref(false,);
    const error = ref('',);
    const isSubmitting = ref(false,);
    const formRef = ref<InstanceType<typeof ReviewForm> | null>(null,);

    const productTitle = computed(() => {
        const locale = useI18n().locale.value;
        return locale === 'ru' ? product.value?.title_ru : product.value?.title_en;
    });

    const productCategory = computed(() => {
        const locale = useI18n().locale.value;
        return locale === 'ru'
            ? product.value?.category?.name_ru
            : product.value?.category?.name_en;
    });

    const difficultyLabel = computed(() => {
        const diff = product.value?.difficulty;
        if (!diff) return t('review.write.unknown',);
        const map: Record<string, string> = {
            beginner: t('review.write.beginner',),
            intermediate: t('review.write.intermediate',),
            advanced: t('review.write.advanced',),
        };
        return map[diff] || diff;
    });

    const productTypeIcon = computed(() => {
        return product.value?.type === 'course' ? '📚' : '🎨';
    });

    onMounted(async () => {
        if (!productId.value || isNaN(productId.value,)) {
            error.value = t('review.write.invalidProduct',);
            return;
        }

        try {
            // Загружаем информацию о продукте
            await productStore.fetchProductById(productId.value,);
            product.value = productStore.currentProduct;

            if (!product.value) {
                error.value = t('review.write.productNotFound',);
                return;
            }

            // Проверяем, есть ли у пользователя покупка этого продукта
            const response = await $fetch(`/api/purchases?product_id=${productId.value}`, {
                headers: { Authorization: `Bearer ${authStore.token}`, },
            }).catch(() => null,);

            if (response && Array.isArray(response,) && response.length > 0) {
                hasPurchase.value = true;
            } else {
                hasPurchase.value = false;
            }
        } catch (err) {
            error.value = t('review.write.loadError',);
            console.error(err,);
        }
    });

    async function handleSubmit(data: { rating: number, title?: string, comment?: string },) {
        if (!productId.value || !hasPurchase.value) return;

        isSubmitting.value = true;
        error.value = '';

        const payload: CreateReviewDto = {
            rating: data.rating,
            title_ru: data.title,
            title_en: data.title,
            comment_ru: data.comment,
            comment_en: data.comment,
        };

        try {
            const response = await $fetch(`/api/products/${productId.value}/reviews`, {
                method: 'POST',
                headers: { Authorization: `Bearer ${authStore.token}`, },
                body: payload,
            });

            // Успешно
            router.push(`/products/${productId.value}?review_submitted=true`,);
        } catch (err: any) {
            error.value = err.data?.error || t('review.write.submitError',);
        } finally {
            isSubmitting.value = false;
        }
    }

    function handleCancel() {
        router.push(`/products/${productId.value}`,);
    }
</script>

<style scoped>
.review-write-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 40px 20px;
}

.review-write-page__container {
    background: white;
    border-radius: 16px;
    padding: 32px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.05);
}

.review-write-page__header {
    text-align: center;
    margin-bottom: 32px;
}

.review-write-page__title {
    font-size: 32px;
    font-weight: 700;
    color: #222;
    margin-bottom: 8px;
}

.review-write-page__subtitle {
    font-size: 16px;
    color: #666;
    line-height: 1.5;
}

.review-write-page__product {
    display: flex;
    gap: 20px;
    align-items: center;
    background: #f9f9f9;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 32px;
}

.review-write-page__product-image {
    flex-shrink: 0;
}

.review-write-page__product-thumbnail {
    width: 100px;
    height: 100px;
    border-radius: 12px;
    object-fit: cover;
}

.review-write-page__product-placeholder {
    width: 100px;
    height: 100px;
    border-radius: 12px;
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36px;
}

.review-write-page__product-info {
    flex: 1;
}

.review-write-page__product-title {
    font-size: 20px;
    font-weight: 600;
    color: #222;
    margin-bottom: 4px;
}

.review-write-page__product-category {
    font-size: 14px;
    color: #6a11cb;
    font-weight: 500;
    margin-bottom: 4px;
}

.review-write-page__product-difficulty {
    font-size: 14px;
    color: #666;
}

.review-write-page__error {
    background: #ffebee;
    color: #d32f2f;
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 24px;
    font-weight: 500;
}

.review-write-page__warning {
    background: #fff3e0;
    color: #e65100;
    padding: 24px;
    border-radius: 12px;
    text-align: center;
    margin-bottom: 24px;
}

.review-write-page__warning-link {
    display: inline-block;
    margin-top: 12px;
    color: #6a11cb;
    font-weight: 600;
    text-decoration: underline;
}
</style>
