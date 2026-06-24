import { useAuthStore, } from '~/stores/AuthStore';

export default defineNuxtRouteMiddleware(async (_to, _from,) => {
    const authStore = useAuthStore();

    if (import.meta.server) {
        // On the server, use cookie-based auth (SSR-safe)
        // This prevents hydration mismatch: if the user has a valid auth cookie,
        // the server renders the admin layout, matching what the client will show.
        const accessTokenCookie = useCookie<string | null>('access_token',).value ?? null;
        const refreshTokenCookie = useCookie<string | null>('refresh_token',).value ?? null;
        authStore.initFromCookie(accessTokenCookie, refreshTokenCookie,);
    } else {
        // On the client, initialize auth state from localStorage before checking
        // This is needed because middleware runs before app.vue's onMounted
        authStore.initFromLocalStorage();

        // Fall back to cookies if localStorage is empty (e.g. first SSR visit).
        // This ensures the client makes the same auth decision as the server,
        // preventing hydration mismatch on page reload.
        if (!authStore.isAuthenticated) {
            const accessTokenCookie = useCookie<string | null>('access_token',).value ?? null;
            const refreshTokenCookie = useCookie<string | null>('refresh_token',).value ?? null;
            authStore.initFromCookie(accessTokenCookie, refreshTokenCookie,);
        }
    }

    const isAuthenticated = authStore.isAuthenticated;
    const isAdmin = authStore.isAdmin;

    // If not authenticated, redirect to login
    if (!isAuthenticated) {
        return navigateTo('/auth/login',);
    }

    // If authenticated but not admin, redirect to home (or show forbidden)
    if (!isAdmin) {
        console.warn('Access denied: user is not an admin',);
        return navigateTo('/',);
    }

    // Allow access if both authenticated and admin
},);
