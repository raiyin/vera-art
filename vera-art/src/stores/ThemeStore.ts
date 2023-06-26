import { defineStore } from 'pinia';

export const useThemeStore = defineStore('themeStore', {
    state: () => ({
        theme: 'light',
    }),
    actions: {
        setTheme(newTheme: string) {
            this.theme = newTheme;
        }
    }
});

// export default {useThemeStore}
