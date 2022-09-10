<script setup lang="ts">

import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import Slogan from "@/components/Slogan.vue"
import SideNewsItem from "@/components/News/SideNewsTrailer.vue";
import NewsItemDescription from "@/components/News/NewsItemDescription.vue"
import { onBeforeMount, onMounted, ref } from "vue";
import { useRoute, useRouter } from 'vue-router'
import axios from "axios";

let newsArray = ref()
let newsObject = ref()
const route = useRoute()
const router = useRouter()

async function fetchData(path: String) {
    try {
        let newsid = path.substring(path.lastIndexOf('/') + 1)
        var response = await axios.get('http://localhost:3001/news', { params: { id: newsid } });
        newsObject.value = response.data[0];

        response = await axios.get('http://localhost:3001/news', { params: { id_ne: newsid } });
        newsArray.value = response.data;

        console.log("2" + Date.now);
    }
    catch (e) {
        alert('Error')
    }
}

onBeforeMount(async () => {
    await router.isReady()
    await fetchData(route.path)
    console.log("1" + Date.now);    
})

onMounted(async () => {
    console.log("3" + Date.now());
})
</script>
    
<template>

    <Header />
    <Slogan />

    <section class="container text-center px-0 main-content">
        <article>
            <div class="news-header">
                <div class="news_img">
                    <img :src="newsObject.img_backfull" />
                </div>
                <div class="other-news">
                    <template v-for="index in 2" :key="index">
                        <SideNewsItem :sideNewsObject="newsArray[index-1]" />
                    </template>
                </div>
            </div>
            <NewsItemDescription :newsObject="newsObject" />
            <div class="news-text">
                <p>{{newsObject.text}}</p>
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
    