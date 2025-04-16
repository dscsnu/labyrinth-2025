<script lang="ts">
    import { PUBLIC_BACKEND_URL, PUBLIC_STREAM_URL } from "$env/static/public";
    import { JwtTokenStore } from "$lib/stores/JwtTokenStore";
    import { TeamStore } from "$lib/stores/TeamStore";
    import { onMount } from "svelte";
    import { get } from "svelte/store";

    let eventSource: EventSource | null = $state(null);

    const cleanup = () => {
        if (eventSource) {
            eventSource.close();
            eventSource = null;
        }
    }

    const connect = () => {
        console.log('Attempting Connection');
        const cleanedUrl = PUBLIC_BACKEND_URL.replace(/\/+$/, '');
        const cleanedEndpoint = PUBLIC_STREAM_URL.replace(/^\/+/, '');
        const jwt = get(JwtTokenStore);

        eventSource = new EventSource(`${cleanedUrl}/${cleanedEndpoint}?team_id=${$TeamStore?.id}`);

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
