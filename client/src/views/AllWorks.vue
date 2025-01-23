<script lang="ts">
import Gallery from '@/components/UI/PicGallery.vue';
import axios from 'axios';

export default {
    inject: ['server'],
    components: {
        Gallery,
    },
    data() {
        return {
            paintingImages: [],
            paitingPage: 0,
            illustrationImages: [],
            illustrationPage: 0,
            threeDImages: [],
            threeDPage: 0,
            limit: 9,
        };
    },
    methods: {
        async loadPainting() {
            try {
                this.paitingPage += 1;
                const response = await axios.get(this.server + 'paintings', {
                    params: {
                        _page: this.paitingPage,
                        _limit: this.limit,
                    },
                });
                this.paintingImages = [
                    ...this.paintingImages,
                    ...response.data,
                ];
            } catch (e) {
                console.error('Error fetching paintings on allworks page ' + e);
            }
        },
        async loadIllustrations() {
            try {
                this.illustrationPage += 1;
                const response = await axios.get(
                    this.server + 'illustrations',
                    {
                        params: {
                            _page: this.illustrationPage,
                            _limit: this.limit,
                        },
                    },
                );
                this.illustrationImages = [
                    ...this.illustrationImages,
                    ...response.data,
                ];
            } catch (e) {
                console.error(
                    'Error fetching imaillustration on all works page ' + e,
                );
            }
        },
        async load3D() {
            try {
                this.threeDPage += 1;
                const response = await axios.get(this.server + 'threeds', {
                    params: {
                        _page: this.threeDPage,
                        _limit: this.limit,
                    },
                });
                this.threeDImages = [...this.threeDImages, ...response.data];
            } catch (e) {
                console.error(
                    'Error fetching 3d images on all works page ' + e,
                );
            }
        },
    },
    mounted() {
        const options = {
            rootMargin: '0px',
            threshold: 1.0,
        };
        const paintingCallback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.loadPainting();
            }
        };
        const illustrationCallback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.loadIllustrations();
            }
        };
        const threeDCallback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.load3D();
            }
        };

        const paintingObserver = new IntersectionObserver(
            paintingCallback,
            options,
        );
        paintingObserver.observe(this.$refs.paintingObserver as Element);

        const illustrationObserver = new IntersectionObserver(
            illustrationCallback,
            options,
        );
        illustrationObserver.observe(
            this.$refs.illustrationObserver as Element,
        );

        const threeDObserver = new IntersectionObserver(
            threeDCallback,
            options,
        );
        threeDObserver.observe(this.$refs.graphics3dObserver as Element);
    },
};
</script>

<template>
    <section class="container main-content">
        <ul class="nav nav-tabs" id="myTab" role="tablist">
            <li class="nav-item" role="presentation">
                <button
                    class="nav-link active"
                    id="painting-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#painting"
                    type="button"
                    role="tab"
                    aria-controls="painting"
                    aria-selected="true"
                >
                    {{ $t('allworks.painting') }}
                </button>
            </li>
            <li class="nav-item" role="presentation">
                <button
                    class="nav-link"
                    id="Illustration-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#Illustration"
                    type="button"
                    role="tab"
                    aria-controls="Illustration"
                    aria-selected="false"
                >
                    {{ $t('allworks.illustration') }}
                </button>
            </li>
            <li class="nav-item" role="presentation">
                <button
                    class="nav-link"
                    id="graphics3d-tab"
                    data-bs-toggle="tab"
                    data-bs-target="#graphics3d"
                    type="button"
                    role="tab"
                    aria-controls="graphics3d"
                    aria-selected="false"
                >
                    {{ $t('allworks.3dgraphics') }}
                </button>
            </li>
        </ul>

        <div class="tab-content" id="myTabContent">
            <div
                class="tab-pane fade show active"
                id="painting"
                role="tabpanel"
                aria-labelledby="painting-tab"
            >
                <Gallery :images="paintingImages" />

                <div ref="paintingObserver" class="observer" />
            </div>

            <div
                class="tab-pane fade"
                id="Illustration"
                role="tabpanel"
                aria-labelledby="Illustration-tab"
            >
                <Gallery :images="illustrationImages" />
                <div ref="illustrationObserver" class="observer" />
            </div>

            <div
                class="tab-pane fade"
                id="graphics3d"
                role="tabpanel"
                aria-labelledby="graphics3d-tab"
            >
                <Gallery :images="threeDImages" />
                <div ref="graphics3dObserver" class="observer" />
            </div>
        </div>
    </section>
</template>

<style scoped>
.tab-content {
    border-left: 1px solid var(--color-surface-secondary);
    border-right: 1px solid var(--color-surface-secondary);
    padding: 10px;
}

.nav-tabs {
    margin-bottom: 0;
    --bs-nav-tabs-border-color: var(--color-surface-secondary);
}

.nav-tabs .nav-item .nav-link {
    color: var(--color-on-surface);
    background-color: var(--color-surface);
}

.nav-tabs .nav-item .nav-link.active {
    color: var(--color-on-surface);
    border-color: var(--color-surface-secondary);
}

.nav-item {
    color: black;
}

.observer {
    height: 0px;
}
</style>
