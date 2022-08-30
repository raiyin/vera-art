<script setup lang="ts">
    import CalendarIcon from "@/components/icons/IconCalendar.vue"
</script>
        
<template>

    <div class="other-news-item-wrapper">
        <router-link to="/news/newsitem" class="other-news-item">
            <div class="other-news-img">
                <img src="@/assets/img/news/2022/08/28/back.jpg" />
            </div>
            <div class="other-news-desc">
                <h6>Учебный год начинается! Запрыгиваем!!!</h6>
                <div class="date">
                    <CalendarIcon />
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
        