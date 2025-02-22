import { createApp } from 'vue';
import { useI18n } from 'vue-i18n';
import i18n from './i18n';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';

const app = createApp(App);

app.provide('server', 'http://localhost:8000/');
app.provide('imagebasedir', '/src/assets');
// app.provide('server', 'https://www.pertsukova.ru/server/');
// app.provide('imagebasedir', '/assets');

app.use(createPinia());
app.use(router);
app.use(i18n);
app.mount('#app');
