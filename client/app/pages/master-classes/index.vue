<script setup lang="ts">
    import { ref, computed, onMounted, } from 'vue';
    import type { MasterClass, } from '~/types/master-class';
    import {
        fetchMasterClasses,
        fetchMyPurchasedProductIds,
    } from '~/api/master-classes';
    import { useAuthStore, } from '~/stores/AuthStore';

    const router = useRouter();
    const authStore = useAuthStore();
    const masterClasses = ref<MasterClass[]>([],);
    const loading = ref(true,);
    const error = ref<string | null>(null,);
    const purchasedIds = ref<number[]>([],);

    onMounted(async () => {
        try {
            masterClasses.value = await fetchMasterClasses();
        } catch (e) {
            error.value = 'Не удалось загрузить мастер-классы';
            console.error(e,);
        } finally {
            loading.value = false;
        }

        // Загружаем купленные мастер-классы, если пользователь авторизован
        if (authStore.accessToken) {
            purchasedIds.value = await fetchMyPurchasedProductIds();
        }
    });

    const selectedMasterClass = ref<MasterClass | null>(null,);
    const showDetailModal = ref(false,);
    const showVideoPlayer = ref(false,);
    const playingMasterClass = ref<MasterClass | null>(null,);
    const videoError = ref<string | null>(null,);

    /**
     * Проверяет, может ли пользователь смотреть мастер-класс:
     * - бесплатные доступны всем
     * - администратор может смотреть всё
     * - пользователь может смотреть только купленные мастер-классы
     */
    function canWatch(mc: MasterClass,): boolean {
        if (mc.is_free) return true;
        if (authStore.isAdmin) return true;
        return purchasedIds.value.includes(mc.id,);
    }

    function openDetailModal(mc: MasterClass,) {
        selectedMasterClass.value = mc;
        showDetailModal.value = true;
    }

    function closeDetailModal() {
        showDetailModal.value = false;
        // Delay clearing so the modal transition plays out
        setTimeout(() => {
            selectedMasterClass.value = null;
        }, 300,);
    }

    function playMasterClass(mc: MasterClass,) {
        if (!canWatch(mc,)) {
            // Если не куплен — перенаправляем на страницу покупки
            router.push(`/products/${mc.id}`,);
            return;
        }
        playingMasterClass.value = mc;
        videoError.value = null;
        showVideoPlayer.value = true;
    }

    function closeVideoPlayer() {
        showVideoPlayer.value = false;
        setTimeout(() => {
            playingMasterClass.value = null;
            videoError.value = null;
        }, 300,);
    }

    function getVideoSrc(mc: MasterClass,): string {
        const url = `${mc.video_url}`;
        if (!mc.is_free && authStore.accessToken) {
            return `${url}?token=${authStore.accessToken}`;
        }
        return url;
    }

    function getThumbnailSrc(mc: MasterClass,): string {
        return `${mc.thumbnail_url}`;
    }

    const activeFilter = ref<string>('all',);

    const filteredClasses = computed(() => {
        if (activeFilter.value === 'all') return masterClasses.value;
        return masterClasses.value.filter((mc,) => {
            // Map difficulty to category for filtering
            if (activeFilter.value === 'watercolor')
                return (
                    mc.difficulty === 'beginner'
                && mc.tags.some(t => t.slug === 'watercolor',)
                );
            if (activeFilter.value === 'oil')
                return mc.difficulty === 'intermediate' || mc.difficulty === 'advanced';
            if (activeFilter.value === 'beginners') return mc.difficulty === 'beginner';
            return true;
        });
    });

    const filterOptions = [
        { value: 'all', label: 'Все мастер-классы', icon: 'i-heroicons-squares-2x2', },
        { value: 'watercolor', label: 'Акварель', icon: 'i-heroicons-paint-brush', },
        { value: 'oil', label: 'Масло', icon: 'i-heroicons-paint-brush', },
        { value: 'beginners', label: 'Для начинающих', icon: 'i-heroicons-sparkles', },
    ];

    const categoryColors: Record<
        string,
        { bg: string, badge: string, gradient: string, icon: string }
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

    function getCategoryForMc(mc: MasterClass,): string {
        if (mc.tags.some(t => t.slug === 'watercolor',)) return 'watercolor';
        if (mc.tags.some(t => t.slug === 'oil',)) return 'oil';
        return 'beginners';
    }

    function formatPrice(price: number,): string {
        if (price === 0) return 'Бесплатно';
        return `${(price / 100).toLocaleString('ru-RU',)} ₽`;
    }

    function reloadPage() {
        window.location.reload();
    }

    function formatDuration(minutes: number,): string {
        if (minutes < 60) return `${minutes} мин`;
        const h = Math.floor(minutes / 60,);
        const m = minutes % 60;
        return m > 0 ? `${h} ч ${m} мин` : `${h} ч`;
    }
