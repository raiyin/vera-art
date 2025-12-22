<script lang="ts">
import Gallery from '@/components/app-ui/PicGallery.vue';
import axios from 'axios';

export default {
    components: {
        Gallery,
    },
    data() {
        return {
            works: [],
            page: -1,
            limit: import.meta.env.VITE_PAGE_SIZE,
            server: import.meta.env.VITE_SERVER_URL,
        };
    },
    methods: {
        async loadWorks() {
            try {
                this.page += 1;
                const response = await axios.get(this.server + 'works', {
                    params: {
                        offset: this.page * this.limit,
                        limit: this.limit,
                    },
                });
                this.works = [...this.works, ...response.data];
            } catch (e) {
                console.error('Error fetching works on allworks page ' + e);
            }
        },
        async loadIllustrations() {
            try {
                this.illustrationPage += 1;
                const response = await axios.get(this.server + 'illustrations', {
                    params: {
                        offset: this.illustrationPage * this.limit,
                        limit: this.limit,
                    },
                });
                this.illustrationImages = [...this.illustrationImages, ...response.data];
            } catch (e) {
                console.error('Error fetching imaillustration on all works page ' + e);
            }
        },
        async load3D() {
            try {
                this.threeDPage += 1;
                const response = await axios.get(this.server + 'threeds', {
                    params: {
                        offset: this.threeDPage * this.limit,
                        limit: this.limit,
                    },
                });
                this.threeDImages = [...this.threeDImages, ...response.data];
            } catch (e) {
                console.error('Error fetching 3d images on all works page ' + e);
            }
        },
    },
    mounted() {
        const options = {
            rootMargin: '0px',
            threshold: 1.0,
        };
        const worksCallback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.loadWorks();
            }
        };
        const worksObserver = new IntersectionObserver(worksCallback, options);
        worksObserver.observe(this.$refs.worksObserver as Element);
    },
};
</script>

<template>
    <section class="container main-content">
        <Gallery :images="works" />

        <div ref="worksObserver" class="observer" />
    </section>
</template>

<style scoped>
.observer {
    height: 0px;
}
</style>
