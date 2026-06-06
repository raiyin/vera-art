<script setup lang="ts">
import type { TypedGetSaleDto } from '~/types/sale';
import Gallery from '@/components/PicGallery.vue';
import { ref, onMounted } from 'vue';

const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;
const limit = config.public.limit;

const sales = ref<TypedGetSaleDto[]>([]);
const page = ref(-1);
const salesObserver = ref<Element | null>(null);

const emit = defineEmits<{
    'work-deleted': [id: string];
}>();

const loadWorks = async () => {
    try {
        page.value += 1;
        const params = new URLSearchParams({
            offset: (page.value * +limit).toString(),
            limit: limit,
        });
        const response = await fetch(`${SERVER_URL}sales?${params}`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const data = await response.json();
        const newWorks = data.map((work: any) => ({
            ...work,
            __type: 'GetSaleDto',
        }));
        sales.value = [...sales.value, ...newWorks];
    } catch (e) {
        console.error('Error fetching works on shop page ' + e);
    }
};

const handleSaleDeleted = (id: string) => {
    sales.value = sales.value.filter(
        (work: TypedGetSaleDto) => work.id !== id && String(work.id) !== id
    );
};

onMounted(() => {
    // Load initial data immediately
    loadWorks();

    // Set up infinite scroll observer
    const options = {
        rootMargin: '0px',
        threshold: 0,
    };
    const salesCallback = (entries: IntersectionObserverEntry[]) => {
        if (entries[0]?.isIntersecting) {
            loadWorks();
        }
    };
    const observer = new IntersectionObserver(salesCallback, options);
    if (salesObserver.value) observer.observe(salesObserver.value);
});
</script>

<template>
    <UContainer class="main-content">
        <Gallery :images="sales" @work-deleted="handleSaleDeleted" />

        <div ref="salesObserver" class="observer" />
    </UContainer>
</template>

<style scoped>
.observer {
    height: 0px;
}
</style>
