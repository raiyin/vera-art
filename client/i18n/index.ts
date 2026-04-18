import { createI18n } from 'vue-i18n';
import ru from './locales/ru.json';
import en from './locales/en.json';

const messages = {
    ru,
    en,
};

const i18n = createI18n({
    locale: import.meta.env.VITE_DEFAULT_LOCALE || 'ru',
    fallbackLocale: import.meta.env.VITE_FALLBACK_LOCALE || 'ru',
    legacy: false,
    messages,
});

export default i18n;
