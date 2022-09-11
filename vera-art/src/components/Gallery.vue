<script>
import axios from 'axios'
import Card from "./Card.vue"
import MySelect from "./MySelect.vue"
import { ref, onMounted, computed, toRefs, onBeforeMount } from 'vue'

export default {
    components: {
        Card,
        MySelect
    },
    props: {
        images: {
            type: Array,
            required: true
        }
    },
    setup() {
    },
    methods: {
    },
    mounted() {
    },
    computed: {
    }
}
</script>

<template>
    <section class="container text-center main-content">
        <div class="row">
            <Card v-for="imgObject in images" :imageObject="imgObject" :key="imgObject.id" />
        </div>
    </section>
</template>

<style scoped>

</style>
