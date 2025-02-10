<script lang="ts">
import type { NewsDesc } from '@/types';
import type { PropType } from 'vue';

export default {
    props: {
        imageObject: {
            type: Object as PropType<NewsDesc>,
            default: {} as NewsDesc,
        },
        selectedIndex: {
            type: Number,
            default: 1,
        },
        setSelectedIndex: {
            type: Function,
        },
    },
    data() {
        return {
            extension: '.jpg',
            imagebasedir: import.meta.env.VITE_IMAGE_DIR,
        };
    },
    methods: {
        makeFileName(dir: string, index: number) {
            return this.imagebasedir + dir + index + this.extension;
        },
        decreaseSelectedIndex() {
            if (this.selectedIndex - 1 === 0) {
                this.setSelectedIndex(this.imageObject.imagescount);
            } else {
                this.setSelectedIndex(this.selectedIndex - 1);
            }
        },
        increaseSelectedIndex() {
            if (this.selectedIndex + 1 > this.imageObject.imagescount) {
                this.setSelectedIndex(0);
            } else {
                this.setSelectedIndex(this.selectedIndex + 1);
            }
        },
    },
};
</script>

<template>
    <div
        id="news_img"
        class="carousel slide carousel-fade"
        data-bs-ride="false"
    >
        <div
            v-if="this.imageObject.imagescount > 1"
            class="carousel-indicators"
        >
            <button
                v-for="index in this.imageObject.imagescount"
                v-bind:key="index"
                type="button"
                data-bs-target="#news_img"
                :class="index === 1 ? 'active' : ''"
                :data-bs-slide-to="index - 1"
                :aria-current="index === 1"
            />
        </div>

        <div class="carousel-inner">
            <div
                v-for="index in this.imageObject.imagescount"
                v-bind:key="index"
                :class="
                    'carousel-item' + (index === selectedIndex ? ' active' : '')
                "
            >
                <img
                    :src="makeFileName(imageObject.dir, index)"
                    class="d-block modal-image"
                    alt="..."
                />
            </div>
        </div>

        <button
            v-if="this.imageObject.imagescount > 1"
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
            v-if="this.imageObject.imagescount > 1"
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
