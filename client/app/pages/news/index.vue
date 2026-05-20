<script lang="ts">
import type { NewsDescDto } from '../../types';
import axios from 'axios';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { ref, onMounted, onUnmounted } from 'vue';

export default {
    setup() {
        const { locale } = useI18n();
        const router = useRouter();
        const news = ref<NewsDescDto[]>([]);
        const page = ref(0);
        const limit = ref(9); // Load 9 news items at a time as requested
        const server = ref(import.meta.env.VITE_SERVER as string);
        const loading = ref(false);
        const hasMore = ref(true);
        const observer = ref<IntersectionObserver | null>(null);
        const observerElement = ref<HTMLElement | null>(null);

        const loadNews = async (initial = false) => {
            if (loading.value || (!hasMore.value && !initial)) return;

            try {
                loading.value = true;
                const response = await axios.get(server.value + 'news', {
                    params: {
                        offset: page.value * limit.value,
                        limit: limit.value,
                    },
                });

                const newNews = response.data || [];

                if (initial) {
                    news.value = newNews;
                } else {
                    news.value = [...news.value, ...newNews];
                }

                // Check if we got fewer items than requested
                if (newNews.length < limit.value) {
                    hasMore.value = false;
                }

                if (!initial) {
                    page.value += 1;
                }
            } catch (e) {
                console.error('Error fetching news', e);
            } finally {
                loading.value = false;
            }
        };

        const loadMoreNews = async () => {
            if (!hasMore.value || loading.value) return;
            await loadNews(false);
        };

        const formatDate = (dateString: string) => {
            const date = new Date(dateString);
            const options: Intl.DateTimeFormatOptions = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            };
            return date.toLocaleDateString(
                locale.value === 'ru' ? 'ru-RU' : 'en-US',
                options
            );
        };

        const getNewsTitle = (newsItem: NewsDescDto) => {
            return locale.value === 'ru' ? newsItem.title_ru : newsItem.title_en;
        };

        const getNewsSubtitle = (newsItem: NewsDescDto) => {
            return locale.value === 'ru' ? newsItem.subTitle_ru : newsItem.subTitle_en;
        };

        const getImageUrl = (newsItem: NewsDescDto) => {
            return newsItem.dir + newsItem.img_back;
        };

        const navigateToNews = (id: string) => {
            router.push(`/news/${id}`);
        };

        onMounted(async () => {
            // Load initial news
            await loadNews(true);

            // Setup intersection observer for infinite scroll
            const options = {
                root: null,
                rootMargin: '100px',
                threshold: 0.1,
            };

            observer.value = new IntersectionObserver((entries) => {
                const entry = entries[0];
                if (entry?.isIntersecting && hasMore.value && !loading.value) {
                    loadMoreNews();
                }
            }, options);

            // Observe the sentinel element
            if (observerElement.value) {
                observer.value.observe(observerElement.value);
            }
        });

        onUnmounted(() => {
            if (observer.value) {
                observer.value.disconnect();
            }
        });

        return {
            news,
            loading,
            hasMore,
            observerElement,
            formatDate,
            getNewsTitle,
            getNewsSubtitle,
            getImageUrl,
            navigateToNews,
            locale,
        };
    },
};
</script>

