<script lang="ts">
    import { onMount } from 'svelte';
    import { token } from '$lib/stores/TokenStore';
    import { goto } from '$app/navigation';
    let { data } = $props();
    let { supabase, user } = $derived(data);

    const handleSignIn = async () => {
        await supabase.auth.signInWithOAuth({
            provider: 'google',
            options: {
                redirectTo: `${window.location.origin}/api/auth/callback`,
                queryParams: {
                    access_type: 'offline',
                    prompt: 'consent',
                },
            }
        })
    }
    onMount(() => {
        if (token) {
            goto('/newplayer');  
        }
    });
</script>

<main class={`h-screen w-screen flex flex-col justify-center items-center`}>
    {user?.email}

    <button onclick={() => handleSignIn()} class={`border-2 px-4 py-2 rounded-lg hover:bg-neutral-200`}>
        Sign In
    </button>
</main>
