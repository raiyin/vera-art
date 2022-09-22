<script>
import axios from 'axios'
import Card from "./Card.vue"
import MySelect from "./MySelect.vue"
import { ref, onMounted, computed, toRefs, onBeforeMount } from 'vue'
import gsap from 'gsap'

export default {
    components: {
        Card,
        MySelect
    },
    props: {
        images: {
            type: Array,
            required: true
        }
    },
    setup() {
        // const myrow = ref(null)
        // return {
        //     myrow
        // }
    },
    methods: {
        animate() {
            let query = gsap.utils.selector(this.$el);
            //let query = gsap.utils.selector(this.$refs.myrow);
            gsap.from(query('.card'), {
                opacity: 0,
                duration: 1
            });
        }
    },
    mounted() {
        this.$nextTick(() => this.animate());
    },
    computed: {
    }
}
</script>

<template>
    <section class="container text-center main-content px-0">
        <!-- <div class="row" ref="myrow"> -->
        <div class="row">
            <Card v-for="imgObject in images" :imageObject="imgObject" :key="imgObject.id" />
        </div>
    </section>
</template>

<style scoped>

</style>
