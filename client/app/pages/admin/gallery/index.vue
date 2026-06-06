<template>
    <div class="admin-page">
        <div class="admin-page__header">
            <h1 class="admin-page__title">Галерея</h1>
            <NuxtLink to="/admin/gallery/add" class="admin-page__btn"
                >+ Добавить работу</NuxtLink
            >
        </div>
        <div class="admin-page__placeholder">
            <p>Управление работами в галерее</p>
            <p class="admin-page__hint">Список работ будет реализован в Шаге 4</p>
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
