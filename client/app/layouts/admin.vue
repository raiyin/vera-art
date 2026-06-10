<template>
    <div class="admin-layout" :class="{ 'admin-layout--dark': isDark && mounted }">
        <!-- Mobile overlay -->
        <div
            v-if="sidebarOpen"
            class="admin-layout__overlay"
            @click="sidebarOpen = false"
        />

        <!-- Sidebar -->
        <AdminSidebar
            :is-open="sidebarOpen"
            :is-dark="isDark && mounted"
            @close="sidebarOpen = false"
        />

        <!-- Main content area -->
        <div class="admin-layout__main">
            <!-- Top header -->
            <AdminHeader
                :is-dark="isDark && mounted"
                :sidebar-open="sidebarOpen"
                @toggle-sidebar="sidebarOpen = !sidebarOpen"
                @toggle-theme="toggleTheme"
            />

            <!-- Page content — use route.fullPath as key to force re-render on navigation -->
            <main class="admin-layout__content">
                <div class="admin-layout__page" :key="route.fullPath">
                    <slot />
                </div>
            </main>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useThemeStore } from '~/stores/ThemeStore';

const themeStore = useThemeStore();
const route = useRoute();
const sidebarOpen = ref(false);
const mounted = ref(false);

const isDark = computed(() => themeStore.theme === 'dark');

onMounted(() => {
    mounted.value = true;
});

function toggleTheme() {
    themeStore.theme = isDark.value ? 'light' : 'dark';
}

// Close sidebar on route change (for mobile)
watch(
    () => route.path,
    () => {
        sidebarOpen.value = false;
    }
);
</script>

<style scoped>
.admin-layout {
    display: flex;
    min-height: 100vh;
    background: var(--admin-bg, #f0f2f5);
    color: var(--admin-text-primary, #2d3436);
    transition: background 0.3s ease, color 0.3s ease;
}

.admin-layout--dark {
    --admin-bg: #0f1117;
    --admin-surface: #1a1d29;
    --admin-sidebar-bg: #0a0c12;
    --admin-sidebar-hover: #1a1d29;
    --admin-text-primary: #e0e0e0;
    --admin-text-secondary: #a0a5b5;
    --admin-border: #2d2d3d;
    --admin-primary: #6c5ce7;
    --admin-primary-hover: #5a4bd1;
}

.admin-layout:not(.admin-layout--dark) {
    --admin-bg: #f0f2f5;
    --admin-surface: #ffffff;
    --admin-sidebar-bg: #1a1d29;
    --admin-sidebar-hover: #2d2d3d;
    --admin-text-primary: #2d3436;
    --admin-text-secondary: #636e72;
    --admin-border: #e0e0e0;
    --admin-primary: #6c5ce7;
    --admin-primary-hover: #5a4bd1;
}

.admin-layout__overlay {
    display: none;
}

@media (max-width: 768px) {
    .admin-layout__overlay {
        display: block;
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
        z-index: 40;
    }
}

.admin-layout__main {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
    margin-left: 260px;
    transition: margin-left 0.3s ease;
}

@media (max-width: 768px) {
    .admin-layout__main {
        margin-left: 0;
    }
}

.admin-layout__content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
}

@media (max-width: 768px) {
    .admin-layout__content {
        padding: 16px;
    }
}

.admin-layout__page {
    animation: adminLayoutPageIn 0.3s ease-out;
}

@keyframes adminLayoutPageIn {
    from {
        opacity: 0;
        transform: translateY(8px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
