import { defineStore, } from 'pinia';
import { ref, watch, } from 'vue';

export const useThemeStore = defineStore('themeStore', () => {
    const theme = ref('light',);
    const DARK_CLASS_NAME = 'body_theme_dark';
    const HTML_DARK_CLASS = 'dark';

    // Client-side only initialization
    if (typeof window !== 'undefined') {
        const themeLocalStorage = localStorage.getItem('theme',);
        if (themeLocalStorage) {
            theme.value = JSON.parse(themeLocalStorage,);
            if (theme.value === 'dark') {
                document.body?.classList.add(DARK_CLASS_NAME,);
                document.documentElement.classList.add(HTML_DARK_CLASS,);
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

    watch(theme, (theme,) => {
        if (typeof window !== 'undefined') {
            localStorage.setItem('theme', JSON.stringify(theme,),);
            const body = document.querySelector('body',);
            const html = document.documentElement;
            if (theme === 'dark') {
                body?.classList.add(DARK_CLASS_NAME,);
                html.classList.add(HTML_DARK_CLASS,);
            } else {
                body?.classList.remove(DARK_CLASS_NAME,);
                html.classList.remove(HTML_DARK_CLASS,);
            }
        }
    }, { immediate: true, },);

    return {
        theme,
    };
},);
