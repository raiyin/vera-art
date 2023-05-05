<script setup>
import { useI18n } from 'vue-i18n';
import Header from '@/components/UI/Header.vue';
import Footer from '@/components/UI/Footer.vue';
const { t } = useI18n({ useScope: 'global' });
</script>

<template>
    <Header />
    <router-view :key="$route.fullPath" />
    <Footer />
</template>
