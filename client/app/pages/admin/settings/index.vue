<template>
    <div class="admin-page">
        <div class="admin-page__breadcrumbs">
            <NuxtLink to="/admin">Dashboard</NuxtLink>
            <span class="admin-page__breadcrumbs-sep">/</span>
            <span class="admin-page__breadcrumbs-current">Настройки</span>
        </div>
        <div class="admin-page__header">
            <h1 class="admin-page__title">
                Настройки
            </h1>
        </div>
        <div class="admin-page__placeholder">
            <p>Управление справочниками: материалы, основы</p>
            <p class="admin-page__hint">
                Будет реализовано в Шаге 14
            </p>
        </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});
</script>

<style scoped>
@import '../_shared.css';
</style>
