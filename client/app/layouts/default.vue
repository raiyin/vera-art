<template>
    <div class="default-layout">
        <Header />
        <UMain>
            <slot :key="route.fullPath" />
        </UMain>
        <AppFooter />
    </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import Header from '../components/Header.vue';
import AppFooter from '../components/AppFooter.vue';

const route = useRoute();
</script>

<style scoped>
.default-layout {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}

.default-layout > .umain {
    flex: 1;
}
</style>
