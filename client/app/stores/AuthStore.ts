import { defineStore, } from 'pinia';
import { ref, watch, computed, } from 'vue';
import { getRoleFromToken, getUserIdFromToken, } from '../utils/jwt';

interface TokenData {
    access_token: string
    refresh_token: string
    access_expires: string
    refresh_expires: string
    token_type: string
}

export const useAuthStore = defineStore('authStore', () => {
    const theme = ref('light',);
    const isAuthenticated = ref(false,);
    const accessToken = ref<string | null>(null,);
    const refreshToken = ref<string | null>(null,);
    const accessTokenExpiry = ref<Date | null>(null,);
    const refreshTokenExpiry = ref<Date | null>(null,);
    const userRole = ref<string | null>(null,);
    const userId = ref<number | null>(null,);
    const DARK_CLASS_NAME = 'body_theme_dark';

    // Client-side only initialization
    if (typeof window !== 'undefined') {
        const themeLocalStorage = localStorage.getItem('theme',);

        if (themeLocalStorage) {
            theme.value = JSON.parse(themeLocalStorage,);
            if (theme.value === 'dark') {
                document.body?.classList.add(DARK_CLASS_NAME,);
            }
        }
    }

    // Watch for theme changes
    watch(theme, (theme,) => {
        if (typeof window !== 'undefined') {
            localStorage.setItem('theme', JSON.stringify(theme,),);
            const body = document.querySelector('body',);
            if (theme === 'dark') {
                body?.classList.add(DARK_CLASS_NAME,);
            } else {
                body?.classList.remove(DARK_CLASS_NAME,);
            }
        }
    },);

    // Save tokens to localStorage
    const saveTokens = (tokenData: TokenData,) => {
        if (typeof window === 'undefined') return;

        accessToken.value = tokenData.access_token;
        refreshToken.value = tokenData.refresh_token;
        accessTokenExpiry.value = new Date(tokenData.access_expires,);
        refreshTokenExpiry.value = new Date(tokenData.refresh_expires,);

        // Extract role and user ID from access token
        const role = getRoleFromToken(tokenData.access_token,);
        const id = getUserIdFromToken(tokenData.access_token,);
        if (role) userRole.value = role;
        if (id) userId.value = id;

        localStorage.setItem('access_token', tokenData.access_token,);
        localStorage.setItem('refresh_token', tokenData.refresh_token,);
        localStorage.setItem('access_expires', tokenData.access_expires,);
        localStorage.setItem('refresh_expires', tokenData.refresh_expires,);
        // Also store as 'token' for backward compatibility with existing code
        localStorage.setItem('token', tokenData.access_token,);

        isAuthenticated.value = true;
    };

    // Clear all tokens (logout)
    const clearTokens = () => {
        if (typeof window === 'undefined') return;

        accessToken.value = null;
        refreshToken.value = null;
        accessTokenExpiry.value = null;
        refreshTokenExpiry.value = null;
        userRole.value = null;
        userId.value = null;

        localStorage.removeItem('access_token',);
        localStorage.removeItem('refresh_token',);
        localStorage.removeItem('access_expires',);
        localStorage.removeItem('refresh_expires',);
        localStorage.removeItem('token',); // Remove old token for backward compatibility

        isAuthenticated.value = false;
    };

    // Update only access token (after refresh)
    const updateAccessToken = (newAccessToken: string, newExpiry: string,) => {
        if (typeof window === 'undefined') return;

        accessToken.value = newAccessToken;
        accessTokenExpiry.value = new Date(newExpiry,);

        // Extract role and user ID from new access token
        const role = getRoleFromToken(newAccessToken,);
        const id = getUserIdFromToken(newAccessToken,);
        if (role) userRole.value = role;
        if (id) userId.value = id;

        localStorage.setItem('access_token', newAccessToken,);
        localStorage.setItem('access_expires', newExpiry,);
        // Also update 'token' for backward compatibility
        localStorage.setItem('token', newAccessToken,);
    };

    // Check if access token is expired
    const isAccessTokenExpired = computed(() => {
        if (!accessTokenExpiry.value) return true;
        return accessTokenExpiry.value <= new Date();
    },);

    // Check if refresh token is expired
    const isRefreshTokenExpired = computed(() => {
        if (!refreshTokenExpiry.value) return true;
        return refreshTokenExpiry.value <= new Date();
    },);

    // Get authorization header for API requests
    const getAuthHeader = computed(() => {
        if (!accessToken.value) return {};
        return { Authorization: `Bearer ${accessToken.value}`, };
    },);

    // For backward compatibility - get token (access token)
    const token = computed(() => accessToken.value,);

    // Role checking computed properties
    const isAdmin = computed(() => userRole.value === 'admin',);
    const isUser = computed(() => userRole.value === 'user',);
    const currentRole = computed(() => userRole.value,);
    const currentUserId = computed(() => userId.value,);

    const setAuthenticated = (auth: boolean,) => {
        isAuthenticated.value = auth;
    };

    // Initialize auth state from localStorage (call on client mount)
    const initFromLocalStorage = () => {
        if (typeof window === 'undefined') return;

        const storedAccessToken = localStorage.getItem('access_token');
        const storedRefreshToken = localStorage.getItem('refresh_token');
        const storedAccessExpiry = localStorage.getItem('access_expires');
        const storedRefreshExpiry = localStorage.getItem('refresh_expires');

        if (storedAccessToken && storedRefreshToken) {
            accessToken.value = storedAccessToken;
            refreshToken.value = storedRefreshToken;

            if (storedAccessExpiry) {
                accessTokenExpiry.value = new Date(storedAccessExpiry);
            }

            if (storedRefreshExpiry) {
                refreshTokenExpiry.value = new Date(storedRefreshExpiry);
            }

            const role = getRoleFromToken(storedAccessToken);
            const id = getUserIdFromToken(storedAccessToken);
            if (role) userRole.value = role;
            if (id) userId.value = id;

            const isAccessValid = accessTokenExpiry.value && accessTokenExpiry.value > new Date();
            const isRefreshValid = refreshTokenExpiry.value && refreshTokenExpiry.value > new Date();

            isAuthenticated.value = !!isRefreshValid;

            if (!isAccessValid && isRefreshValid) {
                console.log('Access token expired, refresh token still valid');
            }
        }
    };

    return {
        theme,
        isAuthenticated,
        accessToken,
        refreshToken,
        accessTokenExpiry,
        refreshTokenExpiry,
        userRole,
        userId,
        isAccessTokenExpired,
        isRefreshTokenExpired,
        getAuthHeader,
        token,
        isAdmin,
        isUser,
        currentRole,
        currentUserId,
        saveTokens,
        clearTokens,
        updateAccessToken,
        setAuthenticated,
        initFromLocalStorage,
    };
},);
