<script lang="ts">
import SideNewsTrailer from '@/components/News/SideNewsTrailer.vue';
import NewsItemDescription from '@/components/News/NewsItemDescription.vue';
import NewsPhotoItem from '@/components/News/NewsPhotoItem.vue';
import { useI18n } from 'vue-i18n';
import { inject } from 'vue';
import type { NewsItemType } from '@/types';
import { fetchCurrentNews, fetchOtherNews } from '@/api/requests';
import VideoSection from '@/components/News/VideoSection.vue';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
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
        NewsPhotoItem,
        VideoSection,
    },
    data() {
        return {
            currentNewsItem: {} as NewsItemType,
            otherNews: [] as NewsItemType[],
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
    },

    async mounted() {
        this.currentNewsItem = await fetchCurrentNews(
            this.$route.path,
            this.jsonserverhost
        );
        this.otherNews = await fetchOtherNews(
            this.$route.path,
            this.jsonserverhost
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
        <article class="container px-0">
            <div class="news-header">
                <div class="news-img">
                    <img :src="background" />
                </div>

                <div class="other-news">
                    <template
                        v-for="newsItem in otherNews"
                        v-bind:key="newsItem"
                    >
                        <div class="other-news-item-wrapper">
                            <SideNewsTrailer :sideNewsObject="newsItem" />
                        </div>
                    </template>
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

        <!-- photo section -->
        <div class="row">
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 1">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 2">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 0">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>
        </div>

        <VideoSection :current-news-item="currentNewsItem" />
    </section>
</template>

<style scoped>
@import '@/assets/css/mdb.min.css';

.news-header {
    display: flex;
    column-gap: 1rem;
}

img {
    width: 100%;
    height: auto;
}

.other-news {
    display: flex;
    flex-direction: column;
    background-color: var(--color-surface);
    padding: 0 15px;
    border-radius: 15px;
}

.other-news-item-wrapper {
    margin-bottom: 0.8rem;
}

.news-text {
    text-align: left;
    margin-bottom: 2rem;
    color: var(--color-on-surface);
}

@media screen and (max-width: 1000px) {
    .news-header {
        flex-direction: column;
    }

    .other-news {
        padding-top: 0.8rem;
    }

    .other-news-item-wrapper {
        margin-bottom: 0.8rem;
    }

    .news-img {
        margin-bottom: 1rem;
    }
}
</style>
