import type { Action } from "svelte/action"

export const clickOutside: Action<HTMLElement, (event: MouseEvent) => void> = (node: HTMLElement, callback) => {
    const handleClick = (event: MouseEvent) => {
        if (!node.contains(event.target as Node)) {
            callback(event);
        }
    }

    $effect(() => {
        window.addEventListener('click', handleClick, true);
        return () => window.removeEventListener('click', handleClick, true);
    })
}