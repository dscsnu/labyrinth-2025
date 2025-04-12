<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchWithAuth } from '$lib/utils/fetchWithAuth';
	import { team as teamStore, setPlayerReadyState, setTeam, type TeamData } from '$lib/stores/TeamStore';
	import { addToast } from '$lib/stores/ToastStore';
	import { goto } from '$app/navigation';
	import { LoadingStore } from '$lib/stores/LoadingStore';
  
	let teamId: string = "810699"; // Hardcoded team ID as requested
	let team: any = null;
	let teamData: TeamData | null;
  
	// Subscribe to team store changes
	$: teamData = $teamStore;
  
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
		
		// Update the TeamStore with the fetched data
		if (team) {
		  setTeam({
			id: teamId,
			name: team.name,
			is_ready: team.is_ready || false,
			members: team.members || []
		  });
		}
	  } catch (err) {
		console.error('Fetch error:', err);
		addToast({ message: 'Unexpected error while fetching team.', type: 'danger' });
	  } finally {
		LoadingStore.set(false);
	  }
	};
  
	onMount(fetchTeamDetails);
  
	// Function to toggle player's ready state
	const togglePlayerReadyState = (playerId: string, currentIsReady: boolean) => {
	  // Toggle player ready state in the store
	  setPlayerReadyState(playerId, !currentIsReady);
	  
	  // You could add an API call here to update the backend if needed
	  // For now just show a toast message
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
			<span>
			  {member.is_ready ? '✅ Ready' : '❌ Not Ready'}
			</span>
			<button
			  on:click={() => togglePlayerReadyState(member.id, member.is_ready)}
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