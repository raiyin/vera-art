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
            return 'Видеослайд ' + index;
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
        <div v-if="currentNewsItem.videoscount > 1" class="carousel-indicators">
            <template v-for="video_index in currentNewsItem.videoscount">
                <button
                    v-if="video_index == 1"
                    v-bind:key="video_index"
                    type="button"
                    data-mdb-target="#carouselVideoExample"
                    :data-mdb-slide-to="video_index"
                    class="active"
                    aria-current="true"
                    :aria-label="makeVideoSlideLabel(video_index)"
                />

                <button
                    v-if="video_index != 1"
                    v-bind:key="video_index"
                    type="button"
                    data-mdb-target="#carouselVideoExample"
                    :data-mdb-slide-to="video_index"
                    :aria-label="makeVideoSlideLabel(video_index)"
                />
            </template>
        </div>

        <div class="carousel-inner">
            <template
                v-for="video_index in currentNewsItem.videoscount"
                v-bind:key="video_index"
            >
                <div v-if="video_index == 1" class="carousel-item active">
                    <video class="img-fluid" controls>
                        <source
                            :src="makeVideoName(video_index, imagebasedir)"
                            type="video/mp4"
                        />
                    </video>
                </div>

                <div
                    v-if="video_index != 1"
                    v-bind:key="video_index"
                    class="carousel-item"
                >
                    <video class="img-fluid" controls>
                        <source
                            :src="makeVideoName(video_index, imagebasedir)"
                            type="video/mp4"
                        />
                    </video>
                </div>
            </template>
        </div>

        <button
            v-if="currentNewsItem.videoscount > 1"
            class="carousel-control-prev"
            type="button"
            data-mdb-target="#carouselVideoExample"
            data-mdb-slide="prev"
        >
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Previous</span>
        </button>

        <button
            v-if="currentNewsItem.videoscount > 1"
            class="carousel-control-next"
            type="button"
            data-mdb-target="#carouselVideoExample"
            data-mdb-slide="next"
        >
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Next</span>
        </button>
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
