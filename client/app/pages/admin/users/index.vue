<template>
    <div class="admin-page">
        <div class="admin-page__header">
            <h1 class="admin-page__title">Пользователи</h1>
        </div>
        <div class="admin-page__placeholder">
            <p>Управление пользователями</p>
            <p class="admin-page__hint">Будет реализовано в Шаге 9</p>
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
