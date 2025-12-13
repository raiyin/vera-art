import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const useAuthStore = defineStore('authStore', () => {
    const theme = ref('light');
    const isAuthenticated = ref(false);
    const DARK_CLASS_NAME = 'body_theme_dark';

    const themeLocalStorage = localStorage.getItem('theme');
    const token = localStorage.getItem('token');
    isAuthenticated.value = !!token;

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

    const setAuthenticated = (auth: boolean) => {
        isAuthenticated.value = auth;
    };

    return {
        theme,
        isAuthenticated,
        setAuthenticated,
    };
});
