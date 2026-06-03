<script setup lang="ts">
import type { TypedGetWorkDto } from '~/types/work';
import Gallery from '@/components/PicGallery.vue';
import { ref, onMounted } from 'vue';

const config = useRuntimeConfig();
const SERVER_URL = config.public.serverUrl;
const limit = config.public.limit;

const works = ref<TypedGetWorkDto[]>([]);
const page = ref(-1);
const worksObserver = ref<Element | null>(null);

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
        const response = await fetch(`${SERVER_URL}works?${params}`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const data = await response.json();

        const newWorks = data.map((work: any) => ({
            ...work,
            __type: 'GetWorkDto',
        }));
        works.value = [...works.value, ...newWorks];
        console.log(works.value);
    } catch (e) {
        console.error('Error fetching works on allworks page ' + e);
    }
};

const handleWorkDeleted = (id: string) => {
    works.value = works.value.filter(
        (work: TypedGetWorkDto) => work.id !== id && String(work.id) !== id
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
    const worksCallback = (entries: IntersectionObserverEntry[]) => {
        if (entries[0]?.isIntersecting) {
            loadWorks();
        }
    };
    const observer = new IntersectionObserver(worksCallback, options);
    if (worksObserver.value) observer.observe(worksObserver.value);
});
</script>

<template>
    <UContainer class="main-content">
        <Gallery :images="works" @work-deleted="handleWorkDeleted" />

        <div ref="worksObserver" class="observer" />
    </UContainer>
</template>

<style scoped>
.observer {
    height: 0px;
}
</style>
