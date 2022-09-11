<script setup lang="ts">
import Modal from "./Modal.vue";
import { ref, onMounted, computed, toRefs } from 'vue'

const props = defineProps<{
    imageObject: any
}>();

var imageObject = props.imageObject;
const imgIdModalToLink = computed(() => {
    return "#" + imageObject.id + "Modal"
})

const mainCardImage = computed(() => {
    return imageObject.dir + '1.jpg'
})

</script>

<template>
    <div class="col d-flex justify-content-center">
        <div class="card shadow-lg p-3 mb-5 bg-body rounded" style="width: 20rem;">
            <img :src=mainCardImage class="card-img-top" alt="..." data-bs-toggle="modal"
                :data-bs-target=imgIdModalToLink>
            <div class="card-body">
                <div class="desc">
                    <h5 class="card-title">{{ imageObject.name_ru }}</h5>
                    <p class="card-text">
                        <span v-if="imageObject.base">{{ imageObject.base }}</span>
                        <span v-if="imageObject.material">, {{ imageObject.material }}</span>
                        <span v-if="imageObject.size">, {{ imageObject.size }}</span>
                        <span v-if="imageObject.year">, {{ imageObject.year }}</span>
                    </p>
                    <p v-if="imageObject.price">Цена: {{ imageObject.price }} р.</p>
                    <a v-if="imageObject.price" href="#" class="btn btn-primary">В корзину</a>
                </div>
            </div>
            <Modal :imageObject=imageObject />
        </div>
    </div>
</template>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
}

.desc {
    align-self: flex-end;
    width: 100%;
}

.card>img:hover {
    cursor: pointer;
}
</style>
