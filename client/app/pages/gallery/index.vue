<script setup lang="ts">
import type { Work } from '~/types/work';
import Gallery from '~/components/PicGallery.vue';
import { useInfiniteScroll } from '~/composables/useInfiniteScroll';
import { getHttpClient } from '~/api/http-client';

const config = useRuntimeConfig();
const relWorksDir = config.public.relWorksDir as string;

const { items: works, sentinelRef, loading } = useInfiniteScroll<Work>(
    async (page) => {
        const limit = config.public.limit as string;
        const { data } = await getHttpClient().get<{
            works: Array<{
                id: number;
                str_id: string;
                width: number;
                height: number;
                year: number;
                name_ru: string;
                name_en: string;
                base_id: number;
                descr_ru: string | null;
                descr_en: string | null;
                work_path: string;
                images: string[];
            }>;
            total: number;
        }>('works', {
            params: {
                page: page + 1,
                limit,
            },
        });
        return data.works.map((work) => ({
            id: work.id,
            str_id: work.str_id,
            work_path: work.work_path,
            dir: relWorksDir + work.work_path,
            name_ru: work.name_ru,
            name_en: work.name_en,
            year: work.year,
            width: work.width,
            height: work.height,
            base_id: work.base_id,
            descr_ru: work.descr_ru ?? '',
            descr_en: work.descr_en ?? '',
            images: work.images ?? [],
            __type: 'GetWorkDto' as const,
        }));
    },
    { immediate: true }
);

const emit = defineEmits<{
    'work-deleted': [id: string];
}>();

const handleWorkDeleted = (id: string) => {
    works.value = works.value.filter((work: Work) => String(work.id) !== id);
};
</script>

<template>
    <UContainer class="main-content">
        <Gallery :images="works" @work-deleted="handleWorkDeleted" />

        <div v-if="loading" class="loading-indicator">
            <span>Загрузка...</span>
        </div>

        <div ref="sentinelRef" class="observer" />
    </UContainer>
</template>

<style scoped>
.observer {
    height: 1px;
}

.loading-indicator {
    text-align: center;
    padding: 1rem;
    color: var(--color-on-surface);
    opacity: 0.7;
}
</style>
