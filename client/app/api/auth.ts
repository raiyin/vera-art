import { getHttpClient, } from '~/api/http-client';
import { useAuthStore, } from '~/stores/AuthStore';

export default {
    // Register new user
    async register(userData: { username: string, password: string, email?: string },) {
        try {
            const response = await getHttpClient().post('register', userData,);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Registration failed', };
        }
    },

    // Verify email with token
    async verifyEmail(token: string,) {
        try {
            const response = await getHttpClient().get('verify-email', { params: { token, }, },);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Email verification failed', };
        }
    },

    // Resend verification email
    async resendVerification(email: string,) {
        try {
            const response = await getHttpClient().post('resend-verification', { email, },);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Failed to resend verification email', };
        }
    },

    // Login user
    async login(credentials: { username: string, password: string },) {
        try {
            const response = await getHttpClient().post('login', credentials,);
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

            const response = await getHttpClient().post('refresh', {
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
            const response = await getHttpClient().get('protected',);
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
            const response = await getHttpClient().get('profile',);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Failed to fetch profile', };
        }
    },

    // Update user profile
    async updateProfile(data: { email?: string, full_name?: string },) {
        try {
            const response = await getHttpClient().put('profile', { name: data.full_name, },);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Failed to update profile', };
        }
    },

    // Upload avatar
    async uploadAvatar(formData: FormData,) {
        try {
            const response = await getHttpClient().post('profile/avatar', formData, {
                headers: {
                    'Content-Type': null,
                },
            },);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Failed to upload avatar', };
        }
    },

    // Delete avatar
    async deleteAvatar() {
        try {
            const response = await getHttpClient().delete('profile/avatar',);
            return response.data;
        } catch (error: any) {
            throw error.response?.data || { error: 'Failed to delete avatar', };
        }
    },
};
