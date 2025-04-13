<script lang="ts">
    import { fetchWithAuth } from '$lib/utils/fetchWithAuth';
    import { team, setPlayerReadyState, setTeam } from '$lib/stores/TeamStore';
    import { addToast } from '$lib/stores/ToastStore';
    import { goto } from '$app/navigation';
    import { LoadingStore } from '$lib/stores/LoadingStore';
    import { onMount } from 'svelte';

    const fetchTeamDetails = async () => {
        if (!$team) {
            addToast({ message: 'No team found.', type: 'warning' });
            goto('/');
            return;
        }

        LoadingStore.set(true);
        try {
            const res = await fetchWithAuth(`api/team?team_id=${$team.id}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (!res.ok) {
                const errorData = await res.json().catch(() => ({}));
                const message = errorData.message || 'Failed to fetch team details.';
                addToast({ message, type: 'danger' });
                goto('/');
                return;
            }

            const data = await res.json();

            if (data) {
                setTeam({
                    id: data.id,
                    name: data.name,
                    is_ready: data.is_ready ?? false,
                    members: data.members ?? []
                });
            }
        } catch (err) {
            console.error('Fetch error:', err);
            addToast({ message: 'Unexpected error while fetching team.', type: 'danger' });
        } finally {
            LoadingStore.set(false);
        }
    };

    onMount(() => {
        fetchTeamDetails();
    })

    const togglePlayerReadyState = (playerId: string, currentIsReady: boolean) => {
        setPlayerReadyState(playerId, !currentIsReady);
        addToast({
            message: `Player status ${!currentIsReady ? 'Ready' : 'Not Ready'} updated!`,
            type: 'success'
        });
    };
</script>

<main class="p-8">
    <h1 class="text-2xl font-bold mb-4">Team Details</h1>

    {#if $team}
        <h2 class="text-lg">Team: {$team.name}</h2>

        {#if $team.is_ready}
            <p class="text-green-500 font-bold mt-4">Everyone is ready!</p>
        {/if}

        <ul class="list-disc list-inside mt-2">
            {#each $team.members as member}
                <li class="flex items-center space-x-2">
                    <span>{member.name}</span>
                    <span>{member.is_ready ? '✅ Ready' : '❌ Not Ready'}</span>
                    <button
                        onclick={() => togglePlayerReadyState(member.id, member.is_ready)}
                        class="ml-2 px-2 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-400"
                    >
                        {member.is_ready ? 'Mark Not Ready' : 'Mark Ready'}
                    </button>
                </li>
            {/each}
        </ul>
    {:else}
        <p class="text-gray-500">Loading team info...</p>
    {/if}
</main>