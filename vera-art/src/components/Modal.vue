<script setup lang="ts">
import Carousel from "./Carousel.vue"
import { ref, onMounted, computed, toRefs } from 'vue'

const props = defineProps<{
    imageObject: any
}>();

var imageObject = props.imageObject;

const imgIdModal = computed(() => {
    return imageObject.id + "Modal"
})

</script>

<template>

    <div class="modal fade" :id=imgIdModal tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <Carousel :imageObject=imageObject :imageId=imageObject.id />
                </div>
            </div>
        </div>
    </div>

</template>

<style scoped>
.modal-dialog {
    position: relative;
    display: table;
    /* This is important */
    overflow-y: auto;
    overflow-x: auto;
    width: auto;
}
</style>
