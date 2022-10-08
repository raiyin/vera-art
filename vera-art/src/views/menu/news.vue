<script setup lang="ts">

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
