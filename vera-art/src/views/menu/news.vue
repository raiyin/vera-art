<script setup lang="ts">

import Header from "../../components/Header.vue"
import Footer from "../../components/Footer.vue"
import Slogan from "../../components/Slogan.vue"

</script>

<template>
    <Header />
    <Slogan />

    <section class="container text-center px-0">
        <div class="row gx-0">

            <div class="col d-flex justify-content-center">
                <div class="news-item">
                    <div class="img-holder">
                        <div class="news-content">
                            <p>Новый год начинается!</p>
                            <p>Запрыгиваем!!!</p>
                            <p>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                                    class="bi bi-calendar" viewBox="0 0 16 16">
                                    <path
                                        d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5zM1 4v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V4H1z" />
                                </svg>
                                <span> 28 сентября, 2022</span>
                            </p>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </section>

    <Footer />
</template>

<style scoped>
.news-item {
    width: 22rem;
    height: 22rem;
    border-radius: 0.3rem;
    border: 1px solid #000000;
    overflow: hidden;
}

.img-holder {
    width: 100%;
    height: 100%;
    background-image: url("../../assets/img/news/2022/08/28/back.jpg");
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100%;
    transition: background-size 1s ease-out;
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
