<script lang="ts">
import SideNewsTrailer from '@/components/News/SideNewsTrailer.vue';
import NewsItemDescription from '@/components/News/NewsItemDescription.vue';
import { inject } from 'vue';
import type { NewsItemType } from '@/types';
import { fetchCurrentNews, fetchOtherNews } from '@/api/requests';
import VideoSection from '@/components/News/VideoSection.vue';
import PhotoSection from './PhotoSection.vue';
import SideNewsTrailerSkeleton from '../UI/Skeletons/SideNewsTrailerSkeleton.vue';

export default {
    setup() {
        const jsonserverhost: string = inject('jsonserverhost') as string;
        const imagebasedir = inject('imagebasedir') as string;
        return {
            jsonserverhost,
            imagebasedir,
        };
    },
    components: {
        SideNewsTrailer,
        NewsItemDescription,
        VideoSection,
        PhotoSection,
        SideNewsTrailerSkeleton,
    },
    data() {
        return {
            currentNewsItem: {} as NewsItemType,
            otherNews: [] as NewsItemType[],
            isImgLoaded: false,
        };
    },
    methods: {
        makeVideoSlideLabel(index: number) {
            return 'Видеослайд ' + index;
        },

        makeVideoName(index: number) {
            return (
                this.imagebasedir + this.currentNewsItem.dir + index + '.mp4'
            );
        },

        onImgLoaded() {
            this.isImgLoaded = true;
        },
    },

    async created() {
        this.currentNewsItem = await fetchCurrentNews(
            this.$route.path,
            this.jsonserverhost,
        );
        this.otherNews = await fetchOtherNews(
            this.$route.path,
            this.jsonserverhost,
        );
        const mdbScript = document.createElement('script');
        mdbScript.setAttribute('src', '/src/assets/js/mdb.min.js');
        document.head.appendChild(mdbScript);
    },
    computed: {
        background() {
            return (
                this.imagebasedir +
                this.currentNewsItem.dir +
                this.currentNewsItem.img_backfull
            );
        },
    },
};
</script>

<template>
    <section class="container text-center px-0 main-content">
        <article class="container">
            <div class="news-header">
                <div class="news-img">
                    <div class="img-mock" v-show="!isImgLoaded" />
                    <img
                        :src="background"
                        width="670px"
                        height="450px"
                        alt="News main image"
                        @load="onImgLoaded"
                        v-show="isImgLoaded"
                    />
                </div>

                <div class="other-news px-sm-0">
                    <div
                        class="other-news-item-wrapper"
                        v-for="newsItem in otherNews"
                        v-bind:key="newsItem.id"
                    >
                        <SideNewsTrailer :sideNewsObject="newsItem" />
                    </div>
                </div>
            </div>

            <NewsItemDescription :newsObject="currentNewsItem" />

            <div class="news-text">
                <p>
                    <span
                        v-html="
                            $i18n.locale === 'RUS'
                                ? currentNewsItem.text_ru
                                : currentNewsItem.text_en
                        "
                    />
                </p>
            </div>
        </article>

        <PhotoSection :current-news-item="currentNewsItem" />

        <VideoSection :current-news-item="currentNewsItem" />
    </section>
</template>

<style scoped>
@import '@/assets/css/mdb.min.css';

.news-header {
    display: flex;
    column-gap: 1rem;
}
.img-mock {
    height: 40rem;
    width: 60rem;
}
img {
    min-width: 100%;
    min-width: auto;
    width: 100%;
    height: auto;
}

.other-news {
    display: flex;
    flex-direction: column;
    background-color: var(--color-surface);
    border-radius: 15px;
    width: 60%;
}

.other-news-item-wrapper {
    margin-bottom: 0.8rem;
}

.news-text {
    text-align: left;
    margin-bottom: 2rem;
    color: var(--color-on-surface);
}

@media screen and (max-width: 991px) {
    .news-header {
        flex-direction: column;
    }

    .other-news {
        padding-top: 0.8rem;
        width: 100%;
        padding: 0 15px;
    }

    .other-news-item-wrapper {
        margin-bottom: 0.8rem;
    }

    .news-img {
        margin-bottom: 1rem;
    }
}
</style>
