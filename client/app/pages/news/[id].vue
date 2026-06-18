<script setup lang="ts">
import { ref, onMounted, computed, } from 'vue';
import { useI18n, } from '#imports';
import { getHttpClient, } from '~/api/http-client';
import { useFormatting, } from '~/composables/useFormatting';
import type { NewsItem, NewsListResponse, } from '~/api/news';
import SideNewsTrailer from '~/components/SideNewsTrailer.vue';
import NewsDescriptionSkeleton from '~/components/NewsDescriptionSkeleton.vue';

const route = useRoute();
const router = useRouter();
const { t, } = useI18n();
const { formatDate, } = useFormatting();

const currentNewsItem = ref<NewsItem | null>(null,);
const otherNews = ref<NewsItem[]>([],);
const loading = ref(true,);
const error = ref<string | null>(null,);
const mainImageError = ref(false,);

const newsId = computed(() => Number(route.params.id),);

const fetchNewsDetail = async (): Promise<void> => {
    try {
        loading.value = true;
        error.value = null;

        const { data: currentData, } = await getHttpClient().get<NewsItem>(
            `news/${newsId.value}`,
            );
        currentNewsItem.value = currentData;

        // Fetch other news for sidebar (excluding current)
        const { data: otherData, } = await getHttpClient().get<NewsListResponse>('news', {
            params: {
                page: 1,
                limit: 6, // Get 6 to potentially exclude current
            },
        });

        // Filter out current news and take first 5
        otherNews.value = (otherData.news || [])
            .filter((news: NewsItem,) => news.id !== newsId.value,)
            .slice(0, 5,);
    } catch (err) {
        console.error('Error fetching news detail:', err,);
        error.value = t('news.detail.error.loadFailed',);
    } finally {
        loading.value = false;
    }
};

const navigateToNews = (id: number,): void => {
    router.push(`/news/${id}`,);
};

const handleMainImageError = (): void => {
    mainImageError.value = true;
};
</script>

