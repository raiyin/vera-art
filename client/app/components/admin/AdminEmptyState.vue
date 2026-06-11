<template>
    <div class="admin-empty" :class="{ 'admin-empty--dark': isDark }">
        <div class="admin-empty__icon">
            <svg
                v-if="icon === 'data'"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                width="48"
                height="48"
            >
                <ellipse cx="12" cy="5" rx="9" ry="3" />
                <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3" />
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" />
            </svg>
            <svg
                v-else-if="icon === 'search'"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                width="48"
                height="48"
            >
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <svg
                v-else-if="icon === 'box'"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                width="48"
                height="48"
            >
                <path
                    d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z"
                />
                <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
                <line x1="12" y1="22.08" x2="12" y2="12" />
            </svg>
            <svg
                v-else
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                width="48"
                height="48"
            >
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="8" x2="12" y2="12" />
                <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
        </div>
        <h3 class="admin-empty__title">{{ title }}</h3>
        <p v-if="description" class="admin-empty__description">{{ description }}</p>
        <div v-if="$slots.default" class="admin-empty__action">
            <slot />
        </div>
    </div>
</template>

<script setup lang="ts">
withDefaults(
    defineProps<{
        icon?: 'data' | 'search' | 'box' | 'info';
        title?: string;
        description?: string;
        isDark?: boolean;
    }>(),
    {
        icon: 'data',
        title: 'Нет данных',
        description: '',
        isDark: false,
    }
);
</script>

<style scoped>
.admin-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    text-align: center;
}

.admin-empty__icon {
    color: var(--admin-text-secondary, #636e72);
    opacity: 0.3;
    margin-bottom: 16px;
}

.admin-empty__title {
    font-size: 18px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin: 0 0 8px;
}

.admin-empty--dark .admin-empty__title {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-empty__description {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    margin: 0 0 20px;
    max-width: 360px;
    line-height: 1.5;
}

.admin-empty__action {
    margin-top: 4px;
}
</style>
