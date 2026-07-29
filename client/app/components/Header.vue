<script setup lang="ts">
import { useI18n, } from '#imports';
import { useAuthStore, } from '~/stores/AuthStore';
import { useNotificationStore, } from '~/stores/NotificationStore';
import { useThemeStore, } from '~/stores/ThemeStore';
import type { NavigationMenuItem, DropdownMenuItem, } from '@nuxt/ui';
import authApi from '~/api/auth';

const { locale, setLocale, t, } = useI18n();
const authStore = useAuthStore();
const notificationStore = useNotificationStore();
const themeStore = useThemeStore();

const route = useRoute();

const avatarUrl = ref('',);

// Fetch avatar URL when authenticated
const fetchAvatar = async () => {
    if (!authStore.isAuthenticated) return;
    try {
        const data = await authApi.getProfile();
        avatarUrl.value = data.avatar_url || '';
    } catch {
        // Silently fail - avatar is optional
    }
};

// Watch auth state to fetch avatar when user logs in
watch(
    () => authStore.isAuthenticated,
    (isAuth,) => {
        if (isAuth) {
            fetchAvatar();
        } else {
            avatarUrl.value = '';
        }
    },
    { immediate: true, },
    );

const switchLocale = (newLocale: 'ru' | 'en',) => {
    setLocale(newLocale,);
};

const toggleTheme = () => {
    themeStore.theme = themeStore.theme === 'light' ? 'dark' : 'light';
};

const logout = () => {
    authStore.clearTokens();
    notificationStore.stopPolling();
    navigateTo('/',);
};

// User dropdown menu items (profile, logout)
const userMenuItems = computed<DropdownMenuItem[]>(() => [
    {
        label: t('profile.menu',),
        icon: 'i-heroicons-user-circle',
        to: '/profile',
        avatar: avatarUrl.value ? { src: avatarUrl.value, } : undefined,
    },
    {
        label: t('auth.logout',),
        icon: 'i-heroicons-arrow-right-on-rectangle',
        onSelect: () => logout(),
    },
]);

// Build navigation items with translations
const navigation = computed<NavigationMenuItem[]>(() => {
    const items: NavigationMenuItem[] = [
        {
            label: t('header.main',),
            to: '/',
        },
        {
            label: t('header.all_works',),
            to: '/gallery',
            active: route.path.startsWith('/gallery',),
        },
        {
            label: t('header.news',),
            to: '/news',
            active: route.path.startsWith('/news',),
        },
        {
            label: t('header.shop',),
            to: '/art-store',
            active: route.path.startsWith('/art-store',),
        },
        {
            label: t('header.payment',),
            to: '/pay-delivery',
            icon: 'i-heroicons-credit-card',
            active: route.path.startsWith('/pay-delivery',),
            value: 'payment',
        },
        {
            label: t('header.services',),
            to: '/services',
            active: route.path.startsWith('/services',),
            icon: 'i-heroicons-document-text',
            value: 'services',
            type: 'trigger',
            children: [
                {
                    label: 'Картины',
                    to: '/art-store',
                    icon: 'i-heroicons-paint-brush',
                    active: route.path.startsWith('/art-store',),
                    value: 'art-store',
                },
                {
                    label: 'Мастер-классы',
                    to: '/master-classes',
                    icon: 'i-heroicons-video-camera',
                    active: route.path.startsWith('/master-classes',),
                    value: 'master-classes',
                },
                {
                    label: 'Онлайн-курсы',
                    to: '/courses/online',
                    icon: 'i-heroicons-computer-desktop',
                    active: route.path.startsWith('/courses/online',),
                    value: 'online-courses',
                },
                {
                    label: 'Индивидуальные занятия',
                    to: '/courses/individual',
                    icon: 'i-heroicons-user',
                    active: route.path.startsWith('/courses/individual',),
                    value: 'individual-courses',
                },
            ],
        },
    ];

    if (authStore.isAuthenticated && authStore.isAdmin) {
        items.push({
            label: t('header.admin',),
            to: '/admin',
            active: route.path.startsWith('/admin',),
        });
    }

    return items;
});
</script>

