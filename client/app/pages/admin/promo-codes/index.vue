<template>
    <div class="admin-page">
        <div class="admin-page__header">
            <h1 class="admin-page__title">Промокоды</h1>
            <NuxtLink to="/admin/promo-codes/add" class="admin-page__btn"
                >+ Создать промокод</NuxtLink
            >
        </div>
        <div class="admin-page__placeholder">
            <p>Управление промокодами</p>
            <p class="admin-page__hint">Будет реализовано в Шаге 12</p>
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
