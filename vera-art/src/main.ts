import { createApp } from 'vue';
import { createI18n, useI18n } from 'vue-i18n'
import { languages } from './i18n';
import { defaultLocale } from './i18n';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';

const messages = Object.assign(languages)
const i18n = createI18n({
    legacy: false,
    locale: defaultLocale,
    fallbackFormat: 'en',
    messages
})
const app = createApp(App, {
    setup() {
        const { t } = useI18n();
        return { t };
    }
});

app.provide('jsonserverhost', 'http://localhost:3001/');
app.provide('imagebasedir', '/src/assets');
// app.provide('jsonserverhost', 'https://www.pertsukova.ru/server/');
// app.provide('imagebasedir', '/assets');

app.use(createPinia());
app.use(router);
app.use(i18n);
app.mount('#app');
