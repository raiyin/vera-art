<script setup lang="ts">
import type { Work } from '~/types/work';
import Gallery from '~/components/PicGallery.vue';
import { useInfiniteScroll } from '~/composables/useInfiniteScroll';
import { getHttpClient } from '~/api/http-client';

const { items: works, sentinelRef, loading } = useInfiniteScroll<Work>(
    async (page) => {
        const limit = useRuntimeConfig().public.limit as string;
        const { data } = await getHttpClient().get<Work[]>('works', {
            params: {
                offset: page * +limit,
                limit,
            },
        });
        return data.map((work) => ({
            ...work,
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
