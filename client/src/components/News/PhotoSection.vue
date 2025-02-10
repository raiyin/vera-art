<script lang="ts">
import NewsPhotoItem from '@/components/News/NewsPhotoItem.vue';
import { PropType } from 'vue';
import type { NewsItemType } from '@/types';
import ModalDialog from '../UI/ModalDialog.vue';
import NewsCarousel from './NewsCarousel.vue';

export default {
    components: {
        NewsPhotoItem,
        NewsCarousel: NewsCarousel,
        Modal: ModalDialog,
    },
    props: {
        currentNewsItem: {
            type: Object as PropType<NewsItemType>,
            default: {} as NewsItemType,
        },
    },
    data() {
        return {
            selectedIndex: 1,
        };
    },
    methods: {
        setSelectedIndex(index: number) {
            this.selectedIndex = index;
        },
    },
};
</script>

<template>
    <section class="container text-center px-0 main-content">
        <div class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
            <template
                v-for="image_index in currentNewsItem.imagescount"
                :key="image_index"
            >
                <NewsPhotoItem
                    :image_index="image_index"
                    :currentNews="currentNewsItem"
                    @click="setSelectedIndex(image_index)"
                />
            </template>

            <Modal modalId="imgNewsModal">
                <NewsCarousel
                    :imageObject="currentNewsItem"
                    :selectedIndex="selectedIndex"
                    :setSelectedIndex="setSelectedIndex"
                />
            </Modal>
        </div>
    </section>
</template>
