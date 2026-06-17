<script setup lang="ts">
import type { PropType } from 'vue';
import type { NewsDesc } from '~/types';
import { ref } from 'vue';

const props = defineProps({
    image_index: {
        type: Number,
        default: 0,
    },
    currentNews: {
        type: Object as PropType<NewsDesc>,
        default: {} as NewsDesc,
    },
    fileName: {
        type: String,
        default: '',
    },
});

const isLoaded = ref(false);

const onImgLoad = () => {
    isLoaded.value = true;
};

const makeImageName = (fileName: string) => {
    return `${props.currentNews?.dir}${fileName}`;
};

const makeModalIdLink = (index: number) => {
    if (index === 0) {
        console.error('Wrong image_index props');
    }
    return '#exampleModal' + index;
};

const makeModalId = (index: number) => {
    if (index === 0) {
        console.error('Wrong image_index props');
    }
    return 'exampleModal' + index;
};
</script>

<template>
    <div :class="{ loading: !isLoaded }">
        <div v-if="!isLoaded" class="image-stub" />

        <img
            v-show="isLoaded"
            v-bind:key="image_index"
            :src="makeImageName(fileName)"
            @load="onImgLoad"
            data-bs-target="#imgNewsModal"
            data-bs-toggle="modal"
            class="w-100 shadow-1-strong rounded mb-4 img-item"
        />
    </div>
</template>

<style scoped>
.image-stub {
    height: 20rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
}

.img-item {
    cursor: pointer;
    object-fit: cover;
    max-width: 100%;
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
    border-radius: 0.4rem;
}

.modal__btn {
    position: absolute;
    right: clamp(0.5rem, 1vw, 1.5rem);
    top: clamp(0.5rem, 1vw, 1.5rem);
    padding: 0;
}

.modal-body {
    padding: clamp(0.5rem, 1vw, 1.5rem);
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}
</style>
