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
    },
    props: {
        currentNewsItem: {
            type: Object as PropType<NewsItemType>,
            default: {} as NewsItemType,
        },
    },
};
</script>

<template>
    <section class="container text-center px-0 main-content">
        <div class="row">
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNewsItem.imagescount">
                    <template v-if="image_index % 3 == 1">
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
        </div>
    </section>
    <NewsCarousel :imageObject="currentNewsItem" />
</template>
