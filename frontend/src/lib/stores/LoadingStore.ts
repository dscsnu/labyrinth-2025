import { writable } from "svelte/store";
export const LoadingStore = writable<boolean>(false);