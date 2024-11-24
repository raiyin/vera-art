<script lang="ts">
export default {
    components: {},
    data() {
        return {
            isLoaded: false,
            loadingGrey: '#adadad',
        };
    },
};
</script>

<template>
    <div>
        <div class="title" />
        <div class="date" />
    </div>
</template>

<style scoped>
.title {
    background-color: grey;
    height: 2rem;
    width: 50%;
}

.date {
    margin-top: 0.7rem;
    background-color: grey;
    height: 1rem;
    width: 15%;
}

.title,
.date {
    background-color: v-bind(loadingGrey);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        v-bind(loadingGrey);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 4px;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}
</style>
