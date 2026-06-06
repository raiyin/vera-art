<template>
    <div class="admin-page">
        <div class="admin-page__header">
            <h1 class="admin-page__title">Курсы</h1>
            <NuxtLink to="/admin/courses/add" class="admin-page__btn"
                >+ Создать курс</NuxtLink
            >
        </div>
        <div class="admin-page__placeholder">
            <p>Управление курсами</p>
            <p class="admin-page__hint">Будет реализовано в Шаге 7</p>
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
@import './_shared.css';
</style>