<template>
    <div class="news-page">
        <div class="container mx-auto px-4 py-8">
            <!-- Page Header -->
            <div class="mb-12 text-center">
                <h1
                    class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4"
                >
                    {{ $t('news.title') }}
                </h1>
                <p class="text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
                    {{ $t('news.subtitle') }}
                </p>
            </div>

            <!-- News Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                <div
                    v-for="newsItem in news"
                    :key="newsItem.id"
                    class="news-card group cursor-pointer"
                    @click="navigateToNews(newsItem.id)"
                >
                    <!-- Image Container -->
                    <div class="news-image-container overflow-hidden rounded-t-2xl">
                        <img
                            :src="getImageUrl(newsItem)"
                            :alt="getNewsTitle(newsItem)"
                            class="news-image w-full h-64 object-cover group-hover:scale-105 transition-transform duration-500"
                            loading="lazy"
                        />
                        <div class="news-image-overlay"></div>
                    </div>

                    <!-- Content Container -->
                    <div class="news-content p-6">
                        <!-- Date -->
                        <div
                            class="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-3"
                        >
                            <svg
                                class="w-4 h-4 mr-2"
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
                                ></path>
                            </svg>
                            {{ formatDate(newsItem.datetime) }}
                        </div>

                        <!-- Title -->
                        <h3
                            class="text-xl font-bold text-gray-900 dark:text-white mb-3 line-clamp-2 group-hover:text-green-600 dark:group-hover:text-green-400 transition-colors"
                        >
                            {{ getNewsTitle(newsItem) }}
                        </h3>

                        <!-- Subtitle/Description -->
                        <p class="text-gray-600 dark:text-gray-300 mb-4 line-clamp-3">
                            {{ getNewsSubtitle(newsItem) }}
                        </p>

                        <!-- Read More Link -->
                        <div
                            class="flex items-center text-green-600 dark:text-green-400 font-medium"
                        >
                            <span class="mr-2">{{ $t('news.readMore') }}</span>
                            <svg
                                class="w-4 h-4 transform group-hover:translate-x-1 transition-transform"
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
                                ></path>
                            </svg>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="mt-12">
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    <div v-for="n in 3" :key="n" class="news-card-skeleton">
                        <div class="skeleton-image h-64 rounded-t-2xl"></div>
                        <div class="p-6">
                            <div class="skeleton-line h-4 w-24 mb-4"></div>
                            <div class="skeleton-line h-6 w-full mb-3"></div>
                            <div class="skeleton-line h-4 w-full mb-2"></div>
                            <div class="skeleton-line h-4 w-3/4"></div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- No More News Message -->
            <div v-if="!hasMore && news.length > 0" class="mt-12 text-center py-8">
                <div
                    class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 mb-4"
                >
                    <svg
                        class="w-8 h-8 text-green-600 dark:text-green-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 13l4 4L19 7"
                        ></path>
                    </svg>
                </div>
                <p class="text-lg text-gray-600 dark:text-gray-300">
                    {{ $t('news.allLoaded') }}
                </p>
            </div>

            <!-- No News Message -->
            <div v-if="!loading && news.length === 0" class="mt-12 text-center py-12">
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-gray-100 dark:bg-gray-800 mb-6"
                >
                    <svg
                        class="w-10 h-10 text-gray-400"
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
                        ></path>
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                    {{ $t('news.noNews') }}
                </h3>
                <p class="text-gray-600 dark:text-gray-300 max-w-md mx-auto">
                    {{ $t('news.noNewsDescription') }}
                </p>
            </div>

            <!-- Intersection Observer Sentinel -->
            <div
                ref="observerElement"
                class="h-1 w-full"
                :class="{ 'opacity-0': !hasMore || loading }"
            ></div>
        </div>
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
    background: white;
    border-radius: 1rem;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    transition: all 0.3s ease;
    overflow: hidden;
    height: 100%;
    display: flex;
    flex-direction: column;
}

.dark .news-card {
    background: #1e293b;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -1px rgba(0, 0, 0, 0.2);
}

.news-card:hover {
    transform: translateY(-8px);
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.dark .news-card:hover {
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.4), 0 10px 10px -5px rgba(0, 0, 0, 0.3);
}

.news-image-container {
    position: relative;
    height: 16rem;
}

.news-image {
    height: 100%;
    width: 100%;
}

.news-image-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(to bottom, transparent 50%, rgba(0, 0, 0, 0.7));
    opacity: 0;
    transition: opacity 0.3s ease;
}

.news-card:hover .news-image-overlay {
    opacity: 1;
}

.news-content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
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
</style>
