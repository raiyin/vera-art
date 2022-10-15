<script>

import Gallery from "@/components/Gallery.vue"
import axios from 'axios'

export default {
    components: {
        Gallery,
    },
    props: [],
    data() {
        return {
            paintingImages: [],
            illustrationImages: [],
            threeDImages: [],
        }
    },
    setup() {
    },
    methods: {
        async fetchData() {
            try {
                var response = await axios.get('http://localhost:3001/painting');
                this.paintingImages = response.data;
                response = await axios.get('http://localhost:3001/illustration');
                this.illustrationImages = response.data;
                response = await axios.get('http://localhost:3001/3d');
                this.threeDImages = response.data;
            }
            catch (e) {
                alert('Error')
                console.log(e);
                var propertyNames = Object.getOwnPropertyNames(e);
                propertyNames.forEach(function (property) {
                    var descriptor = Object.getOwnPropertyDescriptor(e, property);
                    console.log(property + ":" + e[property] + ":" + propsToStr(descriptor));
                });
            }
        }
    },
    mounted() {
        this.fetchData();
    },
    computed: {
    },
    watch: {
    },
}
</script>

<template>

    <section class="container bg-body rounded border-right border-left main-content">
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
                <Gallery :images="paintingImages" />
            </div>
            <div class="tab-pane fade" id="Illustration" role="tabpanel" aria-labelledby="Illustration-tab">
                <Gallery :images="illustrationImages" />
            </div>
            <div class="tab-pane fade" id="graphics3d" role="tabpanel" aria-labelledby="graphics3d-tab">
                <Gallery :images="threeDImages" />
            </div>
        </div>

    </section>
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
