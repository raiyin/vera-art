<script setup lang="ts">
    import { ref, } from 'vue';

    interface CoursePlaylist {
        id: number
        title: string
        description: string
        level: 'beginner' | 'intermediate' | 'advanced'
        duration: string
        lessonsCount: number
        price: number
        discountPrice?: number
        instructor: string
        rating: number
        videos: {
            id: number
            title: string
            duration: string
            thumbnail: string
            isLocked: boolean
        }[]
    }

    const playlists = ref<CoursePlaylist[]>([
        {
            id: 1,
            title: 'Основы акварельной живописи',
            description:
                'Полный курс для начинающих: от базовых техник до создания сложных композиций',
            level: 'beginner',
            duration: '12 часов',
            lessonsCount: 24,
            price: 8900,
            discountPrice: 6900,
            instructor: 'Анна Петрова',
            rating: 4.8,
            videos: [
                {
                    id: 1,
                    title: 'Введение в акварель: материалы и инструменты',
                    duration: '25:10',
                    thumbnail: '/images/courses/watercolor-1.jpg',
                    isLocked: false,
                },
                {
                    id: 2,
                    title: 'Основные техники: заливка, лессировка',
                    duration: '32:45',
                    thumbnail: '/images/courses/watercolor-2.jpg',
                    isLocked: false,
                },
                {
                    id: 3,
                    title: 'Работа с цветом: создание палитры',
                    duration: '28:20',
                    thumbnail: '/images/courses/watercolor-3.jpg',
                    isLocked: true,
                },
                {
                    id: 4,
                    title: 'Пейзаж акварелью: небо и облака',
                    duration: '45:15',
                    thumbnail: '/images/courses/watercolor-4.jpg',
                    isLocked: true,
                },
                {
                    id: 5,
                    title: 'Ботаническая иллюстрация: цветы и листья',
                    duration: '38:30',
                    thumbnail: '/images/courses/watercolor-5.jpg',
                    isLocked: true,
                },
                {
                    id: 6,
                    title: 'Портрет акварелью: основы',
                    duration: '52:10',
                    thumbnail: '/images/courses/watercolor-6.jpg',
                    isLocked: true,
                },
            ],
        },
        {
            id: 2,
            title: 'Масляная живопись: от классики к современности',
            description:
                'Освойте технику масляной живописи от натяжки холста до финального лакирования',
            level: 'intermediate',
            duration: '18 часов',
            lessonsCount: 30,
            price: 11900,
            discountPrice: 9900,
            instructor: 'Иван Смирнов',
            rating: 4.9,
            videos: [
                {
                    id: 1,
                    title: 'Подготовка материалов: холсты, краски, медиумы',
                    duration: '30:15',
                    thumbnail: '/images/courses/oil-1.jpg',
                    isLocked: false,
                },
                {
                    id: 2,
                    title: 'Техника алла прима: быстрая живопись',
                    duration: '35:40',
                    thumbnail: '/images/courses/oil-2.jpg',
                    isLocked: false,
                },
                {
                    id: 3,
                    title: 'Лессировка: создание глубины и объема',
                    duration: '42:20',
                    thumbnail: '/images/courses/oil-3.jpg',
                    isLocked: true,
                },
                {
                    id: 4,
                    title: 'Портрет маслом: анатомия и пропорции',
                    duration: '55:30',
                    thumbnail: '/images/courses/oil-4.jpg',
                    isLocked: true,
                },
                {
                    id: 5,
                    title: 'Натюрморт с драпировкой',
                    duration: '48:15',
                    thumbnail: '/images/courses/oil-5.jpg',
                    isLocked: true,
                },
                {
                    id: 6,
                    title: 'Абстрактная живопись: экспрессия и цвет',
                    duration: '40:50',
                    thumbnail: '/images/courses/oil-6.jpg',
                    isLocked: true,
                },
            ],
        },
        {
            id: 3,
            title: 'Цифровая иллюстрация в Procreate',
            description:
                'Создавайте профессиональные иллюстрации на iPad от скетча до финального рендера',
            level: 'beginner',
            duration: '15 часов',
            lessonsCount: 20,
            price: 9900,
            discountPrice: 7900,
            instructor: 'Мария Козлова',
            rating: 4.7,
            videos: [
                {
                    id: 1,
                    title: 'Знакомство с Procreate: интерфейс и кисти',
                    duration: '22:30',
                    thumbnail: '/images/courses/digital-1.jpg',
                    isLocked: false,
                },
                {
                    id: 2,
                    title: 'Скетчинг: от идеи до эскиза',
                    duration: '28:15',
                    thumbnail: '/images/courses/digital-2.jpg',
                    isLocked: false,
                },
                {
                    id: 3,
                    title: 'Работа со слоями и масками',
                    duration: '35:40',
                    thumbnail: '/images/courses/digital-3.jpg',
                    isLocked: true,
                },
                {
                    id: 4,
                    title: 'Колористика и светотень',
                    duration: '40:20',
                    thumbnail: '/images/courses/digital-4.jpg',
                    isLocked: true,
                },
                {
                    id: 5,
                    title: 'Создание персонажей: от концепта до детализации',
                    duration: '50:10',
                    thumbnail: '/images/courses/digital-5.jpg',
                    isLocked: true,
                },
                {
                    id: 6,
                    title: 'Анимация в Procreate: оживляем иллюстрации',
                    duration: '45:30',
                    thumbnail: '/images/courses/digital-6.jpg',
                    isLocked: true,
                },
            ],
        },
        {
            id: 4,
            title: 'Академический рисунок: от гипсовых голов до натуры',
            description:
                'Классическая школа рисунка для серьезного развития художественных навыков',
            level: 'advanced',
            duration: '25 часов',
            lessonsCount: 35,
            price: 14900,
            discountPrice: 12900,
            instructor: 'Сергей Волков',
            rating: 4.9,
            videos: [
                {
                    id: 1,
                    title: 'Основы композиции и построения',
                    duration: '33:25',
                    thumbnail: '/images/courses/drawing-1.jpg',
                    isLocked: false,
                },
                {
                    id: 2,
                    title: 'Рисунок гипсовой головы: пропорции и объем',
                    duration: '40:15',
                    thumbnail: '/images/courses/drawing-2.jpg',
                    isLocked: false,
                },
                {
                    id: 3,
                    title: 'Анатомия для художников: скелет и мышцы',
                    duration: '48:30',
                    thumbnail: '/images/courses/drawing-3.jpg',
                    isLocked: true,
                },
                {
                    id: 4,
                    title: 'Портрет с натуры: работа с моделью',
                    duration: '55:45',
                    thumbnail: '/images/courses/drawing-4.jpg',
                    isLocked: true,
                },
                {
                    id: 5,
                    title: 'Фигура человека в движении',
                    duration: '52:20',
                    thumbnail: '/images/courses/drawing-5.jpg',
                    isLocked: true,
                },
                {
                    id: 6,
                    title: 'Сложные композиции: групповые портреты',
                    duration: '60:10',
                    thumbnail: '/images/courses/drawing-6.jpg',
                    isLocked: true,
                },
            ],
        },
    ]);

    const selectedPlaylist = ref<CoursePlaylist | null>(null,);

    const openPlaylist = (playlist: CoursePlaylist,) => {
        selectedPlaylist.value = playlist;
    };

    const closePlaylist = () => {
        selectedPlaylist.value = null;
    };

    const getLevelColor = (level: string,) => {
        switch (level) {
    case 'beginner':
        return 'success';
    case 'intermediate':
        return 'warning';
    case 'advanced':
        return 'error';
    default:
        return 'neutral';
        }
    };

    const getLevelText = (level: string,) => {
        switch (level) {
    case 'beginner':
        return 'Начинающий';
    case 'intermediate':
        return 'Средний';
    case 'advanced':
        return 'Продвинутый';
    default:
        return level;
        }
    };