</script>

<template>
    <UContainer class="py-8 md:py-12">
        <!-- Hero Section -->
        <div class="mb-12">
            <div
                class="relative rounded-3xl overflow-hidden bg-linear-to-r from-green-50 to-teal-50 dark:from-gray-800 dark:to-gray-900 p-8 md:p-12"
            >
                <div class="max-w-3xl">
                    <UBreadcrumb
                        :links="[
                            { label: 'Главная', to: '/', },
                            { label: 'Услуги', to: '/services', },
                            { label: 'Мастер-классы', },
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
                        <UBadge
                            color="primary"
                            variant="soft"
                            size="lg"
                        >
                            <UIcon
                                name="i-heroicons-video-camera"
                                class="w-4 h-4 mr-1"
                            />
                            {{ masterClasses.length }} мастер-классов
                        </UBadge>
                        <UBadge
                            color="success"
                            variant="soft"
                            size="lg"
                        >
                            <UIcon
                                name="i-heroicons-lock-open"
                                class="w-4 h-4 mr-1"
                            />
                            {{ masterClasses.filter((m,) => m.is_free,).length }} бесплатных
                        </UBadge>
                        <UBadge
                            color="warning"
                            variant="soft"
                            size="lg"
                        >
                            <UIcon
                                name="i-heroicons-clock"
                                class="w-4 h-4 mr-1"
                            />
                            Разные уровни
                        </UBadge>
                    </div>
                </div>

                <div class="absolute right-8 top-8 hidden lg:block">
                    <div
                        class="w-64 h-64 rounded-full overflow-hidden ring-4 ring-white/50 dark:ring-gray-700/50 shadow-2xl"
                    >
                        <img
                            src="/hero-images/master-hero.jpg"
                            alt="Мастер-класс по рисованию"
                            class="w-full h-full object-cover"
                        >
                    </div>
                </div>
            </div>
        </div>

        <!-- Loading State -->
        <div
            v-if="loading"
            class="text-center py-16"
        >
            <UIcon
                name="i-heroicons-arrow-path"
                class="w-12 h-12 mx-auto text-gray-300 dark:text-gray-600 mb-4 animate-spin"
            />
            <p class="text-gray-500 dark:text-gray-400">
                Загрузка мастер-классов...
            </p>
        </div>

        <!-- Error State -->
        <div
            v-else-if="error"
            class="text-center py-16"
        >
            <UIcon
                name="i-heroicons-exclamation-triangle"
                class="w-16 h-16 mx-auto text-red-300 dark:text-red-600 mb-4"
            />
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                Ошибка загрузки
            </h3>
            <p class="text-gray-500 dark:text-gray-400 mb-6">
                {{ error }}
            </p>
            <UButton
                color="primary"
                variant="outline"
                @click="reloadPage()"
            >
                <UIcon
                    name="i-heroicons-arrow-path"
                    class="w-4 h-4 mr-2"
                />
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
                        <UIcon
                            :name="option.icon"
                            class="w-4 h-4 mr-2"
                        />
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
                                :src="getThumbnailSrc(mc,)"
                                :alt="mc.title_ru"
                                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                            >
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
                                {{ formatPrice(mc.price,) }}
                            </UBadge>
                        </div>

                        <!-- Duration -->
                        <div
                            class="absolute bottom-4 left-4 bg-black/70 text-white px-3 py-1 rounded-full text-sm flex items-center"
                        >
                            <UIcon
                                name="i-heroicons-clock"
                                class="w-3.5 h-3.5 mr-1"
                            />
                            {{ formatDuration(mc.duration_minutes,) }}
                        </div>

                        <!-- Category Badge -->
                        <div class="absolute top-4 right-4">
                            <span
                                :class="categoryColors[getCategoryForMc(mc,)]?.badge"
                                class="px-3 py-1 rounded-full text-xs font-medium"
                            >
                                {{
                                    getCategoryForMc(mc,) === 'watercolor'
                                        ? 'Акварель'
                                        : getCategoryForMc(mc,) === 'oil'
                                            ? 'Масло'
                                            : 'Новичкам'
                                }}
                            </span>
                        </div>

                        <!-- Play overlay (только если можно смотреть) -->
                        <div
                            v-if="canWatch(mc,)"
                            class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 cursor-pointer"
                            @click="playMasterClass(mc,)"
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
                    <div class="p-5 flex flex-col h-full">
                        <!-- Tags -->
                        <div class="flex flex-wrap gap-1.5 mb-3">
                            <span
                                v-for="tag in mc.tags.slice(0, 3,)"
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

                        <!-- Description — 3 lines with fade-out -->
                        <div class="relative mb-3 flex-1 min-h-0">
                            <div class="line-clamp-3">
                                <p class="text-sm text-gray-600 dark:text-gray-400">
                                    {{ mc.short_description_ru }}
                                </p>
                                <p
                                    v-if="mc.description_ru"
                                    class="text-sm text-gray-500 dark:text-gray-500 mt-1"
                                >
                                    {{ mc.description_ru }}
                                </p>
                            </div>
                            <!-- Fade-out gradient at the bottom -->
                            <div
                                class="absolute bottom-0 left-0 right-0 h-8 bg-linear-to-t from-white dark:from-gray-800 to-transparent pointer-events-none"
                            />
                        </div>

                        <!-- Read More Button -->
                        <UButton
                            color="neutral"
                            variant="ghost"
                            size="sm"
                            class="mb-3 self-start group/read"
                            @click="openDetailModal(mc,)"
                        >
                            <span class="text-xs">Читать далее</span>
                            <UIcon
                                name="i-heroicons-arrow-right"
                                class="w-3.5 h-3.5 ml-1 transition-transform duration-200 group-hover/read:translate-x-0.5"
                            />
                        </UButton>

                        <!-- Action Button -->
                        <UButton
                            :color="
                                canWatch(mc,)
                                    ? mc.is_free
                                        ? 'primary'
                                        : 'primary'
                                    : 'warning'
                            "
                            variant="solid"
                            class="w-full shrink-0 mt-auto"
                            @click="playMasterClass(mc,)"
                        >
                            <UIcon
                                :name="
                                    canWatch(mc,)
                                        ? 'i-heroicons-play'
                                        : 'i-heroicons-shopping-cart'
                                "
                                class="w-4 h-4 mr-2"
                            />
                            {{ canWatch(mc,) ? 'Смотреть' : 'Купить' }}
                        </UButton>
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div
                v-if="filteredClasses.length === 0"
                class="text-center py-16"
            >
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
                <UButton
                    color="primary"
                    variant="outline"
                    @click="activeFilter = 'all'"
                >
                    <UIcon
                        name="i-heroicons-squares-2x2"
                        class="w-4 h-4 mr-2"
                    />
                    Показать все
                </UButton>
            </div>
        </template>
    </UContainer>

    <!-- Detail Modal -->
    <UModal
        v-model:open="showDetailModal"
        :transition="true"
        :ui="{
            overlay: 'bg-black/40 backdrop-blur-sm',
            content: 'rounded-3xl shadow-2xl sm:max-w-2xl',
        }"
    >
        <template #header>
            <div class="px-6 pt-6 pb-0">
                <div class="flex items-center gap-2 mb-1">
                    <span
                        v-if="selectedMasterClass"
                        :class="
                            categoryColors[getCategoryForMc(selectedMasterClass,)]?.badge
                        "
                        class="px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                        {{
                            selectedMasterClass
                                ? getCategoryForMc(selectedMasterClass,) === 'watercolor'
                                    ? 'Акварель'
                                    : getCategoryForMc(selectedMasterClass,) === 'oil'
                                        ? 'Масло'
                                        : 'Новичкам'
                                : ''
                        }}
                    </span>
                    <span
                        v-if="selectedMasterClass?.is_free"
                        class="text-xs text-green-600 dark:text-green-400 font-medium"
                    >
                        Бесплатно
                    </span>
                </div>
                <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                    {{ selectedMasterClass?.title_ru }}
                </h2>
            </div>
        </template>

        <template #body>
            <div class="px-6 py-5 space-y-5">
                <!-- Meta info row -->
                <div
                    v-if="selectedMasterClass"
                    class="flex flex-wrap items-center gap-4 text-sm text-gray-500 dark:text-gray-400"
                >
                    <span class="flex items-center gap-1.5">
                        <UIcon
                            name="i-heroicons-clock"
                            class="w-4 h-4"
                        />
                        {{ formatDuration(selectedMasterClass.duration_minutes,) }}
                    </span>
                    <span class="flex items-center gap-1.5">
                        <UIcon
                            name="i-heroicons-eye"
                            class="w-4 h-4"
                        />
                        {{ selectedMasterClass.view_count }} просмотров
                    </span>
                    <span class="flex items-center gap-1.5">
                        <UIcon
                            :name="
                                selectedMasterClass.is_free
                                    ? 'i-heroicons-lock-open'
                                    : 'i-heroicons-lock-closed'
                            "
                            class="w-4 h-4"
                        />
                        {{ formatPrice(selectedMasterClass.price,) }}
                    </span>
                </div>

                <!-- Tags -->
                <div
                    v-if="selectedMasterClass?.tags?.length"
                    class="flex flex-wrap gap-1.5"
                >
                    <span
                        v-for="tag in selectedMasterClass.tags"
                        :key="tag.slug"
                        class="text-xs px-2.5 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400"
                    >
                        {{ tag.name_ru }}
                    </span>
                </div>

                <!-- Full description -->
                <div class="prose prose-sm dark:prose-invert max-w-none">
                    <p class="text-gray-700 dark:text-gray-300 leading-relaxed">
                        {{ selectedMasterClass?.short_description_ru }}
                    </p>
                    <p
                        v-if="selectedMasterClass?.description_ru"
                        class="text-gray-500 dark:text-gray-400 leading-relaxed mt-3"
                    >
                        {{ selectedMasterClass.description_ru }}
                    </p>
                </div>

                <!-- Watch / Buy button in detail modal -->
                <div
                    v-if="selectedMasterClass"
                    class="pt-2"
                >
                    <UButton
                        :color="canWatch(selectedMasterClass,) ? 'primary' : 'warning'"
                        variant="solid"
                        size="lg"
                        class="w-full"
                        @click="
                            closeDetailModal();
                            playMasterClass(selectedMasterClass,);
                        "
                    >
                        <UIcon
                            :name="
                                canWatch(selectedMasterClass,)
                                    ? 'i-heroicons-play'
                                    : 'i-heroicons-shopping-cart'
                            "
                            class="w-5 h-5 mr-2"
                        />
                        {{ canWatch(selectedMasterClass,) ? 'Смотреть' : 'Купить' }}
                    </UButton>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="flex justify-end px-6 pb-6 pt-2 w-full">
                <UButton
                    color="neutral"
                    variant="outline"
                    size="lg"
                    @click="closeDetailModal"
                >
                    Закрыть
                </UButton>
            </div>
        </template>
    </UModal>

    <!-- Video Player Modal -->
    <UModal
        v-model:open="showVideoPlayer"
        :transition="true"
        :ui="{
            overlay: 'bg-black/60 backdrop-blur-sm',
            content: 'rounded-2xl shadow-2xl sm:max-w-4xl overflow-hidden',
        }"
    >
        <template #header>
            <div class="px-6 pt-6 pb-0">
                <h2 class="text-xl font-bold text-gray-900 dark:text-white truncate">
                    {{ playingMasterClass?.title_ru }}
                </h2>
            </div>
        </template>

        <template #body>
            <div class="px-6 py-5">
                <div
                    v-if="videoError"
                    class="bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 p-4 rounded-xl text-center"
                >
                    <UIcon
                        name="i-heroicons-exclamation-triangle"
                        class="w-8 h-8 mx-auto mb-2"
                    />
                    <p>{{ videoError }}</p>
                </div>
                <div
                    v-else-if="playingMasterClass"
                    class="relative bg-black rounded-xl overflow-hidden"
                >
                    <video
                        :key="playingMasterClass.id"
                        controls
                        autoplay
                        class="w-full max-h-[70vh]"
                        :poster="getThumbnailSrc(playingMasterClass,)"
                    >
                        <source
                            :src="getVideoSrc(playingMasterClass,)"
                            type="video/mp4"
                        >
                        Ваш браузер не поддерживает воспроизведение видео.
                    </video>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="flex justify-end px-6 pb-6 pt-2 w-full">
                <UButton
                    color="neutral"
                    variant="soft"
                    size="lg"
                    @click="closeVideoPlayer"
                >
                    <UIcon
                        name="i-heroicons-x-mark"
                        class="w-5 h-5 mr-2"
                    />
                    Закрыть
                </UButton>
            </div>
        </template>
    </UModal>
</template>

<style scoped>
/* Modal body scroll styling */
:deep(.UModal body) {
    scrollbar-width: thin;
    scrollbar-color: #d1d5db transparent;
}

:deep(.UModal body::-webkit-scrollbar) {
    width: 6px;
}

:deep(.UModal body::-webkit-scrollbar-thumb) {
    background-color: #d1d5db;
    border-radius: 3px;
}
</style>
