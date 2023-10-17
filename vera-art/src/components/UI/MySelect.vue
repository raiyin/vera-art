<script lang="ts">
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

export default {
    name: 'my-select',
    components: {
        vSelect,
    },
    props: {
        modelValue: {
            type: String,
        },
        options: {
            type: Array,
            default: () => [],
        },
    },
    methods: {
        changeOption(event) {
            this.$emit('update:modelValue', event.target.value);
        },
        vselectOption(event) {
            this.$emit('update:modelValue', event.target.value);
        },
    },
};
</script>

<template>
    <v-select
        :options="options"
        label="name"
        :reduce="(item) => item.value"
        @input="changeOption"
        :value="selected"
    >
        Выберите из списка
    </v-select>
</template>

<style scoped>
.sort-select {
    max-width: 100px;
}
</style>
