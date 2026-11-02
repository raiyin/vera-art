<template>
    <div
        v-if="showBanner"
        class="cookie-consent-banner fixed bottom-0 left-0 right-0 z-50 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 shadow-lg"
        role="dialog"
        aria-labelledby="cookie-consent-title"
        aria-describedby="cookie-consent-description"
    >
        <div class="container mx-auto px-4 py-6">
            <div
                class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4"
            >
                <div class="flex-1">
                    <h2
                        id="cookie-consent-title"
                        class="text-lg font-semibold text-gray-900 dark:text-white mb-2"
                    >
                        {{ $t('cookieConsent.title',) }}
                    </h2>
                    <p
                        id="cookie-consent-description"
                        class="text-gray-700 dark:text-gray-300 text-sm mb-4"
                    >
                        {{ $t('cookieConsent.description',) }}
                        <NuxtLink
                            :to="localePath('/privacy',)"
                            class="text-teal-600 dark:text-teal-400 hover:underline font-medium"
                        >
                            {{ $t('cookieConsent.privacyPolicy',) }}
                        </NuxtLink>
                        {{ $t('cookieConsent.and',) }}
                        <NuxtLink
                            :to="localePath('/terms',)"
                            class="text-teal-600 dark:text-teal-400 hover:underline font-medium"
                        >
                            {{ $t('cookieConsent.termsOfService',) }} </NuxtLink>.
                    </p>
                    <div class="flex flex-wrap gap-4 mt-4">
                        <div class="flex items-center">
                            <input
                                id="necessary-cookies"
                                type="checkbox"
                                checked
                                disabled
                                class="h-4 w-4 text-teal-600 border-gray-300 rounded focus:ring-teal-500 dark:focus:ring-teal-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                            >
                            <label
                                for="necessary-cookies"
                                class="ml-2 text-sm text-gray-700 dark:text-gray-300"
                            >
                                {{ $t('cookieConsent.necessaryCookies',) }}
                                <span class="text-xs text-gray-500 dark:text-gray-400">({{ $t('cookieConsent.alwaysActive',) }})</span>
                            </label>
                        </div>
                        <div class="flex items-center">
                            <input
                                id="analytics-cookies"
                                v-model="analyticsConsent"
                                type="checkbox"
                                class="h-4 w-4 text-teal-600 border-gray-300 rounded focus:ring-teal-500 dark:focus:ring-teal-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                            >
                            <label
                                for="analytics-cookies"
                                class="ml-2 text-sm text-gray-700 dark:text-gray-300"
                            >
                                {{ $t('cookieConsent.analyticsCookies',) }}
                            </label>
                        </div>
                        <div class="flex items-center">
                            <input
                                id="marketing-cookies"
                                v-model="marketingConsent"
                                type="checkbox"
                                class="h-4 w-4 text-teal-600 border-gray-300 rounded focus:ring-teal-500 dark:focus:ring-teal-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                            >
                            <label
                                for="marketing-cookies"
                                class="ml-2 text-sm text-gray-700 dark:text-gray-300"
                            >
                                {{ $t('cookieConsent.marketingCookies',) }}
                            </label>
                        </div>
                    </div>
                </div>
                <div class="flex flex-col sm:flex-row gap-3">
                    <button
                        type="button"
                        class="px-5 py-2.5 bg-teal-600 hover:bg-teal-700 text-white font-medium rounded-lg text-sm transition-colors duration-200"
                        @click="acceptAll"
                    >
                        {{ $t('cookieConsent.acceptAll',) }}
                    </button>
                    <button
                        type="button"
                        class="px-5 py-2.5 bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 font-medium rounded-lg text-sm transition-colors duration-200"
                        @click="savePreferences"
                    >
                        {{ $t('cookieConsent.savePreferences',) }}
                    </button>
                    <button
                        type="button"
                        class="px-5 py-2.5 bg-red-100 hover:bg-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 text-red-700 dark:text-red-400 font-medium rounded-lg text-sm transition-colors duration-200"
                        @click="rejectAll"
                    >
                        {{ $t('cookieConsent.rejectAll',) }}
                    </button>
                </div>
            </div>
            <div class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                <button
                    type="button"
                    class="text-sm text-teal-600 dark:text-teal-400 hover:underline flex items-center"
                    @click="showDetails = !showDetails"
                >
                    <span>{{ $t('cookieConsent.moreDetails',) }}</span>
                    <svg
                        :class="{ 'rotate-180': showDetails, }"
                        class="w-4 h-4 ml-1 transition-transform duration-200"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 9l-7 7-7-7"
                        />
                    </svg>
                </button>
                <div
                    v-if="showDetails"
                    class="mt-3 text-sm text-gray-600 dark:text-gray-400 space-y-2"
                >
                    <p>{{ $t('cookieConsent.details.necessary',) }}</p>
                    <p>{{ $t('cookieConsent.details.analytics',) }}</p>
                    <p>{{ $t('cookieConsent.details.marketing',) }}</p>
                    <p class="mt-2">
                        {{ $t('cookieConsent.details.rights',) }}
                        <NuxtLink
                            :to="localePath('/privacy#rights',)"
                            class="text-teal-600 dark:text-teal-400 hover:underline"
                        >
                            {{ $t('cookieConsent.details.privacyPolicy',) }} </NuxtLink>.
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, onMounted, watch, } from 'vue';
    import { useI18n, } from '#imports';
    import { useLocalePath, } from '#i18n';
    import { useAuthStore, } from '~/stores/AuthStore';
    import { useCookieConsentStore, } from '~/stores/CookieConsentStore';

    const { locale, } = useI18n();
    const localePath = useLocalePath();
    const authStore = useAuthStore();
    const cookieConsentStore = useCookieConsentStore();

    const showBanner = ref(false,);
    const showDetails = ref(false,);
    const analyticsConsent = ref(false,);
    const marketingConsent = ref(false,);

    // Check if consent is already given
    const checkConsent = () => {
        const consent = cookieConsentStore.getConsent();
        if (consent) {
            // User has already given consent
            showBanner.value = false;
            analyticsConsent.value = consent.analytics;
            marketingConsent.value = consent.marketing;
        } else {
            // Show banner if no consent
            showBanner.value = true;
        }
    };

    // Save consent to cookies and optionally to backend
    const saveConsent = async (analytics: boolean, marketing: boolean,) => {
        const consent = {
            necessary: true,
            analytics,
            marketing,
            timestamp: new Date().toISOString(),
            version: '1.0',
        };

        // Save to cookies
        cookieConsentStore.setConsent(consent,);

        // If user is authenticated, save to backend
        if (authStore.isAuthenticated && authStore.userId) {
            try {
                await $fetch('/api/consent', {
                    method: 'POST',
                    body: {
                        userId: authStore.userId,
                        analyticsConsent: analytics,
                        marketingConsent: marketing,
                        ipAddress: '', // Will be filled by backend
                        userAgent:
                            typeof window !== 'undefined' ? window.navigator.userAgent : '',
                    },
                });
            } catch (error) {
                console.error('Failed to save consent to backend:', error,);
            // Continue anyway - consent is saved in cookies
            }
        }

        showBanner.value = false;

        // Emit event for other components
        window.dispatchEvent(new CustomEvent('cookie-consent-updated', { detail: consent, },),);
    };

    const acceptAll = () => {
        saveConsent(true, true,);
    };

    const savePreferences = () => {
        saveConsent(analyticsConsent.value, marketingConsent.value,);
    };

    const rejectAll = () => {
        saveConsent(false, false,);
    };

    // Watch for authentication changes
    watch(
        () => authStore.isAuthenticated,
        (isAuthenticated,) => {
            if (isAuthenticated && cookieConsentStore.getConsent()) {
                // User logged in and has consent in cookies - sync to backend
                const consent = cookieConsentStore.getConsent();
                if (consent) {
                    saveConsent(consent.analytics, consent.marketing,);
                }
            }
        }
    );

    onMounted(() => {
        // Check consent after component is mounted
        setTimeout(checkConsent, 100,);
    });
</script>

<style scoped>
.cookie-consent-banner {
    animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
    from {
        transform: translateY(100%);
    }
    to {
        transform: translateY(0);
    }
}

.rotate-180 {
    transform: rotate(180deg);
}
</style>
