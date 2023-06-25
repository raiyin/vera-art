<script>
import Gallery from '@/components/UI/Gallery.vue';
import axios from 'axios';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
    },
    inject: ['jsonserverhost'],
    components: {
        Gallery,
    },
    props: [],
    data() {
        return {
            paintingImages: [],
            illustrationImages: [],
            threeDImages: [],
        };
    },
    methods: {
        async fetchData() {
            try {
                let response = await axios.get(this.jsonserverhost + 'painting');
                this.paintingImages = response.data;
                response = await axios.get(this.jsonserverhost + 'illustration');
                this.illustrationImages = response.data;
                response = await axios.get(this.jsonserverhost + '3d');
                this.threeDImages = response.data;
            } catch (e) {
                console.log(e);
            }
        },
    },
    mounted() {
        this.fetchData();
    },
    computed: {},
    watch: {},
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
            </div>
            <div
                class="tab-pane fade"
                id="Illustration"
                role="tabpanel"
                aria-labelledby="Illustration-tab"
            >
                <Gallery :images="illustrationImages" />
            </div>
            <div
                class="tab-pane fade"
                id="graphics3d"
                role="tabpanel"
                aria-labelledby="graphics3d-tab"
            >
                <Gallery :images="threeDImages" />
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
</style>
