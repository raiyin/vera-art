<script lang="ts">
import { useI18n } from 'vue-i18n';
import { ref, computed, type PropType, onMounted, onUnmounted, nextTick } from 'vue';
import type { NewsDesc } from '~/types';

export default {
    props: {
        currentNewsItem: {
            type: Object as PropType<NewsDesc>,
            default: () => ({} as NewsDesc),
        },
    },
    setup(props) {
        const { t, locale } = useI18n({ useScope: 'global' });
        const videoPlayer = ref<HTMLVideoElement | null>(null);
        const videoErrors = ref<Set<number>>(new Set());
        const currentIndex = ref(0);
        const isPlaying = ref(false);
        const isFullscreen = ref(false);
        const showControls = ref(true);
        const controlsTimeout = ref<number | null>(null);

        const handleVideoError = (index: number) => {
            videoErrors.value.add(index);
        };

        const resetVideoErrors = () => {
            videoErrors.value.clear();
        };

        const nextSlide = () => {
            if (!props.currentNewsItem?.videos) return;
            pauseCurrentVideo();
            currentIndex.value = (currentIndex.value + 1) % props.currentNewsItem.videos.length;
            resetControlsTimer();
        };

        const prevSlide = () => {
            if (!props.currentNewsItem?.videos) return;
            pauseCurrentVideo();
            currentIndex.value = currentIndex.value === 0
                ? (props.currentNewsItem.videos.length - 1)
                : currentIndex.value - 1;
            resetControlsTimer();
        };

        const goToSlide = (index: number) => {
            if (!props.currentNewsItem?.videos) return;
            pauseCurrentVideo();
            currentIndex.value = index;
            resetControlsTimer();
        };

        const pauseCurrentVideo = () => {
            if (videoPlayer.value && !videoPlayer.value.paused) {
                videoPlayer.value.pause();
                isPlaying.value = false;
            }
        };

        const togglePlayPause = () => {
            if (!videoPlayer.value) return;

            if (videoPlayer.value.paused) {
                videoPlayer.value.play();
                isPlaying.value = true;
            } else {
                videoPlayer.value.pause();
                isPlaying.value = false;
            }
            resetControlsTimer();
        };

        const toggleFullscreen = () => {
            const container = document.querySelector('.video-carousel-container');
            if (!container) return;

            if (!isFullscreen.value) {
                if (container.requestFullscreen) {
                    container.requestFullscreen();
                }
            } else {
                if (document.exitFullscreen) {
                    document.exitFullscreen();
                }
            }
            resetControlsTimer();
        };

        const resetControlsTimer = () => {
            showControls.value = true;
            if (controlsTimeout.value) clearTimeout(controlsTimeout.value);
            controlsTimeout.value = window.setTimeout(() => {
                showControls.value = false;
            }, 3000);
        };

        const handleFullscreenChange = () => {
            isFullscreen.value = !!document.fullscreenElement;
        };

        onMounted(() => {
            // videoPlayer ref will be set by template
            document.addEventListener('fullscreenchange', handleFullscreenChange);
            resetControlsTimer();
        });

        onUnmounted(() => {
            if (controlsTimeout.value) clearTimeout(controlsTimeout.value);
            document.removeEventListener('fullscreenchange', handleFullscreenChange);
        });

        return {
            t,
            locale,
            videoPlayer,
            videoErrors,
            currentIndex,
            isPlaying,
            isFullscreen,
            showControls,
            controlsTimeout,
            handleVideoError,
            resetVideoErrors,
            nextSlide,
            prevSlide,
            goToSlide,
            togglePlayPause,
            toggleFullscreen,
            resetControlsTimer,
        };
    },
    computed: {
        videoCount() {
            return this.currentNewsItem?.videos?.length || 0;
        },
        hasMultipleVideos() {
            return this.videoCount > 1;
        },
        currentVideoSrc() {
            // Keep for backward compatibility (returns MP4 as default)
            const base = this.getVideoBaseName(this.currentNewsItem?.videos?.[this.currentIndex] || '');
            return base ? `${this.currentNewsItem.dir}videos/${base}/${base}.mp4` : '';
        },
        thumbnailUrls() {
            if (!this.currentNewsItem?.videos) return [];
            return this.currentNewsItem.videos.map(video => {
                const base = this.getVideoBaseName(video);
                return `${this.currentNewsItem.dir}videos/${base}/${base}.mp4`;
            });
        },
        currentVideoSources() {
            const video = this.currentNewsItem?.videos?.[this.currentIndex];
            if (!video) return [];
            const base = this.getVideoBaseName(video);
            const dir = this.currentNewsItem.dir;
            return [
                {
                    src: `${dir}videos/${base}/${base}.m3u8`,
                    type: 'application/vnd.apple.mpegurl',
                },
                {
                    src: `${dir}videos/${base}/${base}.mp4`,
                    type: 'video/mp4',
                },
            ];
        },
    },
    methods: {
        makeVideoSlideLabel(index: number) {
            return this.t('news.videoslide') + ' ' + (index + 1);
        },

        makeVideoName(fileName: string): string {
            return this.currentNewsItem.dir + fileName;
        },

        getVideoBaseName(fileName: string): string {
            if (!fileName) return '';
            // Remove any known video extensions
            return fileName.replace(/\.(mp4|m3u8|webm|ogg|mov|avi)$/i, '');
        },

        handleVideoPlay() {
            this.isPlaying = true;
            this.resetControlsTimer();
        },

        handleVideoPause() {
            this.isPlaying = false;
            this.resetControlsTimer();
        },
    },
};
</script>