</script>

<template>
    <UContainer class="py-8 md:py-12">
        <!-- Breadcrumb -->
        <UBreadcrumb class="mb-8">
            <UBreadcrumbItem to="/">
                Главная
            </UBreadcrumbItem>
            <UBreadcrumbItem to="/services">
                Услуги
            </UBreadcrumbItem>
            <UBreadcrumbItem>Онлайн-курсы</UBreadcrumbItem>
        </UBreadcrumb>

        <!-- Hero Section -->
        <div
            class="relative overflow-hidden rounded-2xl bg-linear-to-r from-blue-600 to-purple-600 p-8 md:p-12 mb-12"
        >
            <div class="relative z-10 max-w-2xl">
                <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">
                    Онлайн-курсы по живописи и рисунку
                </h1>
                <p class="text-xl text-blue-100 mb-6">
                    Профессиональное обучение от ведущих художников. Смотрите первые два
                    видео каждого курса бесплатно!
                </p>
                <div class="flex flex-wrap gap-4">
                    <UButton
                        variant="solid"
                        size="lg"
                        class="font-semibold bg-white text-blue-600 hover:bg-gray-100"
                    >
                        <UIcon
                            name="i-heroicons-play-circle"
                            class="mr-2"
                        />
                        Смотреть демо-урок
                    </UButton>
                    <UButton
                        variant="outline"
                        size="lg"
                        class="font-semibold text-white border-white hover:bg-white/10"
                    >
                        Как приобрести курс
                    </UButton>
                </div>
            </div>
            <div class="absolute right-0 top-0 bottom-0 w-1/3 opacity-20">
                <div class="h-full bg-linear-to-l from-white to-transparent" />
            </div>
        </div>

        <!-- How to Purchase Section -->
        <div
            class="bg-linear-to-br from-gray-50 to-white rounded-2xl p-8 mb-12 border border-gray-200"
        >
            <h2 class="text-3xl font-bold text-gray-900 mb-6">
                Как приобрести видео-курсы
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
                    <div
                        class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4"
                    >
                        <UIcon
                            name="i-heroicons-shopping-cart"
                            class="text-blue-600 text-2xl"
                        />
                    </div>
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Выберите курс
                    </h3>
                    <p class="text-gray-600">
                        Просмотрите доступные курсы и выберите подходящий по уровню и
                        тематике. Первые два видео каждого курса доступны бесплатно.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
                    <div
                        class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-4"
                    >
                        <UIcon
                            name="i-heroicons-credit-card"
                            class="text-green-600 text-2xl"
                        />
                    </div>
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Оплатите онлайн
                    </h3>
                    <p class="text-gray-600">
                        Оплатите курс удобным способом: банковской картой, через PayPal
                        или Яндекс.Кассу. Доступ к курсу открывается мгновенно.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
                    <div
                        class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4"
                    >
                        <UIcon
                            name="i-heroicons-video-camera"
                            class="text-purple-600 text-2xl"
                        />
                    </div>
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Начните обучение
                    </h3>
                    <p class="text-gray-600">
                        Получите пожизненный доступ ко всем материалам курса. Смотрите
                        видео в любое время, задавайте вопросы преподавателю.
                    </p>
                </div>
            </div>
            <div class="mt-8 p-6 bg-blue-50 rounded-xl border border-blue-200">
                <div class="flex items-start">
                    <UIcon
                        name="i-heroicons-information-circle"
                        class="text-blue-600 text-2xl mr-3 mt-1"
                    />
                    <div>
                        <h4 class="text-lg font-semibold text-blue-900 mb-2">
                            Важно знать
                        </h4>
                        <p class="text-blue-800">
                            После покупки курса вы получаете доступ ко всем видео-урокам,
                            дополнительным материалам (PDF-инструкции, референсы) и
                            возможность задавать вопросы преподавателю в закрытом чате.
                            Все курсы обновляются бесплатно для владельцев.
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Course Playlists -->
        <div class="mb-12">
            <div class="flex justify-between items-center mb-8">
                <h2 class="text-3xl font-bold text-gray-900">
                    Доступные курсы
                </h2>
                <div class="flex gap-2">
                    <UButton
                        color="neutral"
                        variant="outline"
                        size="sm"
                    >
                        Все уровни
                    </UButton>
                    <UButton
                        color="neutral"
                        variant="outline"
                        size="sm"
                    >
                        По популярности
                    </UButton>
                    <UButton
                        color="primary"
                        variant="solid"
                        size="sm"
                    >
                        <UIcon
                            name="i-heroicons-funnel"
                            class="mr-2"
                        />
                        Фильтры
                    </UButton>
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <div
                    v-for="playlist in playlists"
                    :key="playlist.id"
                    class="bg-white rounded-2xl border border-gray-200 overflow-hidden shadow-sm hover:shadow-lg transition-shadow duration-300"
                >
                    <div class="p-6">
                        <div class="flex justify-between items-start mb-4">
                            <div>
                                <UBadge
                                    :color="getLevelColor(playlist.level,)"
                                    class="mb-2"
                                >
                                    {{ getLevelText(playlist.level,) }}
                                </UBadge>
                                <h3 class="text-2xl font-bold text-gray-900 mb-2">
                                    {{ playlist.title }}
                                </h3>
                                <p class="text-gray-600 mb-4">
                                    {{ playlist.description }}
                                </p>
                            </div>
                            <div class="text-right">
                                <div class="flex items-center justify-end mb-2">
                                    <UIcon
                                        name="i-heroicons-star"
                                        class="text-yellow-500 mr-1"
                                    />
                                    <span class="font-semibold text-gray-900">{{
                                        playlist.rating
                                    }}</span>
                                </div>
                                <div class="text-sm text-gray-500">
                                    Преподаватель: {{ playlist.instructor }}
                                </div>
                            </div>
                        </div>

                        <div class="flex flex-wrap gap-4 mb-6">
                            <div class="flex items-center text-gray-600">
                                <UIcon
                                    name="i-heroicons-clock"
                                    class="mr-2"
                                />
                                <span>{{ playlist.duration }}</span>
                            </div>
                            <div class="flex items-center text-gray-600">
                                <UIcon
                                    name="i-heroicons-document-text"
                                    class="mr-2"
                                />
                                <span>{{ playlist.lessonsCount }} уроков</span>
                            </div>
                        </div>

                        <!-- Video Preview Grid (First 4 videos) -->
                        <div class="mb-6">
                            <h4 class="text-lg font-semibold text-gray-900 mb-3">
                                Содержание курса
                            </h4>
                            <div class="grid grid-cols-2 gap-3">
                                <div
                                    v-for="video in playlist.videos.slice(0, 4,)"
                                    :key="video.id"
                                    class="relative rounded-lg overflow-hidden border border-gray-200"
                                >
                                    <div
                                        class="aspect-video bg-linear-to-br from-gray-100 to-gray-200 flex items-center justify-center"
                                    >
                                        <div class="text-center">
                                            <div
                                                class="w-10 h-10 mx-auto mb-2 rounded-full bg-white/80 flex items-center justify-center"
                                            >
                                                <UIcon
                                                    v-if="video.isLocked"
                                                    name="i-heroicons-lock-closed"
                                                    class="text-gray-500"
                                                />
                                                <UIcon
                                                    v-else
                                                    name="i-heroicons-play"
                                                    class="text-blue-600"
                                                />
                                            </div>
                                            <div class="text-xs font-medium px-2">
                                                {{ video.title.substring(0, 20,) }}...
                                            </div>
                                        </div>
                                    </div>
                                    <div class="absolute top-2 right-2">
                                        <UBadge
                                            v-if="video.isLocked"
                                            color="neutral"
                                            size="xs"
                                        >
                                            <UIcon
                                                name="i-heroicons-lock-closed"
                                                class="mr-1"
                                            />
                                            Заблокировано
                                        </UBadge>
                                        <UBadge
                                            v-else
                                            color="success"
                                            size="xs"
                                        >
                                            <UIcon
                                                name="i-heroicons-play"
                                                class="mr-1"
                                            />
                                            Доступно
                                        </UBadge>
                                    </div>
                                </div>
                            </div>
                            <div class="mt-3 text-sm text-gray-500">
                                * Первые 2 видео доступны бесплатно. Для доступа к
                                остальным приобретите курс.
                            </div>
                        </div>

                        <div class="flex justify-between items-center">
                            <div>
                                <div class="text-2xl font-bold text-gray-900">
                                    {{ playlist.discountPrice?.toLocaleString('ru-RU',) }}
                                    ₽
                                    <span
                                        v-if="playlist.discountPrice"
                                        class="text-lg text-gray-500 line-through ml-2"
                                    >
                                        {{ playlist.price.toLocaleString('ru-RU',) }} ₽
                                    </span>
                                </div>
                                <div class="text-sm text-gray-500">
                                    Пожизненный доступ
                                </div>
                            </div>
                            <div class="flex gap-3">
                                <UButton
                                    color="neutral"
                                    variant="outline"
                                    @click="openPlaylist(playlist,)"
                                >
                                    <UIcon
                                        name="i-heroicons-eye"
                                        class="mr-2"
                                    />
                                    Подробнее
                                </UButton>
                                <UButton
                                    color="primary"
                                    variant="solid"
                                >
                                    <UIcon
                                        name="i-heroicons-shopping-cart"
                                        class="mr-2"
                                    />
                                    Купить курс
                                </UButton>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- FAQ Section -->
        <div
            class="bg-linear-to-br from-gray-50 to-white rounded-2xl p-8 border border-gray-200"
        >
            <h2 class="text-3xl font-bold text-gray-900 mb-8">
                Часто задаваемые вопросы
            </h2>
            <div class="space-y-4">
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Можно ли смотреть курсы на телефоне?
                    </h3>
                    <p class="text-gray-600">
                        Да, все наши курсы адаптированы для просмотра на любых
                        устройствах: компьютерах, планшетах и смартфонах. Вы можете
                        учиться в любое время и в любом месте.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Есть ли возможность вернуть деньги?
                    </h3>
                    <p class="text-gray-600">
                        Мы предоставляем 14-дневную гарантию возврата средств. Если курс
                        не подошел по какой-либо причине, вы можете запросить возврат в
                        течение двух недель после покупки.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Как долго длится доступ к курсу?
                    </h3>
                    <p class="text-gray-600">
                        Доступ к курсу предоставляется навсегда. Вы можете пересматривать
                        материалы в любое время, даже после завершения курса.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Можно ли задавать вопросы преподавателю?
                    </h3>
                    <p class="text-gray-600">
                        Да, у каждого курса есть закрытый чат с преподавателем, где вы
                        можете задавать вопросы и получать обратную связь по своим
                        работам.
                    </p>
                </div>
            </div>
        </div>

        <!-- CTA Section -->
        <div class="mt-12 text-center">
            <div class="max-w-2xl mx-auto">
                <h2 class="text-3xl font-bold text-gray-900 mb-4">
                    Готовы начать обучение?
                </h2>
                <p class="text-xl text-gray-600 mb-8">
                    Выберите курс, который соответствует вашим интересам и уровню
                    подготовки. Первые два видео каждого курса доступны бесплатно!
                </p>
                <UButton
                    color="primary"
                    variant="solid"
                    size="xl"
                    class="font-bold px-8 py-4"
                >
                    <UIcon
                        name="i-heroicons-rocket-launch"
                        class="mr-3 text-2xl"
                    />
                    Начать бесплатное обучение
                </UButton>
            </div>
        </div>
    </UContainer>
</template>

<style scoped>
.aspect-video {
    aspect-ratio: 16 / 9;
}

.service-card {
    transition: all 0.3s ease;
}

.service-card:hover {
    transform: translateY(-4px);
}
</style>
