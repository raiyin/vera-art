<template>
    <header
        class="admin-header"
        :class="{ 'admin-header--dark': isDark, }"
    >
        <div class="admin-header__left">
            <button
                class="admin-header__menu-btn"
                @click="$emit('toggleSidebar',)"
            >
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="20"
                    height="20"
                >
                    <line
                        x1="3"
                        y1="6"
                        x2="21"
                        y2="6"
                    />
                    <line
                        x1="3"
                        y1="12"
                        x2="21"
                        y2="12"
                    />
                    <line
                        x1="3"
                        y1="18"
                        x2="21"
                        y2="18"
                    />
                </svg>
            </button>

            <div class="admin-header__search">
                <svg
                    class="admin-header__search-icon"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="16"
                    height="16"
                >
                    <circle
                        cx="11"
                        cy="11"
                        r="8"
                    />
                    <line
                        x1="21"
                        y1="21"
                        x2="16.65"
                        y2="16.65"
                    />
                </svg>
                <input
                    v-model="searchQuery"
                    type="text"
                    class="admin-header__search-input"
                    placeholder="Поиск..."
                    @input="onSearch"
                >
            </div>
        </div>

        <div class="admin-header__right">
            <button
                class="admin-header__icon-btn"
                :title="isDark ? 'Светлая тема' : 'Тёмная тема'"
                @click="$emit('toggleTheme',)"
            >
                <svg
                    v-if="isDark"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="20"
                    height="20"
                >
                    <circle
                        cx="12"
                        cy="12"
                        r="5"
                    />
                    <line
                        x1="12"
                        y1="1"
                        x2="12"
                        y2="3"
                    />
                    <line
                        x1="12"
                        y1="21"
                        x2="12"
                        y2="23"
                    />
                    <line
                        x1="4.22"
                        y1="4.22"
                        x2="5.64"
                        y2="5.64"
                    />
                    <line
                        x1="18.36"
                        y1="18.36"
                        x2="19.78"
                        y2="19.78"
                    />
                    <line
                        x1="1"
                        y1="12"
                        x2="3"
                        y2="12"
                    />
                    <line
                        x1="21"
                        y1="12"
                        x2="23"
                        y2="12"
                    />
                    <line
                        x1="4.22"
                        y1="19.78"
                        x2="5.64"
                        y2="18.36"
                    />
                    <line
                        x1="18.36"
                        y1="5.64"
                        x2="19.78"
                        y2="4.22"
                    />
                </svg>
                <svg
                    v-else
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="20"
                    height="20"
                >
                    <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
                </svg>
            </button>

            <NuxtLink
                to="/admin/reviews"
                class="admin-header__icon-btn admin-header__notif-btn"
                title="Ожидающие отзывы"
            >
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="20"
                    height="20"
                >
                    <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" />
                    <path d="M13.73 21a2 2 0 01-3.46 0" />
                </svg>
                <span
                    v-if="pendingReviews > 0"
                    class="admin-header__notif-badge"
                >
                    {{ pendingReviews > 99 ? '99+' : pendingReviews }}
                </span>
            </NuxtLink>

            <div class="admin-header__profile">
                <div class="admin-header__avatar">
                    {{ userInitial }}
                </div>
                <div class="admin-header__profile-info">
                    <span class="admin-header__profile-name">{{ userName }}</span>
                    <span class="admin-header__profile-role">Admin</span>
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, } from 'vue';
import { useAuthStore, } from '~/stores/AuthStore';

const props = defineProps<{
    isDark: boolean
    sidebarOpen: boolean
}>();

defineEmits<{
    toggleSidebar: []
    toggleTheme: []
}>();

const authStore = useAuthStore();
const searchQuery = ref('',);
const mounted = ref(false,);

onMounted(() => {
    mounted.value = true;
});

