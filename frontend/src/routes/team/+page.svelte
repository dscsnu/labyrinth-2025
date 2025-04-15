<script lang="ts">
    import { LoadingStore } from "$lib/stores/LoadingStore";
    import { getPlayerReadyState, setPlayerReadyState, TeamStore } from "$lib/stores/TeamStore";
    import { addToast } from "$lib/stores/ToastStore";
    import { fetchWithAuth } from "$lib/utils/fetchWithAuth";

    const { data } = $props();
    const { user } = $derived(data);

    const handleToggleReady = async () => {
        LoadingStore.set(true)
        try {
            const currentState = getPlayerReadyState(user!.id);

            const response = await fetchWithAuth('api/user/status', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_status: !currentState
                }),
            });

            if (response.status === 200) {
                const data = await response.json();

                if (data.success) {
                    setPlayerReadyState(user!.id, !currentState);
                } else {
                    addToast({
                        message: data.message,
                        type: 'warning'
                    });
                }
            }
        } catch (e) {
            console.error(`Error Toggling Ready State > ${e}`);
            addToast({
                message: 'An unexpected error occured. Please contact helpers.',
                type: 'danger',
            });
        } finally {
            LoadingStore.set(false);
        }
    }
</script>

<main class={`h-screen w-screen flex flex-col justify-center items-center p-4`}>
    <div class={`h-full w-full flex flex-col gap-2 px-2 py-2 border-2 rounded-lg`}>
        <hgroup>
            <h2>{$TeamStore?.name}</h2>
            <p>{$TeamStore?.id}</p>
            <p>COUNTDOWN HERE</p>
        </hgroup>

        <ul class={`flex flex-col gap-2`}>
            {#each $TeamStore?.members! as member (member.id)}
                <li class={`border-2 px-2`}>
                    <p>{member.name}</p>
                    <p>{member.email}</p>
                    <p>
                        {#if member.isReady}
                            Ready
                        {:else}
                            Not Ready
                        {/if}
                    </p>
                </li>
            {/each}
        </ul>

        <button onclick={() => handleToggleReady()} class={`border-2 rounded-lg py-4`}>
            {#if getPlayerReadyState(user!.id)}
                Unready
            {:else}
                Ready
            {/if}
        </button>

        <button class={`border-2 rounded-lg py-4`}>
            Leave Team
        </button>
    </div>
</main>