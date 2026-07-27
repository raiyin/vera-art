<template>
    <div />
</template>

<script setup lang="ts">
definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

// Redirect to the existing shop add page
await navigateTo('/art-store/add', { external: true, },);
</script>