<template>
    <div class="news-detail-page">
        <!-- Loading State -->
        <NewsDescriptionSkeleton v-if="loading && !currentNewsItem" />

        <!-- Error State -->
        <div
            v-else-if="error"
            class="error-container"
        >
            <div class="error-content">
                <svg
                    class="error-icon"
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
                <h2 class="error-title">
                    {{ $t('news.detail.error.title',) }}
                </h2>
                <p class="error-message">
                    {{ error }}
                </p>
                <button
                    class="retry-button"
                    @click="fetchNewsDetail"
                >
                    {{ $t('news.detail.error.retryButton',) }}
                </button>
            </div>
        </div>

        <!-- Main Content -->
        <div
            v-else-if="currentNewsItem"
            class="container mx-auto px-4 py-8"
        >
            <!-- Breadcrumb -->
            <nav class="mb-6">
                <ol class="flex items-center space-x-2 text-sm text-gray-500">
                    <li>
                        <router-link
                            to="/"
                            class="hover:text-green-600 transition-colors"
                        >
                            {{ $t('breadcrumb.home',) }}
                        </router-link>
                    </li>
                    <li class="flex items-center">
                        <svg
                            class="w-4 h-4 mx-1"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M9 5l7 7-7 7"
                            />
                        </svg>
                        <router-link
                            to="/news"
                            class="hover:text-green-600 transition-colors"
                        >
                            {{ $t('breadcrumb.news',) }}
                        </router-link>
                    </li>
                    <li class="flex items-center">
                        <svg
                            class="w-4 h-4 mx-1"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M9 5l7 7-7 7"
                            />
                        </svg>
                        <span
                            class="text-gray-900 dark:text-white font-medium truncate max-w-xs"
                        >
                            {{ currentNewsItem.title }}
                        </span>
                    </li>
                </ol>
            </nav>

            <!-- Top section with main image and sidebar -->
            <div class="flex flex-col lg:flex-row gap-8 mb-8">
                <!-- Main Image (2/3 width on large screens) -->
                <div class="lg:w-2/3">
                    <div
                        class="main-image-container rounded-2xl overflow-hidden shadow-lg min-h-75 md:min-h-100"
                    >
                        <div
                            v-if="mainImageError || !currentNewsItem.image_path"
                            class="image-error-state"
                        >
                            <svg
                                class="w-16 h-16 text-gray-400 mx-auto mb-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                />
                            </svg>
                            <p class="text-gray-500 dark:text-gray-400 text-center">
                                {{ $t('news.detail.image.error',) }}
                            </p>
                        </div>
                        <img
                            v-else
                            :src="currentNewsItem.image_path"
                            :alt="currentNewsItem.title"
                            class="w-full h-auto max-h-150 object-cover"
                            loading="eager"
                            @error="handleMainImageError"
                        >
                    </div>
                </div>

                <!-- Sidebar (1/3 width on large screens) -->
                <div class="lg:w-1/3 flex">
                    <div class="sidebar-container flex-1">
                        <h3
                            class="text-xl font-bold text-gray-900 dark:text-white mb-4 pb-1 border-b border-gray-200 dark:border-gray-700"
                        >
                            {{ $t('news.detail.sidebar.title',) }}
                        </h3>

                        <div
                            v-if="otherNews.length > 0"
                            class="space-y-2 overflow-y-auto sidebar-news-list"
                        >
                            <div
                                v-for="news in otherNews"
                                :key="news.id"
                                class="sidebar-news-item cursor-pointer group"
                                @click="navigateToNews(news.id,)"
                            >
                                <SideNewsTrailer :side-news-object="news" />
                            </div>
                        </div>

                        <div
                            v-else
                            class="text-center py-8"
                        >
                            <svg
                                class="w-12 h-12 mx-auto text-gray-400 mb-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                                />
                            </svg>
                            <p class="text-gray-500 dark:text-gray-400">
                                {{ $t('news.detail.sidebar.empty',) }}
                            </p>
                        </div>

                        <div class="mt-8">
                            <router-link
                                to="/news"
                                class="flex items-center justify-center w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white font-medium rounded-lg transition-colors"
                            >
                                <span>{{ $t('news.viewAllNews',) }}</span>
                                <svg
                                    class="w-5 h-5 ml-2"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M14 5l7 7m0 0l-7 7m7-7H3"
                                    />
                                </svg>
                            </router-link>
                        </div>
                    </div>
                </div>
            </div>

            <!-- News Information -->
            <div class="news-info mb-8">
                <!-- Date -->
                <div
                    class="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-4"
                >
                    <svg
                        class="w-5 h-5 mr-2"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                        />
                    </svg>
                    {{ formatDate(currentNewsItem.created_at,) }}
                </div>

                <!-- Title -->
                <h1
                    class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-4"
                >
                    {{ currentNewsItem.title }}
                </h1>
            </div>

            <!-- News Description -->
            <div
                v-if="currentNewsItem.description"
                class="news-description mb-6"
            >
                <p class="text-lg text-gray-600 dark:text-gray-400 leading-relaxed">
                    {{ currentNewsItem.description }}
                </p>
            </div>

            <!-- News Content -->
            <div
                v-if="currentNewsItem.content"
                class="news-text mb-12"
            >
                <div class="prose prose-lg dark:prose-invert max-w-none">
                    <div
                        class="whitespace-pre-line text-gray-700 dark:text-gray-300 leading-relaxed text-justify"
                        v-html="currentNewsItem.content"
                    />
                </div>
            </div>

            <!-- Video Section -->
            <div
                v-if="currentNewsItem.video_path"
                class="video-section mb-12"
            >
                <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                    {{ $t('news.detail.videos.title',) }}
                </h3>
                <div class="relative aspect-video rounded-xl overflow-hidden shadow-lg bg-black">
                    <video
                        :src="currentNewsItem.video_path"
                        class="w-full h-full object-contain"
                        controls
                        preload="metadata"
                    >
                        <p class="text-white text-center py-8">
                            {{ $t('news.detail.video.notSupported',) }}
                        </p>
                    </video>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.news-detail-page {
    min-height: 100vh;
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

.dark .news-detail-page {
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

.error-container {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 60vh;
    padding: 2rem;
}

.error-content {
    text-align: center;
    max-width: 400px;
}

.error-icon {
    width: 64px;
    height: 64px;
    color: #ef4444;
    margin: 0 auto 1.5rem;
}

.error-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1f2937;
    margin-bottom: 0.5rem;
}

.dark .error-title {
    color: #f9fafb;
}

.error-message {
    color: #6b7280;
    margin-bottom: 1.5rem;
}

.dark .error-message {
    color: #d1d5db;
}

.retry-button {
    background-color: #10b981;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    font-weight: 500;
    transition: background-color 0.2s;
}

.retry-button:hover {
    background-color: #059669;
}

.main-image-container {
    position: relative;
    background: linear-gradient(135deg, #e5e7eb 0%, #d1d5db 100%);
}

.dark .main-image-container {
    background: linear-gradient(135deg, #374151 0%, #4b5563 100%);
}

/* Sidebar container matching main image height */
.sidebar-container {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.sidebar-news-list {
    flex: 1;
    overflow-y: auto;
    max-height: calc(100% - 120px);
    padding-right: 0.5rem;
}

/* Custom scrollbar for sidebar */
.sidebar-news-list::-webkit-scrollbar {
    width: 6px;
}

.sidebar-news-list::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 3px;
}

.dark .sidebar-news-list::-webkit-scrollbar-track {
    background: #374151;
}

.sidebar-news-list::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 3px;
}

.dark .sidebar-news-list::-webkit-scrollbar-thumb {
    background: #4b5563;
}

.sidebar-news-list::-webkit-scrollbar-thumb:hover {
    background: #a1a1a1;
}

.dark .sidebar-news-list::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
}

.sidebar-news-item {
    transition: transform 0.2s ease;
}

.sidebar-news-item:hover {
    transform: translateX(4px);
}

/* Image Error States */
.image-error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 300px;
    padding: 2rem;
    background: linear-gradient(135deg, #e5e7eb 0%, #d1d5db 100%);
    border-radius: 1rem;
}

.dark .image-error-state {
    background: linear-gradient(135deg, #374151 0%, #4b5563 100%);
}

.prose {
    color: inherit;
}

.prose p {
    margin-bottom: 1.5rem;
}

.prose p:last-child {
    margin-bottom: 0;
}

/* Responsive layout adjustments */
@media (max-width: 1024px) {
    .flex-col.lg\:flex-row {
        flex-direction: column;
    }

    .lg\:w-2\/3,
    .lg\:w-1\/3 {
        width: 100%;
    }

    .sidebar-container {
        min-height: auto;
        margin-top: 2rem;
    }

    .sidebar-news-list {
        max-height: 400px;
    }
}
</style>
