<script lang="ts">
    import { onMount } from 'svelte';
    import { fetchWithAuth } from '$lib/utils/fetchWithAuth';
    import { team as teamStore, setPlayerReadyState, setTeam, type TeamData } from '$lib/stores/TeamStore';
    import { addToast } from '$lib/stores/ToastStore';
    import { goto } from '$app/navigation';
    import { LoadingStore } from '$lib/stores/LoadingStore';

    let teamId = $state("810699"); // Using $state for reactive variables
    let team = $state<any>(null);
    let teamData = $derived($teamStore); // Using $derived instead of reactive statement

    const fetchTeamDetails = async () => {
        if (!teamId) {
            addToast({ message: 'No team ID found.', type: 'warning' });
            goto('/');
            return;
        }

        LoadingStore.set(true);
        try {
            const res = await fetchWithAuth(`api/team?team_id=${teamId}`, {
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

            team = await res.json();
            console.log('Fetched team:', team);
            
            if (team) {
                setTeam({
                    id: teamId,
                    name: team.name,
                    is_ready: team.is_ready ?? false,
                    members: team.members ?? []
                });
            }
        } catch (err) {
            console.error('Fetch error:', err);
            addToast({ message: 'Unexpected error while fetching team.', type: 'danger' });
        } finally {
            LoadingStore.set(false);
        }
    };

    // Using the new lifecycle syntax
    $effect(() => {
        fetchTeamDetails();
    });

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

    {#if teamData}
        <h2 class="text-lg">Team: {teamData.name}</h2>
        
        {#if teamData.is_ready}
            <p class="text-green-500 font-bold mt-4">Everyone is ready!</p>
        {/if}

        <ul class="list-disc list-inside mt-2">
            {#each teamData.members as member}
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