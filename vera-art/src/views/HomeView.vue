<script setup lang="ts">

import Card from "../components/Card.vue"
import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"
import Slogan from "../components/Slogan.vue"
import Gallery from "@/components/Gallery.vue"

import { ref, onMounted, computed } from 'vue'

const baseDir = "/src/assets/img/sale/"
let dirs = ref()

function makeFullFileName(dir: string, fileName: string) {
    return baseDir + dir + "/" + fileName
}

onMounted(async () => {
    const res = await fetch("../../src/assets/img/sale/total.json")
    const total = await res.json()
    dirs.value = total.dirs
})
</script>

<template>
    <Header />
    <Slogan />
    <Gallery gallaryDir="../../src/assets/img/sale/" />
    <Footer />
</template>

<style scoped>
</style>