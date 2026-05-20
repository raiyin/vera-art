import { useAuthStore, } from '../stores/AuthStore';

export default defineNuxtRouteMiddleware((_to, _from,) => {
    const authStore = useAuthStore();
    const { isAuthenticated, isAdmin, } = authStore;

    // If not authenticated, redirect to login
    if (!isAuthenticated) {
        return navigateTo('/login',);
    }

    // If authenticated but not admin, redirect to home (or show forbidden)
    if (!isAdmin) {
    // Optionally show a toast or message
        console.warn('Access denied: user is not an admin',);
        return navigateTo('/',);
    }

    // Allow access if both authenticated and admin
},);
