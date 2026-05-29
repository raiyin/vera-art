<script setup>
import { onMounted, onUnmounted, watch } from 'vue';
import { useI18n, useColorMode } from '#imports';
import { useMaterialStore } from './stores/MaterialStore';
import { useThemeStore } from './stores/ThemeStore';
import { useAuthStore } from './stores/AuthStore';
import { useNotificationStore } from './stores/NotificationStore';
import CookieConsent from './components/CookieConsent.vue';
import Header from './components/Header.vue';
const { locale, setLocale } = useI18n();
const materialStore = useMaterialStore();
const themeStore = useThemeStore();
const authStore = useAuthStore();
const notificationStore = useNotificationStore();
const colorMode = useColorMode();

// Sync theme store with Nuxt UI color mode (store -> colorMode)
watch(
    () => themeStore.theme,
    (newTheme) => {
        // Update both preference and value to ensure persistence
        colorMode.preference = newTheme;
        colorMode.value = newTheme;
    },
    { immediate: true }
);

// Sync color mode changes to theme store (colorMode -> store)
watch(
    () => colorMode.value,
    (newColorMode) => {
        if (newColorMode === 'light' || newColorMode === 'dark') {
            if (themeStore.theme !== newColorMode) {
                themeStore.theme = newColorMode;
            }
        }
    }
);

useHead({
    meta: [{ name: 'viewport', content: 'width=device-width, initial-scale=1' }],
    link: [
        { rel: 'icon', href: '/favicon.ico', type: 'image/x-icon' },
        { rel: 'icon', href: '/favicon-16x16.png', type: 'image/png', sizes: '16x16' },
        { rel: 'icon', href: '/favicon-32x32.png', type: 'image/png', sizes: '32x32' },
        { rel: 'apple-touch-icon', href: '/favicon-128x128.png', sizes: '128x128' },
        { rel: 'manifest', href: '/site.webmanifest' },
    ],
    htmlAttrs: {
        lang: locale.value,
    },
});

onMounted(() => {
    // Initialize auth state from localStorage (tokens only available client-side)
    authStore.initFromLocalStorage();

    // Load reference data (materials, bases) on app startup
    materialStore.fetchAll();

    // Start notification polling if user is authenticated
    if (authStore.isAuthenticated) {
        notificationStore.startPolling();
    }
});

onUnmounted(() => {
    // Stop notification polling
    notificationStore.stopPolling();
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
</script>

<template>
    <UApp>
        <Header />

        <UMain>
            <NuxtPage />
        </UMain>

        <AppFooter />

        <!-- Cookie Consent Banner -->
        <CookieConsent />
    </UApp>
</template>
