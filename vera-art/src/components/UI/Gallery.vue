<script lang="ts">
import type { ImageProps } from '@/types';
import type { PropType } from 'vue';
import Card from './Card.vue';

export default {
    components: {
        Card,
    },
    props: {
        images: {
            type: [Array] as PropType<ImageProps[]>,
            required: true,
        },
    },
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
