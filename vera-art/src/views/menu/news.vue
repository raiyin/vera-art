<script setup lang="ts">

import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import Slogan from "@/components/Slogan.vue"
import NewsTrailer from "@/components/News/NewsTrailer.vue";

</script>

<template>
    <Header />
    <Slogan />

    <section class="container text-center main-content">
        <div class="row gx-0">

            <div class="col d-flex justify-content-center">
                <NewsTrailer />
            </div>
            <div class="col d-flex justify-content-center">
                <NewsTrailer />
            </div>
            <div class="col d-flex justify-content-center">
                <NewsTrailer />
            </div>
        </div>
    </section>
    <Footer />
</template>

<style scoped>
</style>
