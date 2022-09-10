<script setup lang="ts">

import Header from "../../components/Header.vue"
import Footer from "../../components/Footer.vue"
import Slogan from "../../components/Slogan.vue"
import Gallery from "@/components/Gallery.vue"

</script>

<template>

    <Header />
    <Slogan />
    <section class="container shadow bg-body rounded border-right border-left main-content">
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
                <Gallery gallaryDir="painting" />
            </div>
            <div class="tab-pane fade" id="Illustration" role="tabpanel" aria-labelledby="Illustration-tab">
                <Gallery gallaryDir="illustration" />
            </div>
            <div class="tab-pane fade" id="graphics3d" role="tabpanel" aria-labelledby="graphics3d-tab">
            </div>
        </div>

    </section>

    <Footer />
</template>

<style scoped>
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
