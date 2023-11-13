<script lang="ts" setup>
import Header from '@/components/UI/Header.vue';
import Footer from '@/components/UI/Footer.vue';
</script>

<template>
    <div>
        <Header />
        <router-view :key="$route.fullPath" />
        <Footer />
    </div>
</template>