<template>
    <div
        class="video-carousel-container"
        @mousemove="resetControlsTimer"
        @touchstart="resetControlsTimer"
    >
        <!-- Main Video Player -->
        <div class="video-main-container">
            <div v-if="videoErrors.has(currentIndex)" class="video-error-state">
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
                    ></path>
                </svg>
                <p class="text-gray-500 dark:text-gray-400 text-center">
                    Video failed to load
                </p>
            </div>

            <div v-else class="video-wrapper relative">
                <video
                    ref="videoPlayer"
                    :key="currentIndex"
                    class="video-player"
                    controls
                    @error="() => handleVideoError(currentIndex)"
                    @play="handleVideoPlay"
                    @pause="handleVideoPause"
                >
                    <source
                        v-for="source in currentVideoSources"
                        :key="source.src"
                        :src="source.src"
                        :type="source.type"
                    />
                    <p>Your browser does not support the video tag.</p>
                </video>

                <!-- Custom Controls Overlay -->
                <transition name="fade">
                    <div v-if="showControls" class="video-controls-overlay">
                        <div class="controls-top">
                            <div class="video-counter">
                                <span class="counter-current">{{
                                    currentIndex + 1
                                }}</span>
                                <span class="counter-separator">/</span>
                                <span class="counter-total">{{ videoCount }}</span>
                            </div>

                            <button
                                @click="toggleFullscreen"
                                class="control-button"
                                aria-label="Toggle fullscreen"
                            >
                                <svg
                                    class="w-6 h-6"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        v-if="isFullscreen"
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M9 9V4.5a.75.75 0 00-1.5 0V9H4.5a.75.75 0 000 1.5H9v4.5a.75.75 0 001.5 0V10.5h4.5a.75.75 0 000-1.5H10.5V4.5a.75.75 0 00-1.5 0V9z"
                                    />
                                    <path
                                        v-else
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M4 8V4m0 0h4M4 4l5 5m11-5h-4m4 0v4m0-4l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5h-4m4 0v-4m0 4l-5-5"
                                    />
                                </svg>
                            </button>
                        </div>

                        <div class="controls-center">
                            <button
                                v-if="hasMultipleVideos"
                                @click="prevSlide"
                                class="nav-button prev"
                                :aria-label="$t('carousel.back')"
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

                            <button
                                @click="togglePlayPause"
                                class="play-pause-button"
                                :aria-label="isPlaying ? 'Pause' : 'Play'"
                            >
                                <svg
                                    v-if="!isPlaying"
                                    class="w-12 h-12"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                                    />
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                                <svg
                                    v-else
                                    class="w-12 h-12"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                            </button>

                            <button
                                v-if="hasMultipleVideos"
                                @click="nextSlide"
                                class="nav-button next"
                                :aria-label="$t('carousel.next')"
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
                        </div>

                        <div class="controls-bottom">
                            <div class="video-title">
                                {{ makeVideoSlideLabel(currentIndex) }}
                            </div>
                        </div>
                    </div>
                </transition>
            </div>
        </div>

        <!-- Thumbnail Navigation -->
        <div v-if="hasMultipleVideos" class="video-thumbnails">
            <div
                v-for="(thumbnail, index) in thumbnailUrls"
                :key="index"
                class="thumbnail-item"
                :class="{ active: index === currentIndex }"
                @click="goToSlide(index)"
                @mouseenter="resetControlsTimer"
            >
                <div class="thumbnail-image-container">
                    <div v-if="videoErrors.has(index)" class="thumbnail-error">
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
                                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                        </svg>
                    </div>
                    <template v-else>
                        <video
                            class="thumbnail-video"
                            :src="thumbnail"
                            muted
                            preload="metadata"
                            playsinline
                            disablePictureInPicture
                            disableRemotePlayback
                            @error="() => handleVideoError(index)"
                        />
                        <div class="thumbnail-play-overlay">
                            <svg
                                class="w-4 h-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                                />
                            </svg>
                        </div>
                    </template>
                </div>
                <div class="thumbnail-label">
                    {{ index + 1 }}
                </div>
            </div>
        </div>

        <!-- Dot Indicators (Mobile) -->
        <div v-if="hasMultipleVideos" class="video-dots">
            <button
                v-for="index in videoCount"
                :key="index"
                class="dot"
                :class="{ active: index - 1 === currentIndex }"
                @click="goToSlide(index - 1)"
                :aria-label="makeVideoSlideLabel(index - 1)"
            />
        </div>
    </div>
