<script>
import CalendarIcon from '@/components/Icons/IconCalendar.vue';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
    },
    inject: ['imagebasedir'],
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
            loadingGrey: '#adadad',
        };
    },
    methods: {
        onImgLoad() {
            this.isLoaded = true;
        },
        getHumanDate(inDate, locale) {
            const options = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            };
            const date = new Date(inDate);
            let stdLocale = locale === 'RUS' ? 'ru-RU' : 'en-EN';
            return date.toLocaleDateString(stdLocale, options);
        },
    },
    computed: {
        newsId() {
            return '/newsitem/' + this.newsObject.id;
        },
        bgImage() {
            return this.imagebasedir + this.newsObject.dir + this.newsObject.img_back;
        },
    },
};
</script>

<template>
    <section class="container text-center px-0">
        <div class="col d-flex justify-content-center">
            <div class="news-item" :class="[!this.isLoaded ? 'loading' : '']">
                <router-link :to="newsId">
                    <div class="img-holder">
                        <img :src="bgImage" @load="onImgLoad" v-show="this.isLoaded" />
                        <div v-show="!this.isLoaded" class="image-stub" />
                    </div>
                    <div class="news-content">
                        <p>
                            {{
                                !this.isLoaded
                                    ? ''
                                    : $i18n.locale === 'RUS'
                                    ? newsObject.title_ru
                                    : newsObject.title_en
                            }}
                        </p>
                        <p>
                            {{
                                !this.isLoaded
                                    ? ''
                                    : $i18n.locale === 'RUS'
                                    ? newsObject.subTitle_ru
                                    : newsObject.subTitle_en
                            }}
                        </p>
                        <p>
                            <CalendarIcon v-if="this.isLoaded" />
                            <span>
                                &nbsp;{{
                                    !this.isLoaded
                                        ? ''
                                        : getHumanDate(newsObject.datetime, $i18n.locale)
                                }}
                            </span>
                        </p>
                    </div>
                </router-link>
            </div>
        </div>
    </section>
</template>

<style scoped>
.news-item {
    width: 22rem;
    height: 27rem;
    border-radius: 0.3rem;
    border: 1px solid var(--color-surface-secondary);
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
    background-color: v-bind(loadingGrey);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        v-bind(loadingGrey);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 4px;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}

.loading .news-content {
    height: 30px;
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
    font-family: Ubuntu;
    font-weight: bold;
    font-size: larger;
    text-align: left;
    text-indent: 1rem;
    line-height: 0.4rem;
    color: var(--color-on-surface);
    letter-spacing: 0.05rem;
}

.news-content > p > span {
    font-size: small;
    font-weight: lighter;
}
</style>
