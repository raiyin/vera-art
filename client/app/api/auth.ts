import axios from 'axios';
import { useAuthStore, } from '../stores/AuthStore';

// Helper to get runtime config — must be called inside Nuxt context
function getApiUrl(): string {
    const config = useRuntimeConfig();
    return config.public.serverUrl;
}

// Create axios instance (baseURL will be set in request interceptor to avoid module-level useRuntimeConfig)
const api = axios.create({
    headers: {
        'Content-Type': 'application/json',
    },
},);

// Request interceptor to add auth token and set base URL
api.interceptors.request.use(
    (config,) => {
        // Set baseURL at request time (inside Nuxt context)
        config.baseURL = getApiUrl();

        const authStore = useAuthStore();
        const token = authStore.accessToken;

        if (token && !config.url?.includes('/refresh',) && !config.url?.includes('/login',) && !config.url?.includes('/register',)) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
    },
    (error,) => {
        return Promise.reject(error,);
    },
);

// Response interceptor to handle token refresh
api.interceptors.response.use(
    response => response,
    async (error,) => {
        const originalRequest = error.config;

        // If error is 401 and we haven't tried refreshing yet
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                const authStore = useAuthStore();

                // Try to refresh the token
                const refreshResponse = await axios.post(`${getApiUrl()}refresh`, {
                    refresh_token: authStore.refreshToken,
                },);

                if (refreshResponse.status === 200) {
                    const { access_token, access_expires, } = refreshResponse.data;

                    // Update the access token in the store
                    authStore.updateAccessToken(access_token, access_expires,);

                    // Update the Authorization header
                    originalRequest.headers.Authorization = `Bearer ${access_token}`;

                    // Retry the original request
                    return api(originalRequest,);
                }
            } catch {
                // Refresh failed, logout the user
                const authStore = useAuthStore();
                authStore.clearTokens();

                // Redirect to login page if we're not already there
                if (typeof window !== 'undefined' && !window.location.pathname.includes('/auth/login',)) {
                    window.location.href = '/auth/login';
                }
            }
        }

        return Promise.reject(error,);
    },
);

export default {
    // Register new user
    async register(userData: { username: string, password: string, email?: string },) {
        try {
            const response = await api.post('register', userData,);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Registration failed', };
        }
    },

    // Login user
    async login(credentials: { username: string, password: string },) {
        try {
            const response = await api.post('login', credentials,);
            const authStore = useAuthStore();

            // Save tokens to store
            authStore.saveTokens(response.data,);

            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Login failed', };
        }
    },

    // Refresh access token
    async refreshToken() {
        try {
            const authStore = useAuthStore();
            const refreshToken = authStore.refreshToken;

            if (!refreshToken) {
                throw new Error('No refresh token available',);
            }

            const response = await axios.post(`${getApiUrl()}refresh`, {
                refresh_token: refreshToken,
            },);

            if (response.status === 200) {
                const { access_token, access_expires, } = response.data;
                authStore.updateAccessToken(access_token, access_expires,);
                return { access_token, access_expires, };
            }
        } catch (error: any) {
            // If refresh fails, clear tokens
            const authStore = useAuthStore();
            authStore.clearTokens();
            throw error.response?.data || { error: 'Token refresh failed', };
        }
    },

    // Logout user
    logout() {
        const authStore = useAuthStore();
        authStore.clearTokens();

        // Redirect to login page
        if (typeof window !== 'undefined') {
            window.location.href = '/auth/login';
        }
    },

    // Get protected content (example)
    async getProtectedContent() {
        try {
            const response = await api.get('protected',);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Failed to fetch protected content', };
        }
    },

    // Check if user is authenticated
    isAuthenticated() {
        const authStore = useAuthStore();
        return authStore.isAuthenticated && !authStore.isRefreshTokenExpired;
    },

    // Get current auth headers
    getAuthHeaders() {
        const authStore = useAuthStore();
        return authStore.getAuthHeader;
    },

    // Get user profile
    async getProfile() {
    	try {
    		const response = await api.get('profile');
    		return response.data;
    	} catch (error: any) {
    		throw error.response?.data || { error: 'Failed to fetch profile' };
    	}
    },

    // Update user profile
    async updateProfile(data: { email?: string; full_name?: string }) {
    	try {
    		const response = await api.put('profile', data);
    		return response.data;
    	} catch (error: any) {
    		throw error.response?.data || { error: 'Failed to update profile' };
    	}
    },

    // Upload avatar
    async uploadAvatar(file: File) {
    	try {
    		const formData = new FormData();
    		formData.append('avatar', file);
    		const response = await api.post('profile/avatar', formData, {
    			headers: {
    				'Content-Type': 'multipart/form-data',
    			},
    		});
    		return response.data;
    	} catch (error: any) {
    		throw error.response?.data || { error: 'Failed to upload avatar' };
    	}
    },

    // Delete avatar
    async deleteAvatar() {
    	try {
    		const response = await api.delete('profile/avatar');
    		return response.data;
    	} catch (error: any) {
    		throw error.response?.data || { error: 'Failed to delete avatar' };
    	}
    },

    // Get axios instance for custom requests
    getApiInstance() {
    	return api;
    },
   };
