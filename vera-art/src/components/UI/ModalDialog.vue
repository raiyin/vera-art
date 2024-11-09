<script lang="ts">
import PicCarousel from '@/components/UI/PicCarousel.vue';
import ButtonClose from '@/components/UI/ButtonClose.vue';
import type { ImageProps } from '@/types';
import type { PropType } from 'vue';

export default {
    components: {
        Carousel: PicCarousel,
        ButtonClose,
    },
    props: {
        imageObject: {
            type: Object as PropType<ImageProps>,
            required: true,
        },
    },
    computed: {
        imgIdModal() {
            return this.imageObject.id + 'Modal';
        },
    },
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
                    <div class="modal__btn">
                        <ButtonClose />
                    </div>
                    <Carousel
                        :imageObject="imageObject"
                        :imageId="imageObject.id"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.modal-content {
    background-color: var(--color-surface);
}
.modal__btn {
    position: absolute;
    right: clamp(5px, 1vw, 15px);
    top: clamp(5px, 1vw, 15px);
    padding: 0;
    z-index: 100;
}
</style>
