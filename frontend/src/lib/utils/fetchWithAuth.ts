import { get } from "svelte/store";
import { token } from "$lib/stores/TokenStore";
import { PUBLIC_BACKEND_URL } from "$env/static/public";

interface FetchOptions extends RequestInit {
    headers?: Record<string, string>;
}

/**
 * Performs an authenticated fetch request to the backend.
 *
 * This function automatically adds a Bearer token from the Svelte store
 * to the Authorization header and ensures the backend URL and route are properly joined
 * without duplicate slashes.
 *
 * @param route - The API route (e.g., "api/team" or "/api/team").
 * @param options - Optional fetch configuration (headers, method, body, etc.).
 * @returns A Promise resolving to the Fetch API Response object.
 */
export const fetchWithAuth = async (route: string, options: FetchOptions = {}): Promise<Response> => {
    const jwt = get(token);
    const cleanedUrl = PUBLIC_BACKEND_URL.replace(/\/+$/, "");
    const cleanedRoute = route.replace(/^\/+/, "");

    return fetch(`${cleanedUrl}/${cleanedRoute}`, {
        ...options,
        headers: {
            ...options.headers,
            'Authorization': `Bearer ${jwt}`
        }
    });
}