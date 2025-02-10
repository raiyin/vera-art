import { createRouter, createWebHistory } from 'vue-router';

// import Home from '@/views/Home.vue';
import Shop from '@/views/Shop.vue';
import AllWorks from '@/views/AllWorks.vue';
import News from '@/views/News.vue';
import PayDelivery from '@/views/PayDeliver.vue';
import Services from '@/views/Services.vue';
import NewsItem from '@/components/News/NewsItem.vue';
import NotFound from '@/views/NotFound.vue';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('@/views/Home.vue'),
        },
        {
            path: '/shop',
            name: 'shop',
            // Без этого не подгружаются стили селекта
            component: Shop,
        },
        {
            path: '/allworks',
            name: 'allworks',
            component: () => import('@/views/AllWorks.vue'),
        },
        {
            path: '/news',
            name: 'news',
            component: () => import('@/views/News.vue')
        },
        {
            path: '/paydelivery',
            name: 'paydelivery',
            component: () => import('@/views/PayDeliver.vue'),
        },
        {
            path: '/services',
            name: 'services',
            component: () => import('@/views/Services.vue'),
        },
        {
            path: '/news/:id',
            name: 'newsitem',
            component: () => import('@/components/News/NewsItem.vue'),
        },
        {
            path: '/:pathMatch(.*)*',
            name: 'notfound',
            component: () => import('@/views/NotFound.vue'),
        },
    ],
});

export default router;
