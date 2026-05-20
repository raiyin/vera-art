<script setup lang="ts">
import NewsPhotoItem from './NewsPhotoItem.vue';
import type { PropType } from 'vue';
import type { NewsDesc } from '../types';
import NewsCarousel from './NewsCarousel.vue';
import { ref } from 'vue';

const props = defineProps({
    currentNewsItem: {
        type: Object as PropType<NewsDesc>,
        default: {} as NewsDesc,
    },
});

const selectedIndex = ref(1);

const setSelectedIndex = (index: number) => {
    selectedIndex.value = index - 1;
};
</script>

<template>
    <UContainer class="text-center main-content">
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

            <UModal modalId="imgNewsModal">
                <NewsCarousel
                    :imageObject="currentNewsItem"
                    :selectedIndex="selectedIndex"
                    :setSelectedIndex="setSelectedIndex"
                />
            </UModal>
        </div>
    </UContainer>
</template>
