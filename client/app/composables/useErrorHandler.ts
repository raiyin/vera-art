/**
 * Composable providing reusable error handling utilities.
 *
 * Consolidates the common pattern of:
 * - Extracting error messages from unknown error types
 * - Logging errors consistently
 * - Providing user-facing error messages
 */
export function useErrorHandler() {
    /**
     * Extract a human-readable error message from an unknown error value.
     *
     * Handles:
     * - Error instances (uses `.message`)
     * - Objects with `.message` property
     * - Axios error responses (uses `.response?.data?.message` or `.response?.statusText`)
     * - Strings
     * - Fallback for anything else
     *
     * @param err - The caught error value
     * @param fallback - Default message if nothing can be extracted (default: 'An unexpected error occurred')
     * @returns Extracted error message string
     */
    function extractErrorMessage(err: unknown, fallback = 'An unexpected error occurred',): string {
        if (typeof err === 'object' && err !== null) {
            // Axios-style error with response data
            const axiosErr = err as { response?: { data?: { message?: string }, statusText?: string }, message?: string };
            if (axiosErr.response?.data?.message) {
                return axiosErr.response.data.message;
            }
            if (axiosErr.response?.statusText) {
                return axiosErr.response.statusText;
            }
            if (axiosErr.message) {
                return axiosErr.message;
            }
        }
        if (typeof err === 'string') {
            return err;
        }
        return fallback;
    }

    /**
     * Log an error to the console with a consistent prefix.
     * Use this instead of raw `console.error()` throughout the codebase.
     *
     * @param context - A short description of where the error occurred
     * @param err - The error value to log
     */
    function logError(context: string, err: unknown,): void {
        console.error(`[${context}]`, err,);
    }

    /**
     * Extract an error message and log it in one call.
     * Returns the extracted message for use in UI state.
     *
     * @param context - A short description of where the error occurred
     * @param err - The error value
     * @param fallback - Default message if nothing can be extracted
     * @returns Extracted error message
     */
    function handleError(context: string, err: unknown, fallback = 'An unexpected error occurred',): string {
        const message = extractErrorMessage(err, fallback,);
        logError(context, err,);
        return message;
    }

    return {
        extractErrorMessage,
        logError,
        handleError,
    };
}
