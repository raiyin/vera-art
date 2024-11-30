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
        <div class="row">
            <template v-for="i in 3">
                <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                    <template v-for="image_index in currentNewsItem.imagescount">
                        <NewsPhotoItem
                            v-if="image_index % 3 === i % 3"
                            :image_index="image_index"
                            :currentNews="currentNewsItem"
                            :key="image_index"
                            @click="setSelectedIndex(image_index)"
                        />
                    </template>
                </div>
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
