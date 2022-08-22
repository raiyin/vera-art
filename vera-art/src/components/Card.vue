<script setup lang="ts">
import Modal from "./Modal.vue";
import { ref, onMounted, computed, toRefs } from 'vue'

const props = defineProps<{
    jsonFile: string;
    filename: string;
}>();

const { jsonFile } = toRefs(props);
let name = ref()
let base = ref()
let material = ref()
let size = ref()
let year = ref()

const imgIdModalToLink = computed(() => {
    return "#" + jsonFile.value.split("/").slice(-2)[0] + "Modal"
})

onMounted(async () => {
    const res = await fetch(jsonFile.value)
    const desc = await res.json()
    name.value = desc.name_ru
    base.value = desc.base
        .replace("canvas on cardboard", "холст на картоне")
        .replace("canvas", "холст")
        .replace("paperboard", "картон")
        .replace("paper", "бумага")

    material.value = desc.material
        .replace("oil", "масло")
        .replace("watercolour", "акварель")
        .replace("texture paste", "текстурная паста")
        .replace("coal", "уголь")
        .replace("gouache", "гуашь")
        .replace("acrylic", "акрил")
    size.value = desc.size
    year.value = desc.year
})
</script>

<template>
    <div class="card shadow-sm p-3 mb-5 bg-body rounded" style="width: 18rem;">
        <img :src=filename class="card-img-top" alt="..." data-bs-toggle="modal" :data-bs-target=imgIdModalToLink>
        <div class="card-body">
            <h5 class="card-title">{{ name }}</h5>
            <p class="card-text">{{ base }}, {{ material }}, {{ size }}, {{ year }}</p>
            <a href="#" class="btn btn-primary">В корзину</a>
        </div>
        <Modal :jsonFile=jsonFile />
    </div>
</template>

<style scoped>
</style>
