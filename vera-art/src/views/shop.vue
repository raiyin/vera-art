<script>
import Gallery from '@/components/UI/Gallery.vue';
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
            images: [],
            page: 1,
            limit: 9,
            selectedSort: '',
            sortOptions: [
                { value: 'name_ru', name: this.$t('shop.byName') },
                { value: 'year', name: this.$t('shop.byNovelty') },
                { value: 'height', name: this.$t('shop.byHeight') },
                { value: 'width', name: this.$t('shop.byWidth') },
            ],
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
        const callback = (entries, observer) => {
            if (entries[0].isIntersecting) {
                this.loadMoreWorks();
            }
        };
        const observer = new IntersectionObserver(callback, options);
        observer.observe(this.$refs.observer);
    },
    computed: {},
    watch: {
        selectedSort(newValue, oldValue) {
            this.images.sort((image_first, image_second) => {
                if (typeof image_first[this.selectedSort] === 'string')
                    return image_first[this.selectedSort]?.localeCompare(
                        image_second[this.selectedSort]
                    );
                if (typeof image_first[this.selectedSort] === 'number')
                    return (
                        image_first[this.selectedSort] - image_second[this.selectedSort]
                    );
            });
        },
    },
};
</script>

<template>
    <section class="container mt-3 px-0">
        <v-select
            :options="sortOptions"
            :reduce="(item) => item.value"
            label="name"
            v-model="selectedSort"
            inputId="value"
            :placeholder="this.$t('mySelect.placeholder')"
        >
            {{ this.$t('mySelect.placeholder') }}
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
