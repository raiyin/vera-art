<script lang="ts" setup>
import AppHeader from '@/components/app-ui/AppHeader.vue';
import AppFooter from '@/components/app-ui/AppFooter.vue';
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
