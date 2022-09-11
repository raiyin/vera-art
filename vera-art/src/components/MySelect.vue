<script>

export default {
    name: 'my-select',
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
            this.$emit("update:modelValue", event.target.value)
        }
    },
    data: () => ({
        value: 'B'
    })
}
</script>

<template>
    <section class="container px-0">

        <select name="" id="" v-model="modelValue" @change="changeOption">
            <option disable value="">Выберите из списка</option>
            <option v-for="option in options" :key="option.value" :value="option.value">{{ option.name }}</option>
        </select>
        
    </section>
</template>

<style scoped>
</style>