import { browser } from "$app/environment";
import { getCookie } from "$lib/utils/getCookie";
import { writable } from "svelte/store";

const TOKEN_NAME = 'labyrinth-gdsc-token'
const initialValue = browser ? getCookie(document.cookie, TOKEN_NAME) : null;
const JwtTokenStore = writable<string | null>(initialValue);
const setToken = (value: string | null) => {
    if (!browser) return;

    if (value) {
        document.cookie = `${TOKEN_NAME}=${value};path=/;max-age=${60 * 60 * 24 * 365}`;
    } else {
        document.cookie = `${TOKEN_NAME}=null;path=/;expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
    }

    JwtTokenStore.set(value);
}

export { JwtTokenStore, setToken, TOKEN_NAME }