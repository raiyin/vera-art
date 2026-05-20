import { defineStore, } from 'pinia';
import { ref, computed, } from 'vue';

export interface CookieConsent {
    necessary: boolean
    analytics: boolean
    marketing: boolean
    timestamp: string
    version: string
}

export const useCookieConsentStore = defineStore('cookieConsent', () => {
    const consent = ref<CookieConsent | null>(null,);
    const COOKIE_NAME = 'cookie_consent';
    const COOKIE_EXPIRY_DAYS = 365;

    // Initialize from cookies
    const initialize = () => {
        if (typeof window === 'undefined') return;

        const cookieValue = getCookie(COOKIE_NAME,);
        if (cookieValue) {
            try {
                consent.value = JSON.parse(cookieValue,);
            } catch (error) {
                console.error('Failed to parse cookie consent:', error,);
                clearConsent();
            }
        }
    };

    // Get current consent
    const getConsent = (): CookieConsent | null => {
        if (!consent.value) {
            initialize();
        }
        return consent.value;
    };

    // Set consent and save to cookie
    const setConsent = (newConsent: CookieConsent,) => {
        consent.value = newConsent;

        if (typeof window !== 'undefined') {
            // Save to cookie
            const expiryDate = new Date();
            expiryDate.setDate(expiryDate.getDate() + COOKIE_EXPIRY_DAYS,);

            const cookieValue = JSON.stringify(newConsent,);
            document.cookie = `${COOKIE_NAME}=${encodeURIComponent(cookieValue,)}; expires=${expiryDate.toUTCString()}; path=/; SameSite=Strict${window.location.protocol === 'https:' ? '; Secure' : ''}`;

            // Also save to localStorage as backup
            localStorage.setItem(COOKIE_NAME, cookieValue,);

            // Initialize analytics based on consent
            initializeAnalytics(newConsent,);
        }
    };

    // Clear consent
    const clearConsent = () => {
        consent.value = null;

        if (typeof window !== 'undefined') {
            // Clear cookie
            document.cookie = `${COOKIE_NAME}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
            localStorage.removeItem(COOKIE_NAME,);

            // Disable analytics
            disableAnalytics();
        }
    };

    // Check if consent is given
    const hasConsent = computed(() => {
        return consent.value !== null;
    },);

    // Check specific consent types
    const hasAnalyticsConsent = computed(() => {
        return consent.value?.analytics === true;
    },);

    const hasMarketingConsent = computed(() => {
        return consent.value?.marketing === true;
    },);

    // Initialize analytics based on consent
    const initializeAnalytics = (consentData: CookieConsent,) => {
        if (typeof window === 'undefined') return;

        // Example: Initialize Google Analytics if consent given
        if (consentData.analytics && window.gtag) {
            // Enable analytics tracking
            window.dataLayer = window.dataLayer || [];
            window.gtag('consent', 'update', {
                analytics_storage: 'granted',
            },);
        } else {
            disableAnalytics();
        }

        // Example: Initialize marketing pixels if consent given
        if (consentData.marketing) {
            // Enable marketing tracking
            // This would typically load Facebook Pixel, etc.
        } else {
            disableMarketing();
        }
    };

    const disableAnalytics = () => {
        if (typeof window === 'undefined') return;

        // Disable analytics tracking
        if (window.gtag) {
            window.gtag('consent', 'update', {
                analytics_storage: 'denied',
            },);
        }

        // Clear analytics cookies
        const analyticsCookies = ['_ga', '_gid', '_gat',];
        analyticsCookies.forEach((cookie,) => {
            document.cookie = `${cookie}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
        },);
    };

    const disableMarketing = () => {
        if (typeof window === 'undefined') return;

        // Clear marketing cookies
        const marketingCookies = ['_fbp', 'fr', 'tr',];
        marketingCookies.forEach((cookie,) => {
            document.cookie = `${cookie}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
        },);
    };

    // Helper function to get cookie value
    const getCookie = (name: string,): string | null => {
        if (typeof window === 'undefined') return null;

        const cookies = document.cookie.split(';',);
        for (const cookie of cookies) {
            const [cookieName, cookieValue,] = cookie.trim().split('=',);
            if (cookieName === name && cookieValue) {
                return decodeURIComponent(cookieValue,);
            }
        }
        return null;
    };

    // Initialize on store creation
    if (typeof window !== 'undefined') {
        initialize();
    }

    return {
        consent,
        getConsent,
        setConsent,
        clearConsent,
        hasConsent,
        hasAnalyticsConsent,
        hasMarketingConsent,
        initialize,
    };
},);

// Type declarations for global analytics objects
declare global {
    interface Window {
        gtag?: (...args: unknown[]) => void
        dataLayer?: unknown[]
        fbq?: (...args: unknown[]) => void
    }
}
