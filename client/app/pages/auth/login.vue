<template>
    <div
        class="min-h-screen flex items-center justify-center bg-linear-to-br from-green-50 to-teal-50 dark:from-gray-900 dark:to-gray-800 p-4"
    >
        <div class="w-full max-w-lg">
            <!-- Logo/Brand -->
            <div class="text-center mb-6">
                <div
                    class="inline-flex items-center justify-center w-14 h-14 bg-teal-500 rounded-2xl shadow-lg mb-3"
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
                        />
                    </svg>
                </div>
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                    {{ $t('auth.login',) }}
                </h1>
                <p class="text-gray-600 dark:text-gray-300 mt-2">
                    {{
                        $t(
                            'auth.welcomeBack',
                            'Welcome back! Please enter your credentials',
                        )
                    }}
                </p>
            </div>

            <!-- Login Card -->
            <UCard class="shadow-2xl border-0">
                <UForm
                    :state="formState"
                    class="space-y-4"
                    @submit="handleLogin"
                >
                    <!-- Username Field -->
                    <UFormField
                        name="username"
                        :label="$t('auth.username',)"
                        required
                    >
                        <UInput
                            v-model="formState.username"
                            type="text"
                            :placeholder="
                                $t('auth.usernamePlaceholder', 'Enter your username',)
                            "
                            icon="i-heroicons-user"
                            size="lg"
                            :disabled="loading"
                            autocomplete="username"
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Password Field -->
                    <UFormField
                        name="password"
                        :label="$t('auth.password',)"
                        required
                    >
                        <UInput
                            v-model="formState.password"
                            :type="passwordVisible ? 'text' : 'password'"
                            :placeholder="
                                $t('auth.passwordPlaceholder', 'Enter your password',)
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
                                    :padded="false"
                                    class="p-1!"
                                    @click="passwordVisible = !passwordVisible"
                                />
                            </template>
                        </UInput>
                    </UFormField>

                    <!-- Remember Me & Forgot Password -->
                    <div class="flex items-center justify-between">
                        <UCheckbox
                            v-model="rememberMe"
                            :label="$t('auth.remember',)"
                            name="remember"
                            :disabled="loading"
                        />
                        <NuxtLink
                            to="/forgot-password"
                            class="text-sm text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300 font-medium transition-colors"
                        >
                            {{ $t('auth.forgot',) }}
                        </NuxtLink>
                    </div>

                    <!-- Error Message -->
                    <UAlert
                        v-if="error"
                        :title="$t('auth.error', 'Error',)"
                        :description="error"
                        icon="i-heroicons-exclamation-triangle"
                        color="error"
                        variant="outline"
                        class="mt-4"
                    />

                    <!-- Email Not Verified Alert -->
                    <UAlert
                        v-if="emailNotVerified"
                        :title="$t('auth.emailNotVerified', 'Email Not Verified',)"
                        :description="emailNotVerified"
                        icon="i-heroicons-envelope"
                        color="warning"
                        variant="outline"
                        class="mt-4"
                    >
                        <template #actions>
                            <UButton
                                color="primary"
                                variant="solid"
                                size="sm"
                                :loading="resending"
                                :disabled="resending"
                                @click="resendVerification"
                            >
                                {{ $t('auth.resendVerification', 'Resend',) }}
                            </UButton>
                        </template>
                    </UAlert>

                    <!-- Resend Success -->
                    <UAlert
                        v-if="resendSuccess"
                        :title="$t('auth.emailSent', 'Email Sent',)"
                        :description="resendSuccess"
                        icon="i-heroicons-check-circle"
                        color="success"
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
                                />
                            </svg>
                        </template>
                        {{ $t('auth.login',) }}
                    </UButton>
                </UForm>

                <!-- Register Link -->
                <div class="mt-6 text-center">
                    <p class="text-gray-600 dark:text-gray-300">
                        {{ $t('auth.noAccount',) }}
                        <NuxtLink
                            to="/auth/register"
                            class="font-semibold text-teal-600 hover:text-teal-700 dark:text-teal-400 dark:hover:text-teal-300 transition-colors"
                        >
                            {{ $t('auth.registerNow',) }}
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
                        />
                    </svg>
                    {{ $t('auth.backToHome', 'Back to home',) }}
                </NuxtLink>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, reactive, onMounted, } from 'vue';
    import { useRouter, } from 'vue-router';
    import { useI18n, } from 'vue-i18n';
    import authApi from '~/api/auth';

    definePageMeta({
        layout: false,
    });

    interface LoginForm {
        username: string
        password: string
    }

    const router = useRouter();
    const { t, } = useI18n();

    const formState = reactive<LoginForm>({
        username: '',
        password: '',
    });

    const rememberMe = ref(false,);
    const loading = ref(false,);
    const error = ref('',);
    const passwordVisible = ref(false,);
    const emailNotVerified = ref('',);
    const unverifiedEmail = ref('',);
    const resending = ref(false,);
    const resendSuccess = ref('',);

    const handleLogin = async () => {
        if (!formState.username.trim() || !formState.password.trim()) {
            error.value = t('auth.fillAllFields', 'Please fill in all fields',);
            return;
        }

        loading.value = true;
        error.value = '';
        emailNotVerified.value = '';
        resendSuccess.value = '';

        try {
            // Use the new auth API which handles token storage automatically
            await authApi.login({
                username: formState.username,
                password: formState.password,
            });

            // Store username for "remember me" functionality
            if (rememberMe.value) {
                if (typeof window !== 'undefined') {
                    localStorage.setItem('rememberedUser', formState.username,);
                }
            } else {
                if (typeof window !== 'undefined') {
                    localStorage.removeItem('rememberedUser',);
                }
            }

            // Redirect to admin page
            await router.push('/admin',);
        } catch (err: any) {
            // Check if it's a 403 email not verified error
            if (
                err.code === 'email_not_verified'
            || err.response?.data?.code === 'email_not_verified'
            ) {
                const email = err.email || err.response?.data?.email || '';
                unverifiedEmail.value = email;
                emailNotVerified.value = t(
                    'auth.emailNotVerifiedMsg',
                    'Please verify your email before logging in. Check your inbox for the verification link.'
                );
                error.value = '';
            }
            // Check if it's a 401 invalid credentials error
            else if (
                err.response?.status === 401
            || err.error?.toLowerCase().includes('invalid credentials',)
                || err.message?.toLowerCase().includes('invalid credentials',)
            ) {
                error.value = t(
                    'auth.invalidCredentials',
                    'Invalid username or password. Please try again.'
                );
            } else {
                error.value
                = err.message
                        || err.error
                    || t('auth.loginError', 'Login failed. Please check your credentials.',);
            }
            console.error('Login error:', err,);
        } finally {
            loading.value = false;
        }
    };

    const resendVerification = async () => {
        if (!unverifiedEmail.value) {
            emailNotVerified.value = t(
                'auth.noEmailForResend',
                'Unable to resend verification. Please register again.'
            );
            return;
        }

        resending.value = true;
        resendSuccess.value = '';

        try {
            const response = await authApi.resendVerification(unverifiedEmail.value,);
            resendSuccess.value
            = response.message
                    || t(
                    'auth.verificationResent',
                    'If this email is registered, a new verification link has been sent.'
                );
        } catch (err: any) {
            emailNotVerified.value
            = err.message
                    || err.error
                || t(
                        'auth.resendError',
                        'Failed to resend verification email. Please try again later.'
                    );
            console.error('Resend error:', err,);
        } finally {
            resending.value = false;
        }
    };

    // Check for remembered username (client-side only, after hydration)
    onMounted(() => {
        const rememberedUser = localStorage.getItem('rememberedUser',);
        if (rememberedUser) {
            formState.username = rememberedUser;
            rememberMe.value = true;
        }
    });
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
    outline: 2px solid var(--color-primary, #4B9E90);
    outline-offset: 2px;
}
</style>
