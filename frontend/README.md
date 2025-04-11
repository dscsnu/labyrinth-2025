# Another Frontend :(

## Documentation

### fetchWithAuth
Automatically append jwt to outgoing request.
```ts
import { fetchWithAuth } from '$lib/utils/fetchWithAuth.js';

const testFetch = async () => {
    const response = await fetchWithAuth('/api/team');
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

### Modal
```svelte
<script lang="ts">
    import Modal from "$lib/components/Modal.svelte";
    let isOpen: boolean = $state(false);
    const toggleModal = () => isOpen = !isOpen;
</script>

<main class={`h-screen w-screen flex flex-col justify-center items-center`}>
    <button onclick={() => toggleModal()}>
        toggle Modal
    </button>

    <!-- anything passed in class with override default styles -->
    <Modal bind:isOpen class={`bg-red-500`}>
        Hello Cro
    </Modal>
</main>
```

### validateInput
Used on Input elements to validate any user input
```svelte
<script lang="ts>
    import { validateInput, ValidationOptions } from "$lib/directives/validateInput.svelte;
</script>

<!--
    Any user will not be able to input non numerice characters into the following.
    The length of the input will also be capped to 6.
 -->
<input
    type={`text`}
    id={`some_id`}
    use:validateInput={{
        allowed: [ValidationOptions.NUMERIC],
        maxLength: 6
    }}
    bind:value
/>
```
