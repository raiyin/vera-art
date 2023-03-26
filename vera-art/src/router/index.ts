import { createRouter, createWebHistory } from 'vue-router';

import Home from '@/views/Home.vue';
import Shop from '@/views/Shop.vue';
import AllWorks from '@/views/AllWorks.vue';
import News from '@/views/News.vue';
import PayDelivery from '@/views/PayDeliver.vue';
import Services from '@/views/Services.vue';
import NewsItem from '@/components/News/NewsItem.vue';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: Home,
        },
        {
            path: '/shop',
            name: 'shop',
            component: Shop,
        },
        {
            path: '/allworks',
            name: 'allworks',
            component: AllWorks,
        },
        {
            path: '/news',
            name: 'news',
            component: News,
        },
        {
            path: '/paydelivery',
            name: 'paydelivery',
            component: PayDelivery,
        },
        {
            path: '/services',
            name: 'services',
            component: Services,
        },
        {
            path: '/newsitem/:id',
            name: 'newsitem',
            component: NewsItem,
        },
    ],
});

export default router;
