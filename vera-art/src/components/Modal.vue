<script setup lang="ts">
import Carousel from "./Carousel.vue"
import { ref, onMounted, computed, toRefs } from 'vue'

const props = defineProps<{
    jsonFile: string;
}>();

const { jsonFile } = toRefs(props);

const imgId = computed(() => {
    return jsonFile.value.split("/").slice(-2)[0]
})

const imgIdModal = computed(() => {
    return jsonFile.value.split("/").slice(-2)[0] + "Modal"
})

</script>

<template>

    <div class="modal fade" :id=imgIdModal tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <Carousel :jsonFile=jsonFile :imageId=imgId />
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
