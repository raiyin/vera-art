<script setup lang="ts">

import Card from "../components/Card.vue"
import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"
import Slogan from "../components/Slogan.vue"

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

    <section class="container text-center works">
        <div class="row">
            <template v-for="value, index in dirs">
                <div class="col d-flex justify-content-center">
                    <Card :jsonFile="makeFullFileName(value, 'desc.json')" :filename="makeFullFileName(value, '1.jpg')" />
                </div>
            </template>
        </div>
    </section>

    <Footer />
</template>

<style scoped>
.works {
    margin-top: 10vh;
}
</style>