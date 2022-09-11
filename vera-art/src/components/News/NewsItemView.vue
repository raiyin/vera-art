<script>

import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import Slogan from "@/components/Slogan.vue"
import SideNewsTrailer from "@/components/News/SideNewsTrailer.vue";
import NewsItemDescription from "@/components/News/NewsItemDescription.vue"
import { useRouter, useRoute } from 'vue-router'
import axios from "axios";


export default {
    components: {
        Header,
        Footer, 
        Slogan, 
        SideNewsTrailer, 
        NewsItemDescription
    },
    data() {
        return {
            news: [],
            currentNews: {}
        }
    },
    setup() {
        const route = useRoute();
        const router = useRouter()
    },
    methods: {
        async fetchData() {
            try {
                let newsid = this.$route.path.substring(this.$route.path.lastIndexOf('/') + 1)
                var response = await axios.get('http://localhost:3001/news', { params: { id: newsid } });
                this.currentNews = response.data[0];

                response = await axios.get('http://localhost:3001/news', { params: { id_ne: newsid } });
                this.news = response.data;
            }
            catch (e) {
                console.log(e);
                var propertyNames = Object.getOwnPropertyNames(e);
                propertyNames.forEach(function (property) {
                    var descriptor = Object.getOwnPropertyDescriptor(e, property);
                    console.log(property + ":" + e[property] + ":" + propsToStr(descriptor));
                });
            }
        },
    },
    mounted() {
        this.fetchData();
    },
    async onBeforeMount() {
        await router.isReady()
    }
}

</script>
    
<template>

    <Header />
    <Slogan />

    <section class="container text-center px-0 main-content">
        <article>
            <div class="news-header">
                <div class="news_img">
                    <img :src="currentNews.img_backfull" />
                </div>
                <div class="other-news">
                    <template v-for="newsItem in news">
                        <SideNewsTrailer :sideNewsObject="newsItem" />
                    </template>
                </div>
            </div>
            <NewsItemDescription :newsObject="currentNews" />
            <div class="news-text">
                <p>{{currentNews.text}}</p>
            </div>
        </article>
    </section>
    <Footer />
</template>
    
<style scoped>
.news-header {
    display: flex;
}

.other-news {
    display: flex;
    flex-direction: column;
}

.news_img {
    display: flex;
    flex-direction: row;
    justify-content: left;
    margin-right: 2rem;
}

.news-text {
    text-align: left;
    margin-bottom: 2rem;
}
</style>
    