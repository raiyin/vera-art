<script lang="ts">
import NewsPhotoItem from '@/components/app-news/NewsPhotoItem.vue';
import { PropType } from 'vue';
import type { NewsDesc } from '@/types';
import ModalDialog from '../app-ui/ModalDialog.vue';
import NewsCarousel from './NewsCarousel.vue';

export default {
    components: {
        NewsPhotoItem,
        NewsCarousel: NewsCarousel,
        Modal: ModalDialog,
    },
    props: {
        currentNewsItem: {
            type: Object as PropType<NewsDesc>,
            default: {} as NewsDesc,
        },
    },
    data() {
        return {
            selectedIndex: 1,
        };
    },
    methods: {
        setSelectedIndex(index: number) {
            this.selectedIndex = index - 1;
        },
    },
};
</script>

<template>
    <section class="container text-center main-content">
        <div class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <!-- Renders 1 to length -->
                <template v-for="image_index in currentNewsItem.images.length">
                    <NewsPhotoItem
                        v-if="image_index % 3 == 1"
                        :key="image_index"
                        :image_index="image_index"
                        :currentNews="currentNewsItem"
                        :fileName="currentNewsItem.images[image_index - 1]"
                        @click="setSelectedIndex(image_index)"
                    />
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.images.length">
                    <NewsPhotoItem
                        v-if="image_index % 3 == 2"
                        :key="image_index"
                        :image_index="image_index"
                        :currentNews="currentNewsItem"
                        :fileName="currentNewsItem.images[image_index - 1]"
                        @click="setSelectedIndex(image_index)"
                    />
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.images.length">
                    <NewsPhotoItem
                        v-if="image_index % 3 == 0"
                        :key="image_index"
                        :image_index="image_index"
                        :currentNews="currentNewsItem"
                        :fileName="currentNewsItem.images[image_index - 1]"
                        @click="setSelectedIndex(image_index)"
                    />
                </template>
            </div>

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
