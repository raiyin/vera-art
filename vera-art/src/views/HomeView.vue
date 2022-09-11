<script>

import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"
import Slogan from "../components/Slogan.vue"
import Gallery from "@/components/Gallery.vue"
import MySelect from "@/components/MySelect.vue"
import axios from 'axios'

export default {
    components: {
        Header,
        Footer,
        Slogan,
        Gallery,
        MySelect
    },
    props: [],
    data() {
        return {
            images: [],
            selectedSort: '',
            sortOptions: [
                { value: "name_ru", name: 'По названию' },
                { value: "year", name: 'По году' },
                { value: "height", name: 'По высоте' },
                { value: "width", name: 'По ширине' }
            ],
        }
    },
    setup() {
    },
    methods: {
        async fetchData() {
            try {
                const response = await axios.get('http://localhost:3001/sale');
                this.images = response.data;
            }
            catch (e) {
                alert('Error')
                console.log(e);
                var propertyNames = Object.getOwnPropertyNames(e);
                propertyNames.forEach(function (property) {
                    var descriptor = Object.getOwnPropertyDescriptor(e, property);
                    console.log(property + ":" + e[property] + ":" + propsToStr(descriptor));
                });
            }
        },
    },
    mounted() {
        this.fetchData();
    },
    computed: {
        sortedImages() {
            return [...this.images].sort((images_first, images_second) => {
                if (typeof images_first[this.selectedSort]==="string")
                    return images_first[this.selectedSort]?.localeCompare(images_second[this.selectedSort])
                if (typeof images_first[this.selectedSort] === "number")
                    return images_first[this.selectedSort] - images_second[this.selectedSort]
            })
        }
    }
}
</script>

<template>
    <Header />
    <Slogan />
    <my-select v-model="selectedSort" :options="sortOptions" />
    <Gallery :images="sortedImages" />
    <Footer />
</template>

<style scoped>

</style>