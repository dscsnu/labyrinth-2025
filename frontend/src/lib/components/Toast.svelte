<script lang="ts">
    import type { IToast } from "$lib/stores/ToastStore";
    import { fly } from "svelte/transition";
    import { twMerge } from "tailwind-merge";

    let { message, type, toastClass }: IToast = $props();

    const toastStyles = {
        default: 'apply-card',
        success: 'bg-green-600 text-amber-50',
        warning: 'bg-amber-400 text-neutral-900',
        danger: 'bg-red-600 text-amber-50',
    }
</script>

<div
    transition:fly={{ y: '-200%' }}
    class={twMerge(`z-[100] fixed top-2 left-1/2 -translate-x-1/2 w-[90%] flex justify-center items-center gap-2 py-2 px-4 overflow-hidden rounded-lg ${toastStyles[type!]}`, toastClass)}
>
    {#if type === 'default'}
        <svelte:element this={null} />
    {:else if type === 'success'}
        {@render successSvg()}
    {:else if type === 'warning'}
        {@render warningSvg()}
    {:else}
        {@render dangerSvg()}
    {/if}

    <p>{message}</p>
</div>

{#snippet successSvg()}
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-badge-check">
        <path d="M3.85 8.62a4 4 0 0 1 4.78-4.77 4 4 0 0 1 6.74 0 4 4 0 0 1 4.78 4.78 4 4 0 0 1 0 6.74 4 4 0 0 1-4.77 4.78 4 4 0 0 1-6.75 0 4 4 0 0 1-4.78-4.77 4 4 0 0 1 0-6.76Z"/>
        <path d="m9 12 2 2 4-4"/>
    </svg>
{/snippet}

{#snippet warningSvg()}
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-circle-alert">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/>
    </svg>
{/snippet}

{#snippet dangerSvg()}
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-triangle-alert">
        <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/>
        <path d="M12 9v4"/><path d="M12 17h.01"/>
    </svg>
{/snippet}