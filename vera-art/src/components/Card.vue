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
let price = ref()

const imgIdModalToLink = computed(() => {
    return "#" + jsonFile.value.split("/").slice(-2)[0] + "Modal"
})

onMounted(async () => {
    const res = await fetch(jsonFile.value)
    const desc = await res.json()
    name.value = desc.name_ru
    base.value = desc.base

    material.value = desc.material
    size.value = desc.size
    year.value = desc.year
    price.value = desc.price
})
</script>

<template>
    <div class="card shadow-lg p-3 mb-5 bg-body rounded" style="width: 20rem;">
        <img :src=filename class="card-img-top" alt="..." data-bs-toggle="modal" :data-bs-target=imgIdModalToLink>
        <div class="card-body">
            <div class="desc">
                <h5 class="card-title">{{ name }}</h5>
                <p class="card-text">
                    <span v-if="base">{{ base }}</span>
                    <span v-if="material">, {{ material }}</span>
                    <span v-if="size">, {{ size }}</span>
                    <span v-if="year">, {{ year }}</span>
                </p>
                <p v-if="price">Цена: {{ price }} р.</p>
                <a v-if="price" href="#" class="btn btn-primary">В корзину</a>
            </div>
        </div>
        <Modal :jsonFile=jsonFile />
    </div>
</template>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
}

.desc {
    /* display: flex; */
    /* flex-direction: column; */
    align-self: flex-end;
    /* justify-content: center; */
    width: 100%;
}

.card>img:hover{
    cursor: pointer;
}
</style>
