import { createI18n } from 'vue-i18n';
import RUS from './locales/RUS.json';
import ENG from './locales/ENG.json';

const messages = {
    RUS,
    ENG,
};

const i18n = createI18n({
    locale: import.meta.env.VITE_DEFAULT_LOCALE,
    fallbackLocale: import.meta.env.VITE_FALLBACK_LOCALE,
    legacy: false,
    messages,
});

export default i18n;
