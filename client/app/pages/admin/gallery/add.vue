<template>
    <div />
</template>

<script setup lang="ts">
definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

// Redirect to the existing gallery add page
await navigateTo('/gallery/add', { external: true, },);
</script>
