# Another Frontend :(

## Documentation

### fetchWithAuth
Automatically append jwt to outgoing request.
```ts
import { fetchWithAuth } from '$lib/utils/fetchWithAuth.js';

const testFetch = async () => {
    const response = await fetchWithAuth('https://github.com/dscsnu/labyrinth-2025');
    const data = await response.json();

    /* ---snip--- */
}
```


### clickOutside
Used to run some function when something is clicked outside the given element.
```svelte
<script lang="ts">
    import { clickOutside } from "$lib/directive/clickOutside.svelte";
    let isOpen: boolean = $state(false);
</script>

<main class={`h-screen w-screen grid place-items-center`}>
    <button class={`border-2 px-4 py-2`} onclick={() => isOpen = true} use:clickOutside={() => isOpen = false}>
        UseDirective
    </button>

    {#if isOpen}
        Component is open
    {/if}
</main>
```