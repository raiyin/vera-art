<template>
    <div class="admin-page">
        <div class="admin-page__header">
            <h1 class="admin-page__title">Новости</h1>
            <NuxtLink to="/admin/news/add" class="admin-page__btn"
                >+ Добавить новость</NuxtLink
            >
        </div>
        <div class="admin-page__placeholder">
            <p>Управление новостями</p>
            <p class="admin-page__hint">Список новостей будет реализован в Шаге 6</p>
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
