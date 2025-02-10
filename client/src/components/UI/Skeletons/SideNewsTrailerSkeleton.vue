<script lang="ts">
export default {};
</script>

<template>
    <div class="side-news">
        <div class="img" />
        <div class="desc">
            <div class="title" />
            <div class="date" />
        </div>
    </div>
</template>

<style scoped>
.side-news {
    display: flex;
    height: 60px;
}

.img {
    min-width: 90px;
}

.desc {
    margin-left: 1rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    width: 100%;
}

.title {
    text-align: left;
    width: 100%;
    height: 1.5rem;
}

.date {
    width: 40%;
    height: 1.5rem;
}

.img,
.title,
.date {
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
    border-radius: 4px;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}
</style>
