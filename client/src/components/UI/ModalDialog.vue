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
        modalId: {
            type: String,
            required: true,
        },
    },
};
</script>

<template>
    <div
        class="modal fade"
        :id="modalId"
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
                    <slot></slot>
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
    z-index: 100;
    box-sizing: border-box;
    right: var(--bs-modal-padding);
    display: flex;
}
</style>
