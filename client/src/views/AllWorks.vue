<script lang="ts">
import Gallery from '@/components/app-ui/PicGallery.vue';
import { TypedGetWorkDto } from '@/types/work';
import axios from 'axios';

export default {
    components: {
        Gallery,
    },
    emits: ['work-deleted'],
    data() {
        return {
            works: [] as TypedGetWorkDto[],
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

                const newWorks = response.data.map((work) => ({
                    __type: 'GetSaleDto',
                    ...work,
                }));
                this.works = [...this.works, ...newWorks];
            } catch (e) {
                console.error('Error fetching works on allworks page ' + e);
            }
        },
        handleWorkDeleted(id: string) {
            // Remove the deleted work from the works array
            this.works = this.works.filter((work: any) => work.id !== id);
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
        <Gallery :images="works" @work-deleted="handleWorkDeleted" />

        <div ref="worksObserver" class="observer" />
    </section>
</template>

<style scoped>
.observer {
    height: 0px;
}
</style>
