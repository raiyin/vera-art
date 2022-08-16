<script setup lang="ts">
</script>

<template>
    <section class="container shadow p-3 mb-5 bg-body rounded">
        <nav class="navbar navbar-expand-md custom-navbar">
            <a class="navbar-brand" href="#">
                <img src="../assets/icons/logo-small.png" alt="" width="130" height="123">
            </a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarToggler"
                aria-controls="navbarToggler" aria-expanded="fasle" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>

            <div class="collapse navbar-collapse" id="navbarToggler">

                <ul class="navbar-nav me-auto my-2 my-lg-">
                    <li class="nav-item dropdown">
                        <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown"
                            aria-expanded="false">
                            Кто&nbsp;я?
                        </a>
                        <ul class="dropdown-menu">
                            <li><a class="dropdown-item" href="./whoami/briefly.html">Кратко</a></li>
                            <li><a class="dropdown-item" href="./whoami/atist.html">Художник</a></li>
                            <li><a class="dropdown-item" href="./whoami/illustrator.html">2D artist/Иллюстратор</a></li>
                            <li><a class="dropdown-item" href="./whoami/volunteer.html">Волонтёр/Наставник</a></li>
                            <li><a class="dropdown-item" href="./whoami/teacher.html">Преподаватель</a></li>
                        </ul>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="./menu/all_works.html">Все&nbsp;работы</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="./menu/services.html">Услуги</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="./menu/news.html">Новости</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="./menu/pay_deliver.html">Оплата/Доставка</a>
                    </li>
                </ul>

                <ul class="nav navbar-nav navbar-right">
                    <li class="me-3">
                        <a href="https://t.me/MilayaV">
                            <i class="fa-brands fa-telegram">
                            </i>
                        </a>
                    </li>
                    <li class="me-3">
                        <a href="https://vk.com/perczukowa">
                            <i class="fa-brands fa-vk"></i>
                        </a>
                    </li>
                    <li>
                        <a href="mailto:perczukowa@yandex.ru">
                            <i class="fa fa-envelope"></i>
                        </a>
                    </li>
                </ul>
            </div>
        </nav>
    </section>
</template>

<style scoped>
</style>