<template>
    <UHeader>
        <template #left>
            <NuxtLink
                to="/"
                class="flex items-center"
                aria-label="Vera site"
            >
                <img
                    src="../assets/icons/favicon-art.svg"
                    alt="Palette"
                    class="h-10 w-10 transition-transform duration-200 hover:scale-110"
                >
            </NuxtLink>
        </template>

        <template #default>
            <UNavigationMenu :items="navigation" />
        </template>

        <template #right>
            <!-- Language switcher -->
            <UTooltip
                :text="locale === 'en' ? 'Switch to Russian' : 'Переключится на Русский'"
            >
                <UButton
                    class="text-grey"
                    variant="ghost"
                    square
                    @click="switchLocale(locale === 'en' ? 'ru' : 'en',)"
                >
                    {{ locale === 'en' ? 'RU' : 'EN' }}
                </UButton>
            </UTooltip>

            <!-- Theme toggle -->
            <UTooltip :text="$t('theme.title',)">
                <UButton
                    class="text-grey"
                    variant="ghost"
                    square
                    @click="toggleTheme"
                >
                    <ClientOnly>
                        <Icon
                            :name="
                                themeStore.theme === 'light'
                                    ? 'i-heroicons-moon'
                                    : 'i-heroicons-sun'
                            "
                            class="w-5 h-5"
                        />
                        <template #fallback>
                            <Icon
                                name="i-heroicons-moon"
                                class="w-5 h-5"
                            />
                        </template>
                    </ClientOnly>
                </UButton>
            </UTooltip>

            <!-- Notifications -->
            <div
                v-if="authStore.isAuthenticated"
                class="relative"
            >
                <UTooltip text="Сообщения">
                    <UButton
                        class="text-grey"
                        variant="ghost"
                        square
                        to="/chat"
                    >
                        <Icon
                            name="i-heroicons-chat-bubble-left-right"
                            class="w-5 h-5"
                        />
                        <span
                            v-if="notificationStore.hasUnread"
                            class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
                        >
                            {{
                                notificationStore.unreadCount > 9
                                    ? '9+'
                                    : notificationStore.unreadCount
                            }}
                        </span>
                    </UButton>
                </UTooltip>
            </div>

            <!-- Auth: User dropdown when authenticated -->
            <UDropdownMenu
                v-if="authStore.isAuthenticated"
                :items="userMenuItems"
            >
                <UButton
                    class="text-grey"
                    variant="ghost"
                    square
                >
                    <img
                        v-if="avatarUrl"
                        :src="avatarUrl"
                        alt="Avatar"
                        class="w-5 h-5 rounded-full object-cover"
                    >
                    <Icon
                        v-else
                        name="i-heroicons-user-circle"
                        class="w-5 h-5"
                    />
                </UButton>
            </UDropdownMenu>

            <UTooltip
                v-else
                :text="$t('auth.login',)"
            >
                <UButton
                    class="text-grey"
                    variant="ghost"
                    square
                    to="/auth/login"
                >
                    <Icon
                        name="i-heroicons-user-circle"
                        class="w-5 h-5"
                    />
                </UButton>
            </UTooltip>

            <!-- Social links (desktop only) -->
            <div class="md:flex items-center space-x-1">
                <UTooltip text="Telegram">
                    <UButton
                        variant="ghost"
                        square
                        href="https://t.me/MilayaV"
                        target="_blank"
                        class="text-grey"
                        external
                    >
                        <Icon
                            name="i-simple-icons-telegram"
                            class="w-5 h-5"
                        />
                    </UButton>
                </UTooltip>

                <UTooltip text="VK">
                    <UButton
                        class="text-grey"
                        variant="ghost"
                        square
                        href="https://vk.com/perczukowa"
                        target="_blank"
                        external
                    >
                        <Icon
                            name="i-simple-icons-vk"
                            class="w-5 h-5"
                        />
                    </UButton>
                </UTooltip>
                <UTooltip text="Email">
                    <UButton
                        class="text-grey"
                        variant="ghost"
                        square
                        href="mailto:perczukowa@yandex.ru"
                        external
                    >
                        <Icon
                            name="i-simple-icons-gmail"
                            class="w-5 h-5"
                        />
                    </UButton>
                </UTooltip>
            </div>
        </template>

        <template #body>
            <UNavigationMenu
                :items="navigation"
                orientation="vertical"
                class="-mx-2.5"
            />
        </template>
    </UHeader>
</template>
