<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
    <section class="container shadow p-3 mb-5 bg-body rounded">
        <nav class="navbar navbar-expand-md custom-navbar">
            <router-link class="navbar-brand" to="/">
                <img src="../assets/icons/logo-small.png" alt="" width="130" height="123">
            </router-link>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarToggler"
                aria-controls="navbarToggler" aria-expanded="fasle" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>

            <div class="collapse navbar-collapse" id="navbarToggler">

                <ul class="navbar-nav me-auto my-2 my-lg-">
                    <li class="nav-item dropdown ">
                        <router-link class="nav-link dropdown-toggle menu-item" to="#" role="button" data-bs-toggle="dropdown"
                            aria-expanded="false">
                            Кто&nbsp;я?..
                        </router-link>
                        <ul class="dropdown-menu">
                            <li>
                                <router-link class="dropdown-item" to="/briefly">Кратко</router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/artist">Художник</router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/illustrator">2D artist/Иллюстратор
                                </router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/volunteer">Волонтёр/Наставник
                                </router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/teacher">Преподаватель
                                </router-link>
                            </li>
                        </ul>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/allworks">Все&nbsp;работы</router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/services">Услуги</router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/news">Новости</router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/paydelivery">Оплата/Доставка</router-link>
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

.nav-item {
    font-family: "Montserrat", "Verdana regular", "Ebrima bold";
    font-size: 14pt;
    font-variant-caps: all-petite-caps;
}

.menu-item {
    --c: #73D1BE;
    --h: 1.5em;

    line-height: var(--h);
    background: linear-gradient(var(--c) 0 0) no-repeat calc(200% - var(--_p, 0%)) 100%/200% var(--_p, .08em);
    color: #0000;
    overflow: hidden;
    text-shadow:
        0 calc(-1*var(--_t, 0em)) var(--c),
        0 calc(var(--h) - var(--_t, 0em)) #fff;
    transition: .3s var(--_s, 0s), background-position .3s calc(.3s - var(--_s, 0s));
}

.menu-item:hover {
    --_t: var(--h);
    --_p: 100%;
    --_s: .3s
}

.dropdown-menu>li.active {
    background-color: #73D1BE;
}

.dropdown-menu>li:hover:active {
    --bs-dropdown-link-active-bg: #4b9e90;
}

.custom-navbar .fa-brands,
.fa {
    font-size: 30px;
    color: #73D1BE;
}

.navbar-brand :hover {
    transform: scale(1.05) rotate(-5deg);
    transition-duration: 0.5s;
}

@media (max-width:768px) {
    .custom-navbar .navbar-right {
        float: right;
        padding-right: 15px;
    }

    .custom-navbar .nav.navbar-nav.navbar-right li {
        float: right;
    }

    .custom-navbar .nav.navbar-nav.navbar-right li>a {
        padding: 8px 5px;
    }

    .custom-navbar .navbar-toggle {
        float: left
    }

    .custom-navbar .navbar-header {
        float: left;
        width: auto !important;
    }

    .custom-navbar .navbar-collapse {
        clear: both;
        float: none;
    }

    .navbar-right {
        visibility: hidden;
    }
}

.navbar-right {
    flex-wrap: nowrap;
}

.fa:hover,
.fa-brands:hover {
    filter: drop-shadow(2px 2px 2px #808080)
}
</style>
