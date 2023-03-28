import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';

const app = createApp(App);
app.provide('jsonserverhost', 'http://localhost:3001/');
app.provide('imagebasedir', '/src/assets');
//app.provide('jsonserverhost', 'https://www.pertsukova.ru/server/');
//app.provide('imagebasedir', '/assets');
app.use(createPinia());
app.use(router);
app.mount('#app');
