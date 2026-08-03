<script setup lang="ts">
    import { useAuthStore, } from '~/stores/AuthStore';
    import { getHttpClient, } from '~/api/http-client';
    import { useI18n, } from 'vue-i18n';
    import { useRouter, } from 'vue-router';
    import { ref, onMounted, onUnmounted, } from 'vue';

    interface NewsItem {
        id: string
        datetime: string
        title_ru: string
        title_en: string
        dir: string
        main_image: string
        text_ru: string
        text_en: string
        images: string[]
        videos: string[]
    }

    const showDeleteModal = ref(false,);
    const newsToDelete = ref<NewsItem | null>(null,);
    const deleteError = ref('',);

    const { locale, t, } = useI18n();
    const router = useRouter();
    const authStore = useAuthStore();

    const news = ref<NewsItem[]>([],);
    const page = ref(0,);
    const limit = ref(9,);
    const loading = ref(false,);
    const initialLoading = ref(true,);
    const hasMore = ref(true,);
    const observer = ref<IntersectionObserver | null>(null,);
    const observerElement = ref<HTMLElement | null>(null,);
    const deletingId = ref<string | null>(null,);
    const error = ref('',);

    const loadNews = async (initial = false,) => {
        if (loading.value || (!hasMore.value && !initial)) return;

        try {
            error.value = '';
            loading.value = true;
            const currentPage = initial ? 1 : page.value;
            const { data, } = await getHttpClient().get<{ news: NewsItem[], total: number }>(
                'news',
                {
                    params: {
                        page: currentPage,
                        limit: limit.value,
                    },
                }
            );

            const newNews = data.news || [];

            if (initial) {
                news.value = newNews;
                page.value = 2;
            } else {
                news.value = [...news.value, ...newNews,];
                page.value += 1;
            }

            if (newNews.length < limit.value) {
                hasMore.value = false;
            }
        } catch (e) {
            console.error('Error fetching news', e,);
            if (news.value.length === 0) {
                error.value = t('news.error.loadFailed',);
            }
        } finally {
            loading.value = false;
            initialLoading.value = false;
        }
    };

    const loadMoreNews = async () => {
        if (!hasMore.value || loading.value) return;
        await loadNews(false,);
    };

    const formatDate = (dateString: string,) => {
        const date = new Date(dateString,);
        const options: Intl.DateTimeFormatOptions = {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        };
        return date.toLocaleDateString(locale.value === 'ru' ? 'ru-RU' : 'en-US', options,);
    };

    const getNewsTitle = (newsItem: NewsItem,) => {
        return locale.value === 'ru' ? newsItem.title_ru : newsItem.title_en;
    };

    const getImageUrl = (newsItem: NewsItem,) => {
        return newsItem.dir + newsItem.main_image;
    };

    const navigateToNews = (id: string,) => {
        router.push(`/news/${id}`,);
    };

    const editNews = (id: string,) => {
        router.push(`/news/edit/${id}/`,);
    };

    const confirmDelete = (newsItem: NewsItem,) => {
        newsToDelete.value = newsItem;
        showDeleteModal.value = true;
    };

    const deleteNews = async () => {
        if (!newsToDelete.value) return;

        deletingId.value = newsToDelete.value.id;
        showDeleteModal.value = false;
        deleteError.value = '';

        try {
            await getHttpClient().delete(`news/${newsToDelete.value.id}`,);

            news.value = news.value.filter(n => n.id !== newsToDelete.value!.id,);
        } catch (error) {
            console.error('Error deleting news:', error,);

            if (error && typeof error === 'object' && 'response' in error) {
                const axiosError = error as { response?: { status?: number } };
                if (axiosError.response?.status === 401) {
                    deleteError.value = t('news.sessionExpired',);
                    authStore.clearTokens();
                    setTimeout(() => {
                        router.push('/auth/login',);
                    }, 2000,);
                } else {
                    deleteError.value = t('news.deleteErrorGeneric',);
                }
            } else {
                deleteError.value = t('news.deleteErrorGeneric',);
            }
        } finally {
            deletingId.value = null;
            newsToDelete.value = null;
        }
    };

    const cancelDelete = () => {
        showDeleteModal.value = false;
        newsToDelete.value = null;
    };

    onMounted(async () => {
        await loadNews(true,);

        const options = {
            root: null,
            rootMargin: '100px',
            threshold: 0.1,
        };

        observer.value = new IntersectionObserver((entries,) => {
            const entry = entries[0];
            if (entry?.isIntersecting && hasMore.value && !loading.value) {
                loadMoreNews();
            }
        }, options,);

        if (observerElement.value) {
            observer.value.observe(observerElement.value,);
        }
    });

    onUnmounted(() => {
        if (observer.value) {
            observer.value.disconnect();
        }
    });
