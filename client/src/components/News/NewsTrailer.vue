<script lang="ts">
import CalendarIcon from '@/components/Icons/IconCalendar.vue';

export default {
    components: {
        CalendarIcon,
    },
    props: {
        newsObject: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            isLoaded: false,
            imagebasedir: import.meta.env.VITE_IMAGE_DIR,
        };
    },
    methods: {
        onImgLoad() {
            this.isLoaded = true;
        },
        getHumanDate(inDate: string, locale: string) {
            const options = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            } as const;
            const date = new Date(inDate);
            const stdLocale = locale === 'RUS' ? 'ru-RU' : 'en-EN';
            return date.toLocaleDateString(stdLocale, options);
        },
    },
    computed: {
        newsId() {
            return '/news/' + this.newsObject.id;
        },
        bgImage() {
            return (
                this.imagebasedir +
                this.newsObject.dir +
                this.newsObject.img_back
            );
        },
    },
};
</script>

<template>
    <section class="container text-center px-0">
        <div class="col d-flex justify-content-center">
            <div class="news-item" :class="{ loading: isLoaded }">
                <router-link :to="newsId">
                    <div class="img-holder">
                        <img
                            :src="bgImage"
                            @load="onImgLoad"
                            v-show="isLoaded"
                        />
                        <div v-show="!isLoaded" class="image-stub" />
                    </div>

                    <div class="news-content">
                        <div>
                            {{
                                !isLoaded
                                    ? ''
                                    : $i18n.locale === 'RUS'
                                      ? newsObject.title_ru
                                      : newsObject.title_en
                            }}
                        </div>
                        <div>
                            {{
                                !isLoaded
                                    ? ''
                                    : $i18n.locale === 'RUS'
                                      ? newsObject.subTitle_ru
                                      : newsObject.subTitle_en
                            }}
                        </div>
                        <div>
                            <CalendarIcon v-if="isLoaded" />
                            <span>
                                &nbsp;{{
                                    !isLoaded
                                        ? ''
                                        : getHumanDate(
                                              newsObject.datetime,
                                              $i18n.locale,
                                          )
                                }}
                            </span>
                        </div>
                    </div>
                </router-link>
            </div>
        </div>
    </section>
</template>

<style scoped>
.news-item {
    width: 25rem;
    height: 30rem;
    border-radius: 0.3rem;
    border: 0.1rem solid var(--color-surface-secondary);
    overflow: hidden;
}

.img-holder {
    width: 100%;
    height: 79%;
    margin-bottom: 1rem;
    overflow: hidden;
}

.img-holder img {
    height: 100%;
    width: 100%;
    transition: transform 0.2s;
}

.img-holder img:hover {
    transform: scale(1.1);
}

.image-stub {
    height: 100%;
    width: 100%;
}

.loading .img-holder,
.loading .news-content {
    background-color: var(--skeleton-gray);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        var(--skeleton-gray);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 0.4rem;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}

.loading .news-content {
    height: 3rem;
    width: 80%;
    margin: auto;
}

.news-item > a {
    text-decoration: none;
}

.img-holder:hover {
    cursor: pointer;
}

.news-content {
    position: relative;
    font-weight: bold;
    font-size: larger;
    text-align: left;
    color: var(--color-on-surface);
    letter-spacing: 0.05rem;
    line-height: 1.4;
    padding-left: 1rem;
}

.news-content > p > span {
    font-size: small;
    font-weight: lighter;
}
</style>
