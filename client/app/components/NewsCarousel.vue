<script setup lang="ts">
import type { NewsDesc } from '../types';
import type { PropType } from 'vue';
import { ref } from 'vue';
import { useI18n } from '#imports';

const props = defineProps({
    imageObject: {
        type: Object as PropType<NewsDesc>,
        required: true,
        default: {} as NewsDesc,
    },
    selectedIndex: {
        type: Number,
        validator: (value: number) => value >= 1,
        default: 1,
    },
    setSelectedIndex: {
        type: Function,
        required: true,
    },
});

const { t } = useI18n();

// Reactive state
const extension = ref('.jpg');

// Methods
const makeFileName = (dir: string, index: number): string => {
    return `${dir}/${index}${extension.value}`;
};

const decreaseSelectedIndex = (): void => {
    const newIndex =
        props.selectedIndex === 1
            ? props.imageObject.images.length
            : props.selectedIndex - 1;
    props.setSelectedIndex(newIndex);
};

const increaseSelectedIndex = () => {
    const newIndex =
        props.selectedIndex === props.imageObject.images.length
            ? 1
            : props.selectedIndex + 1;
    props.setSelectedIndex(newIndex);
};
</script>

<template>
    <div id="news_img" class="carousel slide" data-bs-ride="carousel">
        <div v-if="imageObject.images.length > 1" class="carousel-indicators">
            <UButton
                v-for="i in imageObject.images.length"
                :key="i"
                type="button"
                data-bs-target="#news_img"
                :class="{ active: i === selectedIndex + 1 }"
                :data-bs-slide-to="i - 1"
                :aria-current="i === selectedIndex + 1"
            />
        </div>

        <div class="carousel-inner">
            <div
                v-for="(i, index) in imageObject.images.length"
                :key="index + 1"
                class="carousel-item"
                :class="{ active: index === selectedIndex }"
            >
                <img
                    :src="makeFileName(imageObject.dir, index + 1)"
                    class="d-block modal-image"
                    alt="..."
                />
            </div>
        </div>

        <UButton
            v-if="imageObject.images.length > 1"
            class="carousel-control-prev"
            type="button"
            data-bs-target="#news_img"
            data-bs-slide="prev"
            @click="decreaseSelectedIndex"
        >
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">
                {{ t('carousel.back') }}
            </span>
        </UButton>

        <UButton
            v-if="imageObject.images.length > 1"
            class="carousel-control-next"
            type="button"
            data-bs-target="#news_img"
            data-bs-slide="next"
            @click="increaseSelectedIndex"
        >
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">
                {{ t('carousel.next') }}
            </span>
        </UButton>
    </div>
</template>
