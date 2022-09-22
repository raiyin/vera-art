<script>

import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"
import Gallery from "@/components/Gallery.vue"
import axios from 'axios'
import vSelect from 'vue-select'
import 'vue-select/dist/vue-select.css';

export default {
    components: {
        Header,
        Footer,
        Gallery,
        vSelect
    },
    props: [],
    data() {
        return {
            images: [],
            selectedSort: '',
            sortOptions: [
                { value: "name_ru", name: 'По названию' },
                { value: "year", name: 'По новизне' },
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
        }
    },
    mounted() {
        this.fetchData();
    },
    computed: {
        // sortedImages() {
        //     return [...this.images].sort((image_first, image_second) => {
        //         if (typeof image_first[this.selectedSort] === "string")
        //             return image_first[this.selectedSort]?.localeCompare(image_second[this.selectedSort])
        //         if (typeof image_first[this.selectedSort] === "number")
        //             return image_first[this.selectedSort] - image_second[this.selectedSort]
        //     })
        // }
    },
    watch: {
        selectedSort(newValue, oldValue) {
            this.images.sort((image_first, image_second) => {
                if (typeof image_first[this.selectedSort] === "string")
                    return image_first[this.selectedSort]?.localeCompare(image_second[this.selectedSort])
                if (typeof image_first[this.selectedSort] === "number")
                    return image_first[this.selectedSort] - image_second[this.selectedSort]

            })
        }
    },
}
</script>

<template>
    <Header />
    <section class="container mt-3 px-0">
        <v-select :options="sortOptions" :reduce="item => item.value" label="name" v-model="selectedSort"
            inputId="value" placeholder="Выберите из списка">
            Выберите из списка
        </v-select>
    </section>
    <Gallery :images="images" />
    <Footer />
</template>

<style scoped>
.v-select {
    width: 20rem;
    box-sizing: border-box;
}
</style>