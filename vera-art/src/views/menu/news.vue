<script setup lang="ts">

import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import Slogan from "@/components/Slogan.vue"
import NewsTrailer from "@/components/News/NewsTrailer.vue";
import { ref, onMounted, computed, toRefs, onBeforeMount } from 'vue'
import axios from 'axios'

let news = ref()

async function fetchPosts() {
    try {
        const response = await axios.get('http://localhost:3001/news');
        news.value = response.data;
    }
    catch (e) {
        alert('Error')
    }
}

onBeforeMount(async () => {
    await fetchPosts()
})

</script>

<template>
    <Header />
    <Slogan />

    <section class="container text-center main-content">
        <div class="row gx-0">
            <template v-for="newsObject in news">
                <div class="col d-flex justify-content-center">
                    <NewsTrailer :newsObject="newsObject" />
                </div>
            </template>
        </div>
    </section>
    <Footer />
</template>

<style scoped>

</style>
