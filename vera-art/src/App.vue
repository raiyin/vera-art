<script>
import Header from '@/components/Header.vue';
import Footer from '@/components/Footer.vue';

export default {
    components: {
        Header,
        Footer,
    },
    props: [],
    data() {
        return {};
    },
    methods: {},
    computed: {},
    watch: {},
};
</script>

<template>
    <Header />
    <router-view :key="$route.fullPath"></router-view>
    <Footer />
</template>

<style scoped>
</style>
