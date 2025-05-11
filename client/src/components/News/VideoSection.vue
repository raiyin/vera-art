<script lang="ts">
import { useI18n } from 'vue-i18n';
import { inject, PropType, ref } from 'vue';
import type { NewsItemType } from '@/types';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
        const imagebasedir = inject('imagebasedir') as string;
        const players = ref([]);
        return {
            imagebasedir,
            t,
            locale,
            players,
        };
    },
    props: {
        currentNewsItem: {
            type: Object as PropType<NewsItemType>,
            default: {} as NewsItemType,
        },
    },
    methods: {
        makeVideoSlideLabel(index: number) {
            return this.t('news.videoslide') + ' ' + index;
        },

        makeVideoName(index: number, imagebasedir: string): string {
            return imagebasedir + this.currentNewsItem.dir + index + '.mp4';
        },
    },
    mounted() {
        this.players.value = Array.from(document.querySelectorAll('video'));

        Array.from(
            document.querySelectorAll(
                '.carousel-control-prev, .carousel-control-next, .video-indicator'
            )
        ).forEach((item) => {
            item.addEventListener('click', () => {
                this.players.value.forEach((player) => {
                    player.pause();
                });
            });
        });
    },
};
</script>

<template>
    <section
        id="carouselVideoExample"
        class="carousel slide"
        data-bs-interval="false"
        data-bs-ride="carousel"
    >
        <ol v-if="currentNewsItem.videoscount > 1" class="carousel-indicators">
            <template
                v-for="video_index in currentNewsItem.videoscount"
                :key="video_index"
            >
                <li
                    data-mdb-target="#carouselVideoExample"
                    :data-mdb-slide-to="video_index - 1"
                    class="video-indicator"
                    :class="{ active: video_index === 1 }"
                    :aria-current="video_index === 1"
                    :aria-label="makeVideoSlideLabel(video_index)"
                />
            </template>
        </ol>

        <div class="carousel-inner">
            <div
                v-for="video_index in currentNewsItem.videoscount"
                :key="video_index"
                class="carousel-item"
                data-bs-interval="4000000000"
                :class="{ active: video_index === 1 }"
            >
                <video class="img-fluid" controls>
                    <source
                        :src="makeVideoName(video_index, imagebasedir)"
                        type="video/mp4"
                    />
                    <p>
                        Sorry, there's a problem playing this video. Please try using a
                        different browser.
                    </p>
                </video>
            </div>
        </div>

        <a
            v-if="currentNewsItem.videoscount > 1"
            class="carousel-control-prev"
            type="button"
            data-mdb-target="#carouselVideoExample"
            data-mdb-slide="prev"
        >
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Previous</span>
        </a>

        <a
            v-if="currentNewsItem.videoscount > 1"
            class="carousel-control-next"
            type="button"
            data-mdb-target="#carouselVideoExample"
            data-mdb-slide="next"
        >
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Next</span>
        </a>
    </section>
</template>

<style scoped>
@import '@/assets/css/mdb.min.css';

.carousel-inner {
    object-fit: cover;
}

.carousel-control-prev,
.carousel-control-next {
    height: 15%;
    align-self: center;
}

.img-fluid {
    width: 100%;
}

video {
    cursor: pointer;
}

@media (orientation: landscape) {
    .img-fluid {
        max-height: 90vh;
    }
}

@media (orientation: portrait) {
    .img-fluid {
        max-width: 90vw;
    }
}
</style>
