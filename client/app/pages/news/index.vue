<script setup lang="ts">
import type { NewsDescDto } from '../../types';
import { useAuthStore } from '../../stores/AuthStore';
import axios from 'axios';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { ref, onMounted, onUnmounted } from 'vue';

const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;
const { locale } = useI18n();
const router = useRouter();
const authStore = useAuthStore();

const news = ref<NewsDescDto[]>([]);
const page = ref(0);
const limit = ref(9);
const loading = ref(false);
const hasMore = ref(true);
const observer = ref<IntersectionObserver | null>(null);
const observerElement = ref<HTMLElement | null>(null);
const deletingId = ref<string | null>(null);

const loadNews = async (initial = false) => {
    if (loading.value || (!hasMore.value && !initial)) return;

    try {
        loading.value = true;
        const response = await axios.get(SERVER_URL + 'news', {
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
    return date.toLocaleDateString(locale.value === 'ru' ? 'ru-RU' : 'en-US', options);
};

const getNewsTitle = (newsItem: NewsDescDto) => {
    return locale.value === 'ru' ? newsItem.title_ru : newsItem.title_en;
};

const getImageUrl = (newsItem: NewsDescDto) => {
    return newsItem.dir + newsItem.img_back;
};

const navigateToNews = (id: string) => {
    router.push(`/news/${id}`);
};

const editNews = (id: string) => {
    router.push(`/news/edit/${id}/`);
};

const deleteNews = async (newsItem: NewsDescDto) => {
    if (!window.confirm('Вы уверены, что хотите удалить эту новость?')) {
        return;
    }

    deletingId.value = newsItem.id;

    try {
        const response = await axios.delete(`${SERVER_URL}news/${newsItem.id}`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
        });

        if (response.status === 200) {
            news.value = news.value.filter((n) => n.id !== newsItem.id);
        }
    } catch (error) {
        console.error('Error deleting news:', error);
        alert('Ошибка при удалении новости');
    } finally {
        deletingId.value = null;
    }
};

onMounted(async () => {
    await loadNews(true);

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

    if (observerElement.value) {
        observer.value.observe(observerElement.value);
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
                    {{ $t('news.title') }}
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
                            <UIcon name="i-heroicons-calendar" class="w-4 h-4 mr-2" />
                            {{ formatDate(newsItem.datetime) }}
                        </div>
                        <!-- Title -->
                        <h3
                            class="text-xl font-bold text-gray-900 dark:text-white line-clamp-2 hover:text-green-600 dark:hover:text-green-400 transition-colors cursor-pointer"
                            @click="navigateToNews(newsItem.id)"
                        >
                            {{ getNewsTitle(newsItem) }}
                        </h3>
                    </template>

                    <!-- Default slot: image + content -->
                    <div class="space-y-4">
                        <!-- Image -->
                        <div
                            class="relative overflow-hidden rounded-lg cursor-pointer"
                            @click="navigateToNews(newsItem.id)"
                        >
                            <img
                                :src="getImageUrl(newsItem)"
                                :alt="getNewsTitle(newsItem)"
                                class="w-full h-48 object-cover hover:scale-105 transition-transform duration-500"
                                loading="lazy"
                            />
                        </div>

                        <!-- Read More -->
                        <div
                            class="flex items-center text-green-600 dark:text-green-400 font-medium cursor-pointer"
                            @click="navigateToNews(newsItem.id)"
                        >
                            <span class="mr-2">{{ $t('news.readMore') }}</span>
                            <UIcon name="i-heroicons-arrow-right" class="w-4 h-4" />
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
                                        @click="editNews(newsItem.id)"
                                    >
                                        Редактировать
                                    </UButton>
                                    <UButton
                                        size="sm"
                                        color="error"
                                        variant="outline"
                                        :loading="deletingId === newsItem.id"
                                        :disabled="deletingId === newsItem.id"
                                        @click="deleteNews(newsItem)"
                                    >
                                        Удалить
                                    </UButton>
                                </div>
                            </div>
                        </ClientOnly>
                    </div>
                </UCard>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="mt-12">
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    <div v-for="n in 3" :key="n" class="news-card-skeleton">
                        <div class="skeleton-image h-48 rounded-t-2xl"></div>
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
</style>
