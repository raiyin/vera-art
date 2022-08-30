<script setup lang="ts">
import CalendarIcon from "@/components/icons/IconCalendar.vue"
</script>

<template>
    <section class="container text-center px-0">
        <div class="col d-flex justify-content-center">
            <div class="news-item">
                <router-link to="/news/newsitem">
                    <div class="img-holder">
                        <div class="news-content">
                            <p>Учебный год начинается!</p>
                            <p>Запрыгиваем!!!</p>
                            <p>
                                <CalendarIcon />
                                <span> 28 августа, 2022</span>
                            </p>
                        </div>
                    </div>
                </router-link>
            </div>
        </div>
    </section>
</template>

<style scoped>
.news-item {
    width: 22rem;
    height: 22rem;
    border-radius: 0.3rem;
    border: 1px solid wheat;
    overflow: hidden;
}

.img-holder {
    width: 100%;
    height: 100%;
    background-image: url("@/assets/img/news/2022/08/28/back.jpg");
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100%;
    transition: background-size 1s ease-out;
}

.news-item>a {
    text-decoration: none;
}

.img-holder:hover {
    background-size: 110%;
    cursor: pointer;
}

.news-content {
    position: relative;
    font-family: Ubuntu;
    font-weight: bold;
    font-size: larger;
    text-align: left;
    text-indent: 2rem;
    line-height: 0.4rem;
    top: 17rem;
    color: white;
    letter-spacing: 0.05rem;
}

.news-content>p>span {
    font-size: small;
    font-weight: lighter;
}
</style>
