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
    <section class="container text-center main-content">
        <div class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4">
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 1">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                            @click="setSelectedIndex(image_index)"
                        />
                    </template>
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 2">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                        />
                    </template>
                </template>
            </div>

            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 0">
                        <NewsPhotoItem
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                        />
                    </template>
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
