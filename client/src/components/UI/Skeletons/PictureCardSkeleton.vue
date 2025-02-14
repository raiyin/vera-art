<script lang="ts">
import { useThemeStore } from '../../../stores/ThemeStore';

export default {
    setup() {
        const themeStore = useThemeStore();
        return { themeStore };
    },
};
</script>

<template>
    <div class="card p-3 mb-5 rounded col-lg-4 col-md-6">
        <div class="image" />

        <div class="card-body">
            <div class="desc">
                <h5 class="card-title" />

                <div class="card-text" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
}

.image {
    height: 20rem;
}

h5 {
    width: 80%;
    display: inline-block;
    text-align: center;
}

h5,
.card-text {
    height: 2rem;
    margin-bottom: 0.5rem;
}

.image,
h5,
.card-text {
    background-color: var(--skeleton-gray);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        var(--skeleton-gray);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 0.4rem;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}

.desc {
    align-self: flex-end;
    width: 100%;
}
</style>