</template>

<style scoped>
.video-carousel-container {
    position: relative;
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
    border-radius: 1.5rem;
    overflow: hidden;
    box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.3);
    transition: all 0.3s ease;
}

.dark .video-carousel-container {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.video-main-container {
    position: relative;
    aspect-ratio: 16/9;
    background: #000;
}

.video-wrapper {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.video-player {
    width: 100%;
    height: 100%;
    object-fit: contain;
    outline: none;
}

.video-controls-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(
        to bottom,
        rgba(0, 0, 0, 0.3) 0%,
        rgba(0, 0, 0, 0.2) 20%,
        transparent 40%,
        transparent 60%,
        rgba(0, 0, 0, 0.2) 80%,
        rgba(0, 0, 0, 0.3) 100%
    );
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 1.5rem;
    z-index: 10;
    opacity: 1;
    transition: opacity 0.3s ease;
    pointer-events: none;
}

.video-controls-overlay > * {
    pointer-events: auto;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.controls-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.video-counter {
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 2rem;
    font-size: 0.875rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    backdrop-filter: blur(4px);
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

.control-button {
    background: rgba(0, 0, 0, 0.7);
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
    backdrop-filter: blur(4px);
}

.control-button:hover {
    background: rgba(0, 0, 0, 0.9);
    transform: scale(1.1);
}

.controls-center {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 2rem;
}

.nav-button {
    background: rgba(0, 0, 0, 0.7);
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
    backdrop-filter: blur(4px);
}

.nav-button:hover {
    background: rgba(0, 0, 0, 0.9);
    transform: scale(1.1);
}

.play-pause-button {
    background: rgba(16, 185, 129, 0.9);
    border: none;
    border-radius: 50%;
    width: 72px;
    height: 72px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: white;
    transition: all 0.2s ease;
    backdrop-filter: blur(4px);
}

.play-pause-button:hover {
    background: rgba(16, 185, 129, 1);
    transform: scale(1.1);
}

.controls-bottom {
    display: flex;
    justify-content: center;
}

.video-title {
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 2rem;
    font-size: 0.875rem;
    font-weight: 500;
    backdrop-filter: blur(4px);
}

.video-thumbnails {
    display: flex;
    gap: 0.75rem;
    padding: 1.5rem;
    background: rgba(0, 0, 0, 0.3);
    overflow-x: auto;
    scrollbar-width: thin;
    scrollbar-color: #10b981 transparent;
}

.video-thumbnails::-webkit-scrollbar {
    height: 6px;
}

.video-thumbnails::-webkit-scrollbar-track {
    background: transparent;
}

.video-thumbnails::-webkit-scrollbar-thumb {
    background-color: #10b981;
    border-radius: 3px;
}

.thumbnail-item {
    flex: 0 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    transition: all 0.2s ease;
    opacity: 0.7;
}

.thumbnail-item:hover {
    opacity: 1;
    transform: translateY(-2px);
}

.thumbnail-item.active {
    opacity: 1;
}

.thumbnail-item.active .thumbnail-image-container {
    border-color: #10b981;
    transform: scale(1.05);
}

.thumbnail-image-container {
    width: 80px;
    height: 60px;
    border-radius: 0.75rem;
    overflow: hidden;
    background: linear-gradient(135deg, #374151 0%, #4b5563 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 3px solid transparent;
    transition: all 0.2s ease;
    position: relative;
}

.thumbnail-video {
    width: 100%;
    height: 100%;
    object-fit: cover;
    background: #000;
}

.thumbnail-error {
    color: #ef4444;
}

.thumbnail-play-icon {
    color: white;
    opacity: 0.8;
}

.thumbnail-play-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.3);
    opacity: 0;
    transition: opacity 0.2s ease;
    pointer-events: none;
}

.thumbnail-item:hover .thumbnail-play-overlay {
    opacity: 1;
}

.thumbnail-label {
    color: white;
    font-size: 0.75rem;
    font-weight: 500;
    opacity: 0.8;
}

.thumbnail-item.active .thumbnail-label {
    color: #10b981;
    opacity: 1;
    font-weight: 600;
}

.video-dots {
    display: none;
    justify-content: center;
    gap: 0.5rem;
    padding: 1rem;
    background: rgba(0, 0, 0, 0.2);
}

.dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.3);
    border: none;
    cursor: pointer;
    transition: all 0.2s ease;
    padding: 0;
}

