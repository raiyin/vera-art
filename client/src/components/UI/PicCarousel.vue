<script lang="ts">
import type { ImageParts, ImageProps } from '@/types';
import type { PropType } from 'vue';

export default {
    inject: ['imagebasedir'],
    props: {
        imageObject: {
            type: Object as PropType<ImageProps>,
            default: {} as ImageProps,
        },
        imageId: {
            type: String,
        },
    },
    data() {
        return {
            parts: [] as ImageParts[],
            imgCountGTOne: false,
            extension: '.jpg',
        };
    },
    methods: {
        makeFileName(dir: string, index: number) {
            return this.imagebasedir + dir + index + this.extension;
        },
    },
    computed: {
        imgIdtoLink() {
            return '#' + this.imageObject.id;
        },
    },
    beforeMount() {
        this.parts = this.imageObject.parts;
        this.imgCountGTOne = this.parts.length > 1;
    },
};
</script>

<template>
    <div :id="imageId" class="carousel slide" data-bs-ride="false">
        <div v-if="imgCountGTOne" class="carousel-indicators">
            <button
                v-for="(value, index) in parts"
                v-bind:key="index"
                type="button"
                :data-bs-target="imgIdtoLink"
                :class="index === 0 ? 'active' : null"
                :data-bs-slide-to="index"
                :aria-current="index === 0 ? true : false"
                :aria-label="value.name_ru"
            />
        </div>

        <div class="carousel-inner">
            <template v-for="(value, index) in parts" v-bind:key="index">
                <div
                    :class="
                        index === 0 ? 'carousel-item active' : 'carousel-item'
                    "
                >
                    <img
                        :src="makeFileName(imageObject.dir, index + 1)"
                        class="d-block modal-image"
                        alt="..."
                    />
                    <div class="carousel-caption d-none d-md-block">
                        <h5>{{ value.name_ru }}</h5>
                    </div>
                </div>
            </template>
        </div>

        <button
            v-if="imgCountGTOne"
            class="carousel-control-prev"
            type="button"
            :data-bs-target="imgIdtoLink"
            data-bs-slide="prev"
        >
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">
                {{ $t('carousel.back') }}
            </span>
        </button>

        <button
            v-if="imgCountGTOne"
            class="carousel-control-next"
            type="button"
            :data-bs-target="imgIdtoLink"
            data-bs-slide="next"
        >
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">
                {{ $t('carousel.next') }}
            </span>
        </button>
    </div>
</template>
