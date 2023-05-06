<script>
import Modal from './Modal.vue';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
    },
    inject: ['imagebasedir'],
    components: {
        Modal,
    },
    props: {
        imageObject: {
            type: Object,
        },
    },
    data() {
        return {
            isLoaded: false,
            loadingGrey: '#ededed',
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
    },
    watch: {},
};
</script>

<template>
    <div class="col-12 col-sm-6 col-md-4 d-flex justify-content-center">
        <div
            class="card shadow-lg p-3 mb-5 bg-body rounded"
            :class="[!this.isLoaded ? 'loading' : '']"
        >
            <img
                v-show="this.isLoaded"
                :src="mainCardImage"
                @load="onImgLoad"
                class="card-img-top"
                alt="..."
                data-bs-toggle="modal"
                :data-bs-target="imgIdModalToLink"
            />
            <div v-show="!this.isLoaded" class="image" />

            <div class="card-body">
                <div class="desc">
                    <h5 class="card-title">
                        {{ !this.isLoaded ? '' : imageObject.name_ru }}
                    </h5>
                    <div class="card-text">
                        <span v-if="imageObject.base">
                            {{ !this.isLoaded ? '' : imageObject.base }}
                        </span>
                        <span v-if="imageObject.material">
                            {{ !this.isLoaded ? '' : `, ${imageObject.material}` }}
                        </span>
                        <span v-if="imageObject.size">
                            {{ !this.isLoaded ? '' : `, ${imageObject.size}` }}
                        </span>
                        <span v-if="imageObject.year">
                            {{ !this.isLoaded ? '' : `, ${imageObject.year}` }}
                        </span>
                    </div>
                    <p v-if="imageObject.price">
                        {{
                            !this.isLoaded
                                ? ''
                                : this.$t('card.price') +
                                  ` ${imageObject.price} ` +
                                  this.$t('card.rub')
                        }}
                    </p>
                </div>
            </div>
            <Modal :imageObject="imageObject" />
        </div>
    </div>
</template>

<style scoped>
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
