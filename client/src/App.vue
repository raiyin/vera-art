<script lang="ts" setup>
import AppHeader from '@/components/UI/AppHeader.vue';
import AppFooter from '@/components/UI/AppFooter.vue';
</script>

<template>
    <div class="app-container">
        <AppHeader />
        <div class="content">
            <router-view :key="$route.fullPath" />
        </div>
        <AppFooter />
    </div>
</template>

<style scoped>
.app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
}

.content {
    flex: auto;
}
</style>
