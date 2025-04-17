<script lang="ts">
    import { goto } from "$app/navigation";
    import BorderButton from "$lib/components/BorderButton.svelte";
    import { validateInput, ValidationOptions } from "$lib/directives/validateInput.svelte";
    import { LoadingStore } from "$lib/stores/LoadingStore";
    import { setTeam, type ITeamData } from "$lib/stores/TeamStore";
    import { addToast } from "$lib/stores/ToastStore";
    import { fetchWithAuth } from "$lib/utils/fetchWithAuth";

    const { data } = $props();
    const { user, supabase } = $derived(data);

    type PageState = 'create' | 'join';

    let pageState: PageState = $state('join');
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

<main class={`h-screen w-screen flex flex-col justify-center gap-2 pt-16 pb-8 px-4`}>
    <hgroup>
        <h1 class={`font-neue-machina font-black text-3xl`}>Registration</h1>
        <p class={`text-sm`}>Team up in groups of 4 and dive head first into defeating 7 evil exes!</p>
    </hgroup>

    <div class={`flex flex-col justify-center items-center mt-4`}>
        <div class={`flex w-full justify-between items-center gap-1 p-1 bg-neutral-800 rounded-full text-lg text-black font-medium`}>
            <button
                onclick={() => pageState = 'join'}
                class={`w-1/2 flex justify-center items-center py-2 px-2 rounded-full transition-colors duration-300 ${pageState === 'join' ? 'text-black bg-white' : 'text-white'}`}
            >
                Join
            </button>

            <button
                onclick={() => pageState = 'create'}
                class={`w-1/2 flex justify-center items-center py-2 px-2 rounded-full transition-colors duration-300 ${pageState !== 'join' ? 'text-black bg-white' : 'text-white'}`}
            >
                Create
            </button>
        </div>
    </div>

    <div class={`flex flex-col justify-center items-center gap-2 mt-2`}>
        {#if pageState === 'join'}
            <input
                id={`team_id`}
                name={`team_id`}
                class={`w-[80%] py-4 px-4 text-lg bg-neutral-800 border-2 border-neutral-600 rounded-lg mb-2`}
                use:validateInput={{
                    maxLength: 6,
                    allowed: [ValidationOptions.NUMERIC]
                }}
                bind:value={teamId}
                placeholder={`Team Code`}
            />

            <BorderButton onclick={() => joinTeam()} class={`px-8 active:bg-neutral-400 flex justify-center items-center gap-2`}>
                <p class={`font-semibold`}>Join Team</p>
                <svg xmlns="http://www.w3.org/2000/svg" class={`h-[18px] aspect-square fill-none stroke-3 stroke-current lucide lucide-user-plus-icon lucide-user-plus`} viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                    <line x1="19" x2="19" y1="8" y2="14"/>
                    <line x1="22" x2="16" y1="11" y2="11"/>
                </svg>
            </BorderButton>
        {:else}
            <input
                id={`team_id`}
                name={`team_id`}
                class={`w-[80%] py-4 px-4 text-lg bg-neutral-800 border-2 border-neutral-600 rounded-lg mb-2`}
                use:validateInput={{
                    maxLength: 18,
                    allowed: [ValidationOptions.ALL]
                }}
                bind:value={teamName}
                placeholder={`Team Name`}
            />

            <BorderButton onclick={() => createTeam()} class={`px-8 active:bg-neutral-400 flex justify-center items-center gap-1`}>
                <p class={`font-semibold`}>Create Team</p>
                <svg xmlns="http://www.w3.org/2000/svg" class={`h-[18px] aspect-square fill-none stroke-current stroke-3 lucide lucide-plus-icon lucide-plus`} viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M5 12h14"/>
                    <path d="M12 5v14"/>
                </svg>
            </BorderButton>
        {/if}
    </div>


    <div class={`flex flex-col justify-center items-center text-center gap-2 mt-16`}>
        <div class={`border-t-2 border-white/20 h-[0.5px] w-full`}></div>
        <div class={`mt-1`}>
            <p>Logged in with email</p>
            <p>{user!.email}</p>
        </div>

        <button class={`active:bg-neutral-700 border-2 border-labyrinth-red bg-neutral-800/60 transition-colors duration-200 text-labyrinth-neutral-200 px-4 py-2 rounded-lg`}>
            Sign Out
        </button>
    </div>
</main>

<svelte:head>
    <title>Labyrinth | Team Registration</title>
    <meta name="description" content="Create or join a team and dive into Labyrinth 2025!">
</svelte:head>