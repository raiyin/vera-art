<script setup lang="ts">
</script>
        
<template>

    <div class="other-news-item-wrapper">
        <router-link to="/news/newsitem" class="other-news-item">
            <div class="other-news-img">
                <img src="../assets/img/news/2022/08/28/back.jpg" />
            </div>
            <div class="other-news-desc">
                <h6>Учебный год начинается! Запрыгиваем!!!</h6>
                <div class="date">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                        class="bi bi-calendar fa-lg" viewBox="0 0 16 16">
                        <path d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 
                                    0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 
                                    2-2h1V.5a.5.5 0 0 1 .5-.5zM1 4v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V4H1z" />
                    </svg>
                    <span> 28 августа, 2022</span>
                </div>
            </div>
        </router-link>
    </div>

</template>
        
<style scoped>
.other-news-item-wrapper {
    margin-bottom: 1rem;
}

.other-news-item {
    display: flex;
    text-decoration: none;
    color: black;
}

.other-news-img {
    width: 70px;
    height: 70px;
    margin-right: 1rem;
}

.other-news-img>img {
    max-width: 100%;
    max-height: 100%;
}

.other-news-desc {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 0.1em;
}
</style>
        