</script>

<template>
    <div class="news-page">
        <div class="container mx-auto px-4 py-8">
            <!-- Page Header -->
            <div class="mb-12 text-center">
                <h1
                    class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4"
                >
                    {{ $t('news.title',) }}
                </h1>
            </div>

            <!-- News Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                <UCard
                    v-for="newsItem in news"
                    :key="newsItem.id"
                    class="news-card overflow-hidden hover:shadow-xl transition-shadow duration-300 justify-between"
                >
                    <template #header>
                        <!-- Date -->
                        <div
                            class="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-2"
                        >
                            <UIcon
                                name="i-heroicons-calendar"
                                class="w-4 h-4 mr-2"
                            />
                            {{ formatDate(newsItem.datetime,) }}
                        </div>
                        <!-- Title -->
                        <h3
                            class="text-xl font-bold text-gray-900 dark:text-white line-clamp-2 hover:text-green-600 dark:hover:text-green-400 transition-colors cursor-pointer"
                            @click="navigateToNews(newsItem.id,)"
                        >
                            {{ getNewsTitle(newsItem,) }}
                        </h3>
                    </template>

                    <!-- Default slot: image + content -->
                    <div class="space-y-4">
                        <!-- Image -->
                        <div
                            class="relative overflow-hidden rounded-lg cursor-pointer"
                            @click="navigateToNews(newsItem.id,)"
                        >
                            <img
                                :src="getImageUrl(newsItem,)"
                                :alt="getNewsTitle(newsItem,)"
                                class="w-full h-48 object-cover hover:scale-105 transition-transform duration-500"
                                loading="lazy"
                            >
                        </div>

                        <!-- Read More -->
                        <div
                            class="flex items-center text-green-600 dark:text-green-400 font-medium cursor-pointer"
                            @click="navigateToNews(newsItem.id,)"
                        >
                            <span class="mr-2">{{ $t('news.readMore',) }}</span>
                            <UIcon
                                name="i-heroicons-arrow-right"
                                class="w-4 h-4"
                            />
                        </div>

                        <!-- Admin Controls -->
                        <ClientOnly>
                            <div
                                v-if="authStore.isAuthenticated"
                                class="pt-4 mt-4 border-t border-gray-200 dark:border-gray-700"
                            >
                                <div class="flex gap-2">
                                    <UButton
                                        size="sm"
                                        color="primary"
                                        variant="outline"
                                        @click="editNews(newsItem.id,)"
                                    >
                                        {{ $t('news.edit',) }}
                                    </UButton>
                                    <UButton
                                        size="sm"
                                        color="error"
                                        variant="outline"
                                        :loading="deletingId === newsItem.id"
                                        :disabled="deletingId === newsItem.id"
                                        @click="confirmDelete(newsItem,)"
                                    >
                                        {{ $t('news.delete',) }}
                                    </UButton>
                                </div>
                            </div>
                        </ClientOnly>
                    </div>
                </UCard>
            </div>

            <!-- Initial Loading State (beautiful centered loader) -->
            <div
                v-if="initialLoading"
                class="mt-24"
            >
                <div class="flex flex-col items-center justify-center space-y-6">
                    <!-- Animated spinner -->
                    <div class="news-loader-spinner">
                        <svg
                            class="news-loader-circle"
                            viewBox="0 0 50 50"
                        >
                            <circle
                                class="news-loader-path"
                                cx="25"
                                cy="25"
                                r="20"
                                fill="none"
                                stroke-width="4"
                            />
                        </svg>
                    </div>
                    <div class="flex flex-col items-center space-y-2">
                        <p class="text-lg font-medium text-gray-600 dark:text-gray-300">
                            {{ $t('news.loading',) }}
                        </p>
                        <div class="flex space-x-1.5">
                            <span
                                class="news-loader-dot w-2 h-2 rounded-full bg-green-500 animate-bounce"
                                style="animation-delay: 0s"
                            />
                            <span
                                class="news-loader-dot w-2 h-2 rounded-full bg-green-500 animate-bounce"
                                style="animation-delay: 0.15s"
                            />
                            <span
                                class="news-loader-dot w-2 h-2 rounded-full bg-green-500 animate-bounce"
                                style="animation-delay: 0.3s"
                            />
                        </div>
                    </div>
                </div>
            </div>

            <!-- Error State (initial load failed) -->
            <div
                v-if="error && !initialLoading && !loading && news.length === 0"
                class="mt-12 text-center py-12"
            >
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-red-100 dark:bg-red-900/30 mb-6"
                >
                    <svg
                        class="w-10 h-10 text-red-500"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                    {{ $t('news.error.title',) }}
                </h3>
                <p class="text-gray-600 dark:text-gray-300 max-w-md mx-auto">
                    {{ error }}
                </p>
                <UButton
                    class="mt-6"
                    color="primary"
                    variant="solid"
                    @click="loadNews(true)"
                >
                    {{ $t('news.error.retryButton',) }}
                </UButton>
            </div>

            <!-- Skeleton Loading State (for load more) -->
            <div
                v-if="loading && !initialLoading"
                class="mt-12"
            >
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    <div
                        v-for="n in 3"
                        :key="n"
                        class="news-card-skeleton"
                    >
                        <div class="skeleton-image h-48 rounded-t-2xl" />
                        <div class="p-6">
                            <div class="skeleton-line h-4 w-24 mb-4" />
                            <div class="skeleton-line h-6 w-full mb-3" />
                            <div class="skeleton-line h-4 w-full mb-2" />
                            <div class="skeleton-line h-4 w-3/4" />
                        </div>
                    </div>
                </div>
            </div>

            <!-- No More News Message -->
            <div
                v-if="!hasMore && news.length > 0"
                class="mt-12 text-center py-8"
            >
                <div
                    class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 mb-4"
                >
                    <svg
                        class="w-8 h-8 text-green-600 dark:text-green-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 13l4 4L19 7"
                        />
                    </svg>
                </div>
                <p class="text-lg text-gray-600 dark:text-gray-300">
                    {{ $t('news.allLoaded',) }}
                </p>
            </div>

            <!-- No News Message (only after initial load is complete) -->
            <div
                v-if="!initialLoading && !loading && news.length === 0 && !error"
                class="mt-12 text-center py-12"
            >
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gray-100 dark:bg-gray-800 mb-6"
                >
                    <svg
                        class="w-10 h-10 text-gray-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                        />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                    {{ $t('news.noNews',) }}
                </h3>
                <p class="text-gray-600 dark:text-gray-300 max-w-md mx-auto">
                    {{ $t('news.noNewsDescription',) }}
                </p>
            </div>

            <!-- Delete Error Alert -->
            <UAlert
                v-if="deleteError"
                :title="$t('common.error',)"
                :description="deleteError"
                icon="i-heroicons-exclamation-triangle"
                color="error"
                variant="outline"
                class="mt-6"
                @close="deleteError = ''"
            />

            <!-- Intersection Observer Sentinel -->
            <div
                ref="observerElement"
                class="h-1 w-full"
                :class="{ 'opacity-0': !hasMore || loading, }"
            />
        </div>

        <!-- Delete Confirmation Modal -->
        <UModal
            v-model:open="showDeleteModal"
            :dismissible="false"
            :close="false"
            :transition="true"
            class="delete-modal"
        >
            <template #body>
                <div class="delete-modal-content">
                    <div class="delete-modal-icon-wrapper">
                        <UIcon
                            name="i-heroicons-exclamation-triangle"
                            class="delete-modal-icon"
                        />
                    </div>
                    <div class="delete-modal-main">
                        <h3 class="delete-modal-title">
                            {{ $t('news.deleteConfirmTitle',) }}
                        </h3>
                        <p class="delete-modal-text">
                            <template v-if="newsToDelete">
                                {{ $t('news.deleteConfirmBefore',) }}
                                <span class="delete-modal-highlight">«{{ getNewsTitle(newsToDelete,) }}»</span>
                                {{ $t('news.deleteConfirmAfter',) }}
                            </template>
                        </p>
                        <div class="delete-modal-warning">
                            <UIcon
                                name="i-heroicons-exclamation-circle"
                                class="delete-modal-warning-icon"
                            />
                            <span>{{ $t('news.deleteIrreversible',) }}</span>
                        </div>
                    </div>
                </div>
            </template>

            <template #footer>
                <div class="delete-modal-footer">
                    <UButton
                        size="md"
                        color="neutral"
                        variant="ghost"
                        class="delete-modal-btn"
                        @click="cancelDelete"
                    >
                        {{ $t('common.cancel',) }}
                    </UButton>
                    <div class="delete-modal-footer-divider" />
                    <UButton
                        size="md"
                        color="error"
                        class="delete-modal-btn"
                        :loading="deletingId !== null"
                        :disabled="deletingId !== null"
                        @click="deleteNews"
                    >
                        {{ $t('common.delete',) }}
                    </UButton>
                </div>
            </template>
        </UModal>
    </div>
