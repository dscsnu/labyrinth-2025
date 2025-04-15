<script lang="ts">
    import { goto } from "$app/navigation";
    import { validateInput, ValidationOptions } from "$lib/directives/validateInput.svelte";
    import { LoadingStore } from "$lib/stores/LoadingStore";
    import { UserStore } from "$lib/stores/SupabaseStore";
    import { setTeam, type ITeamData } from "$lib/stores/TeamStore";
    import { addToast } from "$lib/stores/ToastStore";
    import { fetchWithAuth } from "$lib/utils/fetchWithAuth";

    const { data } = $props();
    const { user, supabase } = $derived(data);

    type PageState = 'create' | 'join';

    let pageState: PageState = $state('create');
    let teamName: string = $state('');
    let teamId: string = $state('');

    const createTeam = async () => {
        if (!teamName.trim()) {
            addToast({
                message: 'Please enter a team name',
                type: 'warning',
            });
            return;
        }

        LoadingStore.set(true);
        try {
            const response = await fetchWithAuth('api/team/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    team_name: teamName
                }),
            });

            if (response.status === 200) {
                const data = await response.json();

                if (data.success) {
                    const payload = data.payload;
                    const teamData: ITeamData = {
                        id: payload.team_id,
                        name: teamName,
                        members: [{
                            id: user?.id!,
                            name: user?.user_metadata.full_name!,
                            email: user?.email!,
                            isReady: false
                        }],
                    };

                    setTeam(teamData);
                    goto('/team');
                } else {
                    addToast({
                        message: data.message,
                        type: 'warning'
                    });
                }
            }
        } catch (e) {
            console.error(`Error Creating Team > ${e}`);
            addToast({
                message: 'An unexpected error occured. Please contact helpers.',
                type: 'danger',
            });
        } finally {
            LoadingStore.set(false);
        }
    }

    const joinTeam = async () => {
        if (!teamId.trim() || teamId.trim().length !== 6) {
            addToast({
                message: 'Please enter a valid team code.',
                type: 'warning',
            });
            return;
        }

        LoadingStore.set(true)
        try {
            const response = await fetchWithAuth('api/team/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ team_id: teamId }),
            });

            if (response.status === 200) {
                const data = await response.json()

                if (data.success) {
                    const payload = data.payload;
                    const teamData: ITeamData = {
                        id: payload.id,
                        name: payload.name,
                        members: payload.members.map((m: any) => ({
                            id: m.id,
                            name: m.name,
                            email: m.email,
                            isReady: m.isReady
                        })),
                    };

                    setTeam(teamData);
                    goto('/team');
                } else {
                    addToast({
                        message: data.message,
                        type: 'warning'
                    });
                }
            }
        } catch (e) {
            console.error(`Error Joining Team > ${e}`);
            addToast({
                message: 'An unexpected error occured. Please contact helpers.',
                type: 'danger',
            });
        } finally {
            LoadingStore.set(false);
        }
    }

    const handleSignOut = () => {
        supabase.auth.signOut();
    }
</script>

<main class={`h-screen w-screen flex flex-col justify-center items-center py-4`}>
    <div class={`w-[90%] flex flex-col p-4 border-2 rounded-lg`}>
        <div class={`h-fit w-full gap-3 flex justify-center items-center p-4`}>
            <button class={`px-4 py-2 rounded-lg border-2 ${pageState === 'create' && 'border-green-500'}`} onmousedown={() => pageState = 'create'}>Create Team</button>
            <button class={`px-4 py-2 rounded-lg border-2 ${pageState === 'join' && 'border-green-500'}`} onmousedown={() => pageState = 'join'}>Join Team</button>
        </div>

        {#if pageState === 'create'}
            <div class={`flex flex-col gap-4`}>
                <label for={`team_name`}>Team Name</label>
                <input
                    type={`text`}
                    id={`team_name`}
                    bind:value={teamName}
                    class={`w-full border-2 p-2 rounded-lg`}
                />

                <button onclick={() => createTeam()} class={`w-full border-2 p-2 rounded-lg`}>
                    Create Team
                </button>
            </div>
        {:else if pageState === 'join'}
            <div class={`flex flex-col gap-4`}>
                <label for={`team_name`}>Team Code</label>
                <input
                    type={`text`}
                    id={`team_name`}
                    bind:value={teamId}
                    use:validateInput={{
                        allowed: [ValidationOptions.NUMERIC],
                        maxLength: 6
                    }}
                    class={`w-full border-2 p-2 rounded-lg`}
                    placeholder="e.g, 123456"
                />

                <button onclick={() => joinTeam()} class={`w-full border-2 p-2 rounded-lg`}>
                    Join Team
                </button>
            </div>
        {/if}
    </div>


    <div class={`flex flex-col`}>
        <p>{$UserStore?.email}</p>
        <button onclick={() => handleSignOut()} class={`border-2 px-4 py-2 rounded-lg`}>
            Sign Out
        </button>
    </div>
</main>