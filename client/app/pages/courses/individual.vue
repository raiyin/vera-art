<script setup lang="ts">
import { ref } from 'vue';

interface LessonFormat {
    id: number;
    title: string;
    description: string;
    duration: string;
    price: number;
    features: string[];
    icon: string;
    color: string;
}

interface Instructor {
    id: number;
    name: string;
    specialization: string;
    experience: string;
    avatar: string;
    rating: number;
}

const lessonFormats = ref<LessonFormat[]>([
    {
        id: 1,
        title: 'Разовое занятие',
        description:
            'Идеально для тех, кто хочет попробовать или решить конкретную задачу',
        duration: '2 часа',
        price: 3000,
        features: [
            'Индивидуальный подход',
            'Анализ текущего уровня',
            'Практическое задание',
            'Обратная связь по работе',
            'Рекомендации по материалам',
        ],
        icon: 'i-heroicons-calendar',
        color: 'blue',
    },
    {
        id: 2,
        title: 'Пакет из 4 занятий',
        description: 'Базовый курс для освоения новой техники или темы',
        duration: '2 часа каждое',
        price: 10000,
        features: [
            'Программа из 4 последовательных уроков',
            'Домашние задания с проверкой',
            'Постоянная поддержка в чате',
            'Анализ прогресса',
            'Сертификат о прохождении',
        ],
        icon: 'i-heroicons-academic-cap',
        color: 'green',
    },
    {
        id: 3,
        title: 'Интенсивный курс',
        description: 'Глубокое погружение в тему с быстрым результатом',
        duration: '8 занятий по 3 часа',
        price: 22000,
        features: [
            'Индивидуальная программа',
            'Ежедневная практика',
            'Работа над персональным проектом',
            'Видеозапись всех занятий',
            'Портфолио по окончании',
        ],
        icon: 'i-heroicons-rocket-launch',
        color: 'purple',
    },
]);

const instructors = ref<Instructor[]>([
    {
        id: 1,
        name: 'Анна Петрова',
        specialization: 'Акварель, ботаническая иллюстрация',
        experience: '12 лет преподавания',
        avatar: '/images/instructors/anna.jpg',
        rating: 4.9,
    },
    {
        id: 2,
        name: 'Иван Смирнов',
        specialization: 'Масляная живопись, портрет',
        experience: '15 лет в академической живописи',
        avatar: '/images/instructors/ivan.jpg',
        rating: 4.8,
    },
    {
        id: 3,
        name: 'Мария Козлова',
        specialization: 'Цифровая иллюстрация, скетчинг',
        experience: '8 лет в индустрии дизайна',
        avatar: '/images/instructors/maria.jpg',
        rating: 4.7,
    },
    {
        id: 4,
        name: 'Сергей Волков',
        specialization: 'Академический рисунок, анатомия',
        experience: '20 лет преподавания в художественной школе',
        avatar: '/images/instructors/sergey.jpg',
        rating: 5.0,
    },
]);

const processSteps = ref([
    {
        number: '01',
        title: 'Консультация',
        description:
            'Мы обсуждаем ваши цели, текущий уровень и предпочтения. Определяем направление и формат занятий.',
    },
    {
        number: '02',
        title: 'Составление программы',
        description:
            'Преподаватель разрабатывает индивидуальную программу, учитывающую ваши задачи и сроки.',
    },
    {
        number: '03',
        title: 'Первое занятие',
        description:
            'Знакомство с материалами, базовые упражнения, определение точек роста и дальнейшего плана.',
    },
    {
        number: '04',
        title: 'Регулярные занятия',
        description:
            'Систематическое обучение с домашними заданиями, обратной связью и корректировкой программы.',
    },
    {
        number: '05',
        title: 'Итоговый проект',
        description:
            'Создание завершенной работы, демонстрирующей приобретенные навыки и художественный рост.',
    },
    {
        number: '06',
        title: 'Дальнейшее развитие',
        description:
            'Рекомендации по самостоятельной практике, участие в выставках или продолжение обучения.',
    },
]);

const selectedFormat = ref<LessonFormat | null>(null);
const selectedInstructor = ref<Instructor | null>(null);

const selectFormat = (format: LessonFormat) => {
    selectedFormat.value = format;
};

const selectInstructor = (instructor: Instructor) => {
    selectedInstructor.value = instructor;
};
</script>

