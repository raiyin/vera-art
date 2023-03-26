<script>
import Carousel from '@/components/Carousel.vue';
import ButtonClose from '@/components/UI/ButtonClose.vue';

export default {
    components: {
        Carousel,
        ButtonClose
    },
    props: {
        imageObject: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {};
    },
    methods: {},
    computed: {
        imgIdModal() {
            return this.imageObject.id + 'Modal';
        },
    },
    watch: {},
};
</script>

<template>
    <div
         class="modal fade"
         :id="imgIdModal"
         tabindex="-1"
         aria-labelledby="exampleModalLabel"
         aria-hidden="true">

        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <ButtonClose />
                    <Carousel :imageObject="imageObject" :imageId="imageObject.id" />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
</style>
