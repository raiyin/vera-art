<script lang="ts">
import ButtonClose from '@/components/UI/ButtonClose.vue';

export default {
    inject: ['imagebasedir'],
    components: {
        ButtonClose,
    },
    props: {
        image_index: {
            type: Number,
        },
        currentNews: {
            type: Object,
        },
    },
    data() {
        return {
            isLoaded: false,
            loadingGrey: '#adadad',
        };
    },
    methods: {
        onImgLoad() {
            this.isLoaded = true;
        },
        makeImageName(index: number) {
            return this.imagebasedir + this.currentNews.dir + index + '.jpg';
        },
        makeModalIdLink(index: number) {
            return '#exampleModal' + index;
        },
        makeModalId(index: number) {
            return 'exampleModal' + index;
        },
    },
};
</script>

<template>
    <div :class="[!this.isLoaded ? 'loading' : '']">
        <div v-if="!this.isLoaded" class="image-stub" />

        <img
            v-show="this.isLoaded"
            v-bind:key="image_index"
            :src="makeImageName(image_index)"
            :data-bs-target="makeModalIdLink(image_index)"
            @load="onImgLoad"
            class="w-100 shadow-1-strong rounded mb-4 img-item"
            data-bs-toggle="modal"
        />

        <div
            v-bind:key="image_index"
            class="modal fade"
            :id="makeModalId(image_index)"
            tabindex="-1"
            aria-labelledby="exampleModalLabel"
            aria-hidden="true"
        >
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-body">
                        <ButtonClose />
                        <img class="modal-image" :src="makeImageName(image_index)" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.image-stub {
    width: 100%;
    height: 200px;
    border-radius: 5px;
    margin-bottom: 10px;
}

.loading .image-stub {
    background-color: v-bind(loadingGrey);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        v-bind(loadingGrey);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 4px;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}
</style>
