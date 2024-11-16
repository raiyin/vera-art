<script lang="ts">
import { useI18n } from 'vue-i18n';
import { inject, PropType } from 'vue';
import type { NewsItemType } from '@/types';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
        const imagebasedir = inject('imagebasedir') as string;
        return {
            imagebasedir,
            t,
            locale,
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
};
</script>

<template>
    <div
        id="carouselVideoExample"
        class="carousel slide carousel-fade"
        data-mdb-ride="carousel"
    >
        <ol v-if="currentNewsItem.videoscount > 1" class="carousel-indicators">
            <template
                v-for="video_index in currentNewsItem.videoscount"
                :key="video_index"
            >
                <li
                    data-mdb-target="#carouselVideoExample"
                    :data-mdb-slide-to="video_index - 1"
                    :class="video_index === 1 ? 'active' : ''"
                    :aria-current="video_index === 1"
                    :aria-label="makeVideoSlideLabel(video_index)"
                />
            </template>
        </ol>

        <div class="carousel-inner">
            <div
                v-for="video_index in currentNewsItem.videoscount"
                v-bind:key="video_index"
                :class="'carousel-item ' + (video_index == 1 ? 'active' : '')"
            >
                <video class="img-fluid" controls>
                    <source
                        :src="makeVideoName(video_index, imagebasedir)"
                        type="video/mp4"
                    />
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
    </div>
</template>

<style scoped>
@import '@/assets/css/mdb.min.css';

.carousel-inner {
    object-fit: cover;
}

.img-fluid {
    width: 100%;
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
