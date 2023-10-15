<script>
import 'vue-select/dist/vue-select.css';
import { useI18n } from 'vue-i18n';
import DayIcon from '@/components/Icons/icon_day.vue';
import NightIcon from '@/components/Icons/icon_night.vue';
import { useThemeStore } from '../../stores/ThemeStore';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
        const themeStore = useThemeStore();
        return { themeStore };
    },
    components: {
        DayIcon,
        NightIcon,
        vSelect,
    },
    data() {
        return {
            select: { lang: 'RUS', abbr: this.$t('locale.RUS') },
            languages: [
                { lang: 'RUS', abbr: this.$t('locale.RUS') },
                { lang: 'ENG', abbr: this.$t('locale.ENG') },
            ],
        };
    },
    mounted() {
        const navLinks = document.querySelectorAll('.nav-item');
        const menuToggle = document.getElementById('navbarToggler');
        const bsCollapse = new bootstrap.Collapse(menuToggle, {
            toggle: false,
        });
        navLinks.forEach((l) => {
            l.addEventListener('click', () => {
                if (window.screen.width < 768) {
                    bsCollapse.toggle();
                }
            });
        });
    },
    methods: {
        onChange(event) {
            if (event['lang'] !== null) {
                i18n.global.locale.value = event['lang'];
            }
        },
        handleClick() {
            const newTheme = this.themeStore.theme === 'light' ? 'dark' : 'light';
            this.themeStore.theme = newTheme;
        },
    },
};
</script>

<template>
    <section class="container-xxl header">
        <nav class="navbar navbar-expand-md custom-navbar">
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
                    <li class="nav-item">
                        <router-link class="nav-link menu-item" to="/">
                            {{ $t('header.main') }}
                        </router-link>
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

                <ul class="nav navbar-nav navbar-right d-flex align-items-center">
                    <li class="me-3">
                        <button
                            type="button"
                            @click="handleClick"
                            title="Сменить тему оформления"
                            class="theme-switcher"
                            :class="'theme-switcher_theme_' + this.themeStore.theme"
                        >
                            <DayIcon
                                class="theme-switcher__icon theme-switcher__icon_type_light"
                            />
                            <NightIcon
                                class="theme-switcher__icon theme-switcher__icon_type_dark"
                            />
                        </button>
                    </li>

                    <li class="me-3">
                        <v-select
                            :options="$i18n.availableLocales"
                            :clearable="false"
                            v-model="$i18n.locale"
                            inputId="value"
                        >
                            {{ locale }}
                        </v-select>
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
.header {
    background-color: var(--color-surface);
}

.navbar {
    --bs-navbar-padding-y: 0;
}

.nav-item {
    font-family: 'Montserrat', 'Verdana regular', 'Ebrima bold';
    font-size: 14pt;
    font-variant-caps: all-petite-caps;
}

.nav-link {
    padding-left: 1rem;
}

.menu-item {
    --c: var(--color-on-surface);
    --m: var(--color-header-menu-item);
    --h: 1.85em;

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
    background-color: var(--color-header-menu-item);
}

.dropdown-menu > li:hover:active {
    --bs-dropdown-link-active-bg: #4b9e90;
}

.v-select {
    width: 6rem;
    font-size: 1rem;
}

.vs__dropdown-menu {
    min-width: 0;
}

.header-locale {
    border-width: 1px;
    border-style: solid;
    border-radius: 4px;
    border-color: var(--vs-colors--lightest);
    background-color: var(--color-surface-secondary);
    color: var(--color-on-surface);
}

.custom-navbar .fa-brands,
.fa {
    font-size: 30px;
    color: var(--color-header-menu-item);
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

.theme-switcher {
    margin: 0;
    padding: 0;
    margin-top: -5px;
    border: none;
    background-color: transparent;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 0.15s ease-in-out;
    display: block;
    width: 32px;
    aspect-ratio: 1;
    overflow: hidden;
    border-radius: 50%;
    background-color: var(--color-surface-secondary-solid);
    color: inherit;
    position: relative;
}

.theme-switcher:hover {
    opacity: 1;
}

.theme-switcher__icon {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 50%;
    height: 50%;
    transition: transform 0.5s ease-out;
    transform-origin: 50% 200%;
}

.theme-switcher__icon_type_light {
    width: 70%;
    height: 70%;
}

.theme-switcher_theme_light .theme-switcher__icon_type_light {
    transform: translate(-50%, -50%);
}

.theme-switcher_theme_light .theme-switcher__icon_type_dark {
    transform: translate(-50%, -50%) rotate(180deg);
}

.theme-switcher_theme_dark .theme-switcher__icon_type_light {
    transform: translate(-50%, -50%) rotate(180deg);
}

.theme-switcher_theme_dark .theme-switcher__icon_type_dark {
    transform: translate(-50%, -50%) rotate(360deg);
}
</style>