const userName = computed(() => {
    // During SSR and before client-side mount, show generic name
    // After mount, show user-specific name from localStorage/token
    if (!mounted.value) return 'Admin';
    return authStore.userId ? `Admin #${authStore.userId}` : 'Admin';
});

const userInitial = computed(() => {
    return userName.value.charAt(0,).toUpperCase();
});

const pendingReviews = ref(0,);

function onSearch() {
    // Will be implemented with global search later
}
</script>

<style scoped>
.admin-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
    padding: 0 24px;
    background: var(--admin-surface, #ffffff);
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
    position: sticky;
    top: 0;
    z-index: 30;
    transition: background 0.3s ease, border-color 0.3s ease;
}

.admin-header--dark {
    background: var(--admin-surface, #1a1d29);
    border-color: var(--admin-border, #2d2d3d);
}

.admin-header__left {
    display: flex;
    align-items: center;
    gap: 16px;
    flex: 1;
}

.admin-header__menu-btn {
    display: none;
    background: none;
    border: none;
    color: var(--admin-text-primary, #2d3436);
    cursor: pointer;
    padding: 4px;
    border-radius: 6px;
    transition: background 0.2s;
}

.admin-header__menu-btn:hover {
    background: rgba(0, 0, 0, 0.05);
}

@media (max-width: 768px) {
    .admin-header__menu-btn {
        display: flex;
    }
}

.admin-header__search {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--admin-bg, #f0f2f5);
    border-radius: 8px;
    padding: 8px 12px;
    max-width: 320px;
    width: 100%;
    transition: background 0.3s ease;
}

.admin-header--dark .admin-header__search {
    background: rgba(255, 255, 255, 0.05);
}

.admin-header__search-icon {
    color: var(--admin-text-secondary, #636e72);
    flex-shrink: 0;
}

.admin-header__search-input {
    background: none;
    border: none;
    outline: none;
    color: var(--admin-text-primary, #2d3436);
    font-size: 14px;
    width: 100%;
}

.admin-header__search-input::placeholder {
    color: var(--admin-text-secondary, #636e72);
}

.admin-header__right {
    display: flex;
    align-items: center;
    gap: 8px;
}

.admin-header__icon-btn {
    position: relative;
    background: none;
    border: none;
    color: var(--admin-text-secondary, #636e72);
    cursor: pointer;
    padding: 8px;
    border-radius: 8px;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
}

.admin-header__icon-btn:hover {
    background: rgba(0, 0, 0, 0.05);
    color: var(--admin-text-primary, #2d3436);
}

.admin-header--dark .admin-header__icon-btn:hover {
    background: rgba(255, 255, 255, 0.1);
}

.admin-header__notif-badge {
    position: absolute;
    top: 4px;
    right: 4px;
    background: var(--admin-danger, #e17055);
    color: white;
    font-size: 10px;
    font-weight: 700;
    padding: 1px 5px;
    border-radius: 8px;
    min-width: 16px;
    text-align: center;
    line-height: 1.4;
    animation: adminNotifPulse 2s ease-in-out infinite;
}

@keyframes adminNotifPulse {
    0%,
    100% {
        transform: scale(1);
    }
    50% {
        transform: scale(1.1);
    }
}

.admin-header__profile {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 4px 8px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s;
}

.admin-header__profile:hover {
    background: rgba(0, 0, 0, 0.03);
}

.admin-header--dark .admin-header__profile:hover {
    background: rgba(255, 255, 255, 0.05);
}

.admin-header__avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--admin-primary, #6c5ce7), #a29bfe);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 700;
}

.admin-header__profile-info {
    display: flex;
    flex-direction: column;
    gap: 1px;
}

@media (max-width: 480px) {
    .admin-header__profile-info {
        display: none;
    }
}

.admin-header__profile-name {
    font-size: 13px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.2;
}

.admin-header__profile-role {
    font-size: 11px;
    color: var(--admin-text-secondary, #636e72);
    line-height: 1.2;
}
</style>
