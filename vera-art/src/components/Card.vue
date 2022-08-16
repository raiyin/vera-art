<script setup lang="ts">
import Modal from "./Modal.vue";
defineProps<{
    filename: string;
    name: string;
    base: string;
    material: string;
    size: string;
    year: string;
}>();
</script>

<template>
    <div class="card shadow-sm p-3 mb-5 bg-body rounded" style="width: 18rem;">
        <img :src=filename class="card-img-top" alt="..." data-bs-toggle="modal"
            data-bs-target="#exampleModal">
        <div class="card-body">
            <h5 class="card-title">{{name}}</h5>
            <p class="card-text">{{base}}, {{material}}, {{size}}, {{year}}</p>
            <a href="#" class="btn btn-primary">В корзину</a>
        </div>
        <Modal/>
    </div>
</template>

<style scoped>
</style>
