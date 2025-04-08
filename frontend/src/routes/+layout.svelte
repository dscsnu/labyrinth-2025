<script lang="ts">
	import '../app.css';
    import { onMount } from 'svelte';
    import { invalidateAll } from '$app/navigation';
    import { device, setDevice } from '$lib/stores/DeviceStore';
    import { SupaStore, UserStore } from '$lib/stores/SupabaseStore';
    import MobileOnly from '$lib/components/MobileOnly.svelte';
    import { PUBLIC_ENVIRONMENT } from '$env/static/public';
	
	let { data, children } = $props();
	let { supabase, session, user } = $derived(data);

	onMount(() => {
		const {
			data: { subscription }
		} = supabase.auth.onAuthStateChange((event, newSession) => {
			if (newSession?.expires_at !== session?.expires_at) {
				invalidateAll();
			}

			if (event === 'SIGNED_OUT') {
				invalidateAll();
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