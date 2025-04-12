<script lang="ts">
    import { goto } from "$app/navigation";
    import { fetchWithAuth } from "$lib/utils/fetchWithAuth";
    import { setTeam, type TeamData } from "$lib/stores/TeamStore";
    import { addToast } from "$lib/stores/ToastStore";
    import {
        validateInput,
        ValidationOptions,
    } from "$lib/directives/validateInput.svelte";
    import { LoadingStore } from "$lib/stores/LoadingStore";
    import { SupaStore, UserStore } from "$lib/stores/SupabaseStore";

    let isCreating: boolean = $state(false);
    let loading: boolean = $state(false);
    let teamName: string = $state("");
    let teamCode: string = $state("");

    const toggleMode = () => {
        isCreating = !isCreating;
    };

    const createTeam = async () => {
        if (!teamName.trim()) {
            addToast({
                message: "Please enter a team name",
                type: "warning",
            });
            return;
        }

        LoadingStore.set(true);
        try {
            // Using the backend structure which expects "team_name"
            const res = await fetchWithAuth("api/createteam", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    team_name: teamName,
                }),
            });

            if (!res.ok) {
                const errorData = await res.json().catch(() => ({}));
                const message =
                    errorData.message ||
                    "Failed to create team. Please try again.";
                addToast({
                    message,
                    type: "danger",
                });
                return;
            }
            const data = await res.json();

            // Create a TeamData object with the response data
            const teamData: TeamData = {
                id: data.team_id,
                name: teamName, // Use the name provided by the user
                is_ready: false, // New teams start as not ready
                members: [], // Initialize with an empty array of members
            };

            // Save the complete team data
            setTeam(teamData);

            goto("/team");
        } catch (err) {
            console.error("Error creating team:", err);
            addToast({
                message: "An unexpected error occured. Please try again.",
                type: "danger",
            });
        } finally {
            LoadingStore.set(false);
        }
    };

    const joinTeam = async () => {
        if (!teamCode.trim()) {
            addToast({
                message: "Please enter a team code.",
                type: "warning",
            });
            return;
        }

        LoadingStore.set(true);
        try {
            const res = await fetchWithAuth("api/updateteam", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    team_id: teamCode,
                }),
            });

            if (!res.ok) {
                if (
                    res.status === 500 &&
                    res.statusText.includes("team is full")
                ) {
                    addToast({
                        message: "This team is already full.",
                        type: "warning",
                    });
                } else {
                    const errorData = await res.json().catch(() => ({}));
                    const message =
                        errorData.message ||
                        "Failed to join team. Please check the code and try again.";
                    addToast({
                        message,
                        type: "danger",
                    });
                }
                return;
            }

            // Parse the team data from the response
            const teamResponse = await res.json();
            // Create TeamData object using the response from backend
            const teamData: TeamData = {
                id: teamResponse.id || teamCode,
                name: teamResponse.name || "Team " + teamCode.substring(0, 4),
                is_ready: teamResponse.is_ready || false,
                members: teamResponse.members || [],
            };

            // Save the complete team data
            setTeam(teamData);

            // Redirect to team page
            goto("/team");
        } catch (err) {
            console.error("Error joining team:", err);
            addToast({
                message:
                    "Failed to join team. Please check the code and try again.",
                type: "danger",
            });
        } finally {
            LoadingStore.set(false);
        }
    };

    const handleSignOut = () => $SupaStore.auth.signOut();
</script>

<main class={`h-screen w-screen flex items-center justify-center px-8`}>
    <div class={`rounded-lg border-2 p-8`}>
        <h1 class={`text-3xl font-bold text-center mb-6`}>
            Welcome to Labyrinth
        </h1>

        <div class={`flex justify-center space-x-4 mb-8`}>
            <button
                class={`px-6 py-2 rounded-full border-2 ${!isCreating ? "border-green-500" : "border-black"}`}
                onclick={() => toggleMode()}
            >
                Join Team
            </button>
            <button
                class={`px-6 py-2 rounded-full border-2 ${isCreating ? "border-green-500" : "border-black"}`}
                onclick={() => toggleMode()}
            >
                Create Team
            </button>
        </div>

        {#if isCreating}
            <!-- Create Team Form -->
            <div class={`space-y-4`}>
                <div>
                    <label for={`teamName`}>Team Name</label>
                    <input
                        type={`text`}
                        id={`teamName`}
                        bind:value={teamName}
                        class={`w-full border-2 border-black p-2 rounded-md`}
                        placeholder={`Enter your team name`}
                    />
                </div>

                <button
                    onclick={createTeam}
                    disabled={loading}
                    class={`w-full text-black border-2 py-3 rounded-md font-medium flex items-center justify-center`}
                >
                    Create Team
                </button>
            </div>
        {:else}
            <!-- Join Team Form -->
            <div class={`space-y-4`}>
                <div>
                    <label for={`teamCode`}>Team Code</label>
                    <input
                        type={`text`}
                        id={`teamCode`}
                        use:validateInput={{
                            allowed: [ValidationOptions.NUMERIC],
                            maxLength: 6,
                        }}
                        bind:value={teamCode}
                        class={`w-full border-2 p-2 rounded-md`}
                        placeholder={`Enter team code (e.g., 123456)`}
                    />
                </div>

                <button
                    onclick={joinTeam}
                    disabled={loading}
                    class={`w-full border-2 py-3 rounded-md font-medium flex items-center justify-center`}
                >
                    Join Team
                </button>
            </div>
        {/if}

        <div class={`mt-8 text-center text-sm text-gray-400`}>
            {#if isCreating}
                Already have a team? <button
                    class={`text-purple-400 hover:underline`}
                    onclick={toggleMode}>Join existing team</button
                >
            {:else}
                Need a new team? <button
                    class={`text-purple-400 hover:underline`}
                    onclick={toggleMode}>Create one</button
                >
            {/if}
        </div>

        <div class={`flex flex-col justify-center items-center`}>
            <p>{$UserStore?.email}</p>
            <button onclick={handleSignOut} class={`border-2 rounded-lg p-2`}>
                Sign Out
            </button>
        </div>
    </div>
</main>