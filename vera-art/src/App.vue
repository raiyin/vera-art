<script>
import Header from '@/components/UI/Header.vue';
import Footer from '@/components/UI/Footer.vue';

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

<style scoped></style>
