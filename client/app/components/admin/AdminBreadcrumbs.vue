<template>
    <nav
        class="admin-breadcrumbs"
        :class="{ 'admin-breadcrumbs--dark': isDark, }"
    >
        <NuxtLink
            v-for="(crumb, idx) in items"
            :key="idx"
            :to="crumb.to"
            class="admin-breadcrumbs__item"
            :class="{ 'admin-breadcrumbs__item--last': idx === items.length - 1, }"
        >
            <span
                v-if="idx > 0"
                class="admin-breadcrumbs__sep"
            >/</span>
            {{ crumb.label }}
        </NuxtLink>
    </nav>
</template>

<script setup lang="ts">
    withDefaults(
        defineProps<{
            items: { label: string, to?: string }[]
            isDark?: boolean
        }>(),
        {
            items: () => [],
            isDark: false,
        }
    );
</script>

<style scoped>
.admin-breadcrumbs {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 16px;
}

.admin-breadcrumbs__item {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
    text-decoration: none;
    transition: color 0.2s;
    display: flex;
    align-items: center;
    gap: 4px;
}

.admin-breadcrumbs__item:hover:not(.admin-breadcrumbs__item--last) {
    color: var(--admin-primary, #6c5ce7);
}

.admin-breadcrumbs__item--last {
    color: var(--admin-text-primary, #2d3436);
    font-weight: 600;
    pointer-events: none;
}

.admin-breadcrumbs--dark .admin-breadcrumbs__item--last {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-breadcrumbs__sep {
    color: var(--admin-text-secondary, #636e72);
    opacity: 0.5;
    margin-right: 4px;
}
</style>
