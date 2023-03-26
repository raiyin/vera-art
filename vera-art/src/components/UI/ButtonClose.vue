<script>

export default {
    components: {
    },
    props: {
        imageObject: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {};
    },
    methods: {},
    computed: {},
    watch: {},
};
</script>

<template>
    <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close" />
</template>

<style scoped>
.btn-close {
    position: absolute;
    right: 16px;
    background-color: bisque;
    z-index: 100;
    min-width: 1rem;
    min-height: 1rem;
}
</style>
