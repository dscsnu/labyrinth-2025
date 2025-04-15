import { browser } from "$app/environment";
import { writable, type Writable } from "svelte/store"
import { getCookie } from "./getCookie";

/**
 * Creates a persistent store that saves data to cookies
 *
 * This function returns a Svelte store with additional functionality to persist data in cookies.
 * The store is automatically initialized from existing cookie data if available.
 *
 * @template T - The type of data to be stored
 * @param {string} tokenName - The name of the cookie to use for storage
 * @param {T | null} [initialValue=null] - Initial value to use if no cookie exists
 * @param {number} [maxAge=60 * 60 * 24 * 365] - Cookie max age in seconds (defaults to 1 year)
 * @returns An object containing:
 *  - store: A Svelte writable store of type T | null
 *  - set: A function to update the store value and persist it to cookies
 *
 * @example
 * ```typescript
 * interface UserData {
 *   id: string;
 *   name: string;
 * }
 *
 * const { store: userStore, set: setUser } = createPersistentStore<UserData>('user-data');
 *
 * // Subscribe to store changes
 * userStore.subscribe(data => {
 *   if (data) console.log(`User: ${data.name}`);
 * });
 *
 * // Update store and persist to cookie
 * setUser({ id: '123', name: 'John' });
 *
 * // Clear store and cookie
 * setUser(null);
 * ```
 */
const createPersistentStore = <T>(
    tokenName: string,
    initialvalue: T | null = null,
    maxAge: number = 60 * 60 * 24 * 365
): {
    store: Writable<T | null>;
    set: (data: T | null) => void;
} => {
    const getInitialData = (): T | null => {
        if (!browser) return initialvalue ?? null;

        const cookieData = getCookie(document.cookie, tokenName);
        if (cookieData) {
            try {
                return JSON.parse(cookieData);
            } catch (e) {
                console.error(`Error parsing cookie data for ${tokenName} : `, e);
            }
        }

        return initialvalue ?? null;
    }

    const store = writable<T | null>(getInitialData());

    const setData = (data: T | null) => {
        if (!browser) return;

        if (data) {
            document.cookie = `${tokenName}=${JSON.stringify(data)};path=/;max-age=${maxAge}`;
        } else {
            document.cookie = `${tokenName}=null;path=/;expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
        }

        store.set(data);
    }

    return {
        store,
        set: setData
    };
}

export default createPersistentStore;