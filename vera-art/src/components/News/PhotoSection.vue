<script lang="ts">
import NewsPhotoItem from '@/components/News/NewsPhotoItem.vue';
import { useI18n } from 'vue-i18n';
import { PropType } from 'vue';
import type { NewsItemType } from '@/types';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
        return {};
    },
    components: {
        NewsPhotoItem,
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
</template>

<style scoped></style>
