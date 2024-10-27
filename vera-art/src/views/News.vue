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
                console.log('Error');
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
                console.log('Error');
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
            options
        );
        observer.observe(this.$refs.observer as Element);
    },
};
</script>

<template>
    <section class="container text-center main-content">
        <div class="row gx-0">
            <template v-for="newsObject in news" v-bind:key="newsObject.id">
                <div class="col d-flex justify-content-center mb-5">
                    <NewsTrailer :newsObject="newsObject" />
                </div>
            </template>
        </div>
        <div ref="observer" class="observer"></div>
    </section>
</template>

<style scoped>
.observer {
    height: 0px;
}
</style>
