<script lang="ts">
    import StreamListener from "$lib/components/StreamListener.svelte";

    let { children } = $props();
</script>

{@render children()}

<StreamListener />