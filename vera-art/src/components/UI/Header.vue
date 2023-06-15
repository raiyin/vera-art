<script>
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
    },
    components: {
        vSelect,
    },
    data() {
        return {
            select: { lang: 'ru', abbr: this.$t('locale.ru') },
            languages: [
                { lang: 'ru', abbr: this.$t('locale.ru') },
                { lang: 'en', abbr: this.$t('locale.en') },
            ],
        };
    },
    methods: {
        onChange(event) {
            if (event['lang'] !== null) {
                i18n.global.locale.value = event['lang'];
            }
        },
    },
};
</script>

<template>
    <section class="container-fluid my-shadow bg-body rounded">
        <nav class="navbar navbar-expand-md custom-navbar">
            <router-link class="navbar-brand" to="/">
                <img
                    src="../../assets/icons/logo-small.png"
                    alt=""
                    width="65"
                    height="62"
                />
            </router-link>
            <button
                class="navbar-toggler"
                type="button"
                data-bs-toggle="collapse"
                data-bs-target="#navbarToggler"
                aria-controls="navbarToggler"
                aria-expanded="fasle"
                aria-label="Toggle navigation"
            >
                <span class="navbar-toggler-icon"></span>
            </button>

            <div class="collapse navbar-collapse" id="navbarToggler">
                <ul class="navbar-nav me-auto my-2 my-lg-">
                    <li v-if="false" class="nav-item dropdown">
                        <router-link
                            class="nav-link dropdown-toggle menu-item"
                            to="#"
                            role="button"
                            data-bs-toggle="dropdown"
                            aria-expanded="false"
                        >
                            {{ $t('header.about_me') }}
                        </router-link>
                        <ul class="dropdown-menu">
                            <li>
                                <router-link class="dropdown-item" to="/briefly">
                                    {{ $t('header.briefly') }}
                                </router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/artist">
                                    {{ $t('header.artist') }}
                                </router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/illustrator">
                                    {{ $t('header.illustrator') }}
                                </router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/volunteer">
                                    {{ $t('header.volunteer') }}
                                </router-link>
                            </li>
                            <li>
                                <router-link class="dropdown-item" to="/teacher">
                                    {{ $t('header.teacher') }}
                                </router-link>
                            </li>
                        </ul>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/allworks">
                            {{ $t('header.all_works') }}
                        </router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/news">
                            {{ $t('header.news') }}
                        </router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/shop">
                            {{ $t('header.shop') }}
                        </router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/services">
                            {{ $t('header.services') }}
                        </router-link>
                    </li>
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/paydelivery">
                            {{ $t('header.payment') }}
                        </router-link>
                    </li>
                </ul>

                <!-- <v-select
                    v-model="$i18n.locale"
                    :options="languages"
                    label="abbr"
                    style="width: 170px; margin-right: 3em"
                    @update:modelValue="onChange"
                >
                </v-select> -->

                <ul class="nav navbar-nav navbar-right d-flex align-items-center">
                    <li class="me-3">
                        <select v-model="$i18n.locale" class="header-locale">
                            <option
                                v-for="locale in $i18n.availableLocales"
                                :key="`locale-${locale}`"
                                :value="locale"
                            >
                                {{ locale }}
                            </option>
                        </select>
                    </li>
                    <li class="me-3">
                        <a
                            href="https://t.me/MilayaV"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            <i class="fa-brands fa-telegram"> </i>
                        </a>
                    </li>
                    <li class="me-3">
                        <a
                            href="https://vk.com/perczukowa"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            <i class="fa-brands fa-vk"></i>
                        </a>
                    </li>
                    <li>
                        <a
                            href="mailto:perczukowa@yandex.ru"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            <i class="fa fa-envelope"></i>
                        </a>
                    </li>
                </ul>
            </div>
        </nav>
    </section>
</template>

<style scoped>
.my-shadow {
    box-shadow: 0 0 0.5rem 0.5rem rgba(0, 0, 0, 0.15) !important;
}

.nav-item {
    font-family: 'Montserrat', 'Verdana regular', 'Ebrima bold';
    font-size: 14pt;
    font-variant-caps: all-petite-caps;
}

.menu-item {
    --c: #000000;
    --m: #73d1be;
    --h: 1.75em;

    line-height: var(--h);
    background: linear-gradient(var(--m) 0 0) no-repeat calc(200% - var(--_p, 0%)) 100%/200%
        var(--_p, 0.08em);
    color: #0000;
    overflow: hidden;
    text-shadow: 0 calc(-1 * var(--_t, 0em)) var(--c),
        0 calc(var(--h) - var(--_t, 0em)) #fff;
    transition: 0.3s var(--_s, 0s), background-position 0.3s calc(0.3s - var(--_s, 0s));
}

.menu-item:hover {
    --_t: var(--h);
    --_p: 100%;
    --_s: 0.3s;
}

.dropdown-menu > li.active {
    background-color: #73d1be;
}

.dropdown-menu > li:hover:active {
    --bs-dropdown-link-active-bg: #4b9e90;
}

.header-locale {
    border-width: 1px;
    border-style: solid;
    border-radius: 4px;
    border-color: var(--vs-colors--lightest);
    background-color: #fff;
    color: var(--bs-body-color);
}

.custom-navbar .fa-brands,
.fa {
    font-size: 30px;
    color: #73d1be;
}

.navbar-brand :hover {
    transform: scale(1.05) rotate(-5deg);
    transition-duration: 0.5s;
}

@media (max-width: 768px) {
    .custom-navbar .navbar-right {
        float: right;
        padding-right: 15px;
    }

    .custom-navbar .nav.navbar-nav.navbar-right li {
        float: right;
    }

    .custom-navbar .nav.navbar-nav.navbar-right li > a {
        padding: 8px 5px;
    }

    .custom-navbar .navbar-toggle {
        float: left;
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
    filter: drop-shadow(2px 2px 2px #808080);
}
</style>