.dot:hover {
    background: rgba(255, 255, 255, 0.5);
    transform: scale(1.2);
}

.dot.active {
    background: #10b981;
    transform: scale(1.2);
}

.video-error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    padding: 2rem;
    background: linear-gradient(135deg, #374151 0%, #4b5563 100%);
}

.dark .video-error-state {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
}

/* Responsive Design */
@media (max-width: 768px) {
    .video-carousel-container {
        border-radius: 1rem;
    }

    .video-controls-overlay {
        padding: 1rem;
    }

    .nav-button {
        width: 44px;
        height: 44px;
    }

    .play-pause-button {
        width: 60px;
        height: 60px;
    }

    .play-pause-button svg {
        width: 10px;
        height: 10px;
    }

    .video-thumbnails {
        padding: 1rem;
        gap: 0.5rem;
    }

    .thumbnail-image-container {
        width: 60px;
        height: 45px;
    }

    .video-dots {
        display: flex;
    }
}

@media (max-width: 480px) {
    .video-thumbnails {
        display: none;
    }

    .video-dots {
        display: flex;
    }

    .controls-center {
        gap: 1rem;
    }

    .nav-button {
        width: 36px;
        height: 36px;
    }

    .play-pause-button {
        width: 48px;
        height: 48px;
    }

    .play-pause-button svg {
        width: 8px;
        height: 8px;
    }

    .video-counter {
        font-size: 0.75rem;
        padding: 0.375rem 0.75rem;
    }

    .video-title {
        font-size: 0.75rem;
        padding: 0.5rem 1rem;
    }
}

/* Fullscreen styles */
:fullscreen .video-carousel-container {
    border-radius: 0;
    max-width: 100vw;
    max-height: 100vh;
}

:fullscreen .video-main-container {
    height: 100vh;
    aspect-ratio: unset;
}

:fullscreen .video-player {
    max-height: 100vh;
}

/* Keyboard navigation focus styles */
.video-carousel-container:focus-within {
    outline: 2px solid #10b981;
    outline-offset: 2px;
}

.control-button:focus,
.nav-button:focus,
.play-pause-button:focus,
.dot:focus,
.thumbnail-item:focus {
    outline: 2px solid #10b981;
    outline-offset: 2px;
}
</style>
