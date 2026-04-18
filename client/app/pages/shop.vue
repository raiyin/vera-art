<script lang="ts">
import Gallery from '@/components/app-ui/ShopGallery.vue';
import type { SortOption } from '@/types';
import type { GetSaleDto } from '@/types/sale';
import axios from 'axios';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

export default {
    components: {
        Gallery,
        vSelect,
    },
    data() {
        return {
            images: [] as GetSaleDto[],
            page: 0,
            limit: import.meta.env.VITE_PAGE_SIZE,
            server: import.meta.env.VITE_SERVER_URL,
            selectedSort: '',
        };
    },
    methods: {
        async loadWorks() {
            try {
                this.page += 1;
                const response = await axios.get(this.server + 'sales', {
                    params: {
                        offset: this.page * this.limit,
                        limit: this.limit,
                    },
                });
                const newImages = response.data.map((image) => ({
                    __type: 'GetSaleDto',
                    ...image,
                }));
                this.images = [...this.images, ...newImages];
            } catch (e) {
                console.error('Error fetching images on shop page');
            }
        },
    },
    mounted() {
        const callback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.loadWorks();
            }
        };
        const options = {
            rootMargin: '0px',
            threshold: 1.0,
        };
        const observer = new IntersectionObserver(callback, options);
        observer.observe(this.$refs.observer as Element);
    },
    computed: {
        sortOptions() {
            return [
                { value: 'name_ru', name: this.$t('shop.byName') },
                { value: 'year', name: this.$t('shop.byNovelty') },
                { value: 'height', name: this.$t('shop.byHeight') },
                { value: 'width', name: this.$t('shop.byWidth') },
            ];
        },
    },
    watch: {
        selectedSort() {
            this.images.sort((image_first: GetSaleDto, image_second: GetSaleDto) => {
                if (
                    typeof image_first[this.selectedSort as keyof typeof image_first] ===
                    'string'
                )
                    return (image_first[
                        this.selectedSort as keyof typeof image_first
                    ] as string)?.localeCompare(
                        image_second[
                            this.selectedSort as keyof typeof image_first
                        ] as string
                    );

                if (
                    typeof image_first[this.selectedSort as keyof typeof image_first] ===
                    'number'
                )
                    return (
                        +image_first[this.selectedSort as keyof typeof image_first] -
                        +image_second[this.selectedSort as keyof typeof image_first]
                    );

                return 0;
            });
        },
    },
};
</script>

<template>
    <UContainer class="mt-3">
        <v-select
            :options="sortOptions"
            :reduce="(item: SortOption) => item.value"
            label="name"
            v-model="selectedSort"
            inputId="value"
            :placeholder="$t('mySelect.placeholder')"
        >
        </v-select>
    </UContainer>
    <Gallery :images="images" />
    <div ref="observer" class="observer" />
</template>

<style scoped>
.v-select {
    width: 20rem;
    box-sizing: border-box;
    cursor: pointer;
}

.observer {
    height: 0px;
}
</style>
