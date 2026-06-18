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
                title: string;
                description: string | null;
                image_path: string;
                images: string[];
                year: number | null;
                technique: string | null;
                size: string | null;
                status: string;
                sort_order: number;
                material_ids: number[] | null;
                base_ids: number[] | null;
                created_at: string;
                updated_at: string;
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
            str_id: String(work.id),
            dir: relWorksDir,
            name_ru: work.title,
            name_en: work.title,
            year: work.year ?? 0,
            descr: work.description ?? '',
            base_ru: work.technique ?? '',
            base_en: work.technique ?? '',
            width: 0,
            height: 0,
            type: undefined,
            images: work.images ?? [],
            materials_ids: work.material_ids ?? [],
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
