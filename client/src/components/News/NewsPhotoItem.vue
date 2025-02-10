<script lang="ts">
import ButtonClose from '@/components/UI/ButtonClose.vue';
import type { PropType } from 'vue';
import type { NewsItemType } from '@/types';
import { inject } from 'vue';

export default {
    setup() {
        return {
            imagebasedir: inject('imagebasedir') as string,
        };
    },
    components: {
        ButtonClose,
    },
    props: {
        image_index: {
            type: Number,
            default: 0,
        },
        currentNews: {
            type: Object as PropType<NewsItemType>,
            default: {} as NewsItemType,
        },
    },
    data() {
        return {
            isLoaded: false as boolean,
        };
    },
    methods: {
        onImgLoad() {
            this.isLoaded = true;
        },
        makeImageName(index: number) {
            if (index === 0) {
                console.error('Wrong image_index props');
            }
            return this.imagebasedir + this.currentNews?.dir + index + '.jpg';
        },
        makeModalIdLink(index: number) {
            if (index === 0) {
                console.error('Wrong image_index props');
            }
            return '#exampleModal' + index;
        },
        makeModalId(index: number) {
            if (index === 0) {
                console.error('Wrong image_index props');
            }
            return 'exampleModal' + index;
        },
    },
};
</script>

<template>
    <div :class="[!isLoaded ? 'loading col-lg-4 col-md-6' : 'col-lg-4 col-md-6']">
        <div v-if="!isLoaded" class="image-stub" />

        <img
            v-show="isLoaded"
            v-bind:key="image_index"
            :src="makeImageName(image_index)"
            @load="onImgLoad"
            data-bs-target="#imgNewsModal"
            data-bs-toggle="modal"
            class="w-100 shadow-1-strong rounded mb-4 img-item"
        />
    </div>
</template>

<style scoped>
.image-stub {
    height: 200px;
    border-radius: 5px;
    margin-bottom: 10px;
}

.img-item {
    cursor: pointer;
}

.loading .image-stub {
    background-color: var(--skeleton-gray);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        var(--skeleton-gray);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 4px;
}

.modal__btn {
    position: absolute;
    right: clamp(5px, 1vw, 15px);
    top: clamp(5px, 1vw, 15px);
    padding: 0;
}

.modal-body {
    padding: clamp(5px, 1vw, 15px);
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}
</style>
