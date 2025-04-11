<script>
    import { goto } from "$app/navigation";
    import { fetchWithAuth } from "$lib/utils/fetchWithAuth";
    import { setTeam } from "$lib/stores/TeamStore";
    
    let isCreating = false;
    let loading = false;
    let teamName = '';
    let teamCode = '';
    let error = '';
    
    const toggleMode = () => {
        isCreating = !isCreating;
        error = '';
    };
    
    const createTeam = async () => {
        if (!teamName.trim()) {
            error = 'Please enter a team name';
            return;
        }
        
        loading = true;
        try {
            // Using the backend structure which expects "team_name"
            const res = await fetchWithAuth('http://localhost:3100/api/createteam', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ 
                    team_name: teamName 
                }),
            });
            
            if (!res.ok) {
                const errorData = await res.json().catch(() => ({}));
                error = errorData.message || 'Failed to create team. Please try again.';
                return;
            }
            const data = await res.json();
            
            setTeam(data.team_id);
            goto('/game');
            
        } catch (err) {
            console.error('Error creating team:', err);
            error = 'An unexpected error occurred. Please try again.';
        } finally {
            loading = false;
        }
    };
    
    const joinTeam = async () => {
        if (!teamCode.trim()) {
            error = 'Please enter a team code';
            return;
        }
        
        loading = true;
        try {
            const res = await fetchWithAuth('http://localhost:3100/api/updateteam', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ 
                    team_id: teamCode 
                }),
            });
            console.log(res);
            
            if (!res.ok) {
                if (res.status === 500 && res.statusText.includes("team is full")) {
                    error = 'This team is already full.';
                } else {
                    const errorData = await res.json().catch(() => ({}));
                    error = errorData.message || 'Failed to join team. Please check the code and try again.';
                }
                return;
            }
            
            // The backend handler doesn't return a response body, 
            // so we'll just use the team code as the ID
            setTeam(teamCode);
            
            // Redirect to game page
            goto('/game');
            
        } catch (err) {
            console.error('Error joining team:', err);
            error = 'Failed to join team. Please check the code and try again.';
        } finally {
            loading = false;
        }
    };
</script>


<div class="min-h-screen min-w-screen flex items-center justify-center bg-gray-900 text-white p-4">
    <div class="max-w-md w-full bg-gray-800 rounded-lg shadow-xl p-8">
        <h1 class="text-3xl font-bold text-center mb-6">Welcome to Labyrinth</h1>
        
        <div class="flex justify-center space-x-4 mb-8">
            <button 
                class={`px-6 py-2 rounded-full font-medium transition-colors ${!isCreating ? 'bg-purple-600 text-white' : 'bg-gray-700 text-gray-300'}`} 
                on:click={() => toggleMode()}>
                Join Team
            </button>
            <button 
                class={`px-6 py-2 rounded-full font-medium transition-colors ${isCreating ? 'bg-purple-600 text-white' : 'bg-gray-700 text-gray-300'}`}
                on:click={() => toggleMode()}>
                Create Team
            </button>
        </div>

        {#if error}
            <div class="bg-red-900/50 text-red-200 p-3 rounded-md mb-4 text-sm">
                {error}
            </div>
        {/if}

        {#if isCreating}
            <!-- Create Team Form -->
            <div class="space-y-4">
                <div>
                    <label for="teamName" class="block text-sm font-medium text-gray-300 mb-1">Team Name</label>
                    <input 
                        type="text" 
                        id="teamName" 
                        bind:value={teamName}
                        class="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                        placeholder="Enter your team name" />
                </div>

                <button 
                    on:click={createTeam}
                    disabled={loading}
                    class="w-full bg-purple-600 hover:bg-purple-700 text-white py-3 rounded-md font-medium flex items-center justify-center">
                    {#if loading}
                        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                    {/if}
                    Create Team
                </button>
            </div>
        {:else}
            <!-- Join Team Form -->
            <div class="space-y-4">
                <div>
                    <label for="teamCode" class="block text-sm font-medium text-gray-300 mb-1">Team Code</label>
                    <input 
                        type="text" 
                        id="teamCode" 
                        bind:value={teamCode}
                        class="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                        placeholder="Enter team code (e.g., ABC123)" />
                </div>

                <button 
                    on:click={joinTeam}
                    disabled={loading}
                    class="w-full bg-purple-600 hover:bg-purple-700 text-white py-3 rounded-md font-medium flex items-center justify-center">
                    {#if loading}
                        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                    {/if}
                    Join Team
                </button>
            </div>
        {/if}

        <div class="mt-8 text-center text-sm text-gray-400">
            {#if isCreating}
                Already have a team? <button class="text-purple-400 hover:underline" on:click={toggleMode}>Join existing team</button>
            {:else}
                Need a new team? <button class="text-purple-400 hover:underline" on:click={toggleMode}>Create one</button>
            {/if}
        </div>
    </div>
</div>