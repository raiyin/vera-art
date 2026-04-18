<script setup>
const { locale, setLocale } = useI18n()

useHead({
    meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    link: [
        { rel: 'icon', href: '/favicon.ico', type: 'image/x-icon' },
        { rel: 'icon', href: '/favicon-16x16.png', type: 'image/png', sizes: '16x16' },
        { rel: 'icon', href: '/favicon-32x32.png', type: 'image/png', sizes: '32x32' },
        { rel: 'apple-touch-icon', href: '/favicon-128x128.png', sizes: '128x128' },
        { rel: 'manifest', href: '/site.webmanifest' }
    ],
    htmlAttrs: {
        lang: locale.value,
    },
});

const title = 'Страница художницы Перцуковой Веры';
const description =
    'Это личная страничка художницы и преподавателя Перцуковой Веры Олеговны. Здесь вы можете посмотреть на мои работы, а также купить или заказать картину для себяили в подарок.';

useSeoMeta({
    title,
    description,
    ogTitle: title,
    ogDescription: description,
    ogImage: 'https://ui.nuxt.com/assets/templates/nuxt/starter-light.png',
    twitterCard: 'summary_large_image',
});

const switchLocale = (newLocale) => {
    setLocale(newLocale)
}
</script>

<template>
    <UApp>
        <UHeader class="flex items-center">
            <template #left>
                <div class="flex items-center space-x-6 h-full">
                    <NuxtLink to="/" class="uppercase  py-2">{{ $t('header.main') }}</NuxtLink>
                    <NuxtLink to="/all-works" class="uppercase py-2">{{ $t('header.all_works') }}</NuxtLink>
                    <NuxtLink to="/news" class="uppercase py-2">{{ $t('header.news') }}</NuxtLink>
                    <NuxtLink to="/shop" class="uppercase  py-2">{{ $t('header.shop') }}</NuxtLink>
                    <NuxtLink to="/services" class="uppercase py-2">{{ $t('header.services') }}</NuxtLink>
                    <NuxtLink to="/pay-delivery" class="uppercase py-2">{{ $t('header.payment') }}</NuxtLink>
                </div>
            </template>
            <template #right>
                <div class="flex items-center space-x-4 h-full">
                    <!-- Theme switcher placeholder -->
                    <button class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 flex items-center justify-center">
                        <Icon name="i-simple-icons-sun" class="w-5 h-5" />
                    </button>
                    <!-- Locale switcher -->
                    <button
                        class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 flex items-center justify-center cursor-pointer"
                        @click="switchLocale(locale === 'en' ? 'ru' : 'en')"
                    >
                        {{ locale === 'en' ? 'RU' : 'EN' }}
                    </button>
                    <!-- Social media icons -->
                    <a href="https://t.me/MilayaV" target="_blank" class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 flex items-center justify-center">
                        <Icon name="i-simple-icons-telegram" class="w-5 h-5" />
                    </a>
                    <a href="https://vk.com/perczukowa" target="_blank" class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 flex items-center justify-center">
                        <Icon name="i-simple-icons-vk" class="w-5 h-5" />
                    </a>
                    <a href="mailto:perczukowa@yandex.ru" class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 flex items-center justify-center">
                        <Icon name="i-simple-icons-gmail" class="w-5 h-5" />
                    </a>
                </div>
            </template>
        </UHeader>

        <UMain>
            <NuxtPage />
        </UMain>

        <AppFooter />
    </UApp>
</template>
