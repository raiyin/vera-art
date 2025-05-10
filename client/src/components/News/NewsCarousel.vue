<script lang="ts">
import type { NewsDesc } from '../../types';
import type { PropType } from 'vue';

export default {
    props: {
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
    },
    data() {
        return {
            extension: '.jpg',
            imagebasedir: import.meta.env.VITE_IMAGE_DIR,
        };
    },
    methods: {
        makeFileName(dir: string, index: number): string {
            return `${this.imagebasedir}/${dir}/${index}${this.extension}`;
        },
        decreaseSelectedIndex(): void {
            const newIndex =
                this.selectedIndex === 1
                    ? this.imageObject.imagescount
                    : this.selectedIndex - 1;
            this.setSelectedIndex(newIndex);
        },
        increaseSelectedIndex() {
            const newIndex =
                this.selectedIndex === this.imageObject.imagescount
                    ? 1
                    : this.selectedIndex + 1;
            this.setSelectedIndex(newIndex);
        },
    },
};
</script>

<template>
    <div id="news_img" class="carousel slide carousel-fade" data-bs-ride="carousel">
        <div v-if="imageObject.imagescount > 1" class="carousel-indicators">
            <button
                v-for="i in imageObject.imagescount"
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
                v-for="(i, index) in imageObject.imagescount"
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

        <button
            v-if="imageObject.imagescount > 1"
            class="carousel-control-prev"
            type="button"
            data-bs-target="#news_img"
            data-bs-slide="prev"
            @click="decreaseSelectedIndex"
        >
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">
                {{ $t('carousel.back') }}
            </span>
        </button>

        <button
            v-if="imageObject.imagescount > 1"
            class="carousel-control-next"
            type="button"
            data-bs-target="#news_img"
            data-bs-slide="next"
            @click="increaseSelectedIndex"
        >
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">
                {{ $t('carousel.next') }}
            </span>
        </button>
    </div>
</template>
