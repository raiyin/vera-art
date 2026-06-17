import { useI18n, } from '#imports';

/**
 * Composable providing reusable formatting utilities.
 *
 * Consolidates duplicated formatting logic found across components and pages:
 * - formatPrice: price formatting (kopecks → rubles)
 * - formatDate: date formatting with i18n-aware locale
 * - formatDateTime: date + time formatting
 */
export function useFormatting() {
    const { locale, } = useI18n();

    /**
     * Format a price value.
     *
     * @param price - Price in kopecks (1 RUB = 100 kopecks)
     * @param options - Formatting options
     * @returns Formatted price string (e.g., "1 500 ₽" or "Бесплатно")
     */
    function formatPrice(
        price: number | null | undefined,
        options?: {
            /** If true, input is in kopecks (default: true) */
            inKopecks?: boolean
            /** Custom zero-price label (default: 'Бесплатно') */
            freeLabel?: string
            /** Fallback for null/undefined (default: '—') */
            fallback?: string
        },
    ): string {
        const {
            inKopecks = true,
            freeLabel = 'Бесплатно',
            fallback = '—',
        } = options ?? {};

        if (price == null) return fallback;

        const value = inKopecks ? price / 100 : price;
        if (value === 0) return freeLabel;

        return `${value.toLocaleString('ru-RU',)} ₽`;
    }

    /**
     * Format a date string using the current i18n locale.
     *
     * @param dateStr - ISO date string
     * @param options - Intl.DateTimeFormat options
     * @returns Formatted date string, or fallback for empty input
     */
    function formatDate(
        dateStr: string | null | undefined,
        options?: Intl.DateTimeFormatOptions,
    ): string {
        if (!dateStr) return '—';

        const date = new Date(dateStr,);
        const stdLocale = locale.value === 'ru' ? 'ru-RU' : 'en-US';

        const defaultOptions: Intl.DateTimeFormatOptions = {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            ...options,
        };

        return date.toLocaleDateString(stdLocale, defaultOptions,);
    }

    /**
     * Format a date string with time using the current i18n locale.
     *
     * @param dateStr - ISO date string
     * @returns Formatted date+time string, or fallback for empty input
     */
    function formatDateTime(dateStr: string | null | undefined,): string {
        if (!dateStr) return '—';

        const date = new Date(dateStr,);
        const stdLocale = locale.value === 'ru' ? 'ru-RU' : 'en-US';

        return date.toLocaleString(stdLocale, {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        },);
    }

    /**
     * Format a short date (day + month) for compact displays.
     *
     * @param dateStr - ISO date string
     * @returns Short date string (e.g., "15 марта" or "Mar 15")
     */
    function formatShortDate(dateStr: string | null | undefined,): string {
        if (!dateStr) return '';

        const date = new Date(dateStr,);
        const stdLocale = locale.value === 'ru' ? 'ru-RU' : 'en-US';

        return date.toLocaleDateString(stdLocale, {
            day: 'numeric',
            month: 'short',
        },);
    }

    return {
        formatPrice,
        formatDate,
        formatDateTime,
        formatShortDate,
    };
}
