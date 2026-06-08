import { defineStore, } from 'pinia';
import { ref, watch, } from 'vue';

export const useThemeStore = defineStore('themeStore', () => {
    const theme = ref('light',);
    const DARK_CLASS_NAME = 'body_theme_dark';

    // Client-side only initialization
    if (typeof window !== 'undefined') {
        const themeLocalStorage = localStorage.getItem('theme',);
        if (themeLocalStorage) {
            theme.value = JSON.parse(themeLocalStorage,);
            if (theme.value === 'dark') {
                document.body?.classList.add(DARK_CLASS_NAME,);
            }
        } else {
            // No saved theme, default to light (ignore OS preference)
            theme.value = 'light';
        }

        // Listen for cross-tab synchronization
        const handleStorageChange = (event: StorageEvent,) => {
            if (event.key === 'theme' && event.newValue !== null) {
                try {
                    const newTheme = JSON.parse(event.newValue,);
                    if (newTheme === 'light' || newTheme === 'dark') {
                        // Avoid infinite loop: only update if different
                        if (theme.value !== newTheme) {
                            theme.value = newTheme;
                            // Classes will be updated by the watcher
                        }
                    }
                } catch (e) {
                    console.error('Failed to parse theme from storage', e,);
                }
            }
        };
        window.addEventListener('storage', handleStorageChange,);
    }

    watch(theme, (newTheme,) => {
        if (typeof window !== 'undefined') {
            localStorage.setItem('theme', JSON.stringify(newTheme,),);
            const body = document.querySelector('body',);
            if (newTheme === 'dark') {
                body?.classList.add(DARK_CLASS_NAME,);
            } else {
                body?.classList.remove(DARK_CLASS_NAME,);
            }
        }
    }, { immediate: true, },);

    return {
        theme,
    };
},);
