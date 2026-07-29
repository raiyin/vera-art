import { defineStore, } from 'pinia';
import { ref, computed, } from 'vue';
import { getRoleFromToken, getUserIdFromToken, getTokenExpiry, } from '~/utils/jwt';

interface TokenData {
    access_token: string
    refresh_token: string
    access_expires?: string
    refresh_expires?: string
    token_type?: string
}

export const useAuthStore = defineStore('authStore', () => {
    const isAuthenticated = ref(false,);
    const accessToken = ref<string | null>(null,);
    const refreshToken = ref<string | null>(null,);
    const accessTokenExpiry = ref<Date | null>(null,);
    const refreshTokenExpiry = ref<Date | null>(null,);
    const userRole = ref<string | null>(null,);
    const userId = ref<number | null>(null,);

    // Save tokens to localStorage and cookie (cookie for SSR access)
    const saveTokens = (tokenData: TokenData,) => {
        if (typeof window === 'undefined') return;

        accessToken.value = tokenData.access_token;
        refreshToken.value = tokenData.refresh_token;

        // Derive expiry from JWT token claims if server doesn't send explicit expiry fields
        accessTokenExpiry.value = tokenData.access_expires
            ? new Date(tokenData.access_expires,)
            : getTokenExpiry(tokenData.access_token,);
        refreshTokenExpiry.value = tokenData.refresh_expires
            ? new Date(tokenData.refresh_expires,)
            : getTokenExpiry(tokenData.refresh_token,);

        // Extract role and user ID from access token
        const role = getRoleFromToken(tokenData.access_token,);
        const id = getUserIdFromToken(tokenData.access_token,);
        if (role) userRole.value = role;
        if (id) userId.value = id;

        localStorage.setItem('access_token', tokenData.access_token,);
        localStorage.setItem('refresh_token', tokenData.refresh_token,);
        if (tokenData.access_expires) {
            localStorage.setItem('access_expires', tokenData.access_expires,);
        }
        if (tokenData.refresh_expires) {
            localStorage.setItem('refresh_expires', tokenData.refresh_expires,);
        }
        // Also store as 'token' for backward compatibility with existing code
        localStorage.setItem('token', tokenData.access_token,);

        // Save to cookie for SSR access
        // Using document.cookie directly since useCookie() can only be called
        // within Nuxt setup context (middleware, plugin, component setup)
        setAuthCookieClient(tokenData.access_token, tokenData.refresh_token,);

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

        // Clear auth cookie
        clearAuthCookieClient();

        isAuthenticated.value = false;
    };

    // Update only access token (after refresh)
    const updateAccessToken = (newAccessToken: string, newExpiry?: string,) => {
        if (typeof window === 'undefined') return;

        accessToken.value = newAccessToken;
        accessTokenExpiry.value = newExpiry
            ? new Date(newExpiry,)
            : getTokenExpiry(newAccessToken,);

        // Extract role and user ID from new access token
        const role = getRoleFromToken(newAccessToken,);
        const id = getUserIdFromToken(newAccessToken,);
        if (role) userRole.value = role;
        if (id) userId.value = id;

        localStorage.setItem('access_token', newAccessToken,);
        if (newExpiry) {
            localStorage.setItem('access_expires', newExpiry,);
        }
        // Also update 'token' for backward compatibility
        localStorage.setItem('token', newAccessToken,);

        // Update auth cookie
        if (refreshToken.value) {
            setAuthCookieClient(newAccessToken, refreshToken.value,);
        }
    };

    // Update both tokens (after refresh with rotation)
    const updateTokens = (newAccessToken: string, newRefreshToken: string, newExpiry?: string,) => {
        if (typeof window === 'undefined') return;

        accessToken.value = newAccessToken;
        refreshToken.value = newRefreshToken;
        accessTokenExpiry.value = newExpiry
            ? new Date(newExpiry,)
            : getTokenExpiry(newAccessToken,);
        refreshTokenExpiry.value = getTokenExpiry(newRefreshToken,);

        // Extract role and user ID from new access token
        const role = getRoleFromToken(newAccessToken,);
        const id = getUserIdFromToken(newAccessToken,);
        if (role) userRole.value = role;
        if (id) userId.value = id;

        localStorage.setItem('access_token', newAccessToken,);
        localStorage.setItem('refresh_token', newRefreshToken,);
        if (newExpiry) {
            localStorage.setItem('access_expires', newExpiry,);
        }
        localStorage.setItem('token', newAccessToken,);

        // Update auth cookie
        setAuthCookieClient(newAccessToken, newRefreshToken,);
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

        const storedAccessToken = localStorage.getItem('access_token',);
        const storedRefreshToken = localStorage.getItem('refresh_token',);
        const storedAccessExpiry = localStorage.getItem('access_expires',);
        const storedRefreshExpiry = localStorage.getItem('refresh_expires',);

        if (storedAccessToken && storedRefreshToken) {
            accessToken.value = storedAccessToken;
            refreshToken.value = storedRefreshToken;

            if (storedAccessExpiry) {
                accessTokenExpiry.value = new Date(storedAccessExpiry,);
            } else {
                accessTokenExpiry.value = getTokenExpiry(storedAccessToken,);
            }

            if (storedRefreshExpiry) {
                refreshTokenExpiry.value = new Date(storedRefreshExpiry,);
            } else {
                refreshTokenExpiry.value = getTokenExpiry(storedRefreshToken,);
            }

            const role = getRoleFromToken(storedAccessToken,);
            const id = getUserIdFromToken(storedAccessToken,);
            if (role) userRole.value = role;
            if (id) userId.value = id;

            const isAccessValid = accessTokenExpiry.value && accessTokenExpiry.value > new Date();
            const isRefreshValid = refreshTokenExpiry.value && refreshTokenExpiry.value > new Date();

            isAuthenticated.value = !!isRefreshValid;

            if (!isAccessValid && isRefreshValid) {
                console.log('Access token expired, refresh token still valid',);
            }
        }
    };

    // Initialize auth state from cookie values (SSR-safe, called from middleware)
    // The cookie values are passed in because useCookie() must be called within
    // a Nuxt context (middleware, plugin, component setup), not inside a Pinia store.
    const initFromCookie = (cookieAccessToken: string | null, cookieRefreshToken: string | null,) => {
        if (cookieAccessToken && cookieRefreshToken) {
            const refreshExpiry = getTokenExpiry(cookieRefreshToken,);
            if (!refreshExpiry || refreshExpiry <= new Date()) {
                return; // Refresh token expired
            }

            accessToken.value = cookieAccessToken;
            refreshToken.value = cookieRefreshToken;

            const accessExpiry = getTokenExpiry(cookieAccessToken,);
            if (accessExpiry) {
                accessTokenExpiry.value = accessExpiry;
            }
            if (refreshExpiry) {
                refreshTokenExpiry.value = refreshExpiry;
            }

            const role = getRoleFromToken(cookieAccessToken,);
            const id = getUserIdFromToken(cookieAccessToken,);
            if (role) userRole.value = role;
            if (id) userId.value = id;

            isAuthenticated.value = true;
        }
    };

    return {
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
        updateTokens,
        setAuthenticated,
        initFromLocalStorage,
        initFromCookie,
    };
},);

/**
 * Set auth cookies via document.cookie for SSR access.
 * Using document.cookie directly because useCookie() can only be called
 * within Nuxt setup context (middleware, plugin, component setup).
 */
function setAuthCookieClient(accessToken: string, refreshToken: string,) {
    const maxAge = 60 * 60 * 24 * 7; // 7 days
    document.cookie = `access_token=${encodeURIComponent(accessToken,)}; path=/; max-age=${maxAge}; sameSite=lax`;
    document.cookie = `refresh_token=${encodeURIComponent(refreshToken,)}; path=/; max-age=${maxAge}; sameSite=lax`;
}

/**
 * Clear auth cookies via document.cookie.
 */
function clearAuthCookieClient() {
    document.cookie = 'access_token=; path=/; max-age=0; sameSite=lax';
    document.cookie = 'refresh_token=; path=/; max-age=0; sameSite=lax';
}