<template>
    <UContainer class="py-8 md:py-12">
        <!-- Breadcrumb -->
        <UBreadcrumb class="mb-8">
            <UBreadcrumbItem to="/">Главная</UBreadcrumbItem>
            <UBreadcrumbItem to="/services">Услуги</UBreadcrumbItem>
            <UBreadcrumbItem>Индивидуальные занятия</UBreadcrumbItem>
        </UBreadcrumb>

        <!-- Hero Section -->
        <div
            class="relative overflow-hidden rounded-2xl bg-linear-to-r from-emerald-600 to-teal-600 p-8 md:p-12 mb-12"
        >
            <div class="relative z-10 max-w-3xl">
                <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">
                    Индивидуальные занятия по живописи и рисунку
                </h1>
                <p class="text-xl text-emerald-100 mb-6">
                    Персональный подход, гибкий график и программа, созданная специально
                    для вас. Научитесь рисовать с нуля или совершенствуйте свои навыки под
                    руководством опытного художника.
                </p>
                <div class="flex flex-wrap gap-4">
                    <UButton
                        variant="solid"
                        size="lg"
                        class="font-semibold bg-white text-emerald-600 hover:bg-gray-100"
                    >
                        <UIcon name="i-heroicons-calendar" class="mr-2" />
                        Записаться на пробный урок
                    </UButton>
                    <UButton
                        variant="outline"
                        size="lg"
                        class="font-semibold text-white border-white hover:bg-white/10"
                    >
                        <UIcon name="i-heroicons-chat-bubble-left-right" class="mr-2" />
                        Бесплатная консультация
                    </UButton>
                </div>
            </div>
            <div class="absolute right-0 top-0 bottom-0 w-1/3 opacity-20">
                <div class="h-full bg-linear-to-l from-white to-transparent"></div>
            </div>
        </div>

        <!-- How Lessons Are Organized -->
        <div class="mb-16">
            <div class="text-center mb-12">
                <h2 class="text-3xl font-bold text-gray-900 mb-4">
                    Как организованы индивидуальные занятия
                </h2>
                <p class="text-xl text-gray-600 max-w-3xl mx-auto">
                    Каждое занятие строится вокруг ваших целей и возможностей. Мы создаем
                    комфортную среду для обучения, где вы можете задавать вопросы,
                    экспериментировать и развиваться в своем темпе.
                </p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                <div
                    v-for="step in processSteps"
                    :key="step.number"
                    class="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm hover:shadow-lg transition-shadow duration-300"
                >
                    <div class="flex items-start mb-4">
                        <div
                            class="w-12 h-12 rounded-full bg-emerald-100 flex items-center justify-center mr-4"
                        >
                            <span class="text-2xl font-bold text-emerald-700">{{
                                step.number
                            }}</span>
                        </div>
                        <h3 class="text-xl font-bold text-gray-900 mt-2">
                            {{ step.title }}
                        </h3>
                    </div>
                    <p class="text-gray-600">{{ step.description }}</p>
                </div>
            </div>
        </div>

        <!-- Lesson Formats -->
        <div class="mb-16">
            <div class="text-center mb-12">
                <h2 class="text-3xl font-bold text-gray-900 mb-4">Форматы занятий</h2>
                <p class="text-xl text-gray-600 max-w-3xl mx-auto">
                    Выберите подходящий формат в зависимости от ваших целей, доступного
                    времени и бюджета.
                </p>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div
                    v-for="format in lessonFormats"
                    :key="format.id"
                    class="bg-white rounded-2xl border-2 border-gray-200 overflow-hidden hover:border-emerald-500 transition-all duration-300 hover:shadow-xl"
                >
                    <div class="p-8">
                        <div class="flex items-center mb-6">
                            <div
                                class="w-12 h-12 rounded-lg bg-blue-100 flex items-center justify-center mr-4"
                                v-if="format.color === 'blue'"
                            >
                                <UIcon
                                    :name="format.icon"
                                    class="text-blue-600 text-2xl"
                                />
                            </div>
                            <div
                                class="w-12 h-12 rounded-lg bg-green-100 flex items-center justify-center mr-4"
                                v-else-if="format.color === 'green'"
                            >
                                <UIcon
                                    :name="format.icon"
                                    class="text-green-600 text-2xl"
                                />
                            </div>
                            <div
                                class="w-12 h-12 rounded-lg bg-purple-100 flex items-center justify-center mr-4"
                                v-else
                            >
                                <UIcon
                                    :name="format.icon"
                                    class="text-purple-600 text-2xl"
                                />
                            </div>
                            <div>
                                <h3 class="text-2xl font-bold text-gray-900">
                                    {{ format.title }}
                                </h3>
                                <p class="text-gray-600">{{ format.duration }}</p>
                            </div>
                        </div>

                        <p class="text-gray-700 mb-6">{{ format.description }}</p>

                        <div class="mb-6">
                            <h4 class="font-semibold text-gray-900 mb-3">
                                Что включено:
                            </h4>
                            <ul class="space-y-2">
                                <li
                                    v-for="feature in format.features"
                                    :key="feature"
                                    class="flex items-start"
                                >
                                    <UIcon
                                        name="i-heroicons-check-circle"
                                        class="text-emerald-500 mr-2 mt-1"
                                    />
                                    <span class="text-gray-600">{{ feature }}</span>
                                </li>
                            </ul>
                        </div>

                        <div class="mb-6">
                            <div class="text-3xl font-bold text-gray-900">
                                {{ format.price.toLocaleString('ru-RU') }} ₽
                            </div>
                            <div class="text-gray-500">
                                за {{ format.id === 1 ? 'занятие' : 'пакет' }}
                            </div>
                        </div>

                        <UButton
                            color="primary"
                            variant="solid"
                            class="w-full"
                            @click="selectFormat(format)"
                        >
                            Выбрать этот формат
                        </UButton>
                    </div>
                </div>
            </div>
        </div>

        <!-- Instructors -->
        <div class="mb-16">
            <div class="text-center mb-12">
                <h2 class="text-3xl font-bold text-gray-900 mb-4">Наши преподаватели</h2>
                <p class="text-xl text-gray-600 max-w-3xl mx-auto">
                    Опытные художники-педагоги с академическим образованием и многолетним
                    стажем преподавания.
                </p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <div
                    v-for="instructor in instructors"
                    :key="instructor.id"
                    class="bg-white rounded-2xl border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow duration-300"
                >
                    <div class="p-6">
                        <div class="flex items-center mb-4">
                            <div
                                class="w-16 h-16 rounded-full bg-linear-to-br from-gray-200 to-gray-300 flex items-center justify-center mr-4"
                            >
                                <UIcon
                                    name="i-heroicons-user"
                                    class="text-gray-500 text-2xl"
                                />
                            </div>
                            <div>
                                <h3 class="text-xl font-bold text-gray-900">
                                    {{ instructor.name }}
                                </h3>
                                <div class="flex items-center">
                                    <UIcon
                                        name="i-heroicons-star"
                                        class="text-yellow-500 mr-1"
                                    />
                                    <span class="font-semibold text-gray-900">{{
                                        instructor.rating
                                    }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="mb-4">
                            <div class="text-sm font-semibold text-gray-700 mb-1">
                                Специализация:
                            </div>
                            <p class="text-gray-600">{{ instructor.specialization }}</p>
                        </div>

                        <div class="mb-6">
                            <div class="text-sm font-semibold text-gray-700 mb-1">
                                Опыт:
                            </div>
                            <p class="text-gray-600">{{ instructor.experience }}</p>
                        </div>

                        <UButton
                            color="neutral"
                            variant="outline"
                            class="w-full"
                            @click="selectInstructor(instructor)"
                        >
                            <UIcon name="i-heroicons-eye" class="mr-2" />
                            Посмотреть работы
                        </UButton>
                    </div>
                </div>
            </div>
        </div>

        <!-- What You Will Learn -->
        <div
            class="bg-linear-to-br from-gray-50 to-white rounded-2xl p-8 mb-16 border border-gray-200"
        >
            <h2 class="text-3xl font-bold text-gray-900 mb-8 text-center">
                Чему вы научитесь
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div class="space-y-6">
                    <div class="flex items-start">
                        <div
                            class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center mr-4 shrink-0"
                        >
                            <UIcon name="i-heroicons-paint-brush" class="text-blue-600" />
                        </div>
                        <div>
                            <h3 class="text-xl font-semibold text-gray-900 mb-2">
                                Технические навыки
                            </h3>
                            <p class="text-gray-600">
                                Освоите различные техники живописи и рисунка, работу с
                                материалами, композицию, цветоведение и перспективу.
                            </p>
                        </div>
                    </div>
                    <div class="flex items-start">
                        <div
                            class="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center mr-4 shrink-0"
                        >
                            <UIcon name="i-heroicons-light-bulb" class="text-green-600" />
                        </div>
                        <div>
                            <h3 class="text-xl font-semibold text-gray-900 mb-2">
                                Творческое мышление
                            </h3>
                            <p class="text-gray-600">
                                Разовьете художественное видение, научитесь генерировать
                                идеи и создавать оригинальные работы.
                            </p>
                        </div>
                    </div>
                    <div class="flex items-start">
                        <div
                            class="w-10 h-10 rounded-full bg-purple-100 flex items-center justify-center mr-4 shrink-0"
                        >
                            <UIcon
                                name="i-heroicons-document-text"
                                class="text-purple-600"
                            />
                        </div>
                        <div>
                            <h3 class="text-xl font-semibold text-gray-900 mb-2">
                                Проектный подход
                            </h3>
                            <p class="text-gray-600">
                                Научитесь планировать и реализовывать художественные
                                проекты от эскиза до завершенной работы.
                            </p>
                        </div>
                    </div>
                </div>
                <div class="space-y-6">
                    <div class="flex items-start">
                        <div
                            class="w-10 h-10 rounded-full bg-amber-100 flex items-center justify-center mr-4 shrink-0"
                        >
                            <UIcon name="i-heroicons-chart-bar" class="text-amber-600" />
                        </div>
                        <div>
                            <h3 class="text-xl font-semibold text-gray-900 mb-2">
                                Анализ и критика
                            </h3>
                            <p class="text-gray-600">
                                Научитесь анализировать свои и чужие работы, получать и
                                давать конструктивную обратную связь.
                            </p>
                        </div>
                    </div>
                    <div class="flex items-start">
                        <div
                            class="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center mr-4 shrink-0"
                        >
                            <UIcon name="i-heroicons-briefcase" class="text-red-600" />
                        </div>
                        <div>
                            <h3 class="text-xl font-semibold text-gray-900 mb-2">
                                Профессиональные практики
                            </h3>
                            <p class="text-gray-600">
                                Узнаете о профессиональных стандартах в искусстве,
                                подготовке работ к выставкам и продажам.
                            </p>
                        </div>
                    </div>
                    <div class="flex items-start">
                        <div
                            class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center mr-4 shrink-0"
                        >
                            <UIcon name="i-heroicons-heart" class="text-indigo-600" />
                        </div>
                        <div>
                            <h3 class="text-xl font-semibold text-gray-900 mb-2">
                                Уверенность в себе
                            </h3>
                            <p class="text-gray-600">
                                Приобретете уверенность в своих творческих способностях и
                                научитесь преодолевать творческие кризисы.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- FAQ -->
        <div class="mb-16">
            <h2 class="text-3xl font-bold text-gray-900 mb-8 text-center">
                Частые вопросы об индивидуальных занятиях
            </h2>
            <div class="space-y-4 max-w-3xl mx-auto">
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Нужно ли мне иметь художественное образование или опыт?
                    </h3>
                    <p class="text-gray-600">
                        Нет, не нужно. Мы работаем с учениками любого уровня — от полных
                        новичков до опытных художников. Программа строится исходя из
                        вашего текущего уровня и целей.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Какие материалы нужны для занятий?
                    </h3>
                    <p class="text-gray-600">
                        На первом занятии преподаватель даст список необходимых материалов
                        в зависимости от выбранного направления. Для начала можно
                        использовать базовые наборы, которые есть в большинстве
                        художественных магазинов.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Можно ли заниматься онлайн?
                    </h3>
                    <p class="text-gray-600">
                        Да, мы проводим занятия как очно (в нашей студии или у вас), так и
                        онлайн через Zoom или Skype. Онлайн-занятия включают демонстрацию
                        экрана, обратную связь в реальном времени и запись урока.
                    </p>
                </div>
                <div class="bg-white p-6 rounded-xl border border-gray-200">
                    <h3 class="text-xl font-semibold text-gray-900 mb-2">
                        Как часто нужно заниматься?
                    </h3>
                    <p class="text-gray-600">
                        Рекомендуемая частота — 1-2 раза в неделю. Это позволяет сохранять
                        прогресс и давать время на выполнение домашних заданий. График
                        составляется индивидуально с учетом вашей занятости.
                    </p>
                </div>
            </div>
        </div>

        <!-- CTA -->
        <div class="text-center">
            <div class="max-w-2xl mx-auto">
                <h2 class="text-3xl font-bold text-gray-900 mb-4">Готовы начать?</h2>
                <p class="text-xl text-gray-600 mb-8">
                    Запишитесь на бесплатную консультацию, чтобы обсудить ваши цели и
                    составить индивидуальную программу обучения.
                </p>
                <div class="flex flex-wrap justify-center gap-4">
                    <UButton
                        color="primary"
                        variant="solid"
                        size="xl"
                        class="font-bold px-8 py-4"
                    >
                        <UIcon name="i-heroicons-calendar" class="mr-3 text-2xl" />
                        Записаться на консультацию
                    </UButton>
                    <UButton
                        color="neutral"
                        variant="outline"
                        size="xl"
                        class="font-bold px-8 py-4"
                    >
                        <UIcon name="i-heroicons-phone" class="mr-3 text-2xl" />
                        Позвонить нам
                    </UButton>
                </div>
                <p class="mt-6 text-gray-500">
                    Или напишите нам в WhatsApp / Telegram для быстрого ответа
                </p>
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
