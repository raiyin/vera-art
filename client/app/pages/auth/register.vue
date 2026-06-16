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
                            d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                        ></path>
                    </svg>
                </div>
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                    {{ $t('auth.register') }}
                </h1>
                <p class="text-gray-600 dark:text-gray-300 mt-2">
                    {{ $t('auth.createAccount', 'Create your account to get started') }}
                </p>
            </div>

            <!-- Register Card -->
            <UCard class="shadow-2xl border-0">
                <UForm :state="formState" class="space-y-4" @submit="handleRegister">
                    <!-- Username Field -->
                    <UFormField name="username" :label="$t('auth.username')" required>
                        <UInput
                            v-model="formState.username"
                            type="text"
                            :placeholder="
                                $t('auth.usernamePlaceholder', 'Choose a username')
                            "
                            icon="i-heroicons-user"
                            size="lg"
                            :disabled="loading"
                            autocomplete="username"
                            :rules="[validateUsername]"
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Email Field -->
                    <UFormField name="email" :label="$t('auth.email')" required>
                        <UInput
                            v-model="formState.email"
                            type="email"
                            :placeholder="
                                $t('auth.emailPlaceholder', 'Enter your email address')
                            "
                            icon="i-heroicons-envelope"
                            size="lg"
                            :disabled="loading"
                            autocomplete="email"
                            :rules="[validateEmail]"
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Password Field -->
                    <UFormField name="password" :label="$t('auth.password')" required>
                        <UInput
                            v-model="formState.password"
                            :type="passwordVisible ? 'text' : 'password'"
                            :placeholder="
                                $t('auth.passwordPlaceholder', 'Create a strong password')
                            "
                            icon="i-heroicons-lock-closed"
                            size="lg"
                            :disabled="loading"
                            autocomplete="new-password"
                            :rules="[validatePassword]"
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
                        <template #help>
                            <div
                                class="text-xs text-gray-500 dark:text-gray-400 mt-1 space-y-1"
                            >
                                <p>
                                    {{
                                        $t('auth.passwordRequirements', 'Password must:')
                                    }}
                                </p>
                                <ul class="list-disc pl-4">
                                    <li
                                        :class="
                                            passwordRequirements.minLength
                                                ? 'text-green-600'
                                                : 'text-red-600'
                                        "
                                    >
                                        {{
                                            $t(
                                                'auth.passwordMinLength',
                                                'Be at least 8 characters'
                                            )
                                        }}
                                    </li>
                                    <li
                                        :class="
                                            passwordRequirements.hasUppercase
                                                ? 'text-green-600'
                                                : 'text-red-600'
                                        "
                                    >
                                        {{
                                            $t(
                                                'auth.passwordUppercase',
                                                'Contain an uppercase letter'
                                            )
                                        }}
                                    </li>
                                    <li
                                        :class="
                                            passwordRequirements.hasLowercase
                                                ? 'text-green-600'
                                                : 'text-red-600'
                                        "
                                    >
                                        {{
                                            $t(
                                                'auth.passwordLowercase',
                                                'Contain a lowercase letter'
                                            )
                                        }}
                                    </li>
                                    <li
                                        :class="
                                            passwordRequirements.hasDigit
                                                ? 'text-green-600'
                                                : 'text-red-600'
                                        "
                                    >
                                        {{ $t('auth.passwordDigit', 'Contain a digit') }}
                                    </li>
                                    <li
                                        :class="
                                            passwordRequirements.hasSpecial
                                                ? 'text-green-600'
                                                : 'text-red-600'
                                        "
                                    >
                                        {{
                                            $t(
                                                'auth.passwordSpecial',
                                                'Contain a special character'
                                            )
                                        }}
                                    </li>
                                </ul>
                            </div>
                        </template>
                    </UFormField>

                    <!-- Confirm Password Field -->
                    <UFormField
                        name="confirmPassword"
                        :label="$t('auth.confirmPassword')"
                        required
                    >
                        <UInput
                            v-model="formState.confirmPassword"
                            :type="confirmPasswordVisible ? 'text' : 'password'"
                            :placeholder="
                                $t(
                                    'auth.confirmPasswordPlaceholder',
                                    'Confirm your password'
                                )
                            "
                            icon="i-heroicons-lock-closed"
                            size="lg"
                            :disabled="loading"
                            autocomplete="new-password"
                            :rules="[validateConfirmPassword]"
                            class="w-full"
                        >
                            <template #trailing>
                                <UButton
                                    variant="ghost"
                                    color="neutral"
                                    :icon="
                                        confirmPasswordVisible
                                            ? 'i-heroicons-eye-slash'
                                            : 'i-heroicons-eye'
                                    "
                                    @click="
                                        confirmPasswordVisible = !confirmPasswordVisible
                                    "
                                    :padded="false"
                                    class="!p-1"
                                />
                            </template>
                        </UInput>
                    </UFormField>

                    <!-- Terms Agreement -->
                    <UCheckbox
                        v-model="agreeTerms"
                        name="terms"
                        :disabled="loading"
                        required
                    >
                        <template #label>
                            <span class="text-sm text-gray-700 dark:text-gray-300">
                                {{ $t('auth.agreeTo') }}
                                <NuxtLink
                                    to="/terms"
                                    class="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium"
                                >
                                    {{ $t('auth.termsOfService', 'Terms of Service') }}
                                </NuxtLink>
                                {{ $t('auth.and') }}
                                <NuxtLink
                                    to="/privacy"
                                    class="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium"
                                >
                                    {{ $t('auth.privacyPolicy', 'Privacy Policy') }}
                                </NuxtLink>
                            </span>
                        </template>
                    </UCheckbox>

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

                    <!-- Success Message -->
                    <UAlert
                        v-if="success"
                        :title="$t('auth.success', 'Success!')"
                        :description="success"
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
                        :disabled="loading || !agreeTerms"
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
                                    d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                                ></path>
                            </svg>
                        </template>
                        {{ $t('auth.register') }}
                    </UButton>
                </UForm>

                <!-- Login Link -->
                <div class="mt-6 text-center">
                    <p class="text-gray-600 dark:text-gray-300">
                        {{ $t('auth.haveAccount') }}
                        <NuxtLink
                            to="/auth/login"
                            class="font-semibold text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
                        >
                            {{ $t('auth.login') }}
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
import { ref, reactive, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import authApi from '../../api/auth';

interface RegisterForm {
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
}

export default {
    setup() {
        const router = useRouter();
        const { t } = useI18n();

        const formState = reactive<RegisterForm>({
            username: '',
            email: '',
            password: '',
            confirmPassword: '',
        });

        const agreeTerms = ref(false);
        const loading = ref(false);
        const error = ref('');
        const success = ref('');
        const passwordVisible = ref(false);
        const confirmPasswordVisible = ref(false);

        // Password requirements checker
        const passwordRequirements = computed(() => {
            const password = formState.password;
            return {
                minLength: password.length >= 8,
                hasUppercase: /[A-Z]/.test(password),
                hasLowercase: /[a-z]/.test(password),
                hasDigit: /[0-9]/.test(password),
                hasSpecial: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/.test(password),
                notWeak: ![
                    'password',
                    '12345678',
                    'qwerty',
                    'admin',
                    'letmein',
                ].some((weak) => password.toLowerCase().includes(weak)),
            };
        });

        // Check if all password requirements are met
        const isPasswordValid = computed(() => {
            const req = passwordRequirements.value;
            return (
                req.minLength &&
                req.hasUppercase &&
                req.hasLowercase &&
                req.hasDigit &&
                req.hasSpecial &&
                req.notWeak
            );
        });

        // Validation functions
        const validateUsername = (value: string) => {
            if (!value.trim()) return t('auth.usernameRequired', 'Username is required');
            if (value.length < 3)
                return t(
                    'auth.usernameMinLength',
                    'Username must be at least 3 characters'
                );
            if (value.length > 20)
                return t(
                    'auth.usernameMaxLength',
                    'Username must be less than 20 characters'
                );
            return true;
        };

        const validatePassword = (value: string) => {
            if (!value.trim()) return t('auth.passwordRequired', 'Password is required');

            const req = passwordRequirements.value;
            if (!req.minLength)
                return t(
                    'auth.passwordMinLength',
                    'Password must be at least 8 characters'
                );
            if (!req.hasUppercase)
                return t(
                    'auth.passwordUppercase',
                    'Password must contain an uppercase letter'
                );
            if (!req.hasLowercase)
                return t(
                    'auth.passwordLowercase',
                    'Password must contain a lowercase letter'
                );
            if (!req.hasDigit)
                return t('auth.passwordDigit', 'Password must contain a digit');
            if (!req.hasSpecial)
                return t(
                    'auth.passwordSpecial',
                    'Password must contain a special character'
                );
            if (!req.notWeak)
                return t('auth.passwordWeak', 'Password is too common or weak');

            return true;
        };

        const validateEmail = (value: string) => {
            if (!value.trim()) return t('auth.emailRequired', 'Email is required');
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value.trim()))
                return t('auth.emailInvalid', 'Please enter a valid email address');
            return true;
        };

        const validateConfirmPassword = (value: string) => {
            if (!value.trim())
                return t('auth.confirmPasswordRequired', 'Please confirm your password');
            if (value !== formState.password)
                return t('auth.passwordsDontMatch', 'Passwords do not match');
            return true;
        };

        const isFormValid = computed(() => {
            const emailValid =
                formState.email.trim() &&
                /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formState.email.trim());
            return (
                formState.username.trim() &&
                emailValid &&
                formState.password.trim() &&
                formState.confirmPassword.trim() &&
                formState.password === formState.confirmPassword &&
                isPasswordValid.value &&
                agreeTerms.value
            );
        });

        const handleRegister = async () => {
            if (!isFormValid.value) {
                error.value = t(
                    'auth.fillAllFields',
                    'Please fill in all fields correctly'
                );
                return;
            }

            loading.value = true;
            error.value = '';
            success.value = '';

            try {
                // Use the new auth API
                await authApi.register({
                    username: formState.username,
                    password: formState.password,
                    email: formState.email,
                });

                success.value = t(
                    'auth.verificationEmailSent',
                    'Registration successful! A verification link has been sent to your email. Please check your inbox.'
                );

                // Clear form
                formState.username = '';
                formState.email = '';
                formState.password = '';
                formState.confirmPassword = '';
                agreeTerms.value = false;

                // Redirect to login after 5 seconds
                setTimeout(() => {
                    router.push('/auth/login');
                }, 5000);
            } catch (err: any) {
                error.value =
                    err.message ||
                    err.error ||
                    t(
                        'auth.registrationError',
                        'Registration failed. Username may be taken.'
                    );
                console.error('Registration error:', err);
            } finally {
                loading.value = false;
            }
        };

        return {
            formState,
            agreeTerms,
            loading,
            error,
            success,
            passwordVisible,
            confirmPasswordVisible,
            passwordRequirements,
            validateUsername,
            validateEmail,
            validatePassword,
            validateConfirmPassword,
            isFormValid,
            handleRegister,
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
    outline: 2px solid var(--color-blue-500);
    outline-offset: 2px;
}

/* Password strength indicator */
.password-strength {
    height: 4px;
    border-radius: 2px;
    margin-top: 2px;
    transition: all 0.3s ease;
}

.password-strength.weak {
    background-color: #ef4444;
    width: 33%;
}

.password-strength.medium {
    background-color: #f59e0b;
    width: 66%;
}

.password-strength.strong {
    background-color: #10b981;
    width: 100%;
}
</style>
