<script lang="ts">
import type { ImageProps } from '@/types';
import Modal from './Modal.vue';
import { useThemeStore } from '../../stores/ThemeStore';
import { useI18n } from 'vue-i18n';
import type { PropType } from 'vue';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
        const themeStore = useThemeStore();
        return { themeStore };
    },
    inject: ['imagebasedir'],
    components: {
        Modal,
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
            loadingGrey: '#adadad',
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
        imgIdModalToLink() {
            return '#' + this.imageObject.id + 'Modal';
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
        <div
            class="card p-3 mb-5 rounded"
            :class="[isLoaded ? showCardShadow : 'loading']"
        >
            <img
                v-show="isLoaded"
                :src="mainCardImage"
                @load="onImgLoad"
                class="card-img-top"
                alt="..."
                data-bs-toggle="modal"
                :data-bs-target="imgIdModalToLink"
            />

            <div v-show="!isLoaded" class="image" />

            <div class="card-body">
                <div class="desc">
                    <h5 class="card-title">
                        {{
                            !isLoaded
                                ? ''
                                : $i18n.locale === 'RUS'
                                ? imageObject.name_ru
                                : imageObject.name_en
                        }}
                    </h5>

                    <div class="card-text">
                        <span v-if="imageObject.base_ru">
                            {{
                                !isLoaded
                                    ? ''
                                    : $i18n.locale === 'RUS'
                                    ? imageObject.base_ru
                                    : imageObject.base_en
                            }}
                        </span>
                        <span v-if="imageObject.material_ru">
                            {{
                                !isLoaded
                                    ? ''
                                    : $i18n.locale === 'RUS'
                                    ? `, ${imageObject.material_ru}`
                                    : `, ${imageObject.material_en}`
                            }}
                        </span>
                        <span v-if="imageObject.size">
                            {{ !isLoaded ? '' : `, ${imageObject.size}` }}
                        </span>
                        <span v-if="imageObject.year">
                            {{ !isLoaded ? '' : `, ${imageObject.year}` }}
                        </span>
                    </div>
                    <p v-if="imageObject.price">
                        {{
                            !isLoaded
                                ? ''
                                : $t('card.price') +
                                  ` ${imageObject.price} ` +
                                  $t('card.rub')
                        }}
                    </p>
                </div>
            </div>
            <Modal :imageObject="imageObject" />
        </div>
    </div>
</template>

<style scoped>
.nav-tabs {
    --bs-nav-tabs-border-color: #000 !important;
}

.card {
    --bs-card-bg: var(--color-surface);
}

.card-body {
    display: flex;
    align-items: center;
}

.card {
    opacity: 1;
    width: 100%;
}

.image {
    height: 200px;
}

.loading h5 {
    width: 80%;
    display: inline-block;
    text-align: center;
}

.loading h5,
.loading .card-text,
.loading p {
    height: 20px;
    margin-bottom: 5px;
}

.loading .image,
.loading h5,
.loading .card-text,
.loading p,
.loading .description {
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

.image img {
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
