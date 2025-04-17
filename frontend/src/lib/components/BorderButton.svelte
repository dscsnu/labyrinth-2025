<script lang="ts">
    import { twMerge } from "tailwind-merge";

    let {
        disabled = $bindable(false),
        class: buttonClass = "",
        borderStyle = "",
        onclick = () => {},
        children
    } = $props();
</script>

<div class={`relative h-fit w-fit`}>
    <button
        class={twMerge(
            'bg-white text-black px-6 py-3 font-raleway font-medium rounded-lg',
            buttonClass
        )}
        onclick={() => onclick()}
    >
        {@render children()}
    </button>

    <div class={twMerge(
        `-z-10 absolute top-[6px] left-[6px] h-full w-full rounded-lg border-2 border-white`,
        borderStyle
    )}></div>
</div>