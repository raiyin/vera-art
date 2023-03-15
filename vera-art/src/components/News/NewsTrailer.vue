<script>
import CalendarIcon from '@/components/icons/IconCalendar.vue';

export default {
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
            loadingGrey: '#ededed',
        };
    },
    methods: {
        onImgLoad() {
            this.isLoaded = true;
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
            <div class="news-item" :class="[!this.isLoaded ? 'loading' : '',]">
                <router-link :to="newsId">
                    <div class="img-holder">

                        <img
                             :src="bgImage"
                             @load="onImgLoad"
                             v-show="this.isLoaded" />
                        <div v-show="!this.isLoaded" class="image-stub" />

                    </div>
                    <div class="news-content">
                        <p>{{ !this.isLoaded ? "" : newsObject.title }}</p>
                        <p>{{ !this.isLoaded ? "" : newsObject.subTitle }}</p>
                        <p>
                            <CalendarIcon v-if="this.isLoaded" />
                            <span>&nbsp;{{ !this.isLoaded ? "" : newsObject.datetimehuman }}</span>
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
    border: 1px solid wheat;
    overflow: hidden;
}

.img-holder {
    width: 100%;
    height: 79%;
    margin-bottom: 1rem;
    overflow: hidden
}

.img-holder img {
    height: 100%;
    width: 100%;
    transition: transform .2s;
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
    background: linear-gradient(100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, .5) 50%,
            rgba(255, 255, 255, 0) 60%) v-bind(loadingGrey);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
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

.news-item>a {
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
    color: black;
    letter-spacing: 0.05rem;
}

.news-content>p>span {
    font-size: small;
    font-weight: lighter;
}
</style>
