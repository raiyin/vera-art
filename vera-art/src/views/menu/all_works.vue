<script setup lang="ts">

import Card from "../../components/Card.vue"
import Header from "../../components/Header.vue"
import Footer from "../../components/Footer.vue"
import Slogan from "../../components/Slogan.vue"

import { ref, onMounted, computed } from 'vue'

const extension = ".jpg"
const baseDir = "/src/assets/img/allworks/painting/"
let dirs = ref()

function makeFullFileName(dir: string, fileName: string) {
    return baseDir + dir + "/" + fileName
}

onMounted(async () => {
    const res = await fetch("../../src/assets/img/allworks/painting/total.json")
    const total = await res.json()
    dirs.value = total.dirs
})
</script>

<template>

    <Header />
    <Slogan />


    <section class="container shadow bg-body rounded border-right border-left">
        <ul class="nav nav-tabs" id="myTab" role="tablist">
            <li class="nav-item" role="presentation">
                <button class="nav-link active" id="painting-tab" data-bs-toggle="tab" data-bs-target="#painting"
                    type="button" role="tab" aria-controls="painting" aria-selected="true">ЖИВОПИСЬ</button>
            </li>
            <li class="nav-item" role="presentation">
                <button class="nav-link" id="Illustration-tab" data-bs-toggle="tab" data-bs-target="#Illustration"
                    type="button" role="tab" aria-controls="Illustration" aria-selected="false">ИЛЛЮСТРАЦИЯ</button>
            </li>
            <li class="nav-item" role="presentation">
                <button class="nav-link" id="graphics3d-tab" data-bs-toggle="tab" data-bs-target="#graphics3d"
                    type="button" role="tab" aria-controls="graphics3d" aria-selected="false">3D-ГРАФИКА</button>
            </li>
        </ul>
        <div class="tab-content border-right border-left" id="myTabContent">
            <div class="tab-pane fade show active" id="painting" role="tabpanel" aria-labelledby="painting-tab">
                <section class="container text-center works">
                    <div class="row">
                        <template v-for="value, index in dirs">
                            <div class="col d-flex justify-content-center">
                                <Card :jsonFile="makeFullFileName(value, 'desc.json')"
                                    :filename="makeFullFileName(value, '1.jpg')" />
                            </div>
                        </template>
                    </div>
                </section>
            </div>
            <div class="tab-pane fade" id="Illustration" role="tabpanel" aria-labelledby="Illustration-tab">

            </div>
            <div class="tab-pane fade" id="graphics3d" role="tabpanel" aria-labelledby="graphics3d-tab">

            </div>
        </div>

    </section>

    <Footer />
</template>

<style scoped>
.works {
    margin-top: 10vh;
}

.tab-content {
    border-left: 1px solid #ddd;
    border-right: 1px solid #ddd;
    padding: 10px;
}

.nav-tabs {
    margin-bottom: 0;
}

.nav-tabs .nav-item .nav-link {
  color: black;
}

.nav-tabs .nav-item .nav-link.active {
  color: black;
}

.nav-item {
    color: black;
}
</style>
