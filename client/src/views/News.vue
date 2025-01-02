<script lang="ts">
import NewsTrailer from '@/components/News/NewsTrailer.vue';
import type { NewsDesc } from '@/types';
import axios from 'axios';

export default {
    inject: ['jsonserverhost'],
    components: {
        NewsTrailer,
    },
    props: [],
    data() {
        return {
            news: [] as NewsDesc[],
            page: 1,
            limit: 9,
        };
    },
    methods: {
        async loadNews() {
            try {
                const response = await axios.get(this.jsonserverhost + 'news', {
                    params: {
                        _page: this.page,
                        _limit: this.limit,
                    },
                });
                this.news = response.data;
            } catch (e) {
                console.error('Error fetching news');
            }
        },
        async loadMoreNews() {
            try {
                this.page += 1;
                const response = await axios.get(this.jsonserverhost + 'news', {
                    params: {
                        _page: this.page,
                        _limit: this.limit,
                    },
                });
                this.news = [...this.news, ...response.data];
            } catch (e) {
                console.error('Error fetching news');
            }
        },
    },
    mounted() {
        this.loadNews();
        const options = {
            rootMargin: '0px',
            threshold: 1.0,
        };
        const callback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.loadMoreNews();
            }
        };
        const observer: IntersectionObserver = new IntersectionObserver(
            callback,
            options,
        );
        observer.observe(this.$refs.observer as Element);
    },
};
</script>

<template>
    <section class="container text-center main-content container__news">
        <div class="row row__news">
            <div
                class="col d-flex justify-content-center mb-5"
                v-for="newsObject in news"
                v-bind:key="newsObject.id"
            >
                <NewsTrailer :newsObject="newsObject" />
            </div>
        </div>
        <div ref="observer" class="observer"></div>
    </section>
</template>

<style scoped>
.container__news {
    padding-right: 0;
    padding-left: 0;
}
.observer {
    height: 0px;
}
.row__news {
    justify-content: space-between;
}
</style>
