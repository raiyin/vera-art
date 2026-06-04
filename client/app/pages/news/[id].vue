<script lang="ts">
import { inject, ref, onMounted, computed, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import axios from 'axios';
import type { NewsDesc } from '../../types';
import SideNewsTrailer from '@/components/SideNewsTrailer.vue';
import VideoSection from '@/components/VideoSection.vue';
import NewsDescriptionSkeleton from '@/components/NewsDescriptionSkeleton.vue';

export default {
    components: {
        SideNewsTrailer,
        VideoSection,
        NewsDescriptionSkeleton,
    },
    setup() {
        const config = useRuntimeConfig();
        const SERVER_URL = config.public.serverUrl;

        const route = useRoute();
        const router = useRouter();
        const { locale, t } = useI18n();

        const currentNewsItem = ref<NewsDesc | null>(null);
        const otherNews = ref<NewsDesc[]>([]);
        const loading = ref(true);
        const error = ref<string | null>(null);
        const selectedImageIndex = ref(0);
        const showImageModal = ref(false);
        const imageLoadErrors = ref<Set<number>>(new Set());
        const mainImageError = ref(false);

        const newsId = computed(() => route.params.id as string);

        const getImageUrl = (path: string | undefined, dir: string) => {
            if (!path) return '';
            return dir + path;
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

        const fetchNewsDetail = async () => {
            try {
                loading.value = true;
                error.value = null;

                const response = await axios.get(`${SERVER_URL}news/${newsId.value}`);
                currentNewsItem.value = response.data;

                // Fetch other news for sidebar (excluding current)
                const otherResponse = await axios.get(`${SERVER_URL}news`, {
                    params: {
                        offset: 0,
                        limit: 6, // Get 6 to potentially exclude current
                    },
                });

                // Filter out current news and take first 5
                otherNews.value = otherResponse.data
                    .filter((news: NewsDesc) => news.id.toString() !== newsId.value)
                    .slice(0, 5);
            } catch (err) {
                console.error('Error fetching news detail:', err);
                error.value = t('news.detail.error.loadFailed');
            } finally {
                loading.value = false;
            }
        };

        const openImageModal = (index: number) => {
            selectedImageIndex.value = index;
            showImageModal.value = true;
        };

        const closeImageModal = () => {
            showImageModal.value = false;
        };

        const navigateToNews = (id: string) => {
            router.push(`/news/${id}`);
        };

        const handleMainImageError = () => {
            mainImageError.value = true;
        };

        const handleGalleryImageError = (index: number) => {
            imageLoadErrors.value.add(index);
        };

        const handleModalImageError = () => {
            // If modal image fails, we could show a placeholder
            // For now, just log
            console.error('Modal image failed to load');
        };

        const resetImageErrors = () => {
            imageLoadErrors.value.clear();
            mainImageError.value = false;
        };

        // Keyboard navigation for modal
        const handleKeydown = (event: KeyboardEvent) => {
            if (!showImageModal.value || !currentNewsItem.value) return;

            switch (event.key) {
                case 'Escape':
                    closeImageModal();
                    break;
                case 'ArrowLeft':
                    if (selectedImageIndex.value > 0) {
                        selectedImageIndex.value--;
                    }
                    break;
                case 'ArrowRight':
                    if (
                        currentNewsItem.value.images &&
                        selectedImageIndex.value < currentNewsItem.value.images.length - 1
                    ) {
                        selectedImageIndex.value++;
                    }
                    break;
            }
        };

        onMounted(() => {
            fetchNewsDetail();
            window.addEventListener('keydown', handleKeydown);
        });

        onUnmounted(() => {
            window.removeEventListener('keydown', handleKeydown);
        });

        return {
            currentNewsItem,
            otherNews,
            loading,
            error,
            selectedImageIndex,
            showImageModal,
            imageLoadErrors,
            mainImageError,
            newsId,
            getImageUrl,
            formatDate,
            openImageModal,
            closeImageModal,
            navigateToNews,
            fetchNewsDetail,
            handleMainImageError,
            handleGalleryImageError,
            handleModalImageError,
            resetImageErrors,
            locale,
        };
    },
};
</script>

<template>
    <div class="news-detail-page">
        <!-- Loading State -->
        <NewsDescriptionSkeleton v-if="loading && !currentNewsItem" />

        <!-- Error State -->
        <div v-else-if="error" class="error-container">
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
                    ></path>
                </svg>
                <h2 class="error-title">{{ $t('news.detail.error.title') }}</h2>
                <p class="error-message">{{ error }}</p>
                <button @click="fetchNewsDetail" class="retry-button">
                    {{ $t('news.detail.error.retryButton') }}
                </button>
            </div>
        </div>

        <!-- Main Content -->
        <div v-else-if="currentNewsItem" class="container mx-auto px-4 py-8">
            <!-- Breadcrumb -->
            <nav class="mb-6">
                <ol class="flex items-center space-x-2 text-sm text-gray-500">
                    <li>
                        <router-link
                            to="/"
                            class="hover:text-green-600 transition-colors"
                        >
                            {{ $t('breadcrumb.home') }}
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
                            ></path>
                        </svg>
                        <router-link
                            to="/news"
                            class="hover:text-green-600 transition-colors"
                        >
                            {{ $t('breadcrumb.news') }}
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
                            ></path>
                        </svg>
                        <span
                            class="text-gray-900 dark:text-white font-medium truncate max-w-xs"
                        >
                            {{
                                locale === 'ru'
                                    ? currentNewsItem.title_ru
                                    : currentNewsItem.title_en
                            }}
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
                        ref="mainImageContainer"
                    >
                        <div v-if="mainImageError" class="image-error-state">
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
                                ></path>
                            </svg>
                            <p class="text-gray-500 dark:text-gray-400 text-center">
                                {{ $t('news.detail.image.error') }}
                            </p>
                        </div>
                        <img
                            v-else
                            :src="
                                getImageUrl(
                                    currentNewsItem.img_backfull,
                                    currentNewsItem.dir
                                )
                            "
                            :alt="
                                locale === 'ru'
                                    ? currentNewsItem.title_ru
                                    : currentNewsItem.title_en
                            "
                            class="w-full h-auto max-h-150 object-cover"
                            loading="eager"
                            @error="handleMainImageError"
                        />
                    </div>
                </div>

                <!-- Sidebar (1/3 width on large screens) -->
                <div class="lg:w-1/3 flex">
                    <div class="sidebar-container flex-1">
                        <h3
                            class="text-xl font-bold text-gray-900 dark:text-white mb-4 pb-1 border-b border-gray-200 dark:border-gray-700"
                        >
                            {{ $t('news.detail.sidebar.title') }}
                        </h3>

                        <div
                            v-if="otherNews.length > 0"
                            class="space-y-2 overflow-y-auto sidebar-news-list"
                        >
                            <div
                                v-for="news in otherNews"
                                :key="news.id"
                                class="sidebar-news-item cursor-pointer group"
                                @click="navigateToNews(news.id)"
                            >
                                <SideNewsTrailer :sideNewsObject="news" />
                            </div>
                        </div>

                        <div v-else class="text-center py-8">
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
                                ></path>
                            </svg>
                            <p class="text-gray-500 dark:text-gray-400">
                                {{ $t('news.detail.sidebar.empty') }}
                            </p>
                        </div>

                        <div class="mt-8">
                            <router-link
                                to="/news"
                                class="flex items-center justify-center w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white font-medium rounded-lg transition-colors"
                            >
                                <span>{{ $t('news.viewAllNews') }}</span>
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
                                    ></path>
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
                        ></path>
                    </svg>
                    {{ formatDate(currentNewsItem.datetime) }}
                </div>

                <!-- Title -->
                <h1
                    class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-4"
                >
                    {{
                        locale === 'ru'
                            ? currentNewsItem.title_ru
                            : currentNewsItem.title_en
                    }}
                </h1>
            </div>

            <!-- News Text -->
            <div class="news-text mb-12">
                <div class="prose prose-lg dark:prose-invert max-w-none">
                    <div
                        v-html="
                            locale === 'ru'
                                ? currentNewsItem.text_ru
                                : currentNewsItem.text_en
                        "
                        class="whitespace-pre-line text-gray-700 dark:text-gray-300 leading-relaxed text-justify"
                    ></div>
                </div>
            </div>

            <!-- Image Gallery -->
            <div class="image-gallery mb-12">
                <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                    {{ $t('news.detail.gallery.title') }}
                </h3>
                <div
                    v-if="currentNewsItem.images && currentNewsItem.images.length > 0"
                    class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4"
                >
                    <div
                        v-for="(image, index) in currentNewsItem.images"
                        :key="index"
                        class="gallery-item cursor-pointer group"
                        @click="openImageModal(index)"
                    >
                        <div class="aspect-square overflow-hidden rounded-lg relative">
                            <div
                                v-if="imageLoadErrors.has(index)"
                                class="image-error-placeholder w-full h-full flex items-center justify-center bg-gray-100 dark:bg-gray-800"
                            >
                                <svg
                                    class="w-8 h-8 text-gray-400"
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
                                    ></path>
                                </svg>
                            </div>
                            <img
                                v-else
                                :src="getImageUrl(image, currentNewsItem.dir)"
                                :alt="`Gallery image ${index + 1}`"
                                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                                loading="lazy"
                                @error="() => handleGalleryImageError(index)"
                            />
                        </div>
                        <div class="gallery-overlay" v-if="!imageLoadErrors.has(index)">
                            <svg
                                class="w-8 h-8 text-white"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"
                                ></path>
                            </svg>
                        </div>
                    </div>
                </div>
                <div v-else class="empty-state py-12 text-center">
                    <svg
                        class="w-16 h-16 mx-auto text-gray-400 mb-4"
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
                        ></path>
                    </svg>
                    <p class="text-gray-500 dark:text-gray-400">
                        {{ $t('news.detail.gallery.empty') }}
                    </p>
                </div>
            </div>

            <!-- Video Carousel -->
            <div class="video-carousel mb-12">
                <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                    {{ $t('news.detail.videos.title') }}
                </h3>
                <div
                    v-if="currentNewsItem.videos && currentNewsItem.videos.length > 0"
                    class="relative"
                >
                    <VideoSection :currentNewsItem="currentNewsItem" />
                </div>
                <div v-else class="empty-state py-12 text-center">
                    <svg
                        class="w-16 h-16 mx-auto text-gray-400 mb-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
                        ></path>
                    </svg>
                    <p class="text-gray-500 dark:text-gray-400">
                        {{ $t('news.detail.videos.empty') }}
                    </p>
                </div>
            </div>
        </div>

        <!-- Custom Image Modal -->
        <transition name="modal-fade">
            <div
                v-if="showImageModal && currentNewsItem"
                class="custom-modal-overlay"
                @click.self="closeImageModal"
            >
                <div class="custom-modal-container">
                    <!-- Close Button -->
                    <button
                        class="custom-modal-close"
                        @click="closeImageModal"
                        :aria-label="$t('news.detail.modal.close')"
                    >
                        <svg
                            class="w-6 h-6"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"
                            ></path>
                        </svg>
                    </button>

                    <!-- Image Container -->
                    <div class="custom-modal-image-container">
                        <transition name="image-slide" mode="out-in">
                            <img
                                :key="selectedImageIndex"
                                :src="
                                    getImageUrl(
                                        currentNewsItem.images[selectedImageIndex],
                                        currentNewsItem.dir
                                    )
                                "
                                :alt="`Gallery image ${selectedImageIndex + 1}`"
                                class="custom-modal-image"
                                @click="closeImageModal"
                            />
                        </transition>

                        <!-- Navigation Buttons -->
                        <button
                            v-if="selectedImageIndex > 0"
                            @click="selectedImageIndex--"
                            class="custom-modal-nav-button left"
                            :aria-label="$t('news.detail.modal.previous')"
                        >
                            <svg
                                class="w-8 h-8"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 19l-7-7 7-7"
                                ></path>
                            </svg>
                        </button>
                        <button
                            v-if="selectedImageIndex < currentNewsItem.images.length - 1"
                            @click="selectedImageIndex++"
                            class="custom-modal-nav-button right"
                            :aria-label="$t('news.detail.modal.next')"
                        >
                            <svg
                                class="w-8 h-8"
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
                                ></path>
                            </svg>
                        </button>

                        <!-- Image Counter -->
                        <div class="custom-modal-counter">
                            <span class="counter-current">{{
                                selectedImageIndex + 1
                            }}</span>
                            <span class="counter-separator">/</span>
                            <span class="counter-total">{{
                                currentNewsItem.images.length
                            }}</span>
                        </div>
                    </div>

                    <!-- Thumbnail Strip -->
                    <div
                        v-if="currentNewsItem.images.length > 1"
                        class="custom-modal-thumbnails"
                    >
                        <div
                            v-for="(image, index) in currentNewsItem.images"
                            :key="index"
                            class="thumbnail-item"
                            :class="{ active: index === selectedImageIndex }"
                            @click="selectedImageIndex = index"
                        >
                            <img
                                :src="getImageUrl(image, currentNewsItem.dir)"
                                :alt="`Thumbnail ${index + 1}`"
                                class="thumbnail-image"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </transition>
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
    max-height: calc(100% - 120px); /* Account for title and button */
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

.gallery-item {
    position: relative;
    overflow: hidden;
    border-radius: 0.75rem;
}

.gallery-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
}

.gallery-item:hover .gallery-overlay {
    opacity: 1;
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

.image-error-placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f3f4f6;
}

.dark .image-error-placeholder {
    background: #374151;
}

/* Custom Modal Styles */
.modal-fade-enter-active,
.modal-fade-leave-active {
    transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
    opacity: 0;
}

.image-slide-enter-active,
.image-slide-leave-active {
    transition: transform 0.3s ease, opacity 0.3s ease;
}

.image-slide-enter-from {
    transform: translateX(30px);
    opacity: 0;
}

.image-slide-leave-to {
    transform: translateX(-30px);
    opacity: 0;
}

.custom-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    padding: 1rem;
}

