<script>
import SideNewsTrailer from '@/components/News/SideNewsTrailer.vue';
import NewsItemDescription from '@/components/News/NewsItemDescription.vue';
import NewsPhotoItem from '@/components/News/NewsPhotoItem.vue';
import axios from 'axios';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
    },
    inject: ['jsonserverhost', 'imagebasedir'],
    components: {
        SideNewsTrailer,
        NewsItemDescription,
        NewsPhotoItem,
    },
    data() {
        return {
            news: [],
            currentNews: {},
        };
    },
    methods: {
        async fetchData() {
            try {
                const newsid = this.$route.path.substring(
                    this.$route.path.lastIndexOf('/') + 1
                );
                let response = await axios.get(this.jsonserverhost + 'news', {
                    params: { id: newsid },
                });
                this.currentNews = response.data[0];

                response = await axios.get(this.jsonserverhost + 'news', {
                    params: { id_ne: newsid, _limit: 5 },
                });
                this.news = response.data;
            } catch (e) {
                console.log(e);
            }
        },
        makeVideoSlideLabel(index) {
            return 'Видеослайд ' + index;
        },
        makeVideoName(index) {
            return this.imagebasedir + this.currentNews.dir + index + '.mp4';
        },
    },
    mounted() {
        this.fetchData();
        const mdbScript = document.createElement('script');
        mdbScript.setAttribute('src', '/src/assets/js/mdb.min.js');
        document.head.appendChild(mdbScript);
    },
    computed: {
        background() {
            return (
                this.imagebasedir + this.currentNews.dir + this.currentNews.img_backfull
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
                    <img :src="background" />
                </div>

                <div class="other-news">
                    <template v-for="newsItem in news" v-bind:key="newsItem">
                        <div class="other-news-item-wrapper">
                            <SideNewsTrailer :sideNewsObject="newsItem" />
                        </div>
                    </template>
                </div>
            </div>

            <NewsItemDescription :newsObject="currentNews" />

            <div class="news-text">
                <p>
                    <span
                        v-html="
                            $i18n.locale === 'ru'
                                ? currentNews.text_ru
                                : currentNews.text_en
                        "
                    />
                </p>
            </div>
        </article>

        <!-- photo section -->
        <div class="row">
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNews.imagescount">
                    <template v-if="image_index % 3 == 1">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNews"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNews.imagescount">
                    <template v-if="image_index % 3 == 2">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNews"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNews.imagescount">
                    <template v-if="image_index % 3 == 0">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNews"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>
        </div>

        <!-- video section -->
        <div
            id="carouselVideoExample"
            class="carousel slide carousel-fade"
            data-mdb-ride="carousel"
        >
            <div v-if="currentNews.videoscount > 1" class="carousel-indicators">
                <template v-for="video_index in currentNews.videoscount">
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
                    v-for="video_index in currentNews.videoscount"
                    v-bind:key="video_index"
                >
                    <div v-if="video_index == 1" class="carousel-item active">
                        <video class="img-fluid" controls>
                            <source :src="makeVideoName(video_index)" type="video/mp4" />
                        </video>
                    </div>

                    <div
                        v-if="video_index != 1"
                        v-bind:key="video_index"
                        class="carousel-item"
                    >
                        <video class="img-fluid" controls>
                            <source :src="makeVideoName(video_index)" type="video/mp4" />
                        </video>
                    </div>
                </template>
            </div>

            <button
                v-if="currentNews.videoscount > 1"
                class="carousel-control-prev"
                type="button"
                data-mdb-target="#carouselVideoExample"
                data-mdb-slide="prev"
            >
                <span class="carousel-control-prev-icon" aria-hidden="true"></span>
                <span class="visually-hidden">Previous</span>
            </button>

            <button
                v-if="currentNews.videoscount > 1"
                class="carousel-control-next"
                type="button"
                data-mdb-target="#carouselVideoExample"
                data-mdb-slide="next"
            >
                <span class="carousel-control-next-icon" aria-hidden="true"></span>
                <span class="visually-hidden">Next</span>
            </button>
        </div>
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
    padding: 15px;
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

.img-item {
    cursor: pointer;
}

.btn-close {
    position: absolute;
    right: 16px;
    background-color: bisque;
}

.carousel-inner {
    object-fit: cover;
}

.video-self {
    object-fit: fill;
}

.modal-dialog {
    position: relative;
    display: table;
    /* This is important */
    overflow-y: auto;
    overflow-x: auto;
    width: auto;
}

@media (orientation: landscape) {
    .modal-body > img,
    .img-fluid {
        max-height: 90vh;
    }
}

@media (orientation: portrait) {
    .modal-body > img,
    .img-fluid {
        max-width: 90vw;
    }
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
