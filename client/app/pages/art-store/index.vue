<script setup lang="ts">
    import type { Sale, } from '~/types/work';
    import Gallery from '~/components/PicGallery.vue';
    import { useInfiniteScroll, } from '~/composables/useInfiniteScroll';
    import { getHttpClient, } from '~/api/http-client';

    interface SaleResponseItem {
        id: number
        name_ru: string
        name_en: string
        descr_ru?: string
        descr_en?: string
        image_path: string
        sale_path: string
        dir: string
        images: string[]
        price: number
        year?: number
        technique?: string
        width: number
        height: number
        status: string
        sort_order: number
        sold: boolean
        material_ids?: number[]
        base_ids?: number[]
        created_at: string
        updated_at: string
    }

    interface SalesResponse {
        sales: SaleResponseItem[]
        total: number
    }

    const { items: sales, sentinelRef, loading, } = useInfiniteScroll<Sale>(
        async (page,) => {
            const limit = useRuntimeConfig().public.limit as string;
            const { data, } = await getHttpClient().get<SalesResponse>('sales', {
                params: {
                    page: page + 1,
                    limit,
                },
            });
            return data.sales.map(item => ({
                id: item.id,
                str_id: String(item.id),
                work_path: item.sale_path,
                dir: item.dir,
                name_ru: item.name_ru,
                name_en: item.name_en,
                year: item.year ?? 0,
                descr_ru: item.descr_ru ?? '',
                descr_en: item.descr_en ?? '',
                base_id: item.base_ids?.[0] ?? 0,
                width: item.width,
                height: item.height,
                images: item.images ?? [],
                material_ids: item.material_ids ?? [],
                price: item.price,
                __type: 'GetSaleDto' as const,
            }),);
        },
        { immediate: true, },
);

    const emit = defineEmits<{
        'work-deleted': [id: string,]
    }>();

    const handleSaleDeleted = (id: string,) => {
        sales.value = sales.value.filter((sale: Sale,) => String(sale.id,) !== id,);
    };
</script>

<template>
    <UContainer class="main-content pt-20">
        <Gallery
            :images="sales"
            @work-deleted="handleSaleDeleted"
        />

        <div
            v-if="loading"
            class="loading-indicator"
        >
            <span>Загрузка...</span>
        </div>

        <div
            ref="sentinelRef"
            class="observer"
        />
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
