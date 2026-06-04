<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { MasterClass } from '~/types/master-class';
import { fetchMasterClasses } from '~/api/master-classes';

const masterClasses = ref<MasterClass[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

onMounted(async () => {
    try {
        masterClasses.value = await fetchMasterClasses();
    } catch (e) {
        error.value = 'Не удалось загрузить мастер-классы';
        console.error(e);
    } finally {
        loading.value = false;
    }
});

const activeFilter = ref<string>('all');

const filteredClasses = computed(() => {
    if (activeFilter.value === 'all') return masterClasses.value;
    return masterClasses.value.filter((mc) => {
        // Map difficulty to category for filtering
        if (activeFilter.value === 'watercolor')
            return (
                mc.difficulty === 'beginner' &&
                mc.tags.some((t) => t.slug === 'watercolor')
            );
        if (activeFilter.value === 'oil')
            return mc.difficulty === 'intermediate' || mc.difficulty === 'advanced';
        if (activeFilter.value === 'beginners') return mc.difficulty === 'beginner';
        return true;
    });
});

const filterOptions = [
    { value: 'all', label: 'Все мастер-классы', icon: 'i-heroicons-squares-2x2' },
    { value: 'watercolor', label: 'Акварель', icon: 'i-heroicons-paint-brush' },
    { value: 'oil', label: 'Масло', icon: 'i-heroicons-paint-brush' },
    { value: 'beginners', label: 'Для начинающих', icon: 'i-heroicons-sparkles' },
];

const categoryColors: Record<
    string,
    { bg: string; badge: string; gradient: string; icon: string }
> = {
    watercolor: {
        bg: 'from-blue-50 to-teal-50 dark:from-blue-950/30 dark:to-teal-950/30',
        badge: 'bg-blue-100 text-blue-700 dark:bg-blue-900/50 dark:text-blue-300',
        gradient: 'from-blue-500 to-teal-500',
        icon: 'i-heroicons-paint-brush',
    },
    oil: {
        bg: 'from-amber-50 to-orange-50 dark:from-amber-950/30 dark:to-orange-950/30',
        badge: 'bg-amber-100 text-amber-700 dark:bg-amber-900/50 dark:text-amber-300',
        gradient: 'from-amber-500 to-orange-500',
        icon: 'i-heroicons-paint-brush',
    },
    beginners: {
        bg: 'from-emerald-50 to-green-50 dark:from-emerald-950/30 dark:to-green-950/30',
        badge:
            'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/50 dark:text-emerald-300',
        gradient: 'from-emerald-500 to-green-500',
        icon: 'i-heroicons-sparkles',
    },
};

function getCategoryForMc(mc: MasterClass): string {
    if (mc.tags.some((t) => t.slug === 'watercolor')) return 'watercolor';
    if (mc.tags.some((t) => t.slug === 'oil')) return 'oil';
    return 'beginners';
}

function formatPrice(price: number): string {
    if (price === 0) return 'Бесплатно';
    return `${(price / 100).toLocaleString('ru-RU')} ₽`;
}

function formatDuration(minutes: number): string {
    if (minutes < 60) return `${minutes} мин`;
    const h = Math.floor(minutes / 60);
    const m = minutes % 60;
    return m > 0 ? `${h} ч ${m} мин` : `${h} ч`;
}
</script>

<template>
    <UContainer class="py-8 md:py-12">
        <!-- Hero Section -->
        <div class="mb-12">
            <div
                class="relative rounded-3xl overflow-hidden bg-linear-to-r from-purple-50 to-pink-50 dark:from-gray-800 dark:to-gray-900 p-8 md:p-12"
            >
                <div class="max-w-3xl">
                    <UBreadcrumb
                        :links="[
                            { label: 'Главная', to: '/' },
                            { label: 'Услуги', to: '/services' },
                            { label: 'Мастер-классы' },
                        ]"
                        class="mb-6"
                    />

                    <h1
                        class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4"
                    >
                        Мастер-классы
                    </h1>
                    <p class="text-lg text-gray-700 dark:text-gray-300 mb-6 max-w-2xl">
                        Видеоуроки по рисованию для любого уровня. Каждый мастер-класс —
                        это полноценное занятие с объяснением техник, материалов и
                        пошаговой демонстрацией. Смотрите в удобное время, пересматривайте
                        сколько угодно.
                    </p>

                    <div class="flex flex-wrap gap-3">
                        <UBadge color="primary" variant="soft" size="lg">
                            <UIcon name="i-heroicons-video-camera" class="w-4 h-4 mr-1" />
                            {{ masterClasses.length }} мастер-классов
                        </UBadge>
                        <UBadge color="success" variant="soft" size="lg">
                            <UIcon name="i-heroicons-lock-open" class="w-4 h-4 mr-1" />
                            {{ masterClasses.filter((m) => m.is_free).length }} бесплатных
                        </UBadge>
                        <UBadge color="warning" variant="soft" size="lg">
                            <UIcon name="i-heroicons-clock" class="w-4 h-4 mr-1" />
                            Разные уровни
                        </UBadge>
                    </div>
                </div>

                <div class="absolute right-8 top-8 hidden lg:block">
                    <div
                        class="w-64 h-64 rounded-full bg-linear-to-br from-purple-200 to-pink-200 dark:from-purple-800 dark:to-pink-800 opacity-30"
                    ></div>
                </div>
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="text-center py-16">
            <UIcon
                name="i-heroicons-arrow-path"
                class="w-12 h-12 mx-auto text-gray-300 dark:text-gray-600 mb-4 animate-spin"
            />
            <p class="text-gray-500 dark:text-gray-400">Загрузка мастер-классов...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-16">
            <UIcon
                name="i-heroicons-exclamation-triangle"
                class="w-16 h-16 mx-auto text-red-300 dark:text-red-600 mb-4"
            />
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                Ошибка загрузки
            </h3>
            <p class="text-gray-500 dark:text-gray-400 mb-6">{{ error }}</p>
            <UButton color="primary" variant="outline" @click="window.location.reload()">
                <UIcon name="i-heroicons-arrow-path" class="w-4 h-4 mr-2" />
                Попробовать снова
            </UButton>
        </div>

        <!-- Content -->
        <template v-else>
            <!-- Filter Tabs -->
            <div class="mb-10">
                <div class="flex flex-wrap gap-3">
                    <UButton
                        v-for="option in filterOptions"
                        :key="option.value"
                        :color="activeFilter === option.value ? 'primary' : 'neutral'"
                        :variant="activeFilter === option.value ? 'solid' : 'outline'"
                        @click="activeFilter = option.value"
                    >
                        <UIcon :name="option.icon" class="w-4 h-4 mr-2" />
                        {{ option.label }}
                    </UButton>
                </div>
            </div>

            <!-- Master Classes Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <div
                    v-for="mc in filteredClasses"
                    :key="mc.id"
                    class="group bg-white dark:bg-gray-800 rounded-2xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300 flex flex-col h-full"
                >
                    <!-- Thumbnail -->
                    <div class="relative shrink-0">
                        <div class="w-full h-52 overflow-hidden">
                            <img
                                :src="mc.thumbnail_url"
                                :alt="mc.title_ru"
                                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                            />
                        </div>

                        <!-- Free/Paid Badge -->
                        <div class="absolute top-4 left-4">
                            <UBadge
                                :color="mc.is_free ? 'success' : 'warning'"
                                variant="solid"
                                size="sm"
                            >
                                <UIcon
                                    :name="
                                        mc.is_free
                                            ? 'i-heroicons-lock-open'
                                            : 'i-heroicons-lock-closed'
                                    "
                                    class="w-3.5 h-3.5 mr-1"
                                />
                                {{ formatPrice(mc.price) }}
                            </UBadge>
                        </div>

                        <!-- Duration -->
                        <div
                            class="absolute bottom-4 left-4 bg-black/70 text-white px-3 py-1 rounded-full text-sm flex items-center"
                        >
                            <UIcon name="i-heroicons-clock" class="w-3.5 h-3.5 mr-1" />
                            {{ formatDuration(mc.duration_minutes) }}
                        </div>

                        <!-- Category Badge -->
                        <div class="absolute top-4 right-4">
                            <span
                                :class="categoryColors[getCategoryForMc(mc)]?.badge"
                                class="px-3 py-1 rounded-full text-xs font-medium"
                            >
                                {{
                                    getCategoryForMc(mc) === 'watercolor'
                                        ? 'Акварель'
                                        : getCategoryForMc(mc) === 'oil'
                                        ? 'Масло'
                                        : 'Новичкам'
                                }}
                            </span>
                        </div>

                        <!-- Play overlay -->
                        <div
                            class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                        >
                            <div
                                class="w-16 h-16 rounded-full bg-black/60 flex items-center justify-center shadow-2xl"
                            >
                                <UIcon
                                    name="i-heroicons-play"
                                    class="w-8 h-8 text-white ml-0.5"
                                />
                            </div>
                        </div>
                    </div>

                    <!-- Content -->
                    <div class="p-5 flex flex-col grow">
                        <!-- Tags -->
                        <div class="flex flex-wrap gap-1.5 mb-3">
                            <span
                                v-for="tag in mc.tags.slice(0, 3)"
                                :key="tag.slug"
                                class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400"
                            >
                                {{ tag.name_ru }}
                            </span>
                            <span
                                v-if="mc.tags.length > 3"
                                class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-500"
                            >
                                +{{ mc.tags.length - 3 }}
                            </span>
                        </div>

                        <!-- Title -->
                        <h3
                            class="text-lg font-semibold text-gray-900 dark:text-white mb-2 line-clamp-1"
                        >
                            {{ mc.title_ru }}
                        </h3>

                        <!-- Short Description -->
                        <p
                            class="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2 grow"
                        >
                            {{ mc.short_description_ru }}
                        </p>

                        <!-- Full Description (expandable) -->
                        <p
                            class="text-sm text-gray-500 dark:text-gray-500 mb-4 line-clamp-2 description-text"
                        >
                            {{ mc.description_ru }}
                        </p>

                        <!-- Action Button -->
                        <UButton
                            :to="mc.video_url"
                            :color="mc.is_free ? 'primary' : 'warning'"
                            variant="solid"
                            class="w-full mt-auto"
                        >
                            <UIcon
                                :name="
                                    mc.is_free
                                        ? 'i-heroicons-play'
                                        : 'i-heroicons-shopping-cart'
                                "
                                class="w-4 h-4 mr-2"
                            />
                            {{ mc.is_free ? 'Смотреть бесплатно' : 'Купить доступ' }}
                        </UButton>
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div v-if="filteredClasses.length === 0" class="text-center py-16">
                <UIcon
                    name="i-heroicons-video-camera-slash"
                    class="w-16 h-16 mx-auto text-gray-300 dark:text-gray-600 mb-4"
                />
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                    Нет мастер-классов в этой категории
                </h3>
                <p class="text-gray-500 dark:text-gray-400 mb-6">
                    Попробуйте выбрать другую категорию или вернуться позже.
                </p>
                <UButton color="primary" variant="outline" @click="activeFilter = 'all'">
                    <UIcon name="i-heroicons-squares-2x2" class="w-4 h-4 mr-2" />
                    Показать все
                </UButton>
            </div>
        </template>
    </UContainer>
</template>

<style scoped>
.line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

.line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

.description-text {
    transition: all 0.3s ease;
    cursor: default;
}

.group:hover .description-text {
    -webkit-line-clamp: unset;
    line-clamp: unset;
    overflow: visible;
    white-space: normal;
    background-color: rgba(255, 255, 255, 0.05);
    padding: 0.5rem;
    border-radius: 0.375rem;
    margin: -0.5rem;
}

.dark .group:hover .description-text {
    background-color: rgba(0, 0, 0, 0.1);
}
</style>
