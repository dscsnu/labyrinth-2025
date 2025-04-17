<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';

    let { data } = $props();
    let { supabase } = $derived(data);

    let displaySplash: boolean = $state(true);
    onMount(() => {
        const timeout = setTimeout(() => displaySplash = false, 1000);
        return () => clearTimeout(timeout);
    })

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
</script>

<main class={`h-screen w-screen grid place-items-center`}>
    <div class={`w-[80%] flex flex-col justify-center items-center gap-4`}>
        <p class={`font-philosopher font-bold text-2xl text-white`}>
            Sign in to unlock your
            <span class={`text-labyrinth-blue`}>tragic hero</span>
            act
        </p>

        <button
            onclick={() => handleSignIn()}
            class={`flex justify-center items-center gap-2 px-5 py-4 rounded-lg border-2 border-neutral-400 bg-neutral-800 active:bg-neutral-700`}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class={`h-[24px] aspect-square`} viewBox="0 0 24 24">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                <path d="M1 1h22v22H1z" fill="none"/>
            </svg>
            <p class={`text-lg text-neutral-400`}>Sign In With Google</p>
        </button>
    </div>
</main>

{#if displaySplash}
    <img
        transition:fade
        class={`z-[100] fixed top-0 left-0 h-screen w-screen`}
        fetchpriority={`high`}
        src={`/assets/images/splash.webp`}
        alt={`SplashScreen`}
    />
{/if}