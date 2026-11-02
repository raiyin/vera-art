import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig, } from 'axios';
import { useAuthStore, } from '~/stores/AuthStore';

/**
 * Centralized HTTP client for the application.
 *
 * Features:
 * - Singleton Axios instance (lazy initialization to avoid module-level useRuntimeConfig)
 * - Request interceptor: injects Authorization header from AuthStore
 * - Response interceptor: on 401, queues concurrent requests, performs single token refresh, retries
 * - Unified error handling
 */

interface TokenRefreshResponse {
    access_token: string
    refresh_token: string
    access_expires?: string
}

let isRefreshing = false;
let failedQueue: Array<{
    resolve: (token: string,) => void
    reject: (error: unknown,) => void
}> = [];

function processQueue(error: unknown, token: string | null = null,): void {
    failedQueue.forEach((promise,) => {
        if (error) {
            promise.reject(error,);
        } else {
            promise.resolve(token!,);
        }
    },);
    failedQueue = [];
}

function createHttpClient(): AxiosInstance {
    const config = useRuntimeConfig();
    const baseURL = config.public.serverUrl as string;

    const client = axios.create({
        baseURL,
        timeout: 15000,
    },);

    // Request interceptor: attach auth token (skip for auth endpoints)
    client.interceptors.request.use(
        (config: InternalAxiosRequestConfig,) => {
            const url = config.url || '';

            // Skip auth header for login, register, refresh endpoints
            if (url.includes('/login',) || url.includes('/register',) || url.includes('/refresh',)) {
                return config;
            }

            const authStore = useAuthStore();
            const token = authStore.accessToken;

            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }

            return config;
        },
        error => Promise.reject(error,),
    );

    // Response interceptor: handle 401 with token refresh
    client.interceptors.response.use(
        response => response,
        async (error: AxiosError,) => {
            const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

            // Only attempt refresh on 401, and only once per request
            if (error.response?.status !== 401 || originalRequest._retry) {
                return Promise.reject(error,);
            }

            // If already refreshing, queue this request
            if (isRefreshing) {
                return new Promise<string>((resolve, reject,) => {
                    failedQueue.push({ resolve, reject, },);
                },).then((token,) => {
                    originalRequest.headers.Authorization = `Bearer ${token}`;
                    return client(originalRequest,);
                },);
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                const authStore = useAuthStore();
                const refreshTokenValue = authStore.refreshToken;

                if (!refreshTokenValue) {
                    throw new Error('No refresh token available',);
                }

                const response = await axios.post<TokenRefreshResponse>(
                    `${baseURL}refresh`,
                    { refresh_token: refreshTokenValue, },
                );

                const { access_token, refresh_token, access_expires, } = response.data;

                // Update both tokens in the store (server rotates refresh token too)
                if (refresh_token) {
                    authStore.updateTokens(access_token, refresh_token, access_expires,);
                } else {
                    authStore.updateAccessToken(access_token, access_expires,);
                }

                // Process queued requests with the new token
                processQueue(null, access_token,);

                // Retry the original request
                originalRequest.headers.Authorization = `Bearer ${access_token}`;
                return client(originalRequest,);
            } catch (refreshError) {
                // Refresh failed — clear tokens and redirect to login
                processQueue(refreshError, null,);

                try {
                    const authStore = useAuthStore();
                    authStore.clearTokens();
                } catch {
                    // AuthStore might not be available
                }

                if (typeof window !== 'undefined' && !window.location.pathname.includes('/auth/login',)) {
                    window.location.href = '/auth/login';
                }

                return Promise.reject(refreshError,);
            } finally {
                isRefreshing = false;
            }
        },
    );

    return client;
}

// Singleton — lazy initialization
let httpClientInstance: AxiosInstance | null = null;

/**
 * Get the singleton HTTP client instance.
 * Must be called within Nuxt context (after useRuntimeConfig is available).
 */
export function getHttpClient(): AxiosInstance {
    if (!httpClientInstance) {
        httpClientInstance = createHttpClient();
    }
    return httpClientInstance;
}

/**
 * Reset the HTTP client singleton (useful for testing or after logout).
 */
export function resetHttpClient(): void {
    httpClientInstance = null;
}
