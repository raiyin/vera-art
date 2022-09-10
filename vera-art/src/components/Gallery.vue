<script setup lang="ts">
import axios from 'axios'
import Card from "./Card.vue"
import { ref, onMounted, computed, toRefs, onBeforeMount } from 'vue'

const props = defineProps<{
    gallaryDir: string;
}>();


var { gallaryDir } = toRefs(props);
let images = ref()

async function fetchPosts() {
    try {
        const response = await axios.get('http://localhost:3001/' + gallaryDir.value);
        images.value = response.data;
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
        <div class="row">
            <template v-for="imgObject in images">
                <div class="col d-flex justify-content-center">
                    <Card :imageObject="imgObject" />
                </div>
            </template>
        </div>
    </section>
</template>

<style scoped>

</style>
