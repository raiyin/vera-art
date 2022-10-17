<script>

import NewsTrailer from "@/components/News/NewsTrailer.vue";
import { ref, onMounted, computed, toRefs, onBeforeMount } from 'vue'
import axios from 'axios'

export default {
    inject: ["host"],
    components: {
        NewsTrailer,
    },
    props: [],
    data() {
        return {
            news: [],
        }
    },
    setup() {
    },
    methods: {
        async fetchPosts() {
            try {
                const response = await axios.get('http://' + this.host + ':3001/news');
                this.news = response.data;
            }
            catch (e) {
                console.log('Error')
            }
        }
    },
    mounted() {
        this.fetchPosts();
    },
    onBeforeMount() {
    },
    computed: {
    },
    watch: {
    },
}
</script>

<template>
    <section class="container text-center main-content">
        <div class="row gx-0">
            <template v-for="newsObject in news">
                <div class="col d-flex justify-content-center mb-5">
                    <NewsTrailer :newsObject="newsObject" />
                </div>
            </template>
        </div>
    </section>
</template>

<style scoped>

</style>
