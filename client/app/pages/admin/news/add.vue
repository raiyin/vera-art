<template>
    <div />
</template>

<script setup lang="ts">
definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

// Redirect to the existing news add page
await navigateTo('/news/add', { external: true, },);
</script>
