<script lang="ts" setup>
import AppHeader from '@/components/UI/AppHeader.vue';
import AppFooter from '@/components/UI/AppFooter.vue';
</script>

<template>
    <AppHeader />
    <router-view :key="$route.fullPath" />
    <AppFooter />
</template>
