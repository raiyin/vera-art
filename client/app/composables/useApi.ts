import { getHttpClient, } from '~/api/http-client';

/**
 * Composable wrapper around the centralized HTTP client.
 * Provides typed convenience methods for common HTTP operations.
 */
export function useApi() {
    const client = getHttpClient();

    async function get<T,>(url: string, params?: Record<string, unknown>,): Promise<T> {
        const response = await client.get<T>(url, { params, },);
        return response.data;
    }

    async function post<T,>(url: string, data?: unknown,): Promise<T> {
        const response = await client.post<T>(url, data,);
        return response.data;
    }

    async function put<T,>(url: string, data?: unknown,): Promise<T> {
        const response = await client.put<T>(url, data,);
        return response.data;
    }

    async function del<T,>(url: string,): Promise<T> {
        const response = await client.delete<T>(url,);
        return response.data;
    }

    return { get, post, put, del, };
}
