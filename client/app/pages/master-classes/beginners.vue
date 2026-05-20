<script setup lang="ts">
import { ref } from 'vue';

const videos = ref([
    {
        id: 1,
        title: 'Первые шаги в рисовании: линии и формы',
        description:
            'Научитесь уверенно держать карандаш, создавать базовые формы и понимать пропорции.',
        duration: '28 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1541961017774-22349e4a1262?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: false,
    },
    {
        id: 2,
        title: 'Основы композиции и перспективы',
        description:
            'Поймете, как располагать объекты на листе и создавать иллюзию глубины.',
        duration: '35 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1578301978693-85fa9c0320b9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: false,
    },
    {
        id: 3,
        title: 'Работа с тоном и светотенью',
        description:
            'Научитесь передавать объем с помощью света и тени на простых объектах.',
        duration: '40 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1542744095-fcf48d80b0fd?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: true,
    },
    {
        id: 4,
        title: 'Рисование простых натюрмортов',
        description:
            'Создайте свою первую законченную работу, используя все изученные техники.',
        duration: '48 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1579783902614-a3fb3927b6a5?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: true,
    },
]);

const materials = ref([
    'Бумага для рисования (А4, А3)',
    'Графитные карандаши разной мягкости (2H, HB, 2B, 4B)',
    'Ластик и клячка',
    'Точилка для карандашей',
    'Планшет для рисования (опционально)',
    'Линейка (для начальных упражнений)',
]);
</script>

