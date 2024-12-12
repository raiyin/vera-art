<script lang="ts">
import type { ImageProps } from '@/types';
import ModalDialog from './ModalDialog.vue';
import { useThemeStore } from '../../stores/ThemeStore';
import type { PropType } from 'vue';
import PictureCardSkeleton from '../UI/Skeletons/PictureCardSkeleton.vue';
import PicCarousel from './PicCarousel.vue';

export default {
    setup() {
        const themeStore = useThemeStore();
        return { themeStore };
    },
    inject: ['imagebasedir'],
    components: {
        Modal: ModalDialog,
        CardSkeleton: PictureCardSkeleton,
        Carousel: PicCarousel,
    },
    props: {
        imageObject: {
            type: Object as PropType<ImageProps>,
            default: {} as ImageProps,
        },
    },
    data() {
        return {
            isLoaded: false,
        };
    },
    methods: {
        onImgLoad() {
            setTimeout(() => {
                this.isLoaded = true;
            }, 1000);
        },
    },
    computed: {
        imgIdToModalIdSelector() {
            return '#' + this.imageObject.id + 'Modal';
        },
        imgIdToModalId() {
            return this.imageObject.id + 'Modal';
        },
        mainCardImage() {
            return this.imagebasedir + this.imageObject.dir + '1.jpg';
        },
        showCardShadow() {
            if (this.themeStore.theme === 'light') {
                return 'shadow-lg';
            } else return '';
        },
    },
};
</script>

<template>
    <div class="col-12 col-sm-6 col-md-4 d-flex justify-content-center">
        <div class="card p-3 mb-5 rounded showCardShadow" v-show="isLoaded">
            <img
                :src="mainCardImage"
                @load="onImgLoad"
                class="card-img-top"
                :alt="$i18n.locale === 'RUS' ? imageObject.name_ru : imageObject.name_en"
                data-bs-toggle="modal"
                :data-bs-target="imgIdToModalIdSelector"
            />

            <div class="card-body">
                <div class="desc">
                    <h5 class="card-title">
                        {{
                            $i18n.locale === 'RUS'
                                ? imageObject.name_ru
                                : imageObject.name_en
                        }}
                    </h5>

                    <div class="card-text">
                        <span v-if="imageObject.base_ru">
                            {{
                                $i18n.locale === 'RUS'
                                    ? imageObject.base_ru
                                    : imageObject.base_en
                            }}
                        </span>
                        <span v-if="imageObject.material_ru">
                            {{
                                $i18n.locale === 'RUS'
                                    ? `, ${imageObject.material_ru}`
                                    : `, ${imageObject.material_en}`
                            }}
                        </span>
                        <span v-if="imageObject.size">
                            {{ `, ${imageObject.size}` }}
                        </span>
                        <span v-if="imageObject.year">
                            {{ `, ${imageObject.year}` }}
                        </span>
                    </div>
                    <p v-if="imageObject.price">
                        {{ $t('card.price') + ` ${imageObject.price} ` + $t('card.rub') }}
                    </p>
                </div>
            </div>

            <Modal :modalId="imgIdToModalId">
                <Carousel :imageObject="imageObject" :imageId="imageObject.id" />
            </Modal>
        </div>
        <CardSkeleton v-if="!isLoaded" />
    </div>
</template>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
}

.card {
    --bs-card-bg: var(--color-surface);
    opacity: 1;
    width: 100%;
}

img {
    display: block;
    width: 100%;
    height: inherit;
    object-fit: cover;
}

.desc {
    align-self: flex-end;
    width: 100%;
}

.card > img:hover {
    cursor: pointer;
}
</style>
