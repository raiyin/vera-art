<script setup lang="ts">
import { ref } from 'vue';

const videos = ref([
    {
        id: 1,
        title: 'Основы работы с масляными красками',
        description:
            'Познакомьтесь с материалами, научитесь правильно смешивать краски и наносить их на холст.',
        duration: '32 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1541961017774-22349e4a1262?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: false,
    },
    {
        id: 2,
        title: 'Техника импасто: объемные мазки',
        description:
            'Освойте технику густого нанесения краски для создания выразительной фактуры.',
        duration: '45 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1578301978693-85fa9c0320b9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: false,
    },
    {
        id: 3,
        title: 'Портрет маслом: свет и тень',
        description: 'Научитесь передавать объем и характер в портретной живописи.',
        duration: '68 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1542744095-fcf48d80b0fd?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: true,
    },
    {
        id: 4,
        title: 'Пейзаж с атмосферной перспективой',
        description:
            'Создайте глубину пространства и передайте настроение в пейзажной живописи.',
        duration: '52 мин',
        thumbnail:
            'https://images.unsplash.com/photo-1579783902614-a3fb3927b6a5?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
        videoUrl: '#',
        isLocked: true,
    },
]);

const materials = ref([
    'Холст на подрамнике',
    'Набор масляных красок (основные цвета)',
    'Кисти: щетина, синтетика, плоские и круглые',
    'Масло льняное или пинен',
    'Палитра (деревянная или пластиковая)',
    'Масленка для разбавителя',
    'Тряпки для вытирания кистей',
]);
</script>