<template>
    <UContainer class="py-8 md:py-12">
        <!-- Hero Section -->
        <div class="hero-section mb-12">
            <div
                class="relative rounded-3xl overflow-hidden bg-linear-to-r from-emerald-50 to-green-50 dark:from-gray-800 dark:to-gray-900 p-8 md:p-12"
            >
                <div class="max-w-3xl">
                    <UBreadcrumb
                        :links="[
                            { label: 'Главная', to: '/' },
                            { label: 'Услуги', to: '/services' },
                            { label: 'Мастер-классы', to: '#' },
                            { label: 'Рисование для начинающих' },
                        ]"
                        class="mb-6"
                    />

                    <h1
                        class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4"
                    >
                        Рисование для начинающих
                    </h1>
                    <p class="text-lg text-gray-700 dark:text-gray-300 mb-6 max-w-2xl">
                        Начните свой творческий путь с нуля. Этот курс создан для тех, кто
                        всегда хотел рисовать, но не знал, с чего начать. Мы пройдем весь
                        путь от простых линий до законченных работ.
                    </p>

                    <div class="flex flex-wrap gap-4">
                        <UBadge color="primary" variant="soft" size="lg">
                            <UIcon name="i-heroicons-clock" class="w-4 h-4 mr-1" />
                            12 занятий
                        </UBadge>
                        <UBadge color="success" variant="soft" size="lg">
                            <UIcon name="i-heroicons-user-group" class="w-4 h-4 mr-1" />
                            С нуля
                        </UBadge>
                        <UBadge color="secondary" variant="soft" size="lg">
                            <UIcon name="i-heroicons-video-camera" class="w-4 h-4 mr-1" />
                            Пошаговые уроки
                        </UBadge>
                    </div>
                </div>

                <div class="absolute right-8 top-8 hidden lg:block">
                    <div
                        class="w-64 h-64 rounded-full bg-linear-to-br from-emerald-200 to-green-200 dark:from-emerald-800 dark:to-green-800 opacity-30"
                    ></div>
                </div>
            </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <!-- Main Content -->
            <div class="lg:col-span-2">
                <!-- Video Gallery -->
                <div class="mb-12">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                            Видео уроки
                        </h2>
                        <span class="text-sm text-gray-500 dark:text-gray-400">
                            Доступно {{ videos.filter((v) => !v.isLocked).length }} из
                            {{ videos.length }}
                        </span>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div
                            v-for="video in videos"
                            :key="video.id"
                            class="group bg-white dark:bg-gray-800 rounded-2xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300 flex flex-col h-full card-height"
                        >
                            <div class="relative shrink-0">
                                <div class="image-container">
                                    <img
                                        :src="video.thumbnail"
                                        :alt="video.title"
                                        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                                    />
                                </div>
                                <div class="absolute top-4 right-4">
                                    <UBadge
                                        :color="video.isLocked ? 'neutral' : 'success'"
                                        variant="solid"
                                    >
                                        {{
                                            video.isLocked ? 'Заблокировано' : 'Доступно'
                                        }}
                                    </UBadge>
                                </div>
                                <div
                                    class="absolute bottom-4 left-4 bg-black/70 text-white px-3 py-1 rounded-full text-sm"
                                >
                                    {{ video.duration }}
                                </div>
                                <div
                                    v-if="video.isLocked"
                                    class="absolute inset-0 bg-black/60 flex items-center justify-center"
                                >
                                    <UIcon
                                        name="i-heroicons-lock-closed"
                                        class="w-12 h-12 text-white"
                                    />
                                </div>
                            </div>

                            <div class="p-6 flex flex-col grow">
                                <div class="grow">
                                    <h3
                                        class="text-xl font-semibold text-gray-900 dark:text-white mb-2 line-clamp-1"
                                    >
                                        {{ video.title }}
                                    </h3>
                                    <p
                                        class="text-gray-600 dark:text-gray-400 mb-4 line-clamp-2 description-text"
                                    >
                                        {{ video.description }}
                                    </p>
                                </div>
                                <UButton
                                    :to="video.videoUrl"
                                    :disabled="video.isLocked"
                                    :color="video.isLocked ? 'neutral' : 'primary'"
                                    variant="solid"
                                    class="w-full mt-auto"
                                >
                                    <UIcon
                                        :name="
                                            video.isLocked
                                                ? 'i-heroicons-lock-closed'
                                                : 'i-heroicons-play'
                                        "
                                        class="w-5 h-5 mr-2"
                                    />
                                    {{
                                        video.isLocked
                                            ? 'Доступно после покупки'
                                            : 'Смотреть сейчас'
                                    }}
                                </UButton>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Description -->
                <div class="mb-12">
                    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                        О курсе
                    </h2>
                    <div class="prose prose-lg dark:prose-invert max-w-none">
                        <p class="text-gray-700 dark:text-gray-300 mb-4">
                            Этот курс специально разработан для тех, кто никогда не
                            рисовал или имеет минимальный опыт. Мы не будем использовать
                            сложные термины или требовать специальных навыков — только
                            пошаговые инструкции, практические упражнения и поддержка на
                            каждом этапе.
                        </p>
                        <p class="text-gray-700 dark:text-gray-300 mb-4">
                            Вы научитесь не просто копировать, а понимать, как устроено
                            изображение, развивать глазомер и уверенность в своих
                            движениях. Каждое занятие — это новый шаг к умению выражать
                            свои идеи через рисунок.
                        </p>
                        <ul class="text-gray-700 dark:text-gray-300 mb-4 space-y-2">
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Преодолеете страх чистого листа и научитесь начинать
                                    рисунок</span
                                >
                            </li>
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Освоите базовые техники: линии, штриховка, работа с
                                    тоном</span
                                >
                            </li>
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Поймете принципы композиции, перспективы и
                                    пропорций</span
                                >
                            </li>
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Создадите серию работ, которые покажут ваш
                                    прогресс</span
                                >
                            </li>
                        </ul>
                    </div>
                </div>

                <!-- For Whom -->
                <div class="mb-12">
                    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                        Для кого этот курс
                    </h2>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div
                            class="bg-linear-to-r from-emerald-50 to-green-50 dark:from-gray-800 dark:to-gray-800 rounded-2xl p-6"
                        >
                            <h3
                                class="text-xl font-bold text-gray-900 dark:text-white mb-3 flex items-center"
                            >
                                <UIcon
                                    name="i-heroicons-sparkles"
                                    class="w-6 h-6 mr-2 text-emerald-500"
                                />
                                Абсолютные новички
                            </h3>
                            <p class="text-gray-700 dark:text-gray-300">
                                Никогда не держали карандаш для рисования? Идеально! Мы
                                начнем с самых основ и будем двигаться в комфортном темпе.
                            </p>
                        </div>
                        <div
                            class="bg-linear-to-r from-emerald-50 to-green-50 dark:from-gray-800 dark:to-gray-800 rounded-2xl p-6"
                        >
                            <h3
                                class="text-xl font-bold text-gray-900 dark:text-white mb-3 flex items-center"
                            >
                                <UIcon
                                    name="i-heroicons-arrow-path"
                                    class="w-6 h-6 mr-2 text-emerald-500"
                                />
                                Возвращающиеся к рисованию
                            </h3>
                            <p class="text-gray-700 dark:text-gray-300">
                                Рисовали в детстве, но давно забросили? Курс поможет
                                систематизировать знания и вернуть навыки.
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Sidebar -->
            <div class="lg:col-span-1">
                <!-- Materials -->
                <div
                    class="bg-linear-to-br from-emerald-50 to-green-50 dark:from-gray-800 dark:to-gray-900 rounded-2xl p-6 mb-8"
                >
                    <h3
                        class="text-xl font-bold text-gray-900 dark:text-white mb-4 flex items-center"
                    >
                        <UIcon
                            name="i-heroicons-shopping-bag"
                            class="w-6 h-6 mr-2 text-emerald-500"
                        />
                        Необходимые материалы
                    </h3>
                    <ul class="space-y-3">
                        <li
                            v-for="(material, index) in materials"
                            :key="index"
                            class="flex items-center text-gray-700 dark:text-gray-300"
                        >
                            <UIcon
                                name="i-heroicons-check"
                                class="w-5 h-5 text-green-500 mr-3"
                            />
                            {{ material }}
                        </li>
                    </ul>
                    <UButton color="primary" variant="ghost" class="w-full mt-6">
                        <UIcon name="i-heroicons-shopping-cart" class="w-5 h-5 mr-2" />
                        Где купить материалы
                    </UButton>
                </div>

                <!-- Progress -->
                <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg mb-8">
                    <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">
                        Ваш прогресс
                    </h3>
                    <div class="space-y-4">
                        <div>
                            <div class="flex justify-between mb-1">
                                <span
                                    class="text-sm font-medium text-gray-700 dark:text-gray-300"
                                    >Базовые техники</span
                                >
                                <span class="text-sm font-medium text-emerald-600"
                                    >2/4 уроков</span
                                >
                            </div>
                            <div
                                class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2"
                            >
                                <div
                                    class="bg-emerald-500 h-2 rounded-full"
                                    style="width: 50%"
                                ></div>
                            </div>
                        </div>
                        <div>
                            <div class="flex justify-between mb-1">
                                <span
                                    class="text-sm font-medium text-gray-700 dark:text-gray-300"
                                    >Практические задания</span
                                >
                                <span class="text-sm font-medium text-emerald-600"
                                    >3/6 заданий</span
                                >
                            </div>
                            <div
                                class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2"
                            >
                                <div
                                    class="bg-emerald-500 h-2 rounded-full"
                                    style="width: 50%"
                                ></div>
                            </div>
                        </div>
                        <div>
                            <div class="flex justify-between mb-1">
                                <span
                                    class="text-sm font-medium text-gray-700 dark:text-gray-300"
                                    >Финальный проект</span
                                >
                                <span class="text-sm font-medium text-gray-500"
                                    >Не начат</span
                                >
                            </div>
                            <div
                                class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2"
                            >
                                <div
                                    class="bg-gray-300 dark:bg-gray-600 h-2 rounded-full"
                                    style="width: 0%"
                                ></div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- CTA -->
                <div
                    class="bg-linear-to-r from-emerald-500 to-green-500 rounded-2xl p-6 text-white"
                >
                    <h3 class="text-xl font-bold mb-2">Начните рисовать сегодня</h3>
                    <p class="mb-4 opacity-90">Полный курс для начинающих с нуля</p>
                    <div class="text-3xl font-bold mb-2">8 900 ₽</div>
                    <p class="text-sm opacity-80 mb-6">или 2 970 ₽/мес при рассрочке</p>
                    <UButton color="primary" variant="solid" class="w-full font-bold">
                        <UIcon name="i-heroicons-shopping-cart" class="w-5 h-5 mr-2" />
                        Приобрести курс
                    </UButton>
                    <p class="text-xs opacity-80 mt-4 text-center">
                        Доступ навсегда • Проверка работ • Сертификат
                    </p>
                </div>
            </div>
        </div>
    </UContainer>
</template>

<style scoped>
.hero-section {
    position: relative;
    overflow: hidden;
}

.group:hover .group-hover\:scale-105 {
    transform: scale(1.05);
}

/* Text truncation utilities */
.line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

.line-clamp-2 {
    display: -webkit-box;
    line-clamp: 2;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

/* Description text with hover effect */
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

/* Consistent image container */
.image-container {
    width: 100%;
    height: 12rem;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
}

.image-container img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
}

/* Card height consistency */
.card-height {
    min-height: 24rem;
}

@media (max-width: 768px) {
    .hero-section {
        padding: 2rem 1rem;
    }

    .card-height {
        min-height: 22rem;
    }
}
</style>
