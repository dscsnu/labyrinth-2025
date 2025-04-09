import { get } from "svelte/store";
import { token } from "$lib/stores/TokenStore";

interface FetchOptions extends RequestInit {
    headers?: Record<string, string>;
}

export const fetchWithAuth = async (url: string, options: FetchOptions = {}): Promise<Response> => {
    const jwt = get(token);

    return fetch(url, {
        ...options,
        headers: {
            ...options.headers,
            'Authorization': `Bearer ${jwt}`
        }
    });
}