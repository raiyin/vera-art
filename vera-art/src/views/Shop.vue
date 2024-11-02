<script lang="ts">
import Gallery from '@/components/UI/Gallery.vue';
import type { ImageProps, SortOption } from '@/types';
import axios from 'axios';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
    },
    inject: ['jsonserverhost'],
    components: {
        Gallery,
        vSelect,
    },
    data() {
        return {
            images: [] as ImageProps[],
            page: 1,
            limit: 9,
            selectedSort: '',
            // sortOptions: [
            //     { value: 'name_ru', name: 'По названию' },
            //     { value: 'year', name: 'По новизне' },
            //     { value: 'height', name: 'По высоте' },
            //     { value: 'width', name: 'По ширине' },
            // ],
        };
    },
    methods: {
        async loadWorks() {
            try {
                const response = await axios.get(this.jsonserverhost + 'sale', {
                    params: {
                        _page: this.page,
                        _limit: this.limit,
                    },
                });
                this.images = response.data;
            } catch (e) {
                console.log(e);
            }
        },
        async loadMoreWorks() {
            try {
                this.page += 1;
                const response = await axios.get(this.jsonserverhost + 'sale', {
                    params: {
                        _page: this.page,
                        _limit: this.limit,
                    },
                });
                this.images = [...this.images, ...response.data];
            } catch (e) {
                console.log(e);
            }
        },
    },
    mounted() {
        this.loadWorks();
        const options = {
            rootMargin: '0px',
            threshold: 1.0,
        };
        const callback = (entries: IntersectionObserverEntry[]) => {
            if (entries[0].isIntersecting) {
                this.loadMoreWorks();
            }
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
            this.images.sort(
                (image_first: ImageProps, image_second: ImageProps) => {
                    if (
                        typeof image_first[
                            this.selectedSort as keyof typeof image_first
                        ] === 'string'
                    )
                        return (
                            image_first[
                                this.selectedSort as keyof typeof image_first
                            ] as string
                        )?.localeCompare(
                            image_second[
                                this.selectedSort as keyof typeof image_first
                            ] as string
                        );

                    if (
                        typeof image_first[
                            this.selectedSort as keyof typeof image_first
                        ] === 'number'
                    )
                        return (
                            +image_first[
                                this.selectedSort as keyof typeof image_first
                            ] -
                            +image_second[
                                this.selectedSort as keyof typeof image_first
                            ]
                        );

                    return 0;
                }
            );
        },
    },
};
</script>

<template>
    <section class="container mt-3">
        <v-select
            :options="sortOptions"
            :reduce="(item: SortOption) => item.value"
            label="name"
            v-model="selectedSort"
            inputId="value"
            :placeholder="$t('mySelect.placeholder')"
        >
        </v-select>
    </section>
    <Gallery :images="images" />
    <div ref="observer" class="observer"></div>
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