<template>
    <UContainer class="py-8 md:py-12">
        <!-- Hero Section -->
        <div class="hero-section mb-12">
            <div
                class="relative rounded-3xl overflow-hidden bg-linear-to-r from-amber-50 to-orange-50 dark:from-gray-800 dark:to-gray-900 p-8 md:p-12"
            >
                <div class="max-w-3xl">
                    <UBreadcrumb
                        :links="[
                            { label: 'Главная', to: '/' },
                            { label: 'Услуги', to: '/services' },
                            { label: 'Мастер-классы', to: '#' },
                            { label: 'Масляная живопись' },
                        ]"
                        class="mb-6"
                    />

                    <h1
                        class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4"
                    >
                        Масляная живопись
                    </h1>
                    <p class="text-lg text-gray-700 dark:text-gray-300 mb-6 max-w-2xl">
                        Откройте для себя мир насыщенных цветов и бесконечных
                        возможностей. Масляные краски позволяют создавать работы с
                        глубиной, фактурой и долговечностью.
                    </p>

                    <div class="flex flex-wrap gap-4">
                        <UBadge color="primary" variant="soft" size="lg">
                            <UIcon name="i-heroicons-clock" class="w-4 h-4 mr-1" />
                            10 занятий
                        </UBadge>
                        <UBadge color="success" variant="soft" size="lg">
                            <UIcon name="i-heroicons-user-group" class="w-4 h-4 mr-1" />
                            Средний уровень
                        </UBadge>
                        <UBadge color="secondary" variant="soft" size="lg">
                            <UIcon name="i-heroicons-video-camera" class="w-4 h-4 mr-1" />
                            Онлайн + оффлайн
                        </UBadge>
                    </div>
                </div>

                <div class="absolute right-8 top-8 hidden lg:block">
                    <div
                        class="w-64 h-64 rounded-full bg-linear-to-br from-amber-200 to-orange-200 dark:from-amber-800 dark:to-orange-800 opacity-30"
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
                            Видео мастер-классы
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
                        О мастер-классе
                    </h2>
                    <div class="prose prose-lg dark:prose-invert max-w-none">
                        <p class="text-gray-700 dark:text-gray-300 mb-4">
                            Масляная живопись — классическая техника, которая веками
                            использовалась великими мастерами. Её главные преимущества —
                            медленное высыхание, позволяющее вносить изменения,
                            насыщенность цвета и возможность создавать сложные фактуры.
                        </p>
                        <p class="text-gray-700 dark:text-gray-300 mb-4">
                            На курсе вы пройдете путь от основ до создания сложных
                            композиций. Мы уделим внимание не только техническим аспектам,
                            но и развитию художественного видения и индивидуального стиля.
                        </p>
                        <ul class="text-gray-700 dark:text-gray-300 mb-4 space-y-2">
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Освоите различные техники: лессировка, алла прима,
                                    импасто</span
                                >
                            </li>
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Научитесь работать с цветом, светом и тенью для
                                    создания объема</span
                                >
                            </li>
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Поймете принципы композиции и построения пространства
                                    в картине</span
                                >
                            </li>
                            <li class="flex items-start">
                                <UIcon
                                    name="i-heroicons-check-circle"
                                    class="w-5 h-5 text-green-500 mr-2 mt-1 shrink-0"
                                />
                                <span
                                    >Создадите несколько законченных работ в разных
                                    жанрах</span
                                >
                            </li>
                        </ul>
                    </div>
                </div>

                <!-- Techniques -->
                <div class="mb-12">
                    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                        Техники масляной живописи
                    </h2>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div
                            class="bg-linear-to-r from-amber-50 to-orange-50 dark:from-gray-800 dark:to-gray-800 rounded-2xl p-6"
                        >
                            <h3
                                class="text-xl font-bold text-gray-900 dark:text-white mb-3 flex items-center"
                            >
                                <UIcon
                                    name="i-heroicons-paint-brush"
                                    class="w-6 h-6 mr-2 text-amber-500"
                                />
                                Алла прима
                            </h3>
                            <p class="text-gray-700 dark:text-gray-300">
                                Техника «за один раз» — работа выполняется за один сеанс,
                                пока краска еще влажная. Идеально для этюдов и быстрых
                                набросков.
                            </p>
                        </div>
                        <div
                            class="bg-linear-to-r from-amber-50 to-orange-50 dark:from-gray-800 dark:to-gray-800 rounded-2xl p-6"
                        >
                            <h3
                                class="text-xl font-bold text-gray-900 dark:text-white mb-3 flex items-center"
                            >
                                <UIcon
                                    name="i-heroicons-layers"
                                    class="w-6 h-6 mr-2 text-amber-500"
                                />
                                Лессировка
                            </h3>
                            <p class="text-gray-700 dark:text-gray-300">
                                Нанесение полупрозрачных слоев краски поверх высохшего
                                слоя. Позволяет создавать глубину цвета и сложные световые
                                эффекты.
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Sidebar -->
            <div class="lg:col-span-1">
                <!-- Materials -->
                <div
                    class="bg-linear-to-br from-amber-50 to-orange-50 dark:from-gray-800 dark:to-gray-900 rounded-2xl p-6 mb-8"
                >
                    <h3
                        class="text-xl font-bold text-gray-900 dark:text-white mb-4 flex items-center"
                    >
                        <UIcon
                            name="i-heroicons-shopping-bag"
                            class="w-6 h-6 mr-2 text-amber-500"
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

                <!-- Schedule -->
                <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg mb-8">
                    <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">
                        Расписание занятий
                    </h3>
                    <div class="space-y-4">
                        <div class="flex items-center">
                            <UIcon
                                name="i-heroicons-calendar"
                                class="w-5 h-5 text-amber-500 mr-3"
                            />
                            <div>
                                <div class="font-medium text-gray-900 dark:text-white">
                                    Старт курса
                                </div>
                                <div class="text-sm text-gray-600 dark:text-gray-400">
                                    Каждый понедельник
                                </div>
                            </div>
                        </div>
                        <div class="flex items-center">
                            <UIcon
                                name="i-heroicons-clock"
                                class="w-5 h-5 text-amber-500 mr-3"
                            />
                            <div>
                                <div class="font-medium text-gray-900 dark:text-white">
                                    Длительность
                                </div>
                                <div class="text-sm text-gray-600 dark:text-gray-400">
                                    10 недель, 2 занятия в неделю
                                </div>
                            </div>
                        </div>
                        <div class="flex items-center">
                            <UIcon
                                name="i-heroicons-user-group"
                                class="w-5 h-5 text-amber-500 mr-3"
                            />
                            <div>
                                <div class="font-medium text-gray-900 dark:text-white">
                                    Формат
                                </div>
                                <div class="text-sm text-gray-600 dark:text-gray-400">
                                    Онлайн-трансляции + записи
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- CTA -->
                <div
                    class="bg-linear-to-r from-amber-500 to-orange-500 rounded-2xl p-6 text-white"
                >
                    <h3 class="text-xl font-bold mb-2">Начните рисовать маслом</h3>
                    <p class="mb-4 opacity-90">Полный курс с поддержкой преподавателя</p>
                    <div class="text-3xl font-bold mb-2">15 900 ₽</div>
                    <p class="text-sm opacity-80 mb-6">или 5 300 ₽/мес при рассрочке</p>
                    <UButton color="primary" variant="solid" class="w-full font-bold">
                        <UIcon name="i-heroicons-shopping-cart" class="w-5 h-5 mr-2" />
                        Приобрести курс
                    </UButton>
                    <p class="text-xs opacity-80 mt-4 text-center">
                        Доступ навсегда • Разбор работ • Сертификат
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
