<script>
import MainHeader from '@/components/MainHeader.vue';
import MainFooter from '@/components/MainFooter.vue';

export default {
    components: {
        MainHeader,
        MainFooter,
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
    <MainHeader />
    <router-view :key="$route.fullPath"></router-view>
    <MainFooter />
</template>

<style scoped>
</style>
