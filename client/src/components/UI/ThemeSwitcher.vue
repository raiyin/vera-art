<script lang="ts">
import DayIcon from '@/components/Icons/icon_day.vue';
import NightIcon from '@/components/Icons/icon_night.vue';
import { useThemeStore } from '../../stores/ThemeStore';

export default {
    setup() {
        const themeStore = useThemeStore();
        return { themeStore };
    },
    components: {
        DayIcon,
        NightIcon,
    },
    methods: {
        handleClick() {
            const newTheme =
                this.themeStore.theme === 'light' ? 'dark' : 'light';
            this.themeStore.theme = newTheme;
        },
    },
};
</script>

<template>
    <button
        type="button"
        @click="handleClick"
        :title="$t('theme.title')"
        class="theme-switcher"
        :class="'theme-switcher_theme_' + themeStore.theme"
    >
        <DayIcon class="theme-switcher__icon theme-switcher__icon_type_light" />
        <NightIcon
            class="theme-switcher__icon theme-switcher__icon_type_dark"
        />
    </button>
</template>

<style scoped>
.theme-switcher {
    margin: 0;
    padding: 0;
    margin-top: -0.5rem;
    border: none;
    background-color: transparent;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 0.15s ease-in-out;
    width: 2.2rem;
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
