<template>
    <div
        class="min-h-screen flex items-center justify-center bg-linear-to-br from-indigo-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 p-4"
    >
        <div class="w-full max-w-lg">
            <!-- Logo/Brand -->
            <div class="text-center mb-6">
                <div
                    class="inline-flex items-center justify-center w-14 h-14 bg-indigo-500 rounded-2xl shadow-lg mb-3"
                >
                    <svg
                        class="w-8 h-8 text-white"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                        ></path>
                    </svg>
                </div>
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                    {{ $t('auth.login') }}
                </h1>
                <p class="text-gray-600 dark:text-gray-300 mt-2">
                    {{
                        $t(
                            'auth.welcomeBack',
                            'Welcome back! Please enter your credentials'
                        )
                    }}
                </p>
            </div>

            <!-- Login Card -->
            <UCard class="shadow-2xl border-0">
                <UForm :state="formState" class="space-y-4" @submit="handleLogin">
                    <!-- Username Field -->
                    <UFormField name="username" :label="$t('auth.username')" required>
                        <UInput
                            v-model="formState.username"
                            type="text"
                            :placeholder="
                                $t('auth.usernamePlaceholder', 'Enter your username')
                            "
                            icon="i-heroicons-user"
                            size="lg"
                            :disabled="loading"
                            autocomplete="username"
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Password Field -->
                    <UFormField name="password" :label="$t('auth.password')" required>
                        <UInput
                            v-model="formState.password"
                            :type="passwordVisible ? 'text' : 'password'"
                            :placeholder="
                                $t('auth.passwordPlaceholder', 'Enter your password')
                            "
                            icon="i-heroicons-lock-closed"
                            size="lg"
                            :disabled="loading"
                            autocomplete="current-password"
                            class="w-full"
                        >
                            <template #trailing>
                                <UButton
                                    variant="ghost"
                                    color="neutral"
                                    :icon="
                                        passwordVisible
                                            ? 'i-heroicons-eye-slash'
                                            : 'i-heroicons-eye'
                                    "
                                    @click="passwordVisible = !passwordVisible"
                                    :padded="false"
                                    class="!p-1"
                                />
                            </template>
                        </UInput>
                    </UFormField>

                    <!-- Remember Me & Forgot Password -->
                    <div class="flex items-center justify-between">
                        <UCheckbox
                            v-model="rememberMe"
                            :label="$t('auth.remember')"
                            name="remember"
                            :disabled="loading"
                        />
                        <NuxtLink
                            to="/forgot-password"
                            class="text-sm text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300 font-medium transition-colors"
                        >
                            {{ $t('auth.forgot') }}
                        </NuxtLink>
                    </div>

                    <!-- Error Message -->
                    <UAlert
                        v-if="error"
                        :title="$t('auth.error', 'Error')"
                        :description="error"
                        icon="i-heroicons-exclamation-triangle"
                        color="error"
                        variant="outline"
                        class="mt-4"
                    />

                    <!-- Submit Button -->
                    <UButton
                        type="submit"
                        block
                        size="lg"
                        :loading="loading"
                        :disabled="loading"
                        class="mt-6"
                        color="primary"
                    >
                        <template #leading>
                            <svg
                                class="w-5 h-5"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                                ></path>
                            </svg>
                        </template>
                        {{ $t('auth.login') }}
                    </UButton>
                </UForm>

                <!-- Register Link -->
                <div class="mt-6 text-center">
                    <p class="text-gray-600 dark:text-gray-300">
                        {{ $t('auth.noAccount') }}
                        <NuxtLink
                            to="/auth/register"
                            class="font-semibold text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300 transition-colors"
                        >
                            {{ $t('auth.registerNow') }}
                        </NuxtLink>
                    </p>
                </div>
            </UCard>

            <!-- Back to Home -->
            <div class="mt-6 text-center">
                <NuxtLink
                    to="/"
                    class="inline-flex items-center text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors"
                >
                    <svg
                        class="w-4 h-4 mr-2"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M10 19l-7-7m0 0l7-7m-7 7h18"
                        ></path>
                    </svg>
                    {{ $t('auth.backToHome', 'Back to home') }}
                </NuxtLink>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import authApi from '../../api/auth';

interface LoginForm {
    username: string;
    password: string;
}

export default {
    setup() {
        const router = useRouter();
        const { t } = useI18n();

        const formState = reactive<LoginForm>({
            username: '',
            password: '',
        });

        const rememberMe = ref(false);
        const loading = ref(false);
        const error = ref('');
        const passwordVisible = ref(false);

        const handleLogin = async () => {
            if (!formState.username.trim() || !formState.password.trim()) {
                error.value = t('auth.fillAllFields', 'Please fill in all fields');
                return;
            }

            loading.value = true;
            error.value = '';

            try {
                // Use the new auth API which handles token storage automatically
                await authApi.login({
                    username: formState.username,
                    password: formState.password,
                });

                // Store username for "remember me" functionality
                if (rememberMe.value) {
                    if (typeof window !== 'undefined') {
                        localStorage.setItem('rememberedUser', formState.username);
                    }
                } else {
                    if (typeof window !== 'undefined') {
                        localStorage.removeItem('rememberedUser');
                    }
                }

                // Redirect to admin page
                await router.push('/admin');
            } catch (err: any) {
                // Check if it's a 401 invalid credentials error
                const isInvalidCredentials =
                    err.response?.status === 401 ||
                    err.error?.toLowerCase().includes('invalid credentials') ||
                    err.message?.toLowerCase().includes('invalid credentials');

                if (isInvalidCredentials) {
                    error.value = t(
                        'auth.invalidCredentials',
                        'Invalid username or password. Please try again.'
                    );
                } else {
                    error.value =
                        err.message ||
                        err.error ||
                        t(
                            'auth.loginError',
                            'Login failed. Please check your credentials.'
                        );
                }
                console.error('Login error:', err);
            } finally {
                loading.value = false;
            }
        };

        // Check for remembered user
        const rememberedUser =
            typeof window !== 'undefined' ? localStorage.getItem('rememberedUser') : null;
        if (rememberedUser) {
            formState.username = rememberedUser;
            rememberMe.value = true;
        }

        return {
            formState,
            rememberMe,
            loading,
            error,
            passwordVisible,
            handleLogin,
        };
    },
};
</script>

<style scoped>
/* Custom scrollbar for dark mode */
.dark ::-webkit-scrollbar {
    width: 8px;
}

.dark ::-webkit-scrollbar-track {
    background: #1f2937;
}

.dark ::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 4px;
}

.dark ::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
}

/* Smooth transitions */
* {
    transition: background-color 0.2s ease, border-color 0.2s ease;
}

/* Focus styles */
:focus-visible {
    outline: 2px solid var(--color-green-500);
    outline-offset: 2px;
}
</style>
