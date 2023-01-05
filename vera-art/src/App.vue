<script>
import HeaderBlock from '@/components/Header.vue';
import FooterBlock from '@/components/Footer.vue';

export default {
    components: {
        HeaderBlock,
        FooterBlock,
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
    <HeaderBlock />
    <router-view :key="$route.fullPath"></router-view>
    <FooterBlock />
</template>

<style scoped></style>
