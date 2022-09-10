<script setup lang="ts">

import Card from "../components/Card.vue"
import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"
import Slogan from "../components/Slogan.vue"
import Gallery from "@/components/Gallery.vue"

</script>

<template>
    <Header />
    <Slogan />
    <Gallery gallaryDir="sale" />
    <Footer />
</template>

<style scoped>
</style>