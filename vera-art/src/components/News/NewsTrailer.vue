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
    computed: {
        newsId() {
            return '/newsitem/' + this.newsObject.id;
        },
        bgImage() {
            return `url("${
                this.imagebasedir + this.newsObject.dir + this.newsObject.img_back
            }")`;
        },
    },
};
</script>

<template>
    <section class="container text-center px-0">
        <div class="col d-flex justify-content-center">
            <div class="news-item">
                <router-link :to="newsId">
                    <div class="img-holder"></div>
                    <div class="news-content">
                        <p>{{ newsObject.title }}</p>
                        <p>{{ newsObject.subTitle }}</p>
                        <p>
                            <CalendarIcon />
                            <span>&nbsp;{{ newsObject.datetimehuman }}</span>
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
    background-image: v-bind(bgImage);
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100%;
    transition: background-size 1s ease-out;
    margin-bottom: 1rem;
}

.news-item > a {
    text-decoration: none;
}

.img-holder:hover {
    background-size: 110%;
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

.news-content > p > span {
    font-size: small;
    font-weight: lighter;
}
</style>