.custom-modal-container {
    position: relative;
    background: var(--color-surface);
    border-radius: 1.5rem;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    max-width: 90vw;
    max-height: 90vh;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.dark .custom-modal-container {
    background: #1e293b;
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.custom-modal-close {
    position: absolute;
    top: 1rem;
    right: 1rem;
    z-index: 10;
    background: rgba(0, 0, 0, 0.5);
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: white;
    transition: all 0.2s ease;
}

.custom-modal-close:hover {
    background: rgba(0, 0, 0, 0.8);
    transform: scale(1.1);
}

.custom-modal-image-container {
    position: relative;
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    min-height: 60vh;
}

.custom-modal-image {
    max-width: 100%;
    max-height: 70vh;
    object-fit: contain;
    border-radius: 0.75rem;
    cursor: zoom-out;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3);
}

.custom-modal-nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(0, 0, 0, 0.6);
    border: none;
    border-radius: 50%;
    width: 56px;
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: white;
    transition: all 0.2s ease;
    z-index: 5;
}

.custom-modal-nav-button:hover {
    background: rgba(0, 0, 0, 0.9);
    transform: translateY(-50%) scale(1.1);
}

.custom-modal-nav-button.left {
    left: 1rem;
}

.custom-modal-nav-button.right {
    right: 1rem;
}

