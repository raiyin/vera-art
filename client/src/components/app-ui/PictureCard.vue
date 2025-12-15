<script lang="ts">
import type { ImageProps } from '@/props/image-props';
import ModalDialog from './ModalDialog.vue';
import { useThemeStore } from '../../stores/ThemeStore';
import type { PropType } from 'vue';
import PictureCardSkeleton from '../app-skeletons/PictureCardSkeleton.vue';
import PicCarousel from './PicCarousel.vue';
import { useAuthStore } from '../../stores/AuthStore';
import { storeToRefs } from 'pinia';

export default {
    setup() {
        const themeStore = useThemeStore();
        const authStore = useAuthStore();
        const { isAuthenticated } = storeToRefs(authStore);

        return {
            themeStore,
            isAuthenticated,
        };
    },
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
            imagebasedir: import.meta.env.VITE_IMAGE_DIR,
        };
    },
    methods: {
        onImgLoad() {
            setTimeout(() => {
                this.isLoaded = true;
            }, 1000);
        },
        edit() {
            this.$router.push('/paintings/edit/' + this.imageObject.id + '/');
        },
        onImageDelete(str_id: string) {},
    },
    computed: {
        imgIdToModalIdSelector() {
            return '#' + this.imageObject.str_id + 'Modal';
        },
        imgIdToDeleteIdSelector() {
            return '#' + this.imageObject.str_id + 'DeleteModal';
        },
        imgIdToModalId() {
            return this.imageObject.str_id + 'Modal';
        },
        imgIdToDeleteId() {
            return this.imageObject.str_id + 'DeleteModal';
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
                        <span v-if="imageObject.width && imageObject.width != '0'">
                            {{ `, ${imageObject.width + 'x' + imageObject.height}` }}
                        </span>
                        <span v-if="imageObject.year">
                            {{ `, ${imageObject.year}` }}
                        </span>
                    </div>

                    <p v-if="imageObject.price">
                        {{ $t('card.price') + ` ${imageObject.price} ` + $t('card.rub') }}
                    </p>

                    <div class="image-control" v-if="isAuthenticated">
                        <button
                            class="btn btn-secondary w-100"
                            type="button"
                            v-on:click="edit"
                        >
                            Редактировать
                        </button>
                        <button
                            class="btn btn-secondary w-100"
                            type="button"
                            data-bs-toggle="modal"
                            :data-bs-target="imgIdToDeleteIdSelector"
                        >
                            Удалить
                        </button>
                    </div>
                </div>
            </div>

            <Modal :modalId="imgIdToModalId">
                <Carousel :imageObject="imageObject" :imageId="imageObject.str_id" />
            </Modal>
        </div>
        <CardSkeleton v-if="!isLoaded" />

        <div
            class="modal fade"
            :id="imgIdToDeleteId"
            tabindex="-1"
            aria-labelledby="imageDeleteModalLabel"
            aria-hidden="true"
        >
            <div class="modal-dialog modal-dialog-centered">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="imageDeleteModalLabel">
                            Подтверждение удаления
                        </h5>
                        <button
                            type="button"
                            class="btn-close"
                            data-bs-dismiss="modal"
                            aria-label="Close"
                        ></button>
                    </div>
                    <div class="modal-body">
                        Вы действительно хотите удалить работу "{{
                            imageObject.name_ru
                        }}"?
                    </div>
                    <div class="modal-footer">
                        <button
                            type="button"
                            class="btn btn-secondary"
                            data-bs-dismiss="modal"
                        >
                            Отменить
                        </button>
                        <button type="button" class="btn btn-danger">Удалить</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
    padding: 0;
}

.card {
    --bs-card-bg: var(--color-surface);
    opacity: 1;
    width: 100%;
}

.desc {
    color: var(--color-on-surface);
    align-self: flex-end;
    width: 100%;
}

img {
    display: block;
    width: 100%;
    height: inherit;
    object-fit: cover;
}

.image-control {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 0.5rem;
    margin-top: 1rem;
}

.card > img:hover {
    cursor: pointer;
}
</style>
