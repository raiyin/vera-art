<script setup lang="ts">
import type { TypedGetSaleDto } from '~/types/sale';
import Gallery from '~/components/PicGallery.vue';
import { useInfiniteScroll } from '~/composables/useInfiniteScroll';
import { getHttpClient } from '~/api/http-client';

const { items: sales, sentinelRef, loading } = useInfiniteScroll<TypedGetSaleDto>(
    async (page) => {
        const limit = useRuntimeConfig().public.limit as string;
        const { data } = await getHttpClient().get<TypedGetSaleDto[]>('sales', {
            params: {
                offset: page * +limit,
                limit,
            },
        });
        return data.map((sale) => ({
            ...sale,
            __type: 'GetSaleDto' as const,
        }));
    },
    { immediate: true }
);

const emit = defineEmits<{
    'work-deleted': [id: string];
}>();

const handleSaleDeleted = (id: string) => {
    sales.value = sales.value.filter((sale: TypedGetSaleDto) => String(sale.id) !== id);
};
</script>

<template>
    <UContainer class="main-content">
        <Gallery :images="sales" @work-deleted="handleSaleDeleted" />

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
