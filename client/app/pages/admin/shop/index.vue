<template>
    <div class="admin-page">
        <div class="admin-page__header">
            <h1 class="admin-page__title">Магазин</h1>
            <NuxtLink to="/admin/shop/add" class="admin-page__btn"
                >+ Добавить товар</NuxtLink
            >
        </div>
        <div class="admin-page__placeholder">
            <p>Управление товарами в магазине</p>
            <p class="admin-page__hint">Список товаров будет реализован в Шаге 5</p>
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
