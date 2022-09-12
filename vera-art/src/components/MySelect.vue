<script>
import vSelect from 'vue-select'
import 'vue-select/dist/vue-select.css';


export default {
    name: 'my-select',
    components: {
        vSelect,
    },
    props: {
        modelValue: {
            type: String
        },
        options: {
            type: Array,
            default: () => []
        },
    },
    methods: {
        changeOption(event) {
            console.log(event);
            this.$emit("update:modelValue", event.target.value)
        },
        vselectOption(event) {
            console.log(event);
            this.$emit("update:modelValue", event.target.value)
        }
    }
}
</script>

<template>
    <section class="container px-0">
        <v-select :options="options" label="name" :reduce="item => item.value" @input="changeOption" :value="selected">
            Выберите из списка
        </v-select>

        <!-- <select name="" id="" v-model="modelValue" @change="changeOption">
            <option disable value="">Выберите из списка</option>
            <option v-for="option in options" :key="option.value" :value="option.value">{{ option.name }}</option>
        </select> -->
    </section>
</template>

<style scoped>

</style>