import { nanoid } from "nanoid";
import { writable } from "svelte/store";

export interface IToast {
    id?: string;
    message: string;
    timeout?: number;
    type?: 'default' | 'success' | 'warning' | 'danger';
    toastClass?: string;
}

export const ToastStore = writable<IToast[]>([]);

export const addToast = (toast: IToast) => {
    const id  = nanoid();

    ToastStore.update(all => {
        all.forEach(existingToast => {
            if (existingToast.id) {
                dismissToast(existingToast.id);
            }
        });
        return [{ ...toast, id }];
    });

    if (toast.timeout) setTimeout(() => dismissToast(id), toast.timeout);
    else setTimeout(() => dismissToast(id), 5000);
}

export const dismissToast = (id: string) => {
    ToastStore.update((all) => all.filter((t) => t.id !== id));
};