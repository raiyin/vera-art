<script setup lang="ts">
    import { ref, onMounted, computed, } from 'vue';
    import { useI18n, } from '#imports';
    import { getHttpClient, } from '~/api/http-client';
    import { useFormatting, } from '~/composables/useFormatting';
    import type { NewsItem, NewsListResponse, } from '~/api/news';
    import SideNewsTrailer from '~/components/SideNewsTrailer.vue';
    import SideNewsTrailerSkeleton from '~/components/SideNewsTrailerSkeleton.vue';
    import NewsDescriptionSkeleton from '~/components/NewsDescriptionSkeleton.vue';

    const route = useRoute();
    const router = useRouter();
    const { t, locale, } = useI18n();
    const { formatDate, } = useFormatting();

    const currentNewsItem = ref<NewsItem | null>(null,);
    const otherNews = ref<NewsItem[]>([],);
    const loading = ref(true,);
    const error = ref<string | null>(null,);
    const mainImageError = ref(false,);

    const currentVideoIndex = ref(0,);
    const videoErrors = ref<Set<number>>(new Set(),);

    const selectedGalleryIndex = ref<number | null>(null,);
    const galleryImageErrors = ref<Set<number>>(new Set(),);

    const openGalleryModal = (index: number,): void => {
        selectedGalleryIndex.value = index;
    };

    const closeGalleryModal = (): void => {
        selectedGalleryIndex.value = null;
    };

    const goToPreviousGalleryImage = (): void => {
        if (selectedGalleryIndex.value === null || galleryImages.value.length === 0) return;
        selectedGalleryIndex.value
        = (selectedGalleryIndex.value - 1 + galleryImages.value.length)
                % galleryImages.value.length;
    };

    const goToNextGalleryImage = (): void => {
        if (selectedGalleryIndex.value === null || galleryImages.value.length === 0) return;
        selectedGalleryIndex.value
        = (selectedGalleryIndex.value + 1) % galleryImages.value.length;
    };

    const handleGalleryImageError = (index: number,): void => {
        galleryImageErrors.value.add(index,);
    };

    const newsId = computed(() => String(route.params.id,),);

    const allVideoUrls = computed<string[]>(() => {
        const item = currentNewsItem.value;
        if (!item || !item.videos || item.videos.length === 0) return [];
        return item.videos.map(v => `${item.dir}videos/${v}/${v}.mp4`,);
    });

    const hasMultipleVideos = computed(() => allVideoUrls.value.length > 1,);

    const currentVideoUrl = computed(() => {
        if (allVideoUrls.value.length === 0) return undefined;
        return allVideoUrls.value[currentVideoIndex.value];
    });

    const isCurrentVideoError = computed(() =>
        videoErrors.value.has(currentVideoIndex.value,),
    );

    const resolvedImagePath = computed<string | undefined>(() => {
        const item = currentNewsItem.value;
        if (!item?.main_image) return undefined;
        return item.dir + item.main_image;
    });

    const galleryImages = computed<string[]>(() => {
        const item = currentNewsItem.value;
        if (!item || !item.images || item.images.length === 0) return [];
        return item.images.map(name => item.dir + name,);
    });

    const goToPreviousVideo = (): void => {
        if (allVideoUrls.value.length === 0) return;
        currentVideoIndex.value
        = (currentVideoIndex.value - 1 + allVideoUrls.value.length)
                % allVideoUrls.value.length;
    };

    const goToNextVideo = (): void => {
        if (allVideoUrls.value.length === 0) return;
        currentVideoIndex.value = (currentVideoIndex.value + 1) % allVideoUrls.value.length;
    };

    const goToVideo = (index: number,): void => {
        if (index >= 0 && index < allVideoUrls.value.length) {
            currentVideoIndex.value = index;
        }
    };

    const handleVideoError = (index: number,): void => {
        videoErrors.value.add(index,);
    };

    const handleMainImageError = (): void => {
        mainImageError.value = true;
    };

    const fetchNewsDetail = async (): Promise<void> => {
        try {
            loading.value = true;
            error.value = null;

            const { data: currentData, } = await getHttpClient().get<NewsItem>(
                `news/${newsId.value}`,
        );
            currentNewsItem.value = currentData;

            const { data: otherData, } = await getHttpClient().get<NewsListResponse>('news', {
                params: {
                    page: 1,
                    limit: 50,
                },
            });

            const filtered = (otherData.news || []).filter(
                (news: NewsItem,) => news.id !== newsId.value,);
            otherNews.value = filtered
                .map((n,) => ({ n, sort: Math.random(), }))
                .sort((a, b,) => a.sort - b.sort,)
                .map(({ n, },) => n,)
                .slice(0, 5,);
        } catch (err) {
            console.error('Error fetching news detail:', err,);
            error.value = t('news.detail.error.loadFailed',);
        } finally {
            loading.value = false;
        }
    };

    const navigateToNews = (id: string,): void => {
        router.push(`/news/${id}`,);
    };

    onMounted(() => {
        fetchNewsDetail();
    });
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

        <div
            v-if="currentNewsItem || loading"
            class="container mx-auto px-4 py-8"
        >
            <!-- Breadcrumb -->
            <nav
                v-if="currentNewsItem"
                class="mb-6"
            >
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
                            v-if="currentNewsItem"
                            class="text-gray-900 dark:text-white font-medium truncate max-w-xs"
                        >
                            {{ locale === 'ru' ? currentNewsItem.title_ru : currentNewsItem.title_en }}
                        </span>
                    </li>
                </ol>
            </nav>

            <div class="flex flex-col lg:flex-row gap-8 mb-8">
                <!-- Main Image (2/3 width on large screens) -->
                <div class="lg:w-2/3">
                    <div
                        v-if="currentNewsItem"
                        class="main-image-container rounded-2xl overflow-hidden shadow-lg min-h-75 md:min-h-100"
                    >
                    <div
                        v-if="mainImageError || !currentNewsItem.main_image"
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
                            :src="resolvedImagePath"
                            :alt="locale === 'ru' ? currentNewsItem.title_ru : currentNewsItem.title_en"
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
                            {{ $t('news.detail.sidebar.title',) }}
                        </h3>

                        <div
                            v-if="loading"
                            class="space-y-2"
                        >
                            <SideNewsTrailerSkeleton
                                v-for="i in 5"
                                :key="i"
                            />
                        </div>

                        <div
                            v-else-if="otherNews.length > 0"
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

            <div v-if="currentNewsItem">
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
                        {{ formatDate(currentNewsItem.datetime,) }}
                    </div>

                    <!-- Title -->
                    <h1
                        class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-4"
                    >
                        {{ locale === 'ru' ? currentNewsItem.title_ru : currentNewsItem.title_en }}
                    </h1>
                </div>

                <!-- News Text -->
                <div
                    v-if="currentNewsItem.text_ru || currentNewsItem.text_en"
                    class="news-text mb-12"
                >
                    <div class="prose prose-lg dark:prose-invert max-w-none">
                        <div
                            v-html="locale === 'ru' ? currentNewsItem.text_ru : currentNewsItem.text_en"
                            class="whitespace-pre-line text-gray-700 dark:text-gray-300 leading-relaxed text-justify"
                        ></div>
                    </div>
                </div>

                <!-- Video Carousel Section -->
                <div
                    v-if="allVideoUrls.length > 0"
                    class="video-section mb-12"
                >
                    <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                    {{ $t('news.detail.videos.title',) }}
                    <span
                        v-if="hasMultipleVideos"
                        class="text-lg font-normal text-gray-500 dark:text-gray-400 ml-2"
                    >
                        ({{ currentVideoIndex + 1 }} / {{ allVideoUrls.length }})
                    </span>
                    </h3>
                </div>

                <div class="video-carousel-container">
                    <!-- Video Player -->
                    <div
                        class="relative aspect-video rounded-xl overflow-hidden shadow-lg bg-black"
                    >
                        <!-- Current Video -->
                        <video
                            v-if="currentVideoUrl && !isCurrentVideoError"
                            :key="currentVideoIndex"
                            :src="currentVideoUrl"
                            class="w-full h-full object-contain"
                            controls
                            preload="metadata"
                            @error="handleVideoError(currentVideoIndex,)"
                        >
                            <p class="text-white text-center py-8">
                                {{ $t('news.detail.videos.notSupported',) }}
                            </p>
                        </video>

                        <!-- Video Error Fallback -->
                        <div
                            v-else
                            class="w-full h-full flex items-center justify-center bg-gray-100 dark:bg-gray-800 min-h-48"
                        >
                            <div class="text-center p-8">
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
                                        d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
                                    />
                                </svg>
                                <p class="text-gray-500 dark:text-gray-400">
                                    {{ $t('news.detail.videos.notSupported',) }}
                                </p>
                            </div>
                        </div>

                        <!-- Navigation Arrows (only if multiple videos) -->
                        <button
                            v-if="hasMultipleVideos"
                            class="carousel-arrow carousel-arrow-left"
                            aria-label="Previous video"
                            @click="goToPreviousVideo"
                        >
                            <svg
                                class="w-6 h-6"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 19l-7-7 7-7"
                                />
                            </svg>
                        </button>

                        <button
                            v-if="hasMultipleVideos"
                            class="carousel-arrow carousel-arrow-right"
                            aria-label="Next video"
                            @click="goToNextVideo"
                        >
                            <svg
                                class="w-6 h-6"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 5l7 7-7 7"
                                />
                            </svg>
                        </button>
                    </div>

                    <!-- Carousel Indicators (only if multiple videos) -->
                    <div
                        v-if="hasMultipleVideos"
                        class="flex items-center justify-center gap-2 mt-4"
                    >
                        <button
                            v-for="(_, index) in allVideoUrls"
                            :key="index"
                            class="carousel-dot"
                            :class="{
                                'carousel-dot-active': index === currentVideoIndex,
                            }"
                            :aria-label="`Go to video ${index + 1}`"
                            @click="goToVideo(index,)"
                        />
                    </div>
                </div>

                <!-- Gallery Section -->
                <div
                    v-if="galleryImages.length > 0"
                    class="gallery-section mb-12"
                >
                    <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                        {{ $t('news.detail.gallery.title',) }}
                    </h3>

                    <div class="gallery-grid">
                        <div
                            v-for="(imgUrl, index) in galleryImages"
                            :key="index"
                            class="gallery-item"
                            @click="openGalleryModal(index,)"
                        >
                            <div
                                v-if="galleryImageErrors.has(index,)"
                                class="gallery-item-error"
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
                                        d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                    />
                                </svg>
                            </div>
                            <img
                                v-else
                                :src="imgUrl"
                                :alt="`${$t('news.detail.gallery.title',)} ${index + 1}`"
                                class="gallery-image"
                                loading="lazy"
                                @error="handleGalleryImageError(index,)"
                            />
                        </div>
                    </div>

                    <!-- Empty State -->
                    <div
                        v-if="galleryImages.length === 0"
                        class="text-center py-12 bg-gray-50 dark:bg-gray-800 rounded-xl"
                    >
                        <svg
                            class="w-16 h-16 text-gray-400 mx-auto mb-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                            />
                        </svg>
                        <p class="text-gray-500 dark:text-gray-400">
                            {{ $t('news.detail.gallery.empty',) }}
                        </p>
                    </div>
                </div>

                <!-- Gallery Modal -->
                <Teleport to="body">
                    <div
                        v-if="selectedGalleryIndex !== null"
                        class="gallery-modal-overlay"
                        @click.self="closeGalleryModal"
                    >
                        <div class="gallery-modal-container">
                            <!-- Close button -->
                            <button
                                class="gallery-modal-close"
                                :aria-label="$t('news.detail.modal.close',)"
                                @click="closeGalleryModal"
                            >
                                <svg
                                    class="w-8 h-8"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                </svg>
                            </button>

                            <!-- Previous button -->
                            <button
                                v-if="galleryImages.length > 1"
                                class="gallery-modal-nav gallery-modal-prev"
                                :aria-label="$t('news.detail.modal.previous',)"
                                @click="goToPreviousGalleryImage"
                            >
                                <svg
                                    class="w-8 h-8"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M15 19l-7-7 7-7"
                                    />
                                </svg>
                            </button>

                            <!-- Image -->
                            <div class="gallery-modal-image-wrapper">
                            <img
                                v-if="!galleryImageErrors.has(selectedGalleryIndex,)"
                                :key="selectedGalleryIndex"
                                :src="galleryImages[selectedGalleryIndex]"
                                :alt="`${$t('news.detail.gallery.title',)} ${
                                    selectedGalleryIndex + 1
                                }`"
                                class="gallery-modal-image"
                                @error="handleGalleryImageError(selectedGalleryIndex,)"
                            />
                                <div
                                    v-else
                                    class="gallery-modal-error"
                                >
                                    <svg
                                        class="w-16 h-16 text-gray-400"
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                        />
                                    </svg>
                                </div>
                            </div>

                            <!-- Next button -->
                            <button
                                v-if="galleryImages.length > 1"
                                class="gallery-modal-nav gallery-modal-next"
                                :aria-label="$t('news.detail.modal.next',)"
                                @click="goToNextGalleryImage"
                            >
                                <svg
                                    class="w-8 h-8"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M9 5l7 7-7 7"
                                    />
                                </svg>
                            </button>

                            <!-- Caption -->
                            <div class="gallery-modal-caption">
                                {{
                                    $t('news.detail.modal.caption', {
                                        current: selectedGalleryIndex + 1,
                                        total: galleryImages.length,
                                    },)
                                }}
                            </div>
                        </div>
                    </div>
                </Teleport>
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

/* Video Carousel Styles */
.video-carousel-container {
    position: relative;
    width: 100%;
}

.carousel-arrow {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    border-radius: 50%;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: 2px solid rgba(255, 255, 255, 0.3);
    cursor: pointer;
    transition: all 0.2s ease;
    opacity: 0;
}

.video-carousel-container:hover .carousel-arrow {
    opacity: 1;
}

.carousel-arrow:hover {
    background-color: rgba(0, 0, 0, 0.8);
    border-color: rgba(255, 255, 255, 0.6);
    transform: translateY(-50%) scale(1.1);
}

.carousel-arrow-left {
    left: 12px;
}

.carousel-arrow-right {
    right: 12px;
}

.carousel-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background-color: #d1d5db;
    border: none;
    cursor: pointer;
    transition: all 0.2s ease;
    padding: 0;
}

.dark .carousel-dot {
    background-color: #4b5563;
}

.carousel-dot:hover {
    background-color: #9ca3af;
    transform: scale(1.2);
}

.dark .carousel-dot:hover {
    background-color: #6b7280;
}

.carousel-dot-active {
    background-color: #10b981;
    transform: scale(1.3);
}

.dark .carousel-dot-active {
    background-color: #34d399;
}

.carousel-dot-active:hover {
    background-color: #059669;
    transform: scale(1.3);
}

.dark .carousel-dot-active:hover {
    background-color: #10b981;
}

/* Gallery Section Styles */
.gallery-section {
    scroll-margin-top: 2rem;
}

.gallery-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
}

.gallery-item {
    position: relative;
    border-radius: 0.75rem;
    overflow: hidden;
    cursor: pointer;
    aspect-ratio: 4 / 3;
    background: linear-gradient(135deg, #e5e7eb 0%, #d1d5db 100%);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.dark .gallery-item {
    background: linear-gradient(135deg, #374151 0%, #4b5563 100%);
}

.gallery-item:hover {
    transform: scale(1.03);
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.15);
}

.gallery-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
}

.gallery-item:hover .gallery-image {
    transform: scale(1.08);
}

.gallery-item-error {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    min-height: 200px;
}

/* Gallery Modal Styles */
.gallery-modal-overlay {
    position: fixed;
    inset: 0;
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(4px);
    padding: 1rem;
}

.gallery-modal-container {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    max-width: 90vw;
    max-height: 90vh;
}

.gallery-modal-image-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    max-width: 100%;
    max-height: 80vh;
}

.gallery-modal-image {
    max-width: 100%;
    max-height: 80vh;
    object-fit: contain;
    border-radius: 0.5rem;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.gallery-modal-error {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 400px;
    height: 300px;
    background: #1f2937;
    border-radius: 0.5rem;
}

.gallery-modal-close {
    position: absolute;
    top: -2.5rem;
    right: -0.5rem;
    z-index: 10;
    color: white;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.25rem;
    transition: transform 0.2s ease;
}

.gallery-modal-close:hover {
    transform: scale(1.15);
}

.gallery-modal-nav {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
    color: white;
    background: rgba(0, 0, 0, 0.5);
    border: none;
    border-radius: 50%;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
}

.gallery-modal-nav:hover {
    background: rgba(0, 0, 0, 0.8);
    transform: translateY(-50%) scale(1.1);
}

.gallery-modal-prev {
    left: -3.5rem;
}

.gallery-modal-next {
    right: -3.5rem;
}

.gallery-modal-caption {
    margin-top: 1rem;
    color: #d1d5db;
    font-size: 0.875rem;
    text-align: center;
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

@media (max-width: 640px) {
    .carousel-arrow {
        width: 36px;
        height: 36px;
        opacity: 1;
        background-color: rgba(0, 0, 0, 0.5);
    }

    .carousel-arrow svg {
        width: 18px;
        height: 18px;
    }
}
</style>
