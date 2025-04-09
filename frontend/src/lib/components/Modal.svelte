<script lang="ts">
    import { clickOutside } from "$lib/directives/clickOutside.svelte";
    import { fade } from "svelte/transition";
    import { twMerge } from "tailwind-merge";

    let { isOpen = $bindable(), class: modalClass = "", children} = $props();
    const closeModal = () => isOpen = false;
</script>

{#if isOpen}
    <div
        in:fade={{ duration: 50 }}
        class={`fixed top-0 left-0 z-30 h-screen w-screen grid place-items-center bg-neutral-800/60 backdrop-blur-sm`}
    >
        <div class={twMerge(`relative z-40 grid place-items-center px-8 py-4 bg-white rounded-md`, modalClass)} use:clickOutside={() => closeModal()}>
            <!-- Close modal button -->
            <button class={`absolute top-2 right-2`} aria-label={`Close Modal`} onclick={() => closeModal()}>
                <svg class={`h-[24px] aspect-square fill-none stroke-current stroke-2 lucide lucide-x-icon lucide-x`} viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round" xmlns="http://www.w3.org/2000/svg">
                    <path d="M18 6 6 18"/>
                    <path d="m6 6 12 12"/>
                </svg>
            </button>

            {@render children()}
        </div>
    </div>
{/if}