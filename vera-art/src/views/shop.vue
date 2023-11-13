<script lang="ts">
import Gallery from '@/components/UI/Gallery.vue';
import axios from 'axios';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';
import { useI18n } from 'vue-i18n';
import type { ImageProps, SortOption } from '@/types';

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
            sortingType: '' as string,
            sortOptions: [
                { value: 'name_ru', name: this.$t('shop.byName') },
                { value: 'year', name: this.$t('shop.byNovelty') },
                { value: 'height', name: this.$t('shop.byHeight') },
                { value: 'width', name: this.$t('shop.byWidth') },
            ] as SortOption[],
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
    watch: {
        selectedSort(image_first: ImageProps, image_second: ImageProps) {
            if (
                typeof image_first[
                    this.sortingType as keyof typeof image_first
                ] === 'string' &&
                typeof image_second[
                    this.sortingType as keyof typeof image_second
                ] === 'string'
            ) {
                return (
                    image_first[
                        this.sortingType as keyof typeof image_first
                    ] as string
                )?.localeCompare(
                    image_second[
                        this.sortingType as keyof typeof image_second
                    ] as string
                );
            }
            if (
                typeof image_first[
                    this.sortingType as keyof typeof image_first
                ] === 'number' &&
                typeof image_second[
                    this.sortingType as keyof typeof image_second
                ] === 'number'
            ) {
                return (
                    +image_first[this.sortingType as keyof typeof image_first] -
                    +image_second[this.sortingType as keyof typeof image_second]
                );
            }
            return 0;
        },
    },
};
</script>

<template>
    <section class="container mt-3 px-0">
        <v-select
            :options="sortOptions"
            :reduce="(item: SortOption) => item.value"
            label="name"
            v-model="sortingType"
            inputId="value"
            :placeholder="$t('mySelect.placeholder')"
        >
            {{ $t('mySelect.placeholder') }}
        </v-select>
    </section>
    <Gallery :images="images" />
    <div ref="observer" class="observer"></div>
</template>

<style scoped>
.v-select {
    width: 20rem;
    box-sizing: border-box;
}

.vs__dropdown-option {
    background-color: black;
}

.observer {
    height: 0px;
}
</style>
