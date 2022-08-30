<script setup lang="ts">

import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import Slogan from "@/components/Slogan.vue"
import NewsPreview from "@/components/NewsPreview.vue";

</script>

<template>
    <Header />
    <Slogan />

    <section class="container text-center works">
        <div class="row gx-0">

            <div class="col d-flex justify-content-center">
                <NewsPreview />
            </div>
            <div class="col d-flex justify-content-center">
                <NewsPreview />
            </div>
            <div class="col d-flex justify-content-center">
                <NewsPreview />
            </div>
        </div>
    </section>
    <Footer />
</template>

<style scoped>
</style>
