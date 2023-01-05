<script>
import Carousel from '@/components/Carousel.vue';

export default {
    components: {
        Carousel,
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
        aria-hidden="true"
    >
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <button
                        type="button"
                        class="btn-close"
                        data-bs-dismiss="modal"
                        aria-label="Close"
                    ></button>
                    <Carousel :imageObject="imageObject" :imageId="imageObject.id" />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.btn-close {
    position: absolute;
    right: 16px;
    background-color: bisque;
    z-index: 100;
}

.modal-dialog {
    position: relative;
    display: table;
    /* This is important */
    overflow-y: auto;
    overflow-x: auto;
    width: auto;
}
</style>
