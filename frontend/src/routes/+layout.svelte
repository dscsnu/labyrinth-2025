<script lang="ts">
	import '../app.css';
    import { onMount } from 'svelte';
    import { invalidateAll } from '$app/navigation';
    import { device, setDevice } from '$lib/stores/DeviceStore';
    import { SupaStore, UserStore } from '$lib/stores/SupabaseStore';
    import MobileOnly from '$lib/components/MobileOnly.svelte';
    import { PUBLIC_ENVIRONMENT } from '$env/static/public';
    import { setToken, token } from '$lib/stores/TokenStore';

	let { data, children } = $props();
	let { supabase, session, user } = $derived(data);

	onMount(() => {
		const {
			data: { subscription }
		} = supabase.auth.onAuthStateChange((event, newSession) => {
			if (newSession?.expires_at !== session?.expires_at || event === 'SIGNED_OUT') {
				setToken(null);
				invalidateAll();
			}

			if (['SIGNED_IN', 'TOKEN_REFRESHED'].includes(event) && session?.access_token) {
				setToken(session.access_token);
			}
		});

		SupaStore.set(supabase);
		UserStore.set(user);

		if (window.matchMedia('(max-width: 767px)').matches) setDevice('mobile');
        else setDevice('desktop');

		return () => subscription.unsubscribe();
	})
</script>

<!-- handle device resize -->
<svelte:window onresize={() => setDevice(window.matchMedia('(max-width: 767px').matches ? 'mobile' : 'desktop')}/>

{#if $device === 'mobile' || PUBLIC_ENVIRONMENT === 'development'}
	{@render children()}
{:else}
	<MobileOnly />
{/if}