</template>

<style scoped>
.news-page {
    min-height: 100vh;
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

.dark .news-page {
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

.news-card {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.dark .news-card {
    background: #1e293b;
}

.line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.line-clamp-3 {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

/* Skeleton Loading Styles */
.news-card-skeleton {
    background: white;
    border-radius: 1rem;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    overflow: hidden;
}

.dark .news-card-skeleton {
    background: #1e293b;
}

.skeleton-image {
    background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
    background-size: 200% 100%;
    animation: loading 1.5s infinite;
}

.dark .skeleton-image {
    background: linear-gradient(90deg, #2d3748 25%, #4a5568 50%, #2d3748 75%);
    background-size: 200% 100%;
}

.skeleton-line {
    background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
    background-size: 200% 100%;
    animation: loading 1.5s infinite;
    border-radius: 0.25rem;
}

.dark .skeleton-line {
    background: linear-gradient(90deg, #2d3748 25%, #4a5568 50%, #2d3748 75%);
    background-size: 200% 100%;
}

@keyframes loading {
    0% {
        background-position: 200% 0;
    }
    100% {
        background-position: -200% 0;
    }
}

/* Initial Loader Styles */
.news-loader-spinner {
    width: 64px;
    height: 64px;
}

.news-loader-circle {
    width: 100%;
    height: 100%;
    transform-origin: center;
    animation: news-loader-rotate 2s linear infinite;
}

.news-loader-path {
    stroke: #22c55e;
    stroke-linecap: round;
    animation: news-loader-dash 1.5s ease-in-out infinite;
}

.dark .news-loader-path {
    stroke: #4ade80;
}

@keyframes news-loader-rotate {
    100% {
        transform: rotate(360deg);
    }
}

@keyframes news-loader-dash {
    0% {
        stroke-dasharray: 1, 200;
        stroke-dashoffset: 0;
    }
    50% {
        stroke-dasharray: 90, 200;
        stroke-dashoffset: -35px;
    }
    100% {
        stroke-dasharray: 90, 200;
        stroke-dashoffset: -124px;
    }
}

.news-loader-dot {
    animation: news-loader-bounce 1.4s ease-in-out infinite;
}

@keyframes news-loader-bounce {
    0%,
    80%,
    100% {
        transform: scale(0.6);
        opacity: 0.4;
    }
    40% {
        transform: scale(1);
        opacity: 1;
    }
}
/* Delete Confirmation Modal Styles */
.delete-modal {
    width: min(calc(100vw - 2rem), 520px) !important;
    max-width: 520px !important;
    aspect-ratio: 16 / 9;
    min-height: fit-content;
    border-radius: 18px;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.25);
}

.delete-modal :deep([data-slot="body"]) {
    padding: 0;
}

.delete-modal :deep([data-slot="footer"]) {
    padding: 0;
}

.delete-modal-content {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    width: 100%;
    height: 100%;
    padding: 1.75rem 2rem;
}

.delete-modal-icon-wrapper {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    width: 64px;
    height: 64px;
    border-radius: 50%;
    background: radial-gradient(circle at 30% 25%, #fef2f2 0%, #fee2e2 55%, #fecaca 100%);
    border: 1px solid rgba(220, 38, 38, 0.12);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.6), 0 10px 24px -8px rgba(220, 38, 38, 0.4);
}

.dark .delete-modal-icon-wrapper {
    background: radial-gradient(circle at 30% 25%, #450a0a 0%, #7f1d1d 55%, #991b1b 100%);
    border-color: rgba(248, 113, 113, 0.18);
    box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.08), 0 10px 24px -8px rgba(248, 113, 113, 0.35);
}

.delete-modal-icon {
    width: 30px;
    height: 30px;
    color: #dc2626;
    animation: delete-modal-icon-pulse 2.6s ease-in-out infinite;
}

.dark .delete-modal-icon {
    color: #fca5a5;
}

@keyframes delete-modal-icon-pulse {
    0%,
    100% {
        transform: scale(1);
    }
    50% {
        transform: scale(1.07);
    }
}

.delete-modal-main {
    flex: 1;
    min-width: 0;
}

.delete-modal-title {
    font-size: 1.25rem;
    font-weight: 700;
    letter-spacing: -0.01em;
    color: #1e293b;
    margin: 0 0 0.35rem;
}

.dark .delete-modal-title {
    color: #f1f5f9;
}

.delete-modal-text {
    font-size: 0.95rem;
    color: #475569;
    line-height: 1.6;
    margin: 0 0 0.75rem;
}

.dark .delete-modal-text {
    color: #94a3b8;
}

.delete-modal-highlight {
    display: inline-block;
    font-weight: 600;
    color: #1e293b;
    background: #f1f5f9;
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    padding: 0.05rem 0.5rem;
    margin: 0 0.2rem;
    max-width: 100%;
}

.dark .delete-modal-highlight {
    color: #e2e8f0;
    background: #1e293b;
    border-color: #334155;
}

.delete-modal-warning {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.85rem;
    font-weight: 500;
    color: #ef4444;
    background: #fef2f2;
    border: 1px solid #fecaca;
    border-radius: 999px;
    padding: 0.35rem 0.9rem;
    margin: 0;
}

.dark .delete-modal-warning {
    color: #fca5a5;
    background: rgba(127, 29, 29, 0.35);
    border-color: rgba(239, 68, 68, 0.3);
}

.delete-modal-warning-icon {
    width: 0.95rem;
    height: 0.95rem;
    flex-shrink: 0;
}

.delete-modal-footer {
    display: flex;
    align-items: stretch;
    width: 100%;
    padding: 0;
}

.delete-modal-footer-divider {
    width: 1px;
    flex-shrink: 0;
    background: #e2e8f0;
}

.dark .delete-modal-footer-divider {
    background: #334155;
}

.delete-modal-footer .delete-modal-btn {
    flex: 1 1 0;
    border-radius: 0;
    justify-content: center;
}

@media (max-width: 480px) {
    .delete-modal-content {
        flex-direction: column;
        text-align: center;
        gap: 1rem;
        padding: 1.5rem;
    }

    .delete-modal-icon-wrapper {
        width: 56px;
        height: 56px;
    }
}
</style>
