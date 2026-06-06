<template>
    <aside
        class="admin-sidebar"
        :class="{
            'admin-sidebar--open': isOpen,
            'admin-sidebar--dark': isDark,
        }"
    >
        <div class="admin-sidebar__brand">
            <NuxtLink to="/admin" class="admin-sidebar__logo">
                <svg
                    class="admin-sidebar__logo-icon"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                >
                    <path
                        d="M3 12h2v2H3v-2zm5 0h2v2H8v-2zm5 0h2v2h-2v-2zm5 0h2v2h-2v-2zM3 8h2v2H3V8zm5 0h2v2H8V8zm5 0h2v2h-2V8zm5 0h2v2h-2V8zM3 16h2v2H3v-2zm5 0h2v2H8v-2zm5 0h2v2h-2v-2zm5 0h2v2h-2v-2z"
                    />
                </svg>
                <span class="admin-sidebar__brand-text">Admin Panel</span>
            </NuxtLink>
        </div>

        <nav class="admin-sidebar__nav">
            <NuxtLink
                v-for="item in navItems"
                :key="item.to"
                :to="item.to"
                class="admin-sidebar__link"
                :class="{ 'admin-sidebar__link--active': isActive(item.to) }"
            >
                <span class="admin-sidebar__link-icon" v-html="item.icon" />
                <span class="admin-sidebar__link-text">{{ item.label }}</span>
                <span
                    v-if="item.badge"
                    class="admin-sidebar__badge"
                    :class="`admin-sidebar__badge--${item.badgeType || 'info'}`"
                >
                    {{ item.badge }}
                </span>
            </NuxtLink>
        </nav>

        <div class="admin-sidebar__footer">
            <NuxtLink to="/" class="admin-sidebar__link admin-sidebar__link--back">
                <span class="admin-sidebar__link-icon">
                    <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        width="20"
                        height="20"
                    >
                        <path d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                    </svg>
                </span>
                <span class="admin-sidebar__link-text">{{ $t('header.main') }}</span>
            </NuxtLink>
        </div>
    </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const props = defineProps<{
    isOpen: boolean;
    isDark: boolean;
}>();

const emit = defineEmits<{
    close: [];
}>();

const route = useRoute();

interface NavItem {
    to: string;
    label: string;
    icon: string;
    badge: string | null;
    badgeType?: 'info' | 'warning' | 'danger';
}

const navItems = computed<NavItem[]>(() => [
    {
        to: '/admin',
        label: 'Dashboard',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>',
        badge: null,
    },
    {
        to: '/admin/gallery',
        label: 'Галерея',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="M21 15l-5-5L5 21"/></svg>',
        badge: null,
    },
    {
        to: '/admin/shop',
        label: 'Магазин',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4zM3 6h18"/><path d="M16 10a4 4 0 01-8 0"/></svg>',
        badge: null,
    },
    {
        to: '/admin/news',
        label: 'Новости',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M4 22h16a2 2 0 002-2V4a2 2 0 00-2-2H8a2 2 0 00-2 2v16a2 2 0 01-4 0V7m0 15V7"/></svg>',
        badge: null,
    },
    {
        to: '/admin/courses',
        label: 'Курсы',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M12 6v6l4 2"/><path d="M12 2a10 10 0 1010 10"/></svg>',
        badge: null,
    },
    {
        to: '/admin/master-classes',
        label: 'Мастер-классы',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M14.5 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>',
        badge: null,
    },
    {
        to: '/admin/users',
        label: 'Пользователи',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/></svg>',
        badge: null,
    },
    {
        to: '/admin/reviews',
        label: 'Отзывы',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>',
        badge: null,
    },
    {
        to: '/admin/purchases',
        label: 'Покупки',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4zM3 6h18"/><path d="M16 10a4 4 0 01-8 0"/></svg>',
        badge: null,
    },
    {
        to: '/admin/payments',
        label: 'Платежи',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><rect x="1" y="4" width="22" height="16" rx="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>',
        badge: null,
    },
    {
        to: '/admin/promo-codes',
        label: 'Промокоды',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M20.59 13.41l-7.17 7.17a2 2 0 01-2.83 0L2 12V2h10l8.59 8.59a2 2 0 010 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>',
        badge: null,
    },
    {
        to: '/admin/chats',
        label: 'Чаты',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"/></svg>',
        badge: null,
    },
    {
        to: '/admin/categories',
        label: 'Категории',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M4 4h6v6H4zM14 4h6v6h-6zM4 14h6v6H4zM14 14h6v6h-6z"/></svg>',
        badge: null,
    },
    {
        to: '/admin/tags',
        label: 'Теги',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><path d="M20.59 13.41l-7.17 7.17a2 2 0 01-2.83 0L2 12V2h10l8.59 8.59a2 2 0 010 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>',
        badge: null,
    },
    {
        to: '/admin/settings',
        label: 'Настройки',
        icon:
            '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06A1.65 1.65 0 0019.32 9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>',
        badge: null,
    },
]);

function isActive(path: string): boolean {
    if (path === '/admin') {
        return route.path === '/admin';
    }
    return route.path.startsWith(path);
}
</script>

<style scoped>
.admin-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    width: 260px;
    height: 100vh;
    background: var(--admin-sidebar-bg, #1a1d29);
    color: var(--admin-sidebar-text, #a0a5b5);
    display: flex;
    flex-direction: column;
    z-index: 50;
    transition: transform 0.3s ease;
    overflow-y: auto;
}

@media (max-width: 768px) {
    .admin-sidebar {
        transform: translateX(-100%);
    }

    .admin-sidebar--open {
        transform: translateX(0);
    }
}

.admin-sidebar__brand {
    padding: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.admin-sidebar__logo {
    display: flex;
    align-items: center;
    gap: 12px;
    text-decoration: none;
    color: #ffffff;
}

.admin-sidebar__logo-icon {
    width: 28px;
    height: 28px;
    color: var(--admin-primary, #6c5ce7);
}

.admin-sidebar__brand-text {
    font-size: 18px;
    font-weight: 700;
    letter-spacing: -0.3px;
}

.admin-sidebar__nav {
    flex: 1;
    padding: 12px 0;
    overflow-y: auto;
}

.admin-sidebar__link {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 20px;
    color: var(--admin-sidebar-text, #a0a5b5);
    text-decoration: none;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.2s ease;
    border-left: 3px solid transparent;
    position: relative;
}

.admin-sidebar__link:hover {
    background: var(--admin-sidebar-hover, #2d2d3d);
    color: #ffffff;
}

.admin-sidebar__link--active {
    background: rgba(108, 92, 231, 0.1);
    color: var(--admin-primary, #6c5ce7);
    border-left-color: var(--admin-primary, #6c5ce7);
}

.admin-sidebar__link-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    flex-shrink: 0;
}

.admin-sidebar__link-text {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.admin-sidebar__badge {
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 11px;
    font-weight: 600;
    line-height: 1.4;
}

.admin-sidebar__badge--info {
    background: rgba(108, 92, 231, 0.2);
    color: var(--admin-primary, #6c5ce7);
}

.admin-sidebar__badge--warning {
    background: rgba(253, 203, 110, 0.2);
    color: #fdcb6e;
}

.admin-sidebar__badge--danger {
    background: rgba(225, 112, 85, 0.2);
    color: #e17055;
}

.admin-sidebar__footer {
    padding: 12px 0;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.admin-sidebar__link--back {
    color: rgba(255, 255, 255, 0.5);
}

.admin-sidebar__link--back:hover {
    color: #ffffff;
    background: var(--admin-sidebar-hover, #2d2d3d);
}
</style>
