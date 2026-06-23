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
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        ></path>
                    </svg>
                </div>
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                    {{ $t('auth.emailVerification', 'Email Verification') }}
                </h1>
            </div>

            <!-- Verification Card -->
            <UCard class="shadow-2xl border-0">
                <!-- Loading State -->
                <div v-if="loading" class="text-center py-8">
                    <UIcon
                        name="i-heroicons-arrow-path"
                        class="w-12 h-12 text-indigo-500 mx-auto mb-4 animate-spin"
                    />
                    <p class="text-gray-600 dark:text-gray-300 text-lg">
                        {{ $t('auth.verifyingEmail', 'Verifying your email...') }}
                    </p>
                </div>

                <!-- Success State -->
                <UAlert
                    v-else-if="success"
                    :title="$t('auth.verificationSuccess', 'Email Verified!')"
                    :description="success"
                    icon="i-heroicons-check-circle"
                    color="success"
                    variant="outline"
                    class="mb-4"
                />
                <div v-else-if="success" class="text-center mt-6">
                    <UButton
                        color="primary"
                        size="lg"
                        @click="router.push('/auth/login')"
                    >
                        {{ $t('auth.goToLogin', 'Go to Login') }}
                    </UButton>
                </div>

                <!-- Error State -->
                <template v-else-if="error">
                    <UAlert
                        :title="$t('auth.verificationFailed', 'Verification Failed')"
                        :description="error"
                        icon="i-heroicons-exclamation-triangle"
                        color="error"
                        variant="outline"
                        class="mb-4"
                    />
                    <div class="text-center mt-6 space-y-3">
                        <p class="text-gray-600 dark:text-gray-300">
                            {{
                                $t(
                                    'auth.verificationErrorHint',
                                    'You can request a new verification link below.'
                                )
                            }}
                        </p>
                        <UButton
                            color="primary"
                            variant="outline"
                            size="lg"
                            :loading="resending"
                            :disabled="resending"
                            @click="resendVerification"
                        >
                            {{
                                $t('auth.resendVerification', 'Resend Verification Email')
                            }}
                        </UButton>
                        <div class="mt-2">
                            <UButton
                                color="neutral"
                                variant="ghost"
                                @click="router.push('/auth/login')"
                            >
                                {{ $t('auth.backToLogin', 'Back to Login') }}
                            </UButton>
                        </div>
                    </div>
                </template>

                <!-- Resend Success -->
                <UAlert
                    v-if="resendSuccess"
                    :title="$t('auth.emailSent', 'Email Sent')"
                    :description="resendSuccess"
                    icon="i-heroicons-check-circle"
                    color="success"
                    variant="outline"
                    class="mt-4"
                />
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

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import authApi from '~/api/auth';

definePageMeta({
    layout: false,
});

const router = useRouter();
const route = useRoute();
const { t } = useI18n();

const loading = ref(true);
const error = ref('');
const success = ref('');
const resending = ref(false);
const resendSuccess = ref('');

onMounted(async () => {
    const token = route.query.token as string;

    if (!token) {
        loading.value = false;
        error.value = t(
            'auth.noVerificationToken',
            'No verification token provided. Please use the link from your email.'
        );
        return;
    }

    try {
        const response = await authApi.verifyEmail(token);
        success.value =
            response.message ||
            t(
                'auth.emailVerifiedSuccess',
                'Your email has been verified successfully! You can now log in.'
            );
    } catch (err: any) {
        error.value =
            err.message ||
            err.error ||
            t(
                'auth.verificationError',
                'Email verification failed. The link may be invalid or expired.'
            );
        console.error('Verification error:', err);
    } finally {
        loading.value = false;
    }
});

const resendVerification = async () => {
    const email = route.query.email as string;
    if (!email) {
        error.value = t(
            'auth.noEmailForResend',
            'Unable to resend verification. Please go to login and use the resend option.'
        );
        return;
    }

    resending.value = true;
    resendSuccess.value = '';

    try {
        const response = await authApi.resendVerification(email);
        resendSuccess.value =
            response.message ||
            t(
                'auth.verificationResent',
                'If this email is registered, a new verification link has been sent.'
            );
    } catch (err: any) {
        error.value =
            err.message ||
            err.error ||
            t(
                'auth.resendError',
                'Failed to resend verification email. Please try again later.'
            );
        console.error('Resend error:', err);
    } finally {
        resending.value = false;
    }
};
</script>

<style scoped>
/* Smooth transitions */
* {
    transition: background-color 0.2s ease, border-color 0.2s ease;
}

/* Focus styles */
:focus-visible {
    outline: 2px solid var(--color-blue-500);
    outline-offset: 2px;
}
</style>
