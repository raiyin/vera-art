import { defineStore } from 'pinia';

export const useThemeStore = defineStore('themeStore', {
    state: () => ({
        theme: 'dark',
    }),
    actions: {
        setTheme(newTheme: string) {
            this.theme = newTheme;
        }
    }
});

// export default {useThemeStore}
