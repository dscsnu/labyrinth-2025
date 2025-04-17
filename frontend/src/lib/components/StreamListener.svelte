<script lang="ts">
    import { PUBLIC_BACKEND_URL } from "$env/static/public";
    import { TeamStore } from "$lib/stores/TeamStore";
    import { onMount } from "svelte";

    let eventSource: EventSource | null = $state(null);
    let isConnected: boolean = $state(false);

    const cleanup = () => {
        if (eventSource) {
            alert('closing event source');
            isConnected = false;
            eventSource.close();
            eventSource = null;
        }
    }

    const connect = () => {
        if (isConnected) return;

        console.log('Attempting Connection');
        const cleanedUrl = PUBLIC_BACKEND_URL.replace(/\/+$/, '');
        const params = new URLSearchParams({ team_id: $TeamStore!.id })
        eventSource = new EventSource(`${cleanedUrl}/api/eventlistener?${params.toString()}`);

        eventSource.onopen = () => {
            isConnected = true;
        }

        eventSource.onmessage = (event) => {
            alert(event.data);
        }

        eventSource.onerror = (error) => {
            console.error(`Stream error: ${error}`);
            cleanup();
        }
    }

    onMount(() => {
        const unsubscribe = TeamStore.subscribe(current => {
            if (current) {
                connect();
            } else {
                cleanup();
            }
        })

        return () => unsubscribe();
    })
</script>
