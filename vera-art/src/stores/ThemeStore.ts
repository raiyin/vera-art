import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const useThemeStore = defineStore('themeStore', () => {
    const theme = ref('light');
    const DARK_CLASS_NAME = 'body_theme_dark';

    const themeLocalStorage = localStorage.getItem('theme');
    if (themeLocalStorage) {
        theme.value = JSON.parse(themeLocalStorage);
        if (theme.value === 'dark') {
            document.body?.classList.add(DARK_CLASS_NAME);
        }
    }

    watch(theme, (theme) => {
        localStorage.setItem('theme', JSON.stringify(theme));
        const body = document.querySelector('body');
        if (theme === 'dark') {
            body?.classList.add(DARK_CLASS_NAME);
        } else {
            body?.classList.remove(DARK_CLASS_NAME);
        }
    });

    return {
        theme,
    };
});
