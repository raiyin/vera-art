<script>
import Card from './Card.vue';

export default {
    components: {
        Card,
    },
    props: {
        images: {
            type: Array,
            required: true,
        },
    },
    methods: {},
    computed: {},
};
</script>

<template>
    <section class="container text-center main-content px-0">
        <div class="row">
            <Card
                v-for="imgObject in images"
                :imageObject="imgObject"
                :key="imgObject.id"
            />
        </div>
    </section>
</template>

<style scoped></style>
