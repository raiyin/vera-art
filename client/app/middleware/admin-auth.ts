import { useAuthStore, } from '../stores/AuthStore';

export default defineNuxtRouteMiddleware(async (_to, _from,) => {
    // Skip middleware during SSR — auth depends on localStorage which is client-only
    if (import.meta.server) {
        return;
    }

    const authStore = useAuthStore();

    // Initialize auth state from localStorage before checking
    // This is needed because middleware runs before app.vue's onMounted
    authStore.initFromLocalStorage();

    const { isAuthenticated, isAdmin, } = authStore;

    // If not authenticated, redirect to login
    if (!isAuthenticated) {
        return navigateTo('/auth/login',);
    }

    // If authenticated but not admin, redirect to home (or show forbidden)
    if (!isAdmin) {
        // Optionally show a toast or message
        console.warn('Access denied: user is not an admin',);
        return navigateTo('/',);
    }

    // Allow access if both authenticated and admin
},);