.custom-modal-counter {
    position: absolute;
    top: 1rem;
    left: 1rem;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 2rem;
    font-size: 0.875rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 0.25rem;
}

.counter-current {
    color: #10b981;
}

.counter-separator {
    opacity: 0.7;
}

.counter-total {
    opacity: 0.9;
}

.custom-modal-thumbnails {
    display: flex;
    gap: 0.5rem;
    padding: 1rem;
    background: rgba(0, 0, 0, 0.3);
    overflow-x: auto;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.thumbnail-item {
    flex: 0 0 auto;
    width: 80px;
    height: 80px;
    border-radius: 0.5rem;
    overflow: hidden;
    cursor: pointer;
    border: 3px solid transparent;
    transition: all 0.2s ease;
    opacity: 0.7;
}

.thumbnail-item:hover {
    opacity: 1;
    transform: translateY(-2px);
}

.thumbnail-item.active {
    border-color: #10b981;
    opacity: 1;
    transform: scale(1.05);
}

.thumbnail-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
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

/* Responsive Modal Styles */
@media (max-width: 768px) {
    .custom-modal-container {
        max-width: 95vw;
        max-height: 95vh;
        border-radius: 1rem;
    }

    .custom-modal-image-container {
        padding: 1.5rem;
        min-height: 50vh;
    }

    .custom-modal-image {
        max-height: 60vh;
    }

    .custom-modal-nav-button {
        width: 44px;
        height: 44px;
    }

    .custom-modal-nav-button.left {
        left: 0.5rem;
    }

    .custom-modal-nav-button.right {
        right: 0.5rem;
    }

    .custom-modal-counter {
        top: 0.5rem;
        left: 0.5rem;
        padding: 0.375rem 0.75rem;
        font-size: 0.75rem;
    }

    .custom-modal-thumbnails {
        padding: 0.75rem;
        gap: 0.375rem;
    }

    .thumbnail-item {
        width: 60px;
        height: 60px;
    }

    .custom-modal-close {
        top: 0.5rem;
        right: 0.5rem;
        width: 36px;
        height: 36px;
    }
}

@media (max-width: 480px) {
    .custom-modal-image-container {
        padding: 1rem;
        min-height: 40vh;
    }

    .custom-modal-image {
        max-height: 50vh;
    }

    .thumbnail-item {
        width: 50px;
        height: 50px;
    }

    .custom-modal-nav-button {
        width: 36px;
        height: 36px;
    }

    .custom-modal-nav-button svg {
        width: 20px;
        height: 20px;
    }